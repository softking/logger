package logger

import (
	"bytes"
	"fmt"
)

type formatCacheType struct {
	LastUpdateSeconds    int64
	shortTime, shortDate string
	longTime, longDate   string
	timeStamp            string
}

var formatCache = &formatCacheType{}

// Known format codes:
// %T - Time (15:04:05 MST)
// %t - Time (15:04)
// %D - Date (2006/01/02)
// %d - Date (01/02/06)
// %L - Level (FNST, FINE, DEBG, TRAC, WARN, EROR, CRIT)
// %S - Source
// %M - Message
// Ignores unknown formats
// Recommended: "[%D %T] [%L] (%S) %M"
func FormatLogRecord(format string, rec *LogRecord) string {
	if rec == nil {
		return "<nil>"
	}
	if len(format) == 0 {
		return ""
	}

	out := bytes.NewBuffer(make([]byte, 0, 64))
	secs := rec.Created.UnixNano() / 1e9

	cache := *formatCache
	if cache.LastUpdateSeconds != secs {
		month, day, year := rec.Created.Month(), rec.Created.Day(), rec.Created.Year()
		hour, minute, second := rec.Created.Hour(), rec.Created.Minute(), rec.Created.Second()
		zone, _ := rec.Created.Zone()
		updated := &formatCacheType{
			LastUpdateSeconds: secs,
			shortTime:         fmt.Sprintf("%02d:%02d", hour, minute),
			shortDate:         fmt.Sprintf("%02d/%02d/%02d", day, month, year%100),
			longTime:          fmt.Sprintf("%02d:%02d:%02d %s", hour, minute, second, zone),
			longDate:          fmt.Sprintf("%04d/%02d/%02d", year, month, day),
			timeStamp:         fmt.Sprintf("%d", rec.Created.Unix()),
		}
		cache = *updated
		formatCache = updated
	}

	// Split the string into pieces by % signs
	pieces := bytes.Split([]byte(format), []byte{'%'})

	// Iterate over the pieces, replacing known formats
	for i, piece := range pieces {
		if i > 0 && len(piece) > 0 {
			switch piece[0] {
			case 'T':
				out.WriteString(cache.longTime)
			case 't':
				out.WriteString(cache.shortTime)
			case 'L':
				out.WriteString(levelStrings[rec.Level])
			case 'D':
				out.WriteString(cache.longDate)
			case 'd':
				out.WriteString(cache.shortDate)
			case 'M':
				out.WriteString(rec.Message)
			case 'S':
				out.WriteString(cache.timeStamp)
			}
			if len(piece) > 1 {
				out.Write(piece[1:])
			}
		} else if len(piece) > 0 {
			out.Write(piece)
		}
	}
	out.WriteByte('\n')

	return out.String()
}
