// Package migrations contains DB migration embedded assets.
package migrations

import "embed"

// FS is the embedded filesystem containing migration files.
//
//go:embed *.sql
var FS embed.FS
