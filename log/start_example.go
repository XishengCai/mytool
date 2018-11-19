package main

import (
	"fmt"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

func InitLogger() {
	baseLogPath := path.Join("./", "logrus.log")
	writer, err := rotatelogs.New(
		baseLogPath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %v", errors.WithStack(err))
	}
	switch level := "debug"; level {
	/*
	如果日志级别不是debug就不要打印日志到控制台了
	 */
	case "debug":
		log.SetLevel(log.DebugLevel)
		log.SetOutput(os.Stderr)
	case "info":
		setNull()
		log.SetLevel(log.InfoLevel)
	case "warn":
		setNull()
		log.SetLevel(log.WarnLevel)
	case "error":
		setNull()
		log.SetLevel(log.ErrorLevel)
	default:
		setNull()
		log.SetLevel(log.InfoLevel)
	}

	lfHook := lfshook.NewHook()
	log.AddHook(lfHook)
}

func setNull() {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Printf("err", err)
	}

	writer := bufio.NewWrite(src)
	log.SetOutput(writer)
}
