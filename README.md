# UI Backend

Go backend service for the UI Network application. The project is built on Fiber, GORM, and MySQL, and exposes APIs for authentication, user management, orders, products, payouts, referral workflows, and supporting admin operations.

## Stack

- Go
- Fiber
- GORM
- MySQL
- Viper for environment-based configuration

## Configuration

The service reads configuration from `app.env` in the project root.

1. Copy the example file:

```bash
cp app.env.example app.env
```

2. Set the required values:

- `MYSQL_USER`
- `MYSQL_PASSWORD`
- `MYSQL_DB`
- `MYSQL_HOST`
- `MYSQL_PORT`
- `JWT_SECRET`

Optional values:

- `CHEQUE_DRAW_VALUE` defaults to `4000`
- `APP_PORT` defaults to `8000`

## Run

```bash
go run ./cmd
```

The service fails fast when required configuration is missing.

## Notes

- Local secrets are intentionally not committed.
- `node_modules/` and local build artifacts are excluded from version control.
