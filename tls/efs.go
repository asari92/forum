package tlsecurity

import (
	"embed"
)

//go:embed *.pem
var TlsFiles embed.FS
