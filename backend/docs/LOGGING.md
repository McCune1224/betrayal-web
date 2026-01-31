# Backend Logging Documentation

## Overview

The backend now uses Go 1.21's structured logging package `slog` for comprehensive, configurable logging. This provides excellent observability during development and production.

## Features

### 1. **Structured Logging**
- JSON format for production (LOG_FORMAT=json)
- Text format for development (LOG_FORMAT=text)
- Consistent field naming with snake_case
- Request correlation IDs for tracing

### 2. **Configurable Log Levels**
- `LOG_LEVEL=debug` - Detailed debugging information
- `LOG_LEVEL=info` - Normal operations (default)
- `LOG_LEVEL=warn` - Warnings and recoverable errors
- `LOG_LEVEL=error` - Errors requiring attention

### 3. **HTTP Request Logging**
- Automatic logging of all HTTP requests
- Method, path, status code, duration
- Request ID for tracing
- Client IP address

### 4. **WebSocket Logging**
- Connection/disconnection events
- Message handling (type, size)
- Ping/pong keepalive tracking
- Slow client detection

### 5. **Game Event Logging**
- Room creation and deletion
- Player join/leave events
- Action submissions
- Phase transitions

## Configuration

Add these to your `.env` file:

```bash
# Logging configuration
LOG_LEVEL=info          # debug, info, warn, error
LOG_FORMAT=text         # text (dev) or json (prod)
```

## Log Output Examples

### Development Mode (text)
```
2026-01-30 19:17:49 INFO hub_started
2026-01-30 19:17:49 INFO room_created room_code=IGJ26L host_id=host-123 total_rooms=1
2026-01-30 19:17:49 INFO client_registered room_code=TEST01 player_id=player-123 room_client_count=1 total_rooms=1
2026-01-30 19:17:49 WARN join_attempt_to_nonexistent_room room_code=FAKE99 player_name=Bob
```

### Production Mode (JSON)
```json
{"time":"2026-01-30T19:17:49.618-05:00","level":"INFO","msg":"hub_started"}
{"time":"2026-01-30T19:17:49.618-05:00","level":"INFO","msg":"room_created","room_code":"IGJ26L","host_id":"host-123","total_rooms":1}
{"time":"2026-01-30T19:17:49.618-05:00","level":"WARN","msg":"join_attempt_to_nonexistent_room","room_code":"FAKE99","player_name":"Bob"}
```

## Key Log Events

### Application Lifecycle
- `application starting` - Server startup
- `server starting` - HTTP server startup
- `hub_started` - WebSocket hub initialized
- `database_connection_established` - DB connected
- `server shutting down` - Graceful shutdown

### Room Operations
- `room_created` - New room created
- `room_deleted_empty` - Empty room cleanup
- `player_joined` - Player joined room
- `join_attempt_to_nonexistent_room` - Failed join attempt

### WebSocket Operations
- `websocket_client_connected` - New WebSocket connection
- `client_registered` - Client added to hub
- `client_disconnected` - Connection closed
- `client_unregistered` - Client removed from hub
- `message_received` - Incoming message
- `action_submitted` - Player action

### HTTP Requests
- `http_request_started` - Request begins
- `http_request` - Request completed (info/warn/error based on status)

### Errors
- `websocket_upgrade_failed` - WebSocket upgrade error
- `create_room_failed` - Room creation error
- `join_room_failed` - Room join error
- `health_check_failed` - Health check error

## Using the Logger in Code

### Basic Logging
```go
import "backend/internal/logging"

// Get the global logger
logger := logging.Logger()
logger.Info("something happened", "key", value)
logger.Error("something failed", "error", err)
```

### With Context (Request ID)
```go
// In HTTP handlers
ctx := c.Request().Context()
logger := logging.WithContext(ctx)
logger.Info("handling request") // includes request_id automatically
```

### WebSocket Logging
```go
// In client handlers
logger := logging.WSLogger(roomCode, playerID)
logger.Info("player_action", "action", actionType)
```

### Room Logging
```go
logger := logging.RoomLogger(roomCode)
logger.Info("phase_changed", "from", oldPhase, "to", newPhase)
```

## Debugging Tips

1. **Set LOG_LEVEL=debug** for verbose output during development
2. **Check request_id** to trace a single request through all components
3. **Filter by room_code** to debug specific game rooms
4. **Filter by player_id** to debug specific player issues
5. **Use jq** with JSON logs in production: `docker logs container | jq '. | select(.room_code=="ABC123")'`

## Performance Considerations

- Debug logs are compiled out in production builds when possible
- Use `logger.Debug()` for high-frequency events (every message)
- Use `logger.Info()` for significant events (connections, state changes)
- Buffered channels prevent slow logging from blocking operations
- JSON encoding has minimal overhead with slog

## Testing

Tests use the same logging configuration but output to test logs:
```bash
cd backend
go test ./... -v  # See logs interleaved with test output
```

## Future Enhancements

Potential additions:
- OpenTelemetry integration for distributed tracing
- Log sampling for high-frequency events
- Separate access logs vs application logs
- Log rotation configuration
- Metrics export (Prometheus)
