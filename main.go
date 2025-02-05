package main

import (
	"fmt"
	"gce/pkg/engine"
)

func main() {
	b := engine.NewDefaultBoard()
	for {
		vb := b.VisualBoard()
		fmt.Println(vb.String())

		if b.IsMated() {
			fmt.Println("CHECKMATE BABY!!!!")
			break
		}

		var moveNot string
		fmt.Print("Move: ")
		fmt.Scanln(&moveNot)

		move, err := b.ParseMove(moveNot)
		if err != nil {
			fmt.Println("Invalid move")
			continue
		}
		if !b.MakeMove(move) {
			fmt.Println("Invalid move")
		}
	}
}
