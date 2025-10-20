package main

import (
	"os"

	"github.com/tmc/langchaingo/textsplitter"
)

type Chunk struct {
	Text string
	Page int
}

func splitText(path string) ([]Chunk, error) {
	f, err := os.ReadFile(path)
	ts := textsplitter.NewRecursiveCharacter()
	if err != nil {
		return nil, err
	}
	chunks, _ := ts.SplitText(string(f))
	var out []Chunk
	for _, c := range chunks {
		out = append(out, Chunk{Text: c})
	}
	return out, nil
}
