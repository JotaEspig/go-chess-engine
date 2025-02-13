package engine

import "gce/pkg/chess"

type NodesPerMove struct {
	Move  chess.Move
	Nodes uint64
}

func Perft(board *chess.Board, depth uint) (uint64, []NodesPerMove) {
	if depth == 0 {
		return 1, []NodesPerMove{}
	}

	nodes := uint64(0)
	nodesPerMove := []NodesPerMove{}
	// if board.White.Pawns.Board&e4 != 0 && board.Black.Pawns.Board&(17_179_869_184) != 0 {
	// 	fmt.Println("e4")
	// }
	moves := board.AllLegalMoves()
	// if depth == 1 && initialDepth > 2 {
	// 	fmt.Println(board.VisualBoard().String())
	// 	fmt.Println("Moves: ", len(moves))
	// 	allMovesStr := ""
	// 	for _, move := range moves {
	// 		not := board.MoveToNotation(move)
	// 		allMovesStr += not + " "
	// 	}
	// 	fmt.Println(allMovesStr)
	// }
	for _, move := range moves {
		board.MakeLegalMove(move)
		// if depth == 3 {
		// 	fmt.Println(board.VisualBoard().String())
		// }
		_nodes, _ := Perft(board, depth-1)
		nodes += _nodes
		nodesPerMove = append(nodesPerMove, NodesPerMove{Move: move, Nodes: _nodes})
		board.UndoMove()
	}
	return nodes, nodesPerMove
}
