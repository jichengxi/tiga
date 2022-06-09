package ControllerTypes

import "sync"

var (
	Port       string
	Endpoint   string
	ResultData RespConfig
	WorkGroup  sync.WaitGroup
)
