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
	// b := chess.FenToBoard("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
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
			copyBoard := b.Copy()
			returnCh := make(chan engine.AnalysisReport)
			nodesCountch := make(chan struct{})

			analysisReport := engine.AnalysisByDepth(copyBoard, 5, returnCh, nodesCountch)
			bestBoard := analysisReport.BestBoard
			move, moveNotation := engine.GetEngineMove(copyBoard, bestBoard)
			fmt.Println(moveNotation)
			b.MakeMove(move)

			close(returnCh)
			close(nodesCountch)
			continue
		} else if moveNotation == "eval" {
			copyBoard := b.Copy()
			returnCh := make(chan engine.AnalysisReport)
			nodesCountch := make(chan struct{})

			analysisReport := engine.AnalysisByDepth(copyBoard, 5, returnCh, nodesCountch)
			bestBoard, evaluation := analysisReport.BestBoard, analysisReport.Evaluation
			fmt.Printf("Evaluation: %.2f\n", evaluation)
			fmt.Println(engine.GetEngineLine(&copyBoard, &bestBoard))

			close(returnCh)
			close(nodesCountch)
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
