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
	// fmt.Println(engine.Perft(b, depth))
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
			returnCh := make(chan engine.AnalysisReport)
			nodesCountch := make(chan struct{})

			analysisReport := engine.AnalysisByDepth(b, depth, returnCh, nodesCountch)
			bestBoard, evaluation := analysisReport.BestBoard, analysisReport.Evaluation
			fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")
			fmt.Printf("Evaluation: %.2f\n", evaluation)
			fmt.Println(analysisReport.GetEngineLine())
			fmt.Println("Final position after engine line:")
			fmt.Println(bestBoard.VisualBoard().String())
			fmt.Println("=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=")

			close(returnCh)
			close(nodesCountch)
			continue
		} else if moveNotation == "list" {
			allLegalMoves := engine.MoveSlice(b.AllLegalMoves())
			sort.Sort(allLegalMoves)
			copyBoard := &chess.Board{}
			*copyBoard = *b
			for _, move := range allLegalMoves {
				notation := copyBoard.MoveToNotation(move)
				fmt.Printf("%d -> %s\n", engine.MoveSortingScore(move), notation)
				copyBoard.MakeLegalMove(move)
			}
			fmt.Println("Total moves:", len(allLegalMoves))
			continue
		}

		move, err := b.ParseMove(moveNotation)
		if err != nil {
			fmt.Println("Invalid move")
			continue
		}
		b.MakePseudoLegalMove(move)
		if !b.IsValidPosition() {
			fmt.Println("Illegal move")
			b.UndoMove()
			continue
		}
	}

	fmt.Println(b.GetMoveListInNotation())
}
