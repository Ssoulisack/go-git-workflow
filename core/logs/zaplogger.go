package logs

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Log message structure
type logMessage struct {
	level      zapcore.Level
	message    string
	key        string
	fields     []zap.Field
	callerInfo string // File and line number
}

var log *zap.Logger
var err error

// Channel to queue logs
var logQueue chan logMessage

const CUSTOM_LOG_FORMAT string = "TIME[${time}] PID[${pid}] REQUESTID[${locals:requestid}] RESSTATUS[${status}] - LATENCY[${latency}] METHOD[${method}] PATH[${path}] REFERER[${referer}] PROTOCOL[${protocol}] PORT[${port}] IP[${ip}] IPS[${ips}] HOST[${host}] UA[${ua}] REQHEADERS[${reqHeaders}] REQQUERYPARAMS[${queryParams}] \n URL[${url}]\n REQBODY[${body}] REQHEADER:[${header:}] REQHEADER:[${reqHeader:}] REQQUERY[${query:}] REQFORM[${form:}] REQCOOKIE[${cookie:}] \n RESBODY[${resBody}]\n BYTESSENT[${bytesSent}] BYTESRECEIVED[${bytesReceived}] ROUTE[${route}] ERROR[${error}]  RESPHEADER:[${respHeader:}]  LOCALS:[${locals:}]\n <------------------------------------------------------------------------------------> \n"

// Generate log filename based on date
func getLogFilename() string {
	currentDate := time.Now().Format("2006-01-02") // Format: YYYY-MM-DD
	return fmt.Sprintf("logs/%s.log", currentDate)
}

func init() {
	logQueue = make(chan logMessage, 1000) // Log buffer size

	// Configure lumberjack for daily log rotation
	logFile := &lumberjack.Logger{
		Filename:   getLogFilename(), // Dynamic filename
		MaxSize:    100,              // Max file size in MB
		MaxBackups: 0,                // Unlimited backups (handled by naming convention)
		MaxAge:     0,                // No auto-delete
		Compress:   false,
		LocalTime:  true,
	}

	// Create Zap encoders
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.StacktraceKey = ""

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig) // Console format
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)       // JSON format for files

	// Multi-writer (log to both file and console)
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zapcore.InfoLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
	)

	// Build logger
	log = zap.New(core, zap.AddCaller())

	// Start the background logger
	go processLogQueue()
}

// Goroutine to process logs asynchronously
func processLogQueue() {
	for logMsg := range logQueue {
		fields := make([]zap.Field, len(logMsg.fields))
		copy(fields, logMsg.fields)

		// Add search_key if available
		if logMsg.key != "" {
			fields = append(fields, zap.String("search_key", logMsg.key))
		}

		// Add caller info to log fields (for both console and file)
		fields = append(fields, zap.String("source", logMsg.callerInfo))
		// Add message to log fields
		fields = append(fields, zap.String("message", logMsg.message))

		// Log based on the level (both console and file logging)
		switch logMsg.level {
		case zapcore.InfoLevel:
			log.Info(logMsg.message, fields...) // Log to console and file
		case zapcore.ErrorLevel:
			log.Error(logMsg.message, fields...) // Log to console and file
		}
	}
}

// Helper function to get file and line number
func getCallerInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// Info logs an informational message asynchronously with caller info
func Info(message string, fields ...zap.Field) {
	callerInfo := getCallerInfo(2) // Capture caller location
	select {
	case logQueue <- logMessage{level: zapcore.InfoLevel, message: message, fields: fields, callerInfo: callerInfo}:
	default:
	}
}

// Error logs an error message asynchronously with caller info
func Error(message interface{}, fields ...zap.Field) {
	callerInfo := getCallerInfo(2)
	msgStr := ""

	switch v := message.(type) {
	case error:
		msgStr = v.Error()
	case string:
		msgStr = v
	default:
		msgStr = "Unknown error type"
	}

	select {
	case logQueue <- logMessage{level: zapcore.ErrorLevel, message: msgStr, fields: fields, callerInfo: callerInfo}:
	default:
	}
}

// **ImportantLog logs with a searchable key and caller info**
func ImportantLog(key string, message interface{}, fields ...zap.Field) {
	// Convert message to string
	var msgStr string
	switch v := message.(type) {
	case string:
		msgStr = v
	case error:
		msgStr = v.Error()
	default:
		msgStr = fmt.Sprintf("%v", v) // Convert any other type to string
	}

	// Capture the caller info
	callerInfo := getCallerInfo(2)

	// Send log message to the queue
	select {
	case logQueue <- logMessage{level: zapcore.InfoLevel, key: key, message: msgStr, fields: fields, callerInfo: callerInfo}:
	default:
	}
}

// Sync ensures all logs are flushed before app exit
func Sync() {
	log.Sync() // Ensure logs are written before exit
}
