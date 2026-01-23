package db

import (
	"betrayal-web/internal/db/sqlc"
	"context"
	"database/sql"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func init() {
	_ = godotenv.Load("../../.env") // Load env vars early for all tests; ignore error if .env missing
}

func testDBURL() string {
	if url := os.Getenv("TEST_DATABASE_URL"); url != "" {
		return url
	}
	// Fallback to main DATABASE_URL
	return os.Getenv("DATABASE_URL")
}

func newTestDB(t *testing.T) (*sql.DB, func()) {
	dsn := testDBURL()
	require.NotEmpty(t, dsn, "TEST_DATABASE_URL or DATABASE_URL must be set for DB tests")
	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err, "Failed to open db with dsn: %s", dsn)
	cleanup := func() { db.Close() }
	return db, cleanup
}

func TestDBConnection(t *testing.T) {
	dbObj, cleanup := newTestDB(t)
	defer cleanup()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := dbObj.PingContext(ctx)
	require.NoError(t, err, "Database ping failed, check TEST_DATABASE_URL/DATABASE_URL and DB migrations: %v", err)
}

func TestCreateAndGetRoom(t *testing.T) {
	dbObj, cleanup := newTestDB(t)
	defer cleanup()
	queries := sqlc.New(dbObj)
	ctx := context.Background()

	// Use a truly unique roomCode for every run to avoid key collisions
	roomCode := uuid.NewString()[0:10] // must be <= 10 chars for VARCHAR(10)
	hostID := uuid.New()
	phase := "LOBBY"
	room, err := queries.CreateRoom(ctx, sqlc.CreateRoomParams{
		Code:   roomCode,
		HostID: hostID,
		Phase:  phase,
	})
	require.NoError(t, err, "CreateRoom failed for code %s, hostID %s: %v", roomCode, hostID, err)
	require.Equal(t, roomCode, room.Code, "Room code mismatch (expected %s, got %s)", roomCode, room.Code)
	require.Equal(t, hostID, room.HostID, "HostID mismatch (expected %s, got %s)", hostID, room.HostID)
	require.Equal(t, phase, room.Phase, "Phase mismatch (expected %s, got %s)", phase, room.Phase)

	// Now try to fetch it
	r2, err := queries.GetRoomByCode(ctx, roomCode)
	require.NoError(t, err, "GetRoomByCode failed for code %s: %v", roomCode, err)
	assert.Equal(t, room.Code, r2.Code, "Fetched room code mismatch (expected %s, got %s)", room.Code, r2.Code)
	assert.Equal(t, room.HostID, r2.HostID, "Fetched room HostID mismatch (expected %s, got %s)", room.HostID, r2.HostID)
	assert.Equal(t, room.Phase, r2.Phase, "Fetched room phase mismatch (expected %s, got %s)", room.Phase, r2.Phase)
}
