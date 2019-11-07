package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

//noinspection GoUnusedConst
const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

var (
	consoleBufPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 100))
		},
	}
)

const (
	consoleDefaultTimeFormat = time.Kitchen
)

// Transforms the input into a formatted string.
type formatter func(interface{}) string

// Parses the JSON input and writes it in an (optionally) colorized, human-friendly format to Out.
type prettyConsoleWriter struct {
	Out        io.Writer
	NoColor    bool
	TimeFormat string
}

// Creates a console writer
func newPrettyConsoleWriter() io.Writer {
	return prettyConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05.000",
	}
}

// Write transforms the JSON input with formatters and appends to w.Out.
func (w prettyConsoleWriter) Write(p []byte) (int, error) {

	var buf = consoleBufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		consoleBufPool.Put(buf)
	}()

	var evt map[string]interface{}
	d := json.NewDecoder(bytes.NewReader(p))
	d.UseNumber()
	if err := d.Decode(&evt); err != nil {
		return 0, errors.Wrapf(err, "cannot decode event: %s", err)
	}

	// Print timestamp
	if s := formatTimestampForConsole(w.TimeFormat, w.NoColor)(evt["time"]); len(s) > 0 {
		buf.WriteString(s)
		buf.WriteByte(' ')
	}

	// Print level
	if s := formatLevelForConsole(w.NoColor)(evt["level"]); len(s) > 0 {
		buf.WriteString(s)
		buf.WriteByte(' ')
	}

	// Print caller
	if s := formatCallerForConsole(w.NoColor)(evt["caller"]); len(s) > 0 {
		buf.WriteString(s)
		buf.WriteByte(' ')
	}

	// Print message
	if s := consoleDefaultFormatMessage(evt["message"]); len(s) > 0 {
		buf.WriteString(s)
	}

	// Collect field names, ignoring the ones we just wrote
	// Fields will be sorted alphabetically
	var fields = make([]string, 0, len(evt))
	for field := range evt {
		if field != zerolog.TimestampFieldName &&
			field != zerolog.LevelFieldName &&
			field != zerolog.CallerFieldName &&
			field != zerolog.MessageFieldName &&
			field != zerolog.ErrorFieldName &&
			field != zerolog.ErrorStackFieldName {

			fields = append(fields, field)
		}
	}
	sort.Strings(fields)

	// Space separator
	if len(fields) > 0 {
		buf.WriteByte(' ')
	}

	// Add fields
	for i, field := range fields {
		buf.WriteString(colorize(fmt.Sprintf("%s=", field), colorCyan, w.NoColor))
		switch fValue := evt[field].(type) {
		case string:
			if needsQuote(fValue) {
				buf.WriteString(fmt.Sprintf("%s", strconv.Quote(fValue)))
			} else {
				buf.WriteString(fmt.Sprintf("%s", fValue))
			}
		case json.Number:
			buf.WriteString(fmt.Sprintf("%s", fValue))
		default:
			if b, err := json.Marshal(fValue); err != nil {
				fmt.Fprintf(buf, colorize("[error: %v]", colorRed, w.NoColor), err)
			} else {
				fmt.Fprint(buf, fmt.Sprintf("%s", b))
			}
		}
		if i < len(fields)-1 { // Skip space for last field
			buf.WriteByte(' ')
		}
	}
	buf.WriteByte('\n')

	// Add error
	if errString, ok := evt[zerolog.ErrorFieldName]; ok {
		if stackTrace, ok := evt[zerolog.ErrorStackFieldName]; ok {
			// If we have a stack, it's been marshalled by our own marshaller, which already prints the error itself
			// So printing only the stack is sufficient
			buf.WriteString(fmt.Sprintf("%s\n", stackTrace))
		} else {
			// No stack - print the error on our own; I don't think this can happen (zerolog should always use our marshaller, no?)
			buf.WriteString(fmt.Sprintf("%s\n", errString))
		}
	}

	// Write to output stream
	_, err := buf.WriteTo(w.Out)
	return len(p), err
}

func needsQuote(s string) bool {
	for i := range s {
		if s[i] < 0x20 || s[i] > 0x7e || s[i] == ' ' || s[i] == '\\' || s[i] == '"' {
			return true
		}
	}
	return false
}

func colorize(s interface{}, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

func formatTimestampForConsole(timeFormat string, noColor bool) formatter {
	if timeFormat == "" {
		timeFormat = consoleDefaultTimeFormat
	}
	return func(i interface{}) string {
		t := "<nil>"
		switch tt := i.(type) {
		case string:
			ts, err := time.Parse(zerolog.TimeFieldFormat, tt)
			if err != nil {
				t = tt
			} else {
				t = ts.Format(timeFormat)
			}
		case json.Number:
			i, err := tt.Int64()
			if err != nil {
				t = tt.String()
			} else {
				var sec, nsec int64 = i, 0
				if zerolog.TimeFieldFormat == zerolog.TimeFormatUnixMs {
					nsec = int64(time.Duration(i) * time.Millisecond)
					sec = 0
				}
				ts := time.Unix(sec, nsec).UTC()
				t = ts.Format(timeFormat)
			}
		}
		return colorize(t, colorDarkGray, noColor)
	}
}

func formatLevelForConsole(noColor bool) formatter {
	return func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			switch ll {
			case "debug":
				l = colorize("DBG", colorYellow, noColor)
			case "info":
				l = colorize("INF", colorGreen, noColor)
			case "warn":
				l = colorize("WRN", colorRed, noColor)
			case "error":
				l = colorize(colorize("ERR", colorRed, noColor), colorBold, noColor)
			case "fatal":
				l = colorize(colorize("FTL", colorRed, noColor), colorBold, noColor)
			case "panic":
				l = colorize(colorize("PNC", colorRed, noColor), colorBold, noColor)
			default:
				l = colorize("???", colorBold, noColor)
			}
		} else {
			if i == nil {
				l = colorize("???", colorBold, noColor)
			} else {
				l = strings.ToUpper(fmt.Sprintf("%s", i))[0:3]
			}
		}
		return l
	}
}

func formatCallerForConsole(noColor bool) formatter {
	return func(i interface{}) string {
		var c string
		if cc, ok := i.(string); ok {
			c = cc
		}
		if len(c) > 0 {
			cwd, err := os.Getwd()
			if err == nil {
				c = strings.TrimPrefix(c, cwd)
				c = strings.TrimPrefix(c, "/")
			}
			c = colorize(c, colorBold, noColor) + colorize(" >", colorCyan, noColor)
		}
		return c
	}
}

func consoleDefaultFormatMessage(i interface{}) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%s", i)
}
