package log

import (
	"encoding/json"
	"io"
	"log"
)

type Logger interface {
	Log(msg string, context map[string]interface{})
	Debug(msg string, context map[string]interface{})
}

const flags = log.Ldate | log.Ltime | log.Lmicroseconds

type logst struct {
	*log.Logger
}

func NewLogger(w io.Writer) Logger {
	return &logst{log.New(w, "", flags)}
}

func (l *logst) Log(msg string, context map[string]interface{}) {
	l.logMarsh("", msg, context)
}

func (l *logst) Debug(msg string, context map[string]interface{}) {
	l.logMarsh("Debug", msg, context)
}

func (l *logst) logMarsh(lvl, msg string, ctx map[string]interface{}) {
	json, err := json.Marshal(struct {
		Lvl string                 `json:"lvl"`
		Msg string                 `json:"msg"`
		Ctx map[string]interface{} `json:"context"`
	}{lvl, msg, ctx})

	if err == nil {
		l.Printf(string(json))
	} else {
		l.Printf("unable to marshal log entry: %v", err)
	}
}
