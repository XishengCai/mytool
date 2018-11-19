package main

import (
	log "github.com/sirupsen/logrus"
)

type Hook interface {
	Levels() []log.Level
	Fire(*log.Entry) error
}
type DefaultFieldsHook struct {
}

func (df *DefaultFieldsHook) Fire(entry *log.Entry) error {
	entry.Data["appName"] = "MyAppName"
	return nil
}

func(df *DefaultFieldsHook) Levels() []log.Level {
	return log.AllLevels
}


func main() {
	log.SetFormatter(&log.JSONFormatter{})

	// 添加自己实现的HOOK
	log.AddHook(&DefaultFieldsHook{})


	log.Info("test")
}


