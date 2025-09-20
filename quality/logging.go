package quality

import (
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

type Logger interface {
	Info(msg string, fields map[string]interface{})
	Error(msg string, fields map[string]interface{})
}

type SimpleLogger struct{}

func (l *SimpleLogger) getCallerInfo() string {
	// Get caller file and line number
	_, file, line, ok := runtime.Caller(3) // 3 levels up the stack
	if !ok {
		return ""
	}
	// Shorten the file path to just the last two segments
	parts := strings.Split(filepath.ToSlash(file), "/")
	if len(parts) > 2 {
		file = strings.Join(parts[len(parts)-2:], "/")
	}
	return file + ":" + strconv.Itoa(line) // Fixed: proper int to string conversion
}

func (l *SimpleLogger) Info(msg string, fields map[string]interface{}) {
	caller := l.getCallerInfo()
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["caller"] = caller
	log.Printf("INFO: %s %+v\n", msg, fields)
}

func (l *SimpleLogger) Error(msg string, fields map[string]interface{}) {
	caller := l.getCallerInfo()
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["caller"] = caller
	log.Printf("ERROR: %s %+v\n", msg, fields)
}

var DefaultLogger Logger = &SimpleLogger{}

func LogTestInput(step string, input interface{}) {
	DefaultLogger.Info(step, map[string]interface{}{
		"payload": resolveShallow(input),
	})
}

func LogPayload(msg string, payload interface{}) {
	DefaultLogger.Info(msg, map[string]interface{}{
		"payload": resolveShallow(payload),
	})
}

func LogStart(action string, payload interface{}) {
	DefaultLogger.Info("action started", map[string]interface{}{
		"action":  action,
		"payload": resolveShallow(payload),
	})
}

func LogSuccess(action string, result interface{}) {
	DefaultLogger.Info("action completed successfully", map[string]interface{}{
		"action": action,
		"result": resolveShallow(result),
	})
}

func LogError(action string, err error) {
	DefaultLogger.Error("action failed", map[string]interface{}{
		"action": action,
		"error":  err.Error(),
	})
}

func resolveShallow(payload interface{}) interface{} {
	v := reflect.ValueOf(payload)
	t := reflect.TypeOf(payload)

	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return payload
	}

	result := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" { // unexported
			continue
		}
		result[field.Name] = v.Field(i).Interface()
	}
	return result
}
