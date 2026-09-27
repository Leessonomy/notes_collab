package app

import (
	"log"

	"go.uber.org/fx/fxevent"
)

type errorLogger struct{}

func (errorLogger) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.Invoked:
		if e.Err != nil {
			log.Printf("invoke %s: %v", e.FunctionName, e.Err)
		}
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			log.Printf("start %s: %v", e.FunctionName, e.Err)
		}
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			log.Printf("stop %s: %v", e.FunctionName, e.Err)
		}
	}
}
