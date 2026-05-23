/*
 * output.go
 * capp-parse
 *
 * CLI output functions for parse results.
 *
 * Created by David Richardson on Sunday, April 12, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */
package capp

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/enquora-net/capp-parse/internal/core"
)

const (
	iconOK  = core.IconOK
	iconErr = core.IconErr
)

// MakeEmitFn returns the per-file emit function for the given output target.
func MakeEmitFn(isTTY bool) EmitFn {
	if isTTY {
		return emitTTYResult
	}
	return emitJSONResult
}

// EmitResult writes a single-file parse result to stdout.
func EmitResult(path string, r ParseResult, f Format, isTTY bool) {
	if !isTTY {
		emitJSONResult(path, r)
		return
	}

	if r.HasError {
		for _, e := range r.Errors {
			fmt.Fprintf(os.Stderr, "%s  %s: %s\n", iconErr, path, e.String())
		}
		return
	}

	switch f {
	case FormatSexp, FormatDefault:
		fmt.Println(r.PrettySexp)
	default:
		fmt.Fprintf(os.Stdout, "%s  %s (%d nodes)\n", iconOK, path, r.NodeCount)
	}
}

// EmitWalkSummary writes the walk summary and performance report to stdout.
func EmitWalkSummary(s WalkSummary) {
	fmt.Fprintf(os.Stdout, "\n%d files  %d OK  %d error files  %d bytes  %s\n",
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

// EmitDebugProfile writes the per-file timing summary from a debug walk.
func EmitDebugProfile(p *DebugProfile) {
	if p == nil {
		return
	}

	var maxElapsed time.Duration
	for _, e := range p.Events {
		if e.Elapsed > maxElapsed {
			maxElapsed = e.Elapsed
		}
	}

	fmt.Println("\n" + core.Divider)
	fmt.Println("PARSE PROFILE")
	fmt.Println(core.Divider)

	for _, e := range p.Events {
		icon := iconOK
		if !e.OK {
			icon = iconErr
		}
		fmt.Printf("%s  %8s  %s\n", icon, e.Elapsed.Round(time.Microsecond), e.Path)
	}

	fmt.Println(core.Divider)
	fmt.Printf("   files: %d\n", len(p.Events))
	fmt.Printf("   total: %s\n", p.Total.Round(time.Microsecond))
	if len(p.Events) > 0 {
		fmt.Printf("    mean: %s\n", (p.Total / time.Duration(len(p.Events))).Round(time.Microsecond))
	}
	fmt.Printf("     max: %s\n", maxElapsed.Round(time.Microsecond))
}

func emitTTYResult(path string, r ParseResult) {
	if r.HasError {
		for _, e := range r.Errors {
			fmt.Fprintf(os.Stderr, "%s  %s: %s\n", iconErr, path, e.String())
		}
		return
	}
	fmt.Fprintf(os.Stdout, "%s  %s (%d nodes)\n", iconOK, path, r.NodeCount)
}

func emitJSONResult(path string, r ParseResult) {
	type jsonError struct {
		Row     uint32 `json:"row"`
		Column  uint32 `json:"column"`
		Message string `json:"message"`
	}
	type jsonResult struct {
		Path      string      `json:"path"`
		OK        bool        `json:"ok"`
		NodeCount int         `json:"node_count,omitempty"`
		Errors    []jsonError `json:"errors,omitempty"`
	}

	var jsonErrs []jsonError
	for _, e := range r.Errors {
		jsonErrs = append(jsonErrs, jsonError{
			Row:     e.Row,
			Column:  e.Column,
			Message: e.Message,
		})
	}

	line, _ := json.Marshal(jsonResult{
		Path:      path,
		OK:        !r.HasError,
		NodeCount: r.NodeCount,
		Errors:    jsonErrs,
	})
	fmt.Fprintf(os.Stdout, "%s\n", line)
}
