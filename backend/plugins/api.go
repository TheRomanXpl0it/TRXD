package plugins

import (
	"github.com/tde-nico/log"
	"github.com/yuin/gopher-lua"
)

func Loader(L *lua.LState) int {
    mod := L.SetFuncs(L.NewTable(), exports)
    L.SetField(mod, "trxd", lua.LString("un'api bellissima"))
    L.Push(mod)
    return 1
}

var exports = map[string]lua.LGFunction{
    "on_event": onEvent,
}

func onEvent(L *lua.LState) int{
	eventName := L.ToString(1)
	callback := L.ToFunction(2)
	
	registry := L.Get(lua.RegistryIndex).(*lua.LTable)
	pluginIndex := int(lua.LVAsNumber(registry.RawGetString("trxd_plugin_index")))
	err := registerHandler(eventName, callback, pluginIndex)
	if err != nil {
		L.Panic(L)
		log.Error("Error registering plugin:","err",err)
	}
	return 0
}
