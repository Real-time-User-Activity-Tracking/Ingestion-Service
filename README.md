# Ingestion Service

A simple Go HTTP service that receives user activity events from a frontend application and displays them in the console.

## What it does

- Receives POST requests with user activity data
- Displays the received events in the terminal
- Sends back a success response to the frontend
- Handles CORS for cross-origin requests

## Quick Start

1. **Run the service:**
   ```bash
   go run main.go
   ```

2. **Service will start on:** `http://localhost:9094`

3. **Available endpoints:**
   - `GET /health` - Health check
   - `GET /api/v1/status` - Service status
   - `POST /api/v1/events/track` - Receive events

## Example Event

```json
{
  "event_type": "button_click",
  "timestamp": "2024-01-15T10:30:00Z",
  "user_id": "user123",
  "session_id": "session456",
  "page_url": "https://example.com",
  "event_data": {
    "action": "add_to_cart"
  },
  "client_info": {
    "user_agent": "Mozilla/5.0...",
    "screen_resolution": "1920x1080",
    "language": "en-US"
  }
}
```

## Project Structure

```
ingestion-service/
├── main.go             # Entry point
├── config/             # Configuration
├── handlers/           # Request handlers
├── middleware/         # CORS middleware
├── models/             # Data structures
└── router/             # Route definitions
```

## CORS

The service handles CORS automatically:
- Allows requests from any origin
- Handles preflight OPTIONS requests
- Supports POST requests with JSON data

## Testing

Test with curl:
```bash
curl -X POST http://localhost:9094/api/v1/events/track \
  -H "Content-Type: application/json" \
  -d '{"event_type":"test","user_id":"test123","session_id":"test456","page_url":"http://test.com","event_data":{},"client_info":{"user_agent":"test","screen_resolution":"1920x1080","language":"en-US"}}'
```

## Environment Variables

- `PORT` - Server port (default: 9094)
- `HOST` - Server host (default: 0.0.0.0)