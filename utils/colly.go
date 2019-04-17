package utils

import (
	"github.com/gocolly/colly/debug"
	"github.com/sirupsen/logrus"
)

type LogrusCollectorDebugger struct {
}

func (d LogrusCollectorDebugger) Init() error {
	return nil
}
func (d LogrusCollectorDebugger) Event(e *debug.Event) {
	logrus.WithFields(logrus.Fields{
		"CollectorID": e.CollectorID,
		"RequestID":   e.RequestID,
		"Type":        e.Type,
		"Values":      e.Values,
	}).Debugln()
}

var _ debug.Debugger = LogrusCollectorDebugger{}
