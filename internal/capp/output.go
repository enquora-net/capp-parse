/*
 * internal/capp/output.go
 * capp-parse
 *
 * Created by David Richardson on Sunday, April 12, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 *
 */
package capp

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// MakeEmitFn returns the per-file emit function for the given output target.
func MakeEmitFn(isTTY bool) EmitFn {
	if isTTY {
		return emitTTYResult
	}
	return emitJSONResult
}

// EmitResult writes a single-file parse result to the appropriate output target.
func EmitResult(path string, r ParseResult, isTTY bool) {
	if isTTY {
		emitTTYResult(path, r)
	} else {
		emitJSONResult(path, r)
	}
}

// EmitWalkSummary writes the walk summary and performance report to stdout.
// Called only when stdout is a TTY.
func EmitWalkSummary(s WalkSummary) {
	fmt.Fprintf(os.Stdout, "\n%d files  %d OK  %d errors  %d bytes  %s\n",
		s.TotalFiles, s.OKFiles, s.ErrorFiles, s.TotalBytes,
		s.Elapsed.Round(time.Millisecond))

	if s.Bench == nil {
		return
	}

	b := s.Bench
	mbps := 0.0
	if b.ParseTotal > 0 {
		mbps = float64(b.Bytes) / (1024 * 1024) / b.ParseTotal.Seconds()
	}

	events := make([]BenchEvent, len(b.Events))
	copy(events, b.Events)
	sort.Slice(events, func(i, j int) bool {
		return events[i].Elapsed > events[j].Elapsed
	})
	limit := 10
	if len(events) < limit {
		limit = len(events)
	}

	fmt.Fprintf(os.Stdout, "\n── performance ──────────────────────────────────\n")
	fmt.Fprintf(os.Stdout, "  wall       %s\n", b.WallTime.Round(time.Microsecond))
	fmt.Fprintf(os.Stdout, "  parse cpu  %s (sum across goroutines)\n", b.ParseTotal.Round(time.Microsecond))
	fmt.Fprintf(os.Stdout, "  files      %d\n", b.Files)
	fmt.Fprintf(os.Stdout, "  bytes      %d  (%.1f KB)\n", b.Bytes, float64(b.Bytes)/1024)
	fmt.Fprintf(os.Stdout, "  throughput %.2f MB/s\n", mbps)
	fmt.Fprintf(os.Stdout, "  slowest %d\n", limit)
	for i := range limit {
		e := events[i]
		fmt.Fprintf(os.Stdout, "    %8s  %s\n", e.Elapsed.Round(time.Microsecond), e.Path)
	}
	fmt.Fprintf(os.Stdout, "─────────────────────────────────────────────────\n")
}

func emitTTYResult(path string, r ParseResult) {
	if r.HasError {
		for _, e := range r.Errors {
			fmt.Fprintf(os.Stderr, "%s: %s\n", path, e)
		}
		return
	}
	fmt.Fprintf(os.Stdout, "%s: OK (%d nodes)\n", path, r.NodeCount)
}

func emitJSONResult(path string, r ParseResult) {
	type jsonResult struct {
		Path      string   `json:"path"`
		OK        bool     `json:"ok"`
		NodeCount int      `json:"node_count,omitempty"`
		Errors    []string `json:"errors,omitempty"`
	}
	line, _ := json.Marshal(jsonResult{
		Path:      path,
		OK:        !r.HasError,
		NodeCount: r.NodeCount,
		Errors:    r.Errors,
	})
	fmt.Fprintf(os.Stdout, "%s\n", line)
}
