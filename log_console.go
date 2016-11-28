package glogs

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func init() {
	fmt.Printf("get string: %v", AdapterConsole)
	Register(AdapterConsole, NewConsole)
}

func NewConsole() Logger {
	nLog := &consoleWrite{
		lg:       newLogWriter(os.Stdout),
		Colorful: false,
	}
	return nLog
}

type consoleWrite struct {
	lg       *logWriter
	Colorful bool `json:"color"` // 是否彩色显示
}

// Init 接口实现
func (log *consoleWrite) Init(config string) error {
	if len(config) == 0 {
		return nil
	}
	err := json.Unmarshal([]byte(config), log)
	return err
}

// WriteMsg 接口实现
func (log *consoleWrite) WriteMsg(when time.Time, msg string, level int) error {
	if log.Colorful {
		msg = colors[level](msg)
	}
	log.lg.println(when, msg)
	return nil
}

// 终端颜色显示
type brush func(string) string

func newBrush(colorString string) brush {
	start := "\033["
	end := "\033[0m"
	return func(text string) string {
		return start + colorString + "m" + text + end
	}
}

var colors = []brush{
	newBrush("1;37"),
	newBrush("1;31"), // Error   red
	newBrush("1;36"), // Alert   cyan
	newBrush("1;33"), // Warn    yellow
	newBrush("1;32"), // Notice   green
	newBrush("1;34"), // info  blue
	newBrush("1;35"), // debug  magenta
}
