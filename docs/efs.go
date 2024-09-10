package docs

import (
	"embed"
)

//go:embed *.sql
var DocsFiles embed.FS
