// Usage:
//
// import "github.com/make3782/log"
//
// log := NewLogger(10000)
// log.SetLoger("console", "")

package glogs

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// message levels
const (
	LevelDebug = iota
	LevelInfo
	LevelNotice
	LevelWarn
	LevelAlert
	LevelError
	LevelPanic
)

// 定义一些其他可能用的级别
const (
	LevelTrace   = LevelDebug
	LevelWarning = LevelWarn
)

// 定义adapter name
const (
	AdapterConsole = "console"
	AdapterFile    = "file"
)

type newLoggerFunc func() Logger

// Logger 接口
type Logger interface {
	Init(config string) error                             // 通过Init方法来初始化相关配置
	WriteMsg(when time.Time, msg string, level int) error // adapters的输出，如：console则输出到终端，file则写入到文件
}

var adapters = make(map[string]newLoggerFunc)

// 定义日志输出的开始部分，放在全局里更好点
var levelPrefix = [LevelPanic + 1]string{
	" [Debug] ",
	" [Info] ",
	" [Notice] ",
	" [Warn] ",
	" [Alert] ",
	" [Error] ",
	" [Panic] ",
}

// Register 注册log类型
func Register(name string, log newLoggerFunc) {
	if log == nil { // 是否已经实现该log
		panic("logs: Register provide is nil")
	}
	if _, dupLog := adapters[name]; dupLog { // 是否已经注册
		panic("logs: Register already exists for provider" + name)
	}
	adapters[name] = log
}

// GLogger 是默认的logger
type GLogger struct {
	lock      sync.Mutex
	init      bool          // 是否已被初始化过
	level     int           // 显示的最低级别，高于此级别的才显示
	asyncFlag bool          // 是否是异步输出日志
	msgChan   chan *logMsg  // 具体的一条消息体
	outputs   []*nameLogger // 要输出的目标（eg： 终端+file+...）
}

type nameLogger struct {
	Logger
	name string
}

// 定义消息体
type logMsg struct {
	level int       // 该条消息的级别
	msg   string    // 消息内容
	time  time.Time // 消息发生的时间
}

// NewLogger 创建一个新的logger
// @param channelLens int 用于异步时候的缓冲数量
func NewLogger(channelLens ...int) *GLogger {
	newLog := new(GLogger)
	newLog.level = LevelInfo // 默认的日志级别为Info
	newLog.asyncFlag = false
	return newLog
}

func (gl *GLogger) setLogger(adapterName string, configs ...string) error {
	config := append(configs, "{}")[0]
	// 判断是否已经注册过
	for _, l := range gl.outputs {
		if l.name == adapterName {
			return fmt.Errorf("logs: duplicate adaptername %q", adapterName)
		}
	}

	logFunc, ok := adapters[adapterName]
	if !ok {
		return fmt.Errorf("logs: unknown adaptername %v (forgotten Register?)", adapterName)
	}
	lg := logFunc()
	if err := lg.Init(config); err != nil {
		fmt.Fprintln(os.Stderr, "logs.GLogger.SetLogger error: "+err.Error())
		return err
	}
	gl.outputs = append(gl.outputs, &nameLogger{name: adapterName, Logger: lg})
	return nil
}

// SetLogger 设置logger的输出方式
// configs 为json形式的字符串，用于配置不同的输出方式：如： {"interval":360}.
func (gl *GLogger) SetLogger(adapterName string, configs ...string) error {
	gl.lock.Lock()
	defer gl.lock.Unlock()
	if !gl.init {
		gl.outputs = []*nameLogger{}
		gl.init = true
	}
	return gl.setLogger(adapterName, configs...)
}

// SetLevel 设置日志的最小显示级别
func (gl *GLogger) SetLevel(level int) {
	gl.level = level
}

func (gl *GLogger) writeMsg(level int, msg string, v ...interface{}) error {
	if !gl.init {
		gl.lock.Lock()
		gl.setLogger(AdapterConsole)
		gl.lock.Unlock()
	}

	if len(v) > 0 {
		msg = fmt.Sprintf(msg, v) // 如果有变量替换，则进行替换
	}

	if level == -1 {
		level = LevelError // 设置了不存在的level，则自动恢复为最高级别error
	}
	msg = levelPrefix[level] + msg

	when := time.Now()
	if gl.asyncFlag {
		// 异步写入
	} else {
		gl.wrtiteToLoggers(when, msg, level)
	}
	return nil
}

func (gl *GLogger) wrtiteToLoggers(when time.Time, msg string, level int) {
	for _, l := range gl.outputs {
		// 循环要输出的终端进行输出
		err := l.WriteMsg(when, msg, level)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to write msg to adapter: %v, error: %v\n", l.name, err)
		}
	}
}

func (gl *GLogger) start() {

}

///////////////////////////////////////////////////////////
// 输出的公共方法
///////////////////////////////////////////////////////////

// Panic 输出
func (gl *GLogger) Panic(format string, v ...interface{}) {
	if gl.level <= LevelPanic {
		gl.writeMsg(LevelPanic, format, v...)
	}
}

func (gl *GLogger) Error(format string, v ...interface{}) {
	if gl.level <= LevelError {
		gl.writeMsg(LevelError, format, v...)
	}
}

func (gl *GLogger) Alert(format string, v ...interface{}) {
	if gl.level <= LevelAlert {
		gl.writeMsg(LevelAlert, format, v...)
	}
}

func (gl *GLogger) Warn(format string, v ...interface{}) {
	if gl.level <= LevelWarn {
		gl.writeMsg(LevelWarn, format, v...)
	}
}

func (gl *GLogger) Warning(format string, v ...interface{}) {
	gl.Warn(format, v...)
}

func (gl *GLogger) Notice(format string, v ...interface{}) {
	if gl.level <= LevelNotice {
		gl.writeMsg(LevelNotice, format, v...)
	}
}

func (gl *GLogger) Info(format string, v ...interface{}) {
	if gl.level <= LevelInfo {
		gl.writeMsg(LevelInfo, format, v...)
	}
}

func (gl *GLogger) Debug(format string, v ...interface{}) {
	if gl.level <= LevelDebug {
		gl.writeMsg(LevelDebug, format, v...)
	}
}
