package beelogger

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func SetLevel(level string) {
	switch level {
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
		//Add file and line number to log
		log.Logger = log.With().Caller().Logger()
	default:
		log.Logger = log.Logger.Level(zerolog.InfoLevel)
	}
}

var (
	_logger *zerolog.Logger
)

func init() {
	Color := true
	output := zerolog.ConsoleWriter{
		Out:          os.Stdout,        //os.Stderr:无法重定向文件
		NoColor:      Color,            //日志颜色
		TimeFormat:   "01-02 15:04:05", //日志时间参数
		PartsExclude: []string{},
	}
	output.FormatLevel = func(i interface{}) string {
		prefix := ""
		au := aurora.NewAurora(Color)
		switch i {
		case "fatal":
			prefix = au.Bold(au.Red("[FATAL]")).String()
		case "error":
			prefix = au.Red("[ERR]").String()
		case "warn":
			prefix = au.Yellow("[WRN]").String()
		case "info":
			prefix = au.Blue("[INF]").String()
		case "debug":
			prefix = au.Magenta("[DBG]").String()
		default:
			break
		}
		return prefix
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	_l := log.Output(output)
	_logger = &_l
}

func GetLogger() *zerolog.Logger {
	return _logger
}

func Info() *zerolog.Event {
	return _logger.Info()
}

func Debug() *zerolog.Event {
	return _logger.Debug()
}

func Warn() *zerolog.Event {
	return _logger.Warn()
}

func Error() *zerolog.Event {
	return _logger.Error()
}

func Fatal() *zerolog.Event {
	return _logger.Fatal()
}

func Err(err error) *zerolog.Event {
	return _logger.Err(err)
}

func Panic(err error) *zerolog.Event {
	return _logger.Panic()
}

func Log() *zerolog.Event {
	return _logger.Log()
}
