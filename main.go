package main

import (
	"fmt"
	"gce/pkg/chess"
	"gce/pkg/engine"
	"net/http"
	_ "net/http/pprof"
	"sort"
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

		var moveNotation string
		fmt.Print("Move: ")
		fmt.Scanln(&moveNotation)
		moveNotation = strings.TrimSpace(moveNotation)
		if moveNotation == "q" {
			break
		} else if moveNotation == "move" {
			bestBoard, _ := engine.AnalysisByDepth(*b, 5)
			move, moveNotation := engine.GetEngineMove(*b, bestBoard)
			fmt.Println(moveNotation)
			b.MakeMove(move)
			continue
		} else if moveNotation == "eval" {
			startBoard := b.Copy()
			bestBoard, evaluation := engine.AnalysisByDepth(startBoard, 5)
			fmt.Printf("Evaluation: %.2f\n", evaluation)
			fmt.Println(engine.GetEngineLine(b, &bestBoard))
			continue
		} else if moveNotation == "list" {
			allLegalMoves := engine.MoveSlice(b.AllLegalMoves())
			sort.Sort(allLegalMoves)
			for _, move := range allLegalMoves {
				fmt.Printf("%d -> %s\n", engine.MoveSortingScore(move), move.String())
			}
			continue
		}

		move, err := b.ParseMove(moveNotation)
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
