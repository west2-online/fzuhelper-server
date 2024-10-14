package logger

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// MyFormatter 是一个自定义的日志格式化器，
type MyFormatter struct {
	wd string
}

// Format 将日志条目格式化为字节切片。
// 它根据条目是否有调用者信息，决定写入哪些信息。
func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := entry.Buffer

	// 如果有调用者信息，则写入时间戳、日志级别、简短调用者路径、消息和堆栈跟踪。
	// 否则，只写入时间戳、日志级别、消息和堆栈跟踪。
	if entry.HasCaller() {
		writeTimestamp(b, entry)
		writeLevel(b, entry)
		writeShortCallerPath(b, entry)
		writeMessage(b, entry)
		writeCallerStack(b, entry)
		writeLineBreak(b)
	} else {
		writeTimestamp(b, entry)
		writeLevel(b, entry)
		writeMessage(b, entry)
		writeCallerStack(b, entry)
		writeLineBreak(b)
	}

	return b.Bytes(), nil
}

func getNewFormatter() *MyFormatter {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return &MyFormatter{
		wd,
	}
}

func writeTimestamp(b *bytes.Buffer, entry *logrus.Entry) {
	b.WriteString(entry.Time.Format(time.RFC3339))
	b.WriteRune('	')
}

func writeLevel(b *bytes.Buffer, entry *logrus.Entry) {
	b.WriteRune('[')
	b.WriteString(strings.ToUpper(entry.Level.String()))
	b.WriteRune(']')
	b.WriteRune('	')
}

// writeShortCallerPath 将日志条目的调用者路径写入缓冲区。
// 它使用工作目录来构建一个更短的调用者路径。
func writeShortCallerPath(b *bytes.Buffer, entry *logrus.Entry) {
	parts := strings.Split(entry.Caller.File, "/")
	if len(parts) > 2 {
		b.WriteString(filepath.Join(parts[len(parts)-2], parts[len(parts)-1])) // // 只使用文件名和最后一个目录。
	} else {
		b.WriteString(entry.Caller.File) // 如果路径太短，则直接使用完整路径。
	}

	b.WriteRune(':')
	b.WriteString(strconv.Itoa(entry.Caller.Line))
	b.WriteRune(' ')
}

func writeMessage(b *bytes.Buffer, entry *logrus.Entry) {
	b.WriteString(entry.Message)
}

// writeCallerStack 如果日志级别大于等于 ErrorLevel，则将调用堆栈写入缓冲区。
func writeCallerStack(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.Level > logrus.ErrorLevel {
		return
	}
	b.WriteString("\nStack Trace:\n")

	pc := make([]uintptr, 10)
	n := runtime.Callers(10, pc)
	frames := runtime.CallersFrames(pc[:n-1])

	for {
		frame, more := frames.Next()
		b.WriteString(fmt.Sprintf("%s\n", frame.Function))
		b.WriteString(fmt.Sprintf("	%s:%d\n", frame.File, frame.Line))
		if !more {
			break
		}
	}
}

func writeLineBreak(b *bytes.Buffer) {
	b.WriteRune('\n')
}
