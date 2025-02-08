package engine

import "gce/pkg/chess"

type AnalysisReport struct {
	BestBoard  chess.Board
	Evaluation float64
}
