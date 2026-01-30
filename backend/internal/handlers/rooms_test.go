package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/internal/game"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoomHandler_CreateRoom(t *testing.T) {
	hub := game.NewHub()
	go hub.Run()

	handler := &RoomHandler{Hub: hub}

	t.Run("creates room successfully", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/rooms", strings.NewReader(`{"hostName":"TestHost"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateRoom(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "roomCode")
		assert.Contains(t, rec.Body.String(), "hostId")
	})

	t.Run("handles invalid JSON", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/rooms", strings.NewReader(`invalid json`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateRoom(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "error")
	})
}

func TestRoomHandler_JoinRoom(t *testing.T) {
	hub := game.NewHub()
	go hub.Run()

	handler := &RoomHandler{Hub: hub}

	roomCode := hub.GetRoomManager().CreateRoom("host-123")

	t.Run("joins existing room", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/rooms/"+roomCode+"/join",
			strings.NewReader(`{"playerName":"Alice"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPathValues(echo.PathValues{
			{Name: "code", Value: roomCode},
		})

		err := handler.JoinRoom(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "playerId")
		assert.Contains(t, rec.Body.String(), "phase")
	})

	t.Run("fails to join non-existing room", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/rooms/FAKE99/join",
			strings.NewReader(`{"playerName":"Bob"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPathValues(echo.PathValues{
			{Name: "code", Value: "FAKE99"},
		})

		err := handler.JoinRoom(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "error")
	})

	t.Run("handles invalid JSON", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/rooms/"+roomCode+"/join",
			strings.NewReader(`invalid json`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPathValues(echo.PathValues{
			{Name: "code", Value: roomCode},
		})

		err := handler.JoinRoom(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestHealthHandler(t *testing.T) {
	t.Run("returns unhealthy when DB not initialized", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/health/db", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := HealthHandler(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	})
}
