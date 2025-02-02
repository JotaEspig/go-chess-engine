package main

import (
	"fmt"
	"gce/pkg/engine"
)

func main() {
	engine.Bla()

	b := engine.FenToBoard("r1bq1rk1/pp2bppp/2n2n2/3p4/3p4/2N1P1P1/PP2NPBP/R1BQ1RK1 w - - 0 1")
	originalBoard := b
	vb := b.ToVisualBoard()
	fmt.Println(vb.String())
	fmt.Printf("Evaluation: %.1f\n", engine.Evaluate(b))

	fmt.Println("Knight")
	knightMoves := b.White.Knights.AllPossibleMoves(b, true)
	for _, move := range knightMoves {
		fmt.Println("Is white move: ", b.Ctx.WhiteToMove)
		fmt.Println("Is capture: ", move.IsCapture)
		b.MakeMove(move)
		vb := b.ToVisualBoard()
		fmt.Println(vb.String())
		// ask for one key input
		fmt.Scanln()

		// reset board
		b = originalBoard
	}

	fmt.Println("Rook")
	rookMoves := b.White.Rooks.AllPossibleMoves(b, true)
	for _, move := range rookMoves {
		fmt.Println("Is white move: ", b.Ctx.WhiteToMove)
		fmt.Println("Is capture: ", move.IsCapture)
		b.MakeMove(move)
		vb := b.ToVisualBoard()
		fmt.Println(vb.String())
		// ask for one key input
		fmt.Scanln()

		// reset board
		b = originalBoard
	}

	fmt.Println("Bishop")
	bishopMoves := b.White.Bishops.AllPossibleMoves(b, true)
	for _, move := range bishopMoves {
		fmt.Println("Is white move: ", b.Ctx.WhiteToMove)
		fmt.Println("Is capture: ", move.IsCapture)
		b.MakeMove(move)
		vb := b.ToVisualBoard()
		fmt.Println(vb.String())
		// ask for one key input
		fmt.Scanln()

		// reset board
		b = originalBoard
	}

	fmt.Println("Queen")
	queenMoves := b.White.Queens.AllPossibleMoves(b, true)
	for _, move := range queenMoves {
		fmt.Println("Is white move: ", b.Ctx.WhiteToMove)
		fmt.Println("Is capture: ", move.IsCapture)
		b.MakeMove(move)
		vb := b.ToVisualBoard()
		fmt.Println(vb.String())
		// ask for one key input
		fmt.Scanln()

		// reset board
		b = originalBoard
	}

	fmt.Println("King")
	kingMoves := b.White.King.AllPossibleMoves(b, true)
	for _, move := range kingMoves {
		fmt.Println("Is white move: ", b.Ctx.WhiteToMove)
		fmt.Println("Is capture: ", move.IsCapture)
		b.MakeMove(move)
		vb := b.ToVisualBoard()
		fmt.Println(vb.String())
		// ask for one key input
		fmt.Scanln()

		// reset board
		b = originalBoard
	}
}
