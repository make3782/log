// 实现log的file输出

package glogs

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// init初始化进行注册
func init() {
	Register(AdapterFile, NewFile)
}

// NewFile 新的文件日志记录器
func NewFile() Logger {
	nLog := &fileWrite{
		Rotate: true,
	}
	return nLog
}

type fileWrite struct {
	sync.RWMutex
	Filename   string `json:"filename"`
	fileWriter *os.File

	// 按最大行数分割文件
	MaxLines int64 `json:"maxlines"`

	// 当前已经写入的行
	currLines int64

	// 按最大大小分割文件
	MaxSize int64 `json:"maxsize"`
	// 当前文件大小
	currSize int64

	// 按时间日期分割文件
	MaxDays int64 `json:"maxdays"`

	// 是否分割文件
	Rotate bool `json:"rotate"`

	// 文件后缀名 默认为.log
	ext string
}

// Init 接口实现
func (w *fileWrite) Init(config string) error {
	err := json.Unmarshal([]byte(config), w)
	if err != nil {
		return err
	}
	if len(w.Filename) == 0 {
		return errors.New("jsonconfig must have filename")
	}
	if w.ext == "" {
		w.ext = ".log"
	}
	return nil
}

// WriteMsg 接口实现
func (w *fileWrite) WriteMsg(when time.Time, msg string, level int) error {
	t := formatTime(when)
	msg = t + msg + "\n"
	if w.Rotate {
		// 需要分割文件的话

	}
	w.Lock() // 写文件时锁定
	s, err := w.fileWriter.WriteString(msg)
	if err == nil {
		w.currLines++
		w.currSize += s
	}
	w.Unlock()
	return err
}
