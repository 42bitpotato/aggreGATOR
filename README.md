# aggreGATOR 🐊
A small RSS → PostgreSQL CLI written in Go. 


## Prerequisites

- Go (1.21+ recommended)
- PostgreSQL

## Install

```bash
go install github.com/42bitpotato/aggreGATOR/cmd/gator@latest

```

This builds a `gator` binary in `$GOPATH/bin` (also in `$HOME/go/bin` if `GOPATH` unset). Ensure that directory is on your `PATH`.

## Configure

Create `~/.gatorconfig.json`:

```json
{"db_url":"postgres://<username>:@localhost:5432/gator?sslmode=disable","current_user_name":"","date_format":"2006-01-02 15:04:05"}
```

- `db_url`: your Postgres DSN.
- `current_user_name`: set after you register or login.
- `date_format`: output format used when displaying dates.

## Run

General form:

```bash
gator <command> [args]
```

Common commands:

- `register <user>`  
  Create a new user.

- `login <user>`  
  Set `<user>` as the current user (updates `~/.gatorconfig.json`).

- `addfeed <title> <url>`  
  Add a feed to the catalog.

- `follow <url>`  
  Follow a feed as the current user.

- `following`  
  List feeds the current user follows.

- `feeds`  
  List all feeds and who added them.

- `browse [n]`  
  Show latest posts from followed feeds. Defaults to a small number if `n` omitted.

- `agg <duration>`  
  Continuously aggregate posts on an interval, e.g. `gator agg 5m`.

- `users`  
  List users.

- `unfollow <url>`  
  Unfollow a feed.

- `reset`  
  Reset users (destructive).

## Notes

- You need a running Postgres instance reachable by `db_url`.
- If you change the config while the app is running, restart your `gator` process.
- If you see “permission denied” on install, ensure `$GOBIN` or `$GOPATH/bin` is on your `PATH`.

## License

MIT

If your import path for the CLI differs, replace the `go install` line with the exact module path to your `cmd/gator`.