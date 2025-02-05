package main

import (
	"fmt"
	"gce/pkg/chess"
	"gce/pkg/engine"
	"strings"
)

func main() {
	b := chess.NewDefaultBoard()
	for {
		vb := b.VisualBoard()
		fmt.Println(vb.String())

		evaluation := engine.EvaluatePosition(*b)
		fmt.Printf("Evaluation: %.2f\n", evaluation)

		if b.IsMated() {
			fmt.Println("CHECKMATE BABY!!!!")
			break
		}

		var moveNot string
		fmt.Print("Move: ")
		fmt.Scanln(&moveNot)
		if strings.TrimSpace(moveNot) == "q" {
			break
		}

		move, err := b.ParseMove(moveNot)
		if err != nil {
			fmt.Println("Invalid move")
			continue
		}
		if !b.MakeMove(move) {
			fmt.Println("Illegal move")
		}
	}

	fmt.Println(b.GetMoveListInNotation())
}
