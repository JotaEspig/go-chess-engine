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

	b := chess.FenToBoard("rn1qkbnr/1pp2ppp/p2p4/4N3/2B1P3/2N5/PPPP1PPP/R1BbK2R w KQkq - 0 6")
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
		} else if moveNotation == "eval" {
			bestBoard, evaluation := engine.AnalysisByDepth(*b, 5)
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
