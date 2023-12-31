package liblogger

import (
	"os"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var initialize = false
var logger *zap.Logger

// loadConfig -
func loadConfig() {
	if !initialize {
		file := os.Getenv("log_file")
		if file == "" {
			file = "log.log"
		}

		dir := os.Getenv("dir_log")
		if dir == "" {
			file = "./" + file
		} else {
			file = dir + "/" + file
		}

		var ws []zapcore.WriteSyncer
		ws = append(ws, zapcore.AddSync(&lumberjack.Logger{
			Filename:   file,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		}))

		if os.Getenv("log_silent") != "1" {
			ws = append(ws, zapcore.AddSync(os.Stdout))
		}

		cfg := zap.NewProductionEncoderConfig()
		cfg.EncodeTime = zapcore.ISO8601TimeEncoder

		// Custom encoder to skip lengthy stacktrace
		customEncoder, _ := newCustomEncoder(cfg)

		core := zapcore.NewCore(
			//	zapcore.NewJSONEncoder(cfg), // Default JSON encoder
			customEncoder, // Use custom encoder
			zapcore.NewMultiWriteSyncer(ws...),
			zap.InfoLevel,
		)

		// Load caller skip from .env
		envCallerSkip, err := strconv.Atoi(os.Getenv("log_caller_skip"))
		if err != nil || envCallerSkip < 0 {
			envCallerSkip = 0
		}

		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(envCallerSkip))
		defer logger.Sync()

		initialize = true
	}
}

// Sync
func Sync() {
	logger.Sync()
}

// Infow -
func Infow(msg string, values ...interface{}) {
	loadConfig()
	logger.Sugar().Infow(msg, values...)
}

// Infof -
func Infof(template string, args ...interface{}) {
	loadConfig()
	logger.Sugar().Infof(template, args...)
}

// Errorf -
func Errorf(template string, args ...interface{}) {
	loadConfig()
	logger.Sugar().Errorf(template, args...)
}
