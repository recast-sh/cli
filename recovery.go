package cli

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	for i := skip; ; i++ {   // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fmt.Fprintf(buf, "%s:%d ", file, line)
		fmt.Fprintf(buf, "%s\n", function(pc))
	}

	b := buf.Bytes()

	// TODO: This is very very wrong!
	if r := os.Getenv("GOPATH"); r != "" {
		b = bytes.Replace(b, []byte(r), []byte("."), -1)
	}
	if r := os.Getenv("GOROOT"); r != "" {
		b = bytes.Replace(b, []byte(r), []byte("."), -1)
	}
	return b
}

func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	//	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
