package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"bubbleRouge/pkg/testharness"
)

type Step struct {
	Key   string `json:"key"`
	Ticks int    `json:"ticks"`
}

type Script struct {
	Seed   int64   `json:"seed"`
	Width  int     `json:"width"`
	Height int     `json:"height"`
	DT     float64 `json:"dt"`
	Steps  []Step  `json:"steps"`
}

func main() {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var s Script
	if err := json.Unmarshal(b, &s); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	c := testharness.NewController(testharness.Options{Seed: s.Seed, Width: s.Width, Height: s.Height})
	for _, st := range s.Steps {
		if st.Key != "" {
			c.InjectKey(st.Key)
		}
		n := st.Ticks
		if n <= 0 {
			n = 1
		}
		dt := s.DT
		if dt == 0 {
			dt = 1.0
		}
		c.Tick(n, dt)
	}
	out, err := c.Snapshot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Stdout.Write(out)
}
