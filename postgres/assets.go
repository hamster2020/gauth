package postgres

import "embed"

//go:embed sql/init.sql
var sqlInit string

//go:embed sql/migration-*.sql
var migrations embed.FS
