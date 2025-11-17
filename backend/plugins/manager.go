package plugins

import (
	"fmt"
	"context"
	"path/filepath"

	"github.com/tde-nico/log"
	"github.com/yuin/gopher-lua"
)

type Plugin struct {
	path string
	state *lua.LState
}

type Handler struct {
	pluginIndex uint
	callback *lua.LFunction
}

type Manager struct {
	plugins [] Plugin
	handlers map[string] []Handler
}

var manager Manager

func InitManager() error {
	manager = Manager{
		plugins: [] Plugin {},
		handlers: map[string] []Handler {},
	}
	
	matches, err := filepath.Glob("plugins/*.lua")
	if err != nil {
	    return err
	}
	
	// register a new plugin
	for index, path := range matches {
	    log.Info("Found plugin: ","path",path)
		l := lua.NewState()
		l.PreloadModule("trxd",Loader)
		// l.SetGlobal("_trxd_plugin_index",index)
		registry := l.Get(lua.RegistryIndex).(*lua.LTable)
		registry.RawSetString("trxd_plugin_index",lua.LNumber(index))
		plugin := Plugin{
			path: path,
			state: l,
		}
		manager.plugins = append(manager.plugins, plugin)
		err = l.DoFile(path)
		if err != nil {
			return err
		}
	}
	
	pluginsLoaded := map[string] any {}
	pluginsLoaded["plugins"] = []string {}
	for _, value := range (manager.plugins) {
		pluginsLoaded["plugins"] = value.path
	}
	
	pluginsLoaded, err = DispatchEvent(context.TODO(),"pluginsLoaded",pluginsLoaded)
	if err != nil {
		log.Error("Error executing plugins:","err",err)
	}
	
	return nil
}

func DestroyManager() {
	for _, plugin := range manager.plugins {
		plugin.state.Close()
	}
}


func registerHandler(eventName string, luaFunction *lua.LFunction, pluginIndex int) error {
	
	if luaFunction.Proto.NumParameters != 1 {
		return fmt.Errorf("Function num parameters, expected: %d, got: %d",1,luaFunction.Proto.NumParameters)
	}
	
	value, ok := manager.handlers[eventName]
	handler := Handler{pluginIndex: uint(pluginIndex), callback: luaFunction}
	
	if ok {
		manager.handlers[eventName] = append(value, handler)
	} else {
		manager.handlers[eventName] = [] Handler {handler}
	}
	
	return nil
}


// Send the information about the event to the plugins that manage that event
func DispatchEvent(c context.Context, event string, parameters map [string] any) (map[string] any, error) {
	handlers, ok := manager.handlers[event]
	if ok {
		for _, handler := range handlers {
			plugin := manager.plugins[handler.pluginIndex]
			interpreter := plugin.state
			interpreter.SetContext(c)
			parameterTable := interpreter.NewTable()
			backupParameterTable := interpreter.NewTable()
			
			for key, value := range parameters {
				value, err := goToLua(interpreter,value)
				if err != nil {
					return nil, err
				}
				parameterTable.RawSetString(key,value)
				backupParameterTable.RawSetString(key,value)
			}
			
			err := interpreter.CallByParam(
				lua.P{
					Fn: handler.callback,
					NRet: 1,
				    Protect: true,
				}, 
				parameterTable,
			)
			if err != nil {
				log.Error("Error executing plugin:","path",plugin.path,"err",err)
				continue
			}
			
			ret := interpreter.Get(-1).(*lua.LTable)
			interpreter.Pop(1)
			
			oldLuaParams := parameters

			for key := range parameters {
				value, err := luaToGo(interpreter,ret.RawGetString(key),backupParameterTable.RawGetString(key))
				if err != nil {
					log.Error("Error converting plugin result:","path",plugin.path,"err",err)
					parameters = oldLuaParams
					break
				}
				parameters[key] = value
			}
		}
	}
	return parameters, nil
}