package digcloudlog

import (
	"errors"
	"fmt"
)

type IDCLogger interface {
	Info(log string)
	Error(log string)
}

type CDCLogAdapter struct {
	instList []IDCLogger
	dblog    bool
}

var g_SingleLogAdapter *CDCLogAdapter = &CDCLogAdapter{}

func GetLogAdapter() *CDCLogAdapter {
	return g_SingleLogAdapter
}
func (g *CDCLogAdapter) Initialize() error {
	g.dblog = false
	return nil
}
func (g *CDCLogAdapter) Info(format string, args ...any) {
	strLog := fmt.Sprintf(format, args...)
	for _, inst := range g.instList {
		inst.Info(strLog)
	}
}
func (g *CDCLogAdapter) Error(format string, args ...any) {
	strLog := fmt.Sprintf(format, args...)
	for _, inst := range g.instList {
		inst.Error(strLog)
	}
}
func (g *CDCLogAdapter) AddStdout() {
	getLogStdout().Initialize()
	g.instList = append(g.instList, getLogStdout())
}
func (g *CDCLogAdapter) AddLocalFile(basePath, infoFileName, errorFileName string) error {
	localLog := newLogLocalFile()
	err := localLog.Initialize(basePath, infoFileName, errorFileName)
	if err != nil {
		return err
	}
	g_defaultFileLog = localLog
	g.instList = append(g.instList, localLog)

	return nil
}
func (g *CDCLogAdapter) AddLogflare(sourceid, apiKey string) error {
	logflare := newCLogFlare()
	err := logflare.Initialize(sourceid, apiKey)
	if err != nil {
		return err
	}

	g.instList = append(g.instList, logflare)
	return nil
}

func (g *CDCLogAdapter) AddNewLogService(service interface{}) error {
	tmp, ok := service.(IDCLogger)
	if !ok {
		return errors.New("interface is not IDCLogger")
	}
	g.instList = append(g.instList, tmp)
	return nil
}
