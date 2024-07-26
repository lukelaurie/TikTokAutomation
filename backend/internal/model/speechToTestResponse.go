package model

type Word struct {
	Word            string `json:"Word"`
	Offset          int    `json:"Offset"`
	Duration        int    `json:"Duration"`
	Confidence      float64 `json:"Confidence"`
	AccuracyScore   float64 `json:"AccuracyScore"`
}

type NBest struct {
	Confidence float64 `json:"Confidence"`
	Lexical    string  `json:"Lexical"`
	ITN        string  `json:"ITN"`
	MaskedITN  string  `json:"MaskedITN"`
	Display    string  `json:"Display"`
	Words      []Word  `json:"Words"`
}

type TranscriptResponse struct {
	RecognitionStatus string  `json:"RecognitionStatus"`
	Offset            int     `json:"Offset"`
	Duration          int     `json:"Duration"`
	NBest             []NBest `json:"NBest"`
	DisplayText       string  `json:"DisplayText"`
}