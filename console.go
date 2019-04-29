package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// ConsoleLogWriter 终端输出
type ConsoleLogWriter struct {
	format string
	record chan *LogRecord
}

// NewConsoleLogWriter 创建
func NewConsoleLogWriter(format string) *ConsoleLogWriter {
	w := &ConsoleLogWriter{
		format: format,
		record: make(chan *LogRecord, LogBufferLength),
	}
	go w.run(os.Stdout)
	return w
}

func (c *ConsoleLogWriter) run(out io.Writer) {
	for rec := range c.record {
		fmt.Fprint(out, FormatLogRecord(c.format, rec))
	}
}

// LogWrite 接口
func (c *ConsoleLogWriter) LogWrite(rec *LogRecord) {
	c.record <- rec
}

// Close 接口
func (c *ConsoleLogWriter) Close() {
	close(c.record)
	time.Sleep(50 * time.Millisecond) // Try to give console I/O time to complete
}
