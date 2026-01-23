package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"betrayal-web/internal"
	"github.com/labstack/echo/v4"
	pgxmock "github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

// NOTE: multi-import statement fix for Go linting and IDEs. Only single import block used.

func TestHealthHandler_DBConnected_OK(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/health/db", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	origConn := internal.Conn
	defer func() { internal.Conn = origConn }()

	mock, err := pgxmock.NewConn()
	assert.NoError(t, err)
	defer mock.Close(context.Background())
	mock.ExpectPing().WillReturnError(nil)
	internal.Conn = mock // pgxmock.PgxConnIface satisfies Pingable

	if assert.NoError(t, HealthHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "\"status\":\"ok\"")
		assert.NotContains(t, rec.Body.String(), "postgresql://")
		assert.NotContains(t, rec.Body.String(), "DATABASE_URL")
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHealthHandler_DBConnected_PingError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/health/db", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	origConn := internal.Conn
	defer func() { internal.Conn = origConn }()

	mock, err := pgxmock.NewConn()
	assert.NoError(t, err)
	defer mock.Close(context.Background())
	mock.ExpectPing().WillReturnError(errors.New("simulated DB ping failure"))
	internal.Conn = mock

	if assert.NoError(t, HealthHandler(c)) {
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
		assert.Contains(t, rec.Body.String(), "db error")
		assert.Contains(t, rec.Body.String(), "simulated DB ping failure")
		assert.NotContains(t, rec.Body.String(), "postgresql://")
		assert.NotContains(t, rec.Body.String(), "DATABASE_URL")
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHealthHandler_DBNotConnected(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/health/db", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	origConn := internal.Conn
	defer func() { internal.Conn = origConn }()
	internal.Conn = nil

	if assert.NoError(t, HealthHandler(c)) {
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
		assert.Contains(t, rec.Body.String(), "db not connected")
		assert.NotContains(t, rec.Body.String(), "postgresql://")
		assert.NotContains(t, rec.Body.String(), "DATABASE_URL")
	}
}
