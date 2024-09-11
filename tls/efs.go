package tlsecurity

import (
	"embed"
)

//go:embed *.pem
var Files embed.FS
