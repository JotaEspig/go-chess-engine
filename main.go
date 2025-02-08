package main

import (
	"bufio"
	"fmt"
	"gce/pkg/chess"
	"gce/pkg/engine"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sort"
	"strings"
)

func main() {
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	var depth uint
	fmt.Print("Depth: ")
	fmt.Scanln(&depth)

	var fen string
	fmt.Print("FEN: ")
	// Read entire line not using Scanln
	reader := bufio.NewReader(os.Stdin)
	fen, _ = reader.ReadString('\n')
	fen = strings.ReplaceAll(fen, "\n", "")
	fen = strings.ReplaceAll(fen, "\r", "")
	fen = strings.TrimSpace(fen)
	if fen == "" {
		fen = chess.DefaultStartFen
	}

	b := chess.FenToBoard(fen)
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
			nodesCountch := make(chan chess.Move)

			analysisReport := engine.AnalysisByDepth(copyBoard, depth, returnCh, nodesCountch)
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
			nodesCountch := make(chan chess.Move)

			analysisReport := engine.AnalysisByDepth(copyBoard, depth, returnCh, nodesCountch)
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
