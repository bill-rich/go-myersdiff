package myersdiff

// This implementation of the Myers Diff Algorythm was written using http://simplygenius.net/Article/DiffTutorial1,
// and https://github.com/cj1128/myers-diff as references.

import (
	"bytes"
)

const (
	ADD    opType = 1
	DELETE        = 2
	NOOP          = 3
)

type opType int

func (op opType) String() string {
	switch op {
	case ADD:
		return "ADD"
	case DELETE:
		return "DEL"
	case NOOP:
		return "NOP"
	}
	return ""
}

// DiffOptions changes the behavior of how the diff is run or outputted.
type DiffOptions struct {
	PrintAdd    bool
	PrintNoOp   bool
	PrintDelete bool
}

type trace struct {
	v      []map[int]int
	srcLen int
	dstLen int
}

func (t *trace) length() int {
	return len(t.v)
}

func (t *trace) append(v map[int]int) {
	t.v = append(t.v, v)
}

// GenerateDiff provides the diff of two string slices. To diff two files, split them by lines and provide each as a
// slice.
func GenerateDiff(src, dst []string, opts *DiffOptions) *bytes.Buffer {
	script := shortestEditScript(src, dst)
	return writeDiff(src, dst, script, opts)
}

// NewOptions returns a default set of options.
func NewOptions() *DiffOptions {
	return &DiffOptions{
		PrintAdd:    true,
		PrintNoOp:   true,
		PrintDelete: true,
	}
}

func writeDiff(src, dst []string, script []opType, opts *DiffOptions) *bytes.Buffer {
	buffer := bytes.Buffer{}
	srcIndex, dstIndex := 0, 0
	for _, op := range script {
		switch op {
		case ADD:
			if opts.PrintAdd {
				buffer.Write([]byte("+ "))
				buffer.Write([]byte(dst[dstIndex]))
				buffer.Write([]byte("\n"))
			}
			dstIndex += 1

		case NOOP:
			if opts.PrintNoOp {
				buffer.Write([]byte("  "))
				buffer.Write([]byte(src[srcIndex]))
				buffer.Write([]byte("\n"))
			}
			srcIndex += 1
			dstIndex += 1

		case DELETE:
			if opts.PrintDelete {
				buffer.Write([]byte("- "))
				buffer.Write([]byte(src[srcIndex]))
				buffer.Write([]byte("\n"))
			}
			srcIndex += 1
		}
	}
	return &buffer
}

func createTrace(src, dst []string) *trace {
	srcLen := len(src)
	dstLen := len(dst)
	maxLen := srcLen + dstLen
	var x, y int
	trace := trace{
		srcLen: len(src),
		dstLen: len(dst),
	}

	for d := 0; d <= maxLen; d++ {
		v := make(map[int]int, d+2)
		trace.append(v)

		// Find the first difference
		if d == 0 {
			firstDiff := 0
			// Keep looking until the files differ
			for len(src) > firstDiff && len(dst) > firstDiff && src[firstDiff] == dst[firstDiff] {
				firstDiff++
			}
			// The first diff is at line firstDiff
			v[0] = firstDiff

			// If firstDiff is the end of the file, there is no diff
			if firstDiff == len(src) && firstDiff == len(dst) {
				return &trace
			}
			continue
		}

		lastV := trace.v[d-1]

		for k := -d; k <= d; k += 2 {
			switch {
			// Go down (insert dest) if at the lowest k-line
			case k == -d:
				x = lastV[k+1]
			// Go down (insert dest) lower k-line x is behind the higher k-line x. This comparison can't be made at the
			// highest k-line (k==d).
			case k != d && lastV[k-1] < lastV[k+1]:
				x = lastV[k+1]
			// Move right (del source) if at the highest k-line (k==d) or if the lower k-line x is further along.
			default:
				x = lastV[k-1] + 1
			}

			// Get y using the slope function
			y = x - k

			// Look for next diff along the diagonal (snake)
			for x < srcLen && y < dstLen && src[x] == dst[y] {
				x, y = x+1, y+1
			}

			// Set the k-line/x intercept
			v[k] = x

			// Reached the end of the source or dest.
			if x == srcLen && y == dstLen {
				return &trace
			}
		}
	}
	return &trace
}

func createScript(trace *trace) []opType {
	var x, y int
	var script []opType

	x = trace.srcLen
	y = trace.dstLen
	var k, prevK, prevX, prevY int

	for d := trace.length() - 1; d > 0; d-- {
		k = x - y
		lastV := trace.v[d-1]

		switch {
		case k == -d:
			prevK = k + 1
		case k != d && lastV[k-1] < lastV[k+1]:
			prevK = k + 1
		default:
			prevK = k - 1
		}

		prevX = lastV[prevK]
		prevY = prevX - prevK

		for x > prevX && y > prevY {
			script = append(script, NOOP)
			x -= 1
			y -= 1
		}

		if x == prevX {
			script = append(script, ADD)
		} else {
			script = append(script, DELETE)
		}
		x, y = prevX, prevY
	}

	if trace.v[0][0] != 0 {
		for i := 0; i < trace.v[0][0]; i++ {
			script = append(script, NOOP)
		}
	}

	return reverse(script)
}

func shortestEditScript(src, dst []string) []opType {
	trace := createTrace(src, dst)
	return createScript(trace)
}

func reverse(s []opType) []opType {
	result := make([]opType, len(s))
	end := len(s) - 1

	for i, v := range s {
		result[end-i] = v
	}

	return result
}
