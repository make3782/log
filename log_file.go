// 实现log的file输出

package glogs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

// init初始化进行注册
func init() {
	Register(AdapterFile, NewFile)
}

// NewFile 新的文件日志记录器
func NewFile() Logger {
	nLog := &fileWrite{
		Rotate: true,
		Perm:   "0777",
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
	Daily         bool  `json:"daily"`
	MaxDays       int64 `json:"maxdays"`
	dailyOpenDate int
	dailyOpenTime time.Time

	// 是否分割文件
	Rotate bool `json:"rotate"`

	// 是否将文件按错误level分割
	RotateLevel bool `json:"rotatelevel"`

	// 文件后缀名 默认为.log
	ext string

	// 文件权限
	Perm string `json:"perm"`
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

	return w.startLogger()
}

func (w *fileWrite) startLogger() error {
	// 检查是否创建文件
	file, err := w.createLogFile()
	if err != nil {
		return err
	}

	if w.fileWriter != nil {
		w.fileWriter.Close()
	}
	w.fileWriter = file

	// 初始化文件相关
	fInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("get stat err: %s\n", err)
	}
	w.currSize = int64(fInfo.Size())
	w.currLines = 0
	if w.Daily {
		go w.dailyRotate(w.dailyOpenTime) // 启动后，启动一个协程去处理按日期的文件分割
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
		w.currSize += int64(s)
	}
	w.Unlock()
	return err
}

// createLogFile 创建日志文件
func (w *fileWrite) createLogFile() (*os.File, error) {
	perm, err := strconv.ParseInt(w.Perm, 8, 64)
	if err != nil {
		return nil, err
	}
	fd, err := os.OpenFile(w.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
	if err == nil {
		os.Chmod(w.Filename, os.FileMode(perm))
	}
	return fd, err
}

// dailyRotate 这里用计时器协程，效率高于每次检查文件
func (w *fileWrite) dailyRotate(openTime time.Time) {
	y, m, d := openTime.Add(time.Hour * 24).Date()
	nextDay := time.Date(y, m, d, 0, 0, 0, 0, openTime.Location())
	tm := time.NewTimer(time.Duration(nextDay.UnixNano() - openTime.UnixNano() + 100))
	select {
	case <-tm.C:
		w.Lock()
		if w.needRotate(time.Now().Day()) {
			if err := w.doRotate(time.Now()); err != nil {
				fmt.Fprintf(os.Stderr, "FileLogWrite(%q): %s\n", w.Filename, err)
			}
		}
		w.Unlock()
	}
}

func (w *fileWrite) needRotate(day int) bool {
	return (w.MaxLines > 0 && w.currLines >= w.MaxLines) || (w.MaxSize > 0 && w.currSize >= w.MaxSize) || (w.Daily && day != w.dailyOpenDate)
}

//doRotate 将文件分割，用于写入新的文件
func (w *fileWrite) doRotate(logTime time.Time) error {
	//num := 1
	//fName := ""

	_, err := os.Lstat(w.Filename)
	if err != nil {

	}

	if w.MaxLines > 0 || w.MaxSize > 0 {

	}

	// 关闭原来的旧文件
	w.fileWriter.Close()

	// 重新启动logger

	startLoggerErr := w.startLogger()
	return startLoggerErr
}
