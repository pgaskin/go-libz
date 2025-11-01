package embed

import (
	_ "embed"

	"github.com/pgaskin/go-libz"
)

//go:generate docker build --platform amd64 --pull --no-cache --progress plain --output . .
//go:embed libz.wasm
var binary []byte

func init() {
	libz.Binary = binary
}
