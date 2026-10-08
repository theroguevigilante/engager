// Package migrations embeds the SQL schema migrations so both the API and the
// seed command can bring a database up to date.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
