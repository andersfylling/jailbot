package jailbotque

import (
	"errors"
	"sync"
)

type Alert struct {
}

type AlertStack struct {
	sync.RWMutex
	stack []*Alert
	index uint
}

func NewAlertStack() *AlertStack {
	array := []*Alert{}

	return &AlertStack{
		stack: array,
		index: 0,
	}
}

// TODO: an ID should be used to make sure it's the correct stack entry and not a new one
// func (a *AlertStack) GetFirstAlert() (*Alert, error)      {}
// func (a *AlertStack) GetAlert(index uint) (*Alert, error) {}
// func (a *AlertStack) RemoveAlert(alert *Alert)            {}
// func (a *AlertStack) RemoveAlertByIndex(index uint)       {}

func (a *AlertStack) Add(alert *Alert) {
	a.Lock()
	defer a.Unlock()

	if uint(len(a.stack)) == a.index {
		a.stack = append(a.stack, alert)
	} else {
		a.stack[a.index] = alert
		a.index++
	}
}
func (a *AlertStack) Pop(index uint) (*Alert, error) {
	a.Lock()
	defer a.Unlock()

	if index > 0 {
		index--
		return a.stack[a.index], nil
	}

	return nil, errors.New("No alerts in stack")
}
