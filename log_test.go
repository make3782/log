package glogs

import (
	"testing"
)

func TestNewLoger(t *testing.T) {
	log := NewLogger()

	ok := log.SetLogger("console", `{"color": true}`)
	log.SetLevel(LevelDebug)
	if ok != nil {
		t.Errorf("get error: %v", ok)
	}
	log.Error("wzx test")
	log.Alert("wzx test2")
	log.Info("wzx test3")
	log.Debug("wzx test4")
	log.Notice("wzx test5")
	log.Warn("wzx test6")
	log.Panic("wzx test7")
}
