package main

import (
	"fmt"
	"gce/pkg/chess"
	"gce/pkg/engine"
	"net/http"
	_ "net/http/pprof"
	"strings"
)

func main() {
	go func() {
		fmt.Println("Starting pprof server")
		http.ListenAndServe("localhost:6060", nil)
	}()

	b := chess.NewDefaultBoard()
	for {
		vb := b.VisualBoard()
		fmt.Println(vb.String())

		if b.IsMated() {
			fmt.Println("CHECKMATE BABY!!!!")
			break
		} else if b.IsDraw() {
			fmt.Println("DRAW")
			break
		}

		var moveNot string
		fmt.Print("Move: ")
		fmt.Scanln(&moveNot)
		if strings.TrimSpace(moveNot) == "q" {
			break
		} else if strings.TrimSpace(moveNot) == "eval" {
			evaluation := engine.AnalysisByDepth(*b, 4)
			fmt.Printf("Evaluation: %.2f\n", evaluation)
			continue
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
