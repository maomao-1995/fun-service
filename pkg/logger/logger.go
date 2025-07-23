package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Logger 全局日志实例
var Logger *logrus.Logger

// Config 日志配置结构体
type Config struct {
	Path  string // 日志文件路径（如 "storage/logs/app.log"）
	Level string // 日志级别（debug, info, warn, error, fatal, panic）
}

// Init 初始化日志
func Init(cfg Config) {
	// 1. 创建日志实例
	Logger = logrus.New()

	// 2. 设置日志级别
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		// 级别无效时默认使用 info
		level = logrus.InfoLevel
		fmt.Printf("日志级别设置失败，使用默认级别: %v, 错误: %v\n", level, err)
	}
	Logger.SetLevel(level)

	// 3. 设置日志格式（带文件名和行号）
	Logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05", // 时间格式
		FullTimestamp:   true,                  // 显示完整时间
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// 格式化调用者信息（文件名:行号）
			filename := filepath.Base(f.File)
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})

	// 4. 启用调用者信息（文件名和行号）
	Logger.SetReportCaller(true)

	// 5. 设置输出（同时输出到控制台和文件）
	if cfg.Path != "" {
		// 创建日志目录（若不存在）
		logDir := filepath.Dir(cfg.Path)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Printf("创建日志目录失败: %v\n", err)
			// 目录创建失败时仅输出到控制台
			return
		}

		// 创建日志文件（追加模式，权限 0644）
		file, err := os.OpenFile(cfg.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Printf("打开日志文件失败: %v\n", err)
			return
		}

		// 同时输出到文件和控制台
		Logger.SetOutput(os.Stdout) // 控制台
		Logger.AddHook(&fileHook{   // 文件（通过 hook 实现多输出）
			file: file,
		})
	} else {
		// 未指定路径时仅输出到控制台
		Logger.SetOutput(os.Stdout)
	}
}

// fileHook 用于将日志同时输出到文件的 hook
type fileHook struct {
	file *os.File
}

// Fire 处理日志输出
func (h *fileHook) Fire(entry *logrus.Entry) error {
	// 格式化日志内容（与控制台输出一致）
	line, err := entry.String()
	if err != nil {
		return err
	}
	// 写入文件
	_, err = h.file.WriteString(line)
	return err
}

// Levels 指定该 hook 生效的日志级别
func (h *fileHook) Levels() []logrus.Level {
	return logrus.AllLevels // 所有级别都输出到文件
}

// 以下为简化日志调用的包装函数（避免直接使用 Logger 时重复写调用位置）

// Debug 调试日志
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Debugf 格式化调试日志
func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

// Info 信息日志
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

// Warn 警告日志
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// Warnf 格式化警告日志
func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

// Error 错误日志
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Errorf 格式化错误日志
func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

// Fatal 致命错误日志（会触发程序退出）
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// Fatalf 格式化致命错误日志
func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

// Panic 恐慌日志（会触发 panic）
func Panic(args ...interface{}) {
	Logger.Panic(args...)
}

// Panicf 格式化恐慌日志
func Panicf(format string, args ...interface{}) {
	Logger.Panicf(format, args...)
}
