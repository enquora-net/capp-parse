/*
 * internal/core/bench.go
 * capp-parse
 *
 * Created by David Richardson on Sunday, April 12, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 */
package core

import (
	"time"

	"github.com/enquora-net/capp-parse/internal/types"
)

type benchAccumulator struct {
	enabled bool
	events  []types.BenchEvent
	total   time.Duration
	bytes   int
}

func newBenchAccumulator(enabled bool) *benchAccumulator {
	return &benchAccumulator{enabled: enabled}
}

func (b *benchAccumulator) record(path string, elapsed time.Duration, bytes int) {
	if !b.enabled {
		return
	}
	b.events = append(b.events, types.BenchEvent{Path: path, Elapsed: elapsed, Bytes: bytes})
	b.total += elapsed
	b.bytes += bytes
}

func (b *benchAccumulator) report(wall time.Duration) *types.BenchReport {
	if !b.enabled {
		return nil
	}
	return &types.BenchReport{
		WallTime:   wall,
		ParseTotal: b.total,
		Files:      len(b.events),
		Bytes:      b.bytes,
		Events:     b.events,
	}
}
