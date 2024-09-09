package docs

import (
	"embed"
)

//go:embed *.sql
var FilesDocs embed.FS
