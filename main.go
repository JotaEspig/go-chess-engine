package main

import (
	"fmt"
	"gce/pkg/engine"
)

func main() {
	engine.Bla()

	b := engine.NewDefaultBoard()
	vb := b.ToVisualBoard()
	fmt.Println(vb.String())
}
