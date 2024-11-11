# PipeTG

A Telegram channel scraper using MTProto with session management and connection pooling.

## What it does

- Manages Telegram MTProto sessions
- Handles auth lifecycle and persistence
- Provides connection pooling and validation
- Scrapes channel messages on intervals

## Quick Start

```bash
# Set up env vars
export TELEGRAM_APP_ID=your_app_id
export TELEGRAM_APP_HASH=your_app_hash
export TELEGRAM_USER_ID=your_user_id
export SCRAPE_WINDOW=5m
export PORT=8080

# Optional: Skip auth flow by providing existing session
export SESSION_CONFIG=your_session_json

# Run it
go run cmd/main.go
```

## Modules

### Session Management
```go
// Inject and persist Telegram sessions
session := NewInjectedSessionStorage(rawSession)
client := NewTelegramClient(uid, appID, appHash, rawSession)

// Managed connections with validation
client.WithContext(ctx, func(ctx context.Context) error {
    conn, err := client.GetTgConnFromCtx(ctx)
    if err != nil {
        return err
    }

    return scraper.Run(ctx)
})
```

### Auth Flow
```go
// Phone auth initiation
client.SendCode(ctx, conn)

// Code verification
client.AuthenticateWithCode(ctx, code, conn)

// Automatic validation
client.Validate(ctx, conn)
```

### Connection Options
```go
// Connection with session
conn := client.Raw() 

// Fresh connection without session
conn := client.RawWithoutSession()

// Echo handler integration
client.WithEchoContext(e, handler)
```

## Architecture

The client wrapper provides:
- Session storage and injection
- Connection pooling and validation
- Context-based connection management
- Echo framework integration

## Development Notes

- Use `WithContext` for automatic validation
- Handle session persistence with `InjectedSessionStorage`
- Implement proper error handling for auth states
- Validate connections before operations
- Optionally inject existing sessions via env