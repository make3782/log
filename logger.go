package glogs

import (
	"bytes"
	"io"
	"sync"
	"time"
)

type logWriter struct {
	sync.Mutex // 这里是否需要锁呢？
	writer     io.Writer
}

func newLogWriter(wr io.Writer) *logWriter {
	return &logWriter{writer: wr}
}

func (lg *logWriter) println(when time.Time, msg string) {
	lg.Lock()
	var buf bytes.Buffer
	timeString := formatTime(when)
	buf.WriteString(timeString)
	buf.WriteString(msg)
	buf.WriteString("\n")
	lg.writer.Write(buf.Bytes())
	lg.Unlock()
}

// formatTime 格式化时间显示
func formatTime(when time.Time) string {
	return when.Format("2006/01/02 - 15:04:05.000")
}
