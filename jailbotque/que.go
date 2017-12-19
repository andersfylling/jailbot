package jailbotque

import (
	"sync"
)

type JailbotQue struct {
	Alerts *AlertStack
}

var singleton *JailbotQue
var setupIssue error
var once sync.Once

// https://console.bluemix.net/docs/services/ComposeForMongoDB/connecting-external.html#connecting-external-app
func setup() {
	singleton = &JailbotQue{Alerts: NewAlertStack()}
}

// GetInstance singleton pattern
func GetInstance() (*JailbotQue, error) {
	once.Do(setup)

	return singleton, setupIssue
}
