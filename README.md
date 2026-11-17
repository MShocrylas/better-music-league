Requires sqlite3
`go tool sqlc generate` generates the db package go code from sql schema and queries. Config in `sqlc.yaml`

`goose -dir migrations/ create <filename> sql` to generate a new migration file
