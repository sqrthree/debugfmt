// Package debugfmt implements a development-friendly textual handler.
package debugfmt

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/sqrthree/debugfmt/colors"
)

// color function.
type colorFunc func(string) string

// Colors mapping.
var Colors = [...]colorFunc{
	log.DebugLevel: colors.Purple,
	log.InfoLevel:  colors.Blue,
	log.WarnLevel:  colors.Yellow,
	log.ErrorLevel: colors.Magenta,
	log.FatalLevel: colors.Red,
}

// Strings mapping.
var Strings = [...]string{
	log.DebugLevel: "DEBUG",
	log.InfoLevel:  " INFO",
	log.WarnLevel:  " WARN",
	log.ErrorLevel: "ERROR",
	log.FatalLevel: "FATAL",
}

// Handler implementation.
type Handler struct {
	mu     sync.Mutex
	Writer io.Writer
}

// New returns a new handle.
func New(w io.Writer) *Handler {
	return &Handler{
		Writer: w,
	}
}

// HandleLog implements log.Handler.
func (h *Handler) HandleLog(e *log.Entry) error {
	var buf bytes.Buffer

	color := Colors[e.Level]
	level := Strings[e.Level]

	names := e.Fields.Names()

	len := len(names)

	time := formatDateString(e.Timestamp.Local())

	fmt.Fprintf(&buf, "%s %s %s", colors.Gray(time), color(level), color(e.Message))

	if len != 0 {
		fmt.Fprintf(&buf, ": ")
	}

	for _, name := range names {
		fmt.Fprintf(&buf, "%s%s%v ", color(name), "=", e.Fields.Get(name))
	}

	fmt.Fprintln(&buf)

	b := buf.Bytes()

	h.mu.Lock()
	defer h.mu.Unlock()

	h.Writer.Write(b)

	return nil
}

// formatDateString formats t to a human-friendly string.
func formatDateString(t time.Time) string {
	hour, minute, second := t.Clock()

	return fmt.Sprintf("%02v:%02v:%02v", hour, minute, second)
}
