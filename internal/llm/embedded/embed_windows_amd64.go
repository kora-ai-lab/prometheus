//go:build windows && amd64

package embedded

import _ "embed"

//go:embed llama-server-windows-amd64.exe
var ServerBinary []byte

const ServerName = "llama-server-windows-amd64.exe"
