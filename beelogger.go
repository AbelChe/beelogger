package beelogger

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/logrusorgru/aurora"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	_logger     *zerolog.Logger
	fnCache     sync.Map
	allowCaller = true
)

// SetLevel 仅设置日志级别；函数信息在 Debug / Error 系列调用时自动添加（如果 allowCaller=true）
func SetLevel(level string) {
	switch strings.ToLower(level) {
	case "off":
		log.Logger = log.Logger.Level(zerolog.Disabled)
	case "fatal":
		log.Logger = log.Logger.Level(zerolog.FatalLevel)
	case "error":
		log.Logger = log.Logger.Level(zerolog.ErrorLevel)
	case "warn":
		log.Logger = log.Logger.Level(zerolog.WarnLevel)
	case "info":
		log.Logger = log.Logger.Level(zerolog.InfoLevel)
	case "debug":
		log.Logger = log.Logger.Level(zerolog.DebugLevel)
	default:
		log.Logger = log.Logger.Level(zerolog.InfoLevel)
	}
}

func init() {
	colorEnabled := os.Getenv("TERM") != "" && os.Getenv("NO_COLOR") == ""
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    !colorEnabled,
		TimeFormat: "01-02 15:04:05",
	}

	output.FormatLevel = func(i interface{}) string {
		au := aurora.NewAurora(colorEnabled)
		switch i {
		case "fatal":
			return au.Bold(au.Red("[FATAL]")).String()
		case "error":
			return au.Red("[ERR]").String()
		case "warn":
			return au.Yellow("[WRN]").String()
		case "info":
			return au.Blue("[INF]").String()
		case "debug":
			return au.Magenta("[DBG]").String()
		default:
			return fmt.Sprintf("[%s]", i)
		}
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%v", i)
	}

	_l := log.Output(output)
	_logger = &_l
}

// ---- 对外基础获取 ----
func GetLogger() *zerolog.Logger { return _logger }

// ---- 公共级别封装 ----
// Info / Warn 不自动加函数信息（按你的需求只在 Debug 和 Error 系列加）
// 如果以后想给 Info 也加，改成调用 addCallerFields 即可。

func Info() *zerolog.Event {
	return _logger.Info()
}

func Warn() *zerolog.Event {
	return _logger.Warn()
}

func Debug() *zerolog.Event {
	e := _logger.Debug()
	// 如果 Debug 级别没启用(e.Disabled)则不做任何 caller 开销
	return addCallerFieldsIfEnabled(e, 2) // skip=2: addCallerFieldsIfEnabled -> Debug -> 用户代码
}

func Error() *zerolog.Event {
	e := _logger.Error()
	return addCallerFieldsForce(e, 2)
}

// Err(err) 依然要附带函数信息
func Err(err error) *zerolog.Event {
	e := _logger.Err(err)
	return addCallerFieldsForce(e, 2)
}

func Fatal() *zerolog.Event {
	e := _logger.Fatal()
	return addCallerFieldsForce(e, 2)
}

func Panic() *zerolog.Event {
	e := _logger.Panic()
	return addCallerFieldsForce(e, 2)
}

// 通用 Log() 不加（保持灵活）
func Log() *zerolog.Event {
	return _logger.Log()
}

// ---- 采集调用者信息的内部方法 ----
func addCallerFieldsIfEnabled(e *zerolog.Event, skip int) *zerolog.Event {
	if !allowCaller {
		return e
	}
	if e == nil || !e.Enabled() { // 未启用级别，不做开销
		return e
	}
	fn, file, line := callerInfo(skip + 1) // 再 +1 跳过本函数
	return e.
		Str("func", fn).
		Str("caller", fmt.Sprintf("%s:%d", file, line))
}

func addCallerFieldsForce(e *zerolog.Event, skip int) *zerolog.Event {
	if !allowCaller || e == nil {
		return e
	}
	// 对于 Error 系列，即使当前级别允许或不允许，这里 e.Enabled() 基本是 true（因为调用的是已通过级别过滤的方法）。
	// 保险起见仍判断：
	if !e.Enabled() {
		return e
	}
	fn, file, line := callerInfo(skip + 1)
	return e.
		Str("func", fn).
		Str("caller", fmt.Sprintf("%s:%d", file, line))
}

func callerInfo(skip int) (fnShort, fileShort string, line int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", "unknown", 0
	}
	fileShort = path.Base(file)
	fnShort = resolveFuncName(pc)
	return
}

func resolveFuncName(pc uintptr) string {
	if v, ok := fnCache.Load(pc); ok {
		return v.(string)
	}
	name := "unknown"
	if f := runtime.FuncForPC(pc); f != nil {
		n := f.Name()
		if idx := strings.LastIndex(n, "/"); idx >= 0 {
			n = n[idx+1:]
		}
		if i := strings.Index(n, "["); i >= 0 { // 去掉泛型实例化部分
			n = n[:i]
		}
		name = n
	}
	fnCache.Store(pc, name)
	return name
}
