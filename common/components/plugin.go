package components

import (
	"plugin"
	"pro2d/common/logger"
	"sync"
)

type PluginOption func(*Plugin)

type Plugin struct {
	pluginPath string

	Actions sync.Map
}

func NewPlugin(path string, options ...PluginOption) IPlugin {
	if path == "" {
		return nil
	}
	p := &Plugin{
		pluginPath: path,
		Actions:    sync.Map{},
	}
	for _, option := range options {
		option(p)
	}
	return p
}

func (p *Plugin) LoadPlugin() error {
	plu, err := plugin.Open(p.pluginPath)

	if err != nil {
		return err
	}
	logger.Debug("func LoadPlugin open success...")

	f, err := plu.Lookup("GetActionMap")
	if err != nil {
		return err
	}
	logger.Debug("func LoadPlugin Lookup success...")

	if x, ok := f.(func() map[interface{}]interface{}); ok {
		logger.Debug("func LoadPlugin GetActionMap success...")
		p.SetActions(x())
	}

	return nil
}

func (p *Plugin) GetAction(cmd uint32) interface{} {
	f, ok := p.Actions.Load(cmd)
	if !ok {
		return nil
	}
	return f
}

func (p *Plugin) SetActions(am map[interface{}]interface{}) {
	p.Actions.Range(func(key, value interface{}) bool {
		p.Actions.Delete(key.(uint32))
		return true
	})
	for k, v := range am {
		cmd := k.(uint32)
		p.Actions.Store(cmd, v)
	}
}
