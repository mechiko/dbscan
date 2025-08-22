package main

import (
	"go.uber.org/zap"
)

type app struct {
	loger *zap.SugaredLogger
	pwd   string
}

// var _ dbscan.Apper = (*app)(nil)

func NewApp(logger *zap.SugaredLogger, pwd string) *app {
	newApp := &app{}
	newApp.pwd = pwd
	newApp.loger = logger
	return newApp
}

func (a *app) Pwd() string {
	return a.pwd
}

func (a *app) Logger() *zap.SugaredLogger {
	return a.loger
}

func (a *app) ConfigPath() string {
	return `.nevakod\4zupper`
}

func (a *app) DbPath() string {
	return `.`
}
