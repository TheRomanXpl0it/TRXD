package plugins

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/tde-nico/log"
	"github.com/vadv/gopher-lua-libs"
	"github.com/yuin/gopher-lua"
)

type Plugin struct {
	path string
	states []*lua.LState
 	mus []*sync.Mutex
}

type Handler struct {
	pluginIndex uint
	callback string
}

type Manager struct {
	plugins [] Plugin
	handlers map[string] []Handler
}

var manager Manager
var INTERPRETER_NUMBER int = 1

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
		lA := []* lua.LState {}
		muA := []* sync.Mutex {}
		for interpreterIndex := range INTERPRETER_NUMBER {
			l := lua.NewState()
			l.PreloadModule("trxd",Loader)
			registry := l.Get(lua.RegistryIndex).(*lua.LTable)
			registry.RawSetString("trxd_plugin_index",lua.LNumber(index))
			registry.RawSetString("trxd_interpreter_index",lua.LNumber(interpreterIndex))
			libs.Preload(l)
			err = l.DoFile(path)
			if err != nil {
				return err
			}
			lA = append(lA, l)
			mutex := sync.Mutex{}
			muA = append(muA, &mutex)
		}
		// load the built in functions
		plugin := Plugin{
			path: path,
			states: lA,
			mus: muA,
		}
		manager.plugins = append(manager.plugins, plugin)
	}
	
	pluginsLoaded := map[string] any {}
	pluginsLoaded["plugins"] = []string {}
	for _, value := range (manager.plugins) {
		pluginsLoaded["plugins"] = append(pluginsLoaded["plugins"].([]string),value.path)
	}
	
	pluginsLoaded, err = DispatchEvent(context.TODO(),"pluginsLoaded",pluginsLoaded)
	if err != nil {
		log.Error("Error executing plugins:","err",err)
	}
	
	return nil
}

func DestroyManager() {
	for _, plugin := range manager.plugins {
		for _, interpreters := range plugin.states {	
			interpreters.Close()
		}
	}
}


func registerHandler(eventName string, luaFunction string, pluginIndex int, l *lua.LState) error {
	
	luaFunctionPointer, ok := l.GetGlobal(luaFunction).(*lua.LFunction)
	if ! ok{
		return fmt.Errorf("Expected function, got:",luaFunctionPointer.Type().String())
		
	}
	if luaFunctionPointer.Proto.NumParameters != 1 {
		return fmt.Errorf("Function num parameters, expected: %d, got: %d",1,luaFunctionPointer.Proto.NumParameters)
	}
	
	value, ok := manager.handlers[eventName]
	handler := Handler{pluginIndex: uint(pluginIndex), callback: luaFunction}
	
	if ok {
		skipAppend:= false
		for _,existingHandler := range manager.handlers[eventName] {
			if existingHandler.callback == handler.callback && existingHandler.pluginIndex == handler.pluginIndex {
				skipAppend = true
				break
			} 
		}
		if !skipAppend {
			manager.handlers[eventName] = append(value, handler)
		}
	} else {
		manager.handlers[eventName] = [] Handler {handler}
	}
	
	return nil
}


// // Send the information about the event to the plugins that manage that event
// func DispatchEvent(c context.Context, event string, parameters map [string] any) (map[string] any, error) {
// 	handlers, ok := manager.handlers[event]
// 	if ok {
// 		for _, handler := range handlers {
// 			plugin := manager.plugins[handler.pluginIndex]
// 			interpreter := plugin.state
// 			interpreter.SetContext(c)
// 			parameterTable := interpreter.NewTable()
// 			backupParameterTable := interpreter.NewTable()
			
// 			for key, value := range parameters {
// 				value, err := goToLua(interpreter,value)
// 				if err != nil {
// 					return nil, err
// 				}
// 				parameterTable.RawSetString(key,value)
// 				backupParameterTable.RawSetString(key,value)
// 			}
			
// 			err := interpreter.CallByParam(
// 				lua.P{
// 					Fn: handler.callback,
// 					NRet: 1,
// 				    Protect: true,
// 				}, 
// 				parameterTable,
// 			)
// 			if err != nil {
// 				log.Error("Error executing plugin:","path",plugin.path,"err",err)
// 				continue
// 			}
			
// 			ret := interpreter.Get(-1).(*lua.LTable)
// 			interpreter.Pop(1)
			
// 			oldLuaParams := parameters

// 			for key := range parameters {
// 				value, err := luaToGo(interpreter,ret.RawGetString(key),backupParameterTable.RawGetString(key))
// 				if err != nil {
// 					log.Error("Error converting plugin result:","path",plugin.path,"err",err)
// 					parameters = oldLuaParams
// 					break
// 				}
// 				parameters[key] = value
// 			}
// 		}
// 	}
// 	return parameters, nil
// }

func DispatchEvent[T any](c context.Context, event string, payload T) (T, error) {

    out, err := dispatchEventRaw(c, event, payload, reflect.TypeOf(payload))
    if err != nil {
        return payload, err
    }

    v, ok := out.(T)
    if !ok {
        return payload, fmt.Errorf("plugin '%s' returned incompatible type %T, expected %T",
            event, out, payload)
    }

    return v, nil
}

func retrieveLock(locks []*sync.Mutex) int {
	index := 0
	for true {
		lock := locks[index]
		success := lock.TryLock()
		if success {
			return index
		}
		index++
		if index == len(locks){
			time.Sleep(1*time.Millisecond)
		}
		index %= len(locks)
	}
	return 0 //cannot reach
}


func dispatchEventRaw(
    c context.Context,
    event string,
    payload any,
    expectedType reflect.Type,
) (any, error) {
    handlers, ok := manager.handlers[event]
    if !ok {
        return payload, nil
    }

    for _, handler := range handlers {
        plugin := manager.plugins[handler.pluginIndex]
        interpreterIdx := retrieveLock(plugin.mus)
         	
        interpreter := plugin.states[interpreterIdx]
        interpreter.SetContext(c)
        luaPayload, err := goToLua(interpreter, payload)
        // Error on the go side, the type given was not convertible to Lua type
        if err != nil {
            return nil, err
        }
        backupPayload, err := goToLua(interpreter, payload)
        if err != nil {
            return nil, err
        }

        err = interpreter.CallByParam(
            lua.P{
                Fn:      interpreter.GetGlobal(handler.callback).(*lua.LFunction),
                NRet:    1,
                Protect: true,
            },
            luaPayload,
        )
        if err != nil {
            log.Error("Error executing plugin:", "path", plugin.path, "err", err)
            continue
        }


        ret := interpreter.Get(-1)
        interpreter.Pop(1)

        converted, err := luaToGo(
            interpreter,
            ret,
            backupPayload,
            expectedType,
        )
        if err != nil {
            log.Error("Error converting plugin result:", "path", plugin.path, "err", err)
            break
        }

        plugin.mus[interpreterIdx].Unlock()
        payload = converted
    }

    return payload, nil
}