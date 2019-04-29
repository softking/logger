package logger

import (
	"fmt"
	"os"
	"time"
)

// 不建议使用文件作为日志存储介质,慢慢慢,慢的不行不行的,而且只能一个协程同步写入,建议用 tcp 或者udp
type FileLogWriter struct {
	rec      chan *LogRecord
	filename string
	file     *os.File
	format   string
	openTime string
}

// This is the FileLogWriter's output method
func (w *FileLogWriter) LogWrite(rec *LogRecord) {
	w.rec <- rec
}

func (w *FileLogWriter) Close() {
	close(w.rec)
	w.file.Sync()
}

func NewFileLogWriter(format, fname string) *FileLogWriter {
	w := &FileLogWriter{
		rec:      make(chan *LogRecord, LogBufferLength),
		filename: fname,
		format:   format,
	}
	w.run()
	return w
}

func (w *FileLogWriter) run() {
	if w.file == nil {
		nowTime := time.Now().Format("2006-01-02")
		fd, err := os.OpenFile(fmt.Sprintf("%s_%s", w.filename, nowTime), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
		if err == nil {
			w.file = fd
			w.openTime = nowTime
		}
	}

	go func() {
		for {
			rec, ok := <-w.rec
			if !ok {
				return
			}

			nowTime := time.Now().Format("2006-01-02")
			if w.openTime != nowTime {
				w.file.Close()
				fd, err := os.OpenFile(fmt.Sprintf("%s_%s", w.filename, nowTime), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
				if err == nil {
					w.file = fd
					w.openTime = nowTime
				}
			}

			_, err := fmt.Fprint(w.file, FormatLogRecord(w.format, rec))
			if err != nil {
				fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
				return
			}
		}
	}()

}
