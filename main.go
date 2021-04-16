//go:generate statik -src=./static -f
package main

import (
	"github.com/bzp2010/webvpn/cmd"
)

func main() {
	cmd.Execute()
}
