// Copyright © 2018 Chad Netzer <chad.netzer@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"fmt"
	"strings"
	"time"
)

type Progress interface {
	ShowDirsFilesFound()
	Clear()
}

// A simple progress meter while scanning directories and performing linking
type TTYProgress struct {
	lastLineLen     int
	lastTime        time.Time
	lastFPSTime     time.Time
	updateDelay     time.Duration
	updateFPSDelay  time.Duration
	dirFilesCounter int
	counterMin      int
	lastFPS         float64
	lastFPSDiff     float64

	stats   *LinkingStats
	options *Options
}

type DisabledProgress struct{}

// Initialize TTYProgress and return pointer to it
func NewTTYProgress(stats *LinkingStats, options *Options) *TTYProgress {
	now := time.Now()
	return &TTYProgress{
		lastTime:       now,
		lastFPSTime:    now,
		updateDelay:    60 * time.Millisecond,
		updateFPSDelay: 180 * time.Millisecond,
		counterMin:     151, // Prime number makes output more dynamic
		stats:          stats,
		options:        options,
	}
}

// Output a line (without a newline at the end) of progress on directory
// scanning and inode linking (ie. which inodes have identical content and
// matching inode parameters).  Call in the main directory walk/link
// calculation loop.
func (p *TTYProgress) ShowDirsFilesFound() {
	p.dirFilesCounter += 1
	if p.dirFilesCounter < p.counterMin {
		return
	} else {
		p.dirFilesCounter = 0
	}

	now := time.Now()
	timeSinceLast := now.Sub(p.lastTime)
	if timeSinceLast < p.updateDelay {
		return
	}
	p.lastTime = now

	numDirs := p.stats.DirCount
	numFiles := p.stats.FileCount

	duration := now.Sub(p.stats.StartTime)
	durStr := duration.Round(time.Second).String()

	var fps, fpsDiff float64
	timeSinceLastFPS := now.Sub(p.lastFPSTime)
	if timeSinceLastFPS > p.updateFPSDelay {
		fps = float64(numFiles) / duration.Seconds()
		fpsDiff = fps - p.lastFPS

		p.lastFPS = fps
		p.lastFPSDiff = fpsDiff
		p.lastFPSTime = now
	} else {
		fps = p.lastFPS
		fpsDiff = p.lastFPSDiff
	}

	fmtStr := "\r%d files in %d dirs, elapsed time: %s  files/sec: %.0f (%+.0f)"
	s := fmt.Sprintf(fmtStr, numFiles, numDirs, durStr, fps, fpsDiff)
	p.line(s)
}

// Call to erase the progress loop (before 'normal' program post-processing
// output)
func (p *TTYProgress) Clear() {
	p.line("\r")
	p.lastLineLen = 1
	p.line("\r")
}

func (p *TTYProgress) line(s string) {
	numSpaces := p.lastLineLen - len(s)
	p.lastLineLen = len(s)
	if numSpaces > 0 {
		s += strings.Repeat(" ", numSpaces)
	}
	fmt.Print(s)
}

func (p *DisabledProgress) ShowDirsFilesFound() {}
func (p *DisabledProgress) Clear()              {}
