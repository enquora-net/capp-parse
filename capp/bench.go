/*
 * capp/bench.go
 * cappuccino
 *
 * Created by David Richardson on Sunday, April 12, 2026.
 * Copyright (c) 2026 David Richardson. All rights reserved.
 * All responsibility for usage rests with the user.
 * The author bears no liability for damages arising from usage,
 * whether direct or indirect.
 */

package capp

import "time"

type benchAccumulator struct {
	enabled bool
	events  []BenchEvent
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
	b.events = append(b.events, BenchEvent{Path: path, Elapsed: elapsed, Bytes: bytes})
	b.total += elapsed
	b.bytes += bytes
}

func (b *benchAccumulator) report(wall time.Duration) *BenchReport {
	if !b.enabled {
		return nil
	}
	return &BenchReport{
		WallTime:   wall,
		ParseTotal: b.total,
		Files:      len(b.events),
		Bytes:      b.bytes,
		Events:     b.events,
	}
}
