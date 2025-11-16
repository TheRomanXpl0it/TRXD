package plugins

import (
	"fmt"

	"github.com/yuin/gopher-lua"
)

func goToLua(L *lua.LState, v any) (lua.LValue, error) {
	switch x := v.(type) {
	case nil:
		return lua.LNil, nil
	case string:
		return lua.LString(x), nil
	case float64:
		return lua.LNumber(x), nil
	case bool:
		return lua.LBool(x), nil
	case []any:
		t := L.NewTable()
		for _, item := range x {
			val, err := goToLua(L, item)
			if err != nil {
				return nil, err
			}
			t.Append(val)
		}
		return t, nil
	case map[string]any:
		t := L.NewTable()
		for k, item := range x {
			val, err := goToLua(L, item)
			if err != nil {
				return nil, err
			}
			t.RawSetString(k, val)
		}
		return t, nil
	default:
		return nil, fmt.Errorf("Unsupported type: %T", v)
	}
}

func luaToGo(L *lua.LState, v lua.LValue, oldV lua.LValue) (any, error) {
	if v.Type() != oldV.Type() {
		return nil, fmt.Errorf("Type mismatch between prototype and returned value: got [%s] expected [%s]",v.Type().String(),oldV.Type().String())
	}
	switch x := v.(type) {
	case *lua.LNilType:
		return nil, nil
	case lua.LString:
		return string(x), nil
	case lua.LNumber:
		return float64(x), nil
	case lua.LBool:
		return bool(x), nil
	case *lua.LTable:
		// Check if the table is an array or a map
		isList := true
		x.ForEach(func(k, v lua.LValue) {
			// list if keys are 1..n
			if _, ok := k.(lua.LNumber); !ok {
				isList = false
			}
		})

		if !isList {
			// Check if all keys are strings
			neitherMapNorList := false
			x.ForEach(func(k, v lua.LValue) {
				if _, ok := k.(lua.LString); !ok {
					neitherMapNorList = true
				}
			})
			if neitherMapNorList {
				return nil, fmt.Errorf("Unsupported table with mixed keys")
			}
		}

		var errDuringForEach error = nil
		if isList {
			var result []any
			x.ForEach(func(k, v lua.LValue) {
				if errDuringForEach != nil {
					return
				}
				oldInnerV := oldV.(*lua.LTable).RawGet(k)
				val, err := luaToGo(L, v, oldInnerV)
				if err != nil {
					errDuringForEach = err
					return
				}
				result = append(result, val)
			})
			if errDuringForEach != nil {
				return nil, errDuringForEach
			} else {
				return result, nil
			}
		} else {
			result := make(map[string]any)
			x.ForEach(func(k, v lua.LValue) {
				if errDuringForEach != nil {
					return
				}
				keyStr, ok := k.(lua.LString)
				if !ok {
					// Unreachable due to earlier check
					errDuringForEach = fmt.Errorf("Non-string key in map")
					return
				}
				oldInnerV := oldV.(*lua.LTable).RawGet(k)
				val, err := luaToGo(L, v, oldInnerV)
				if err != nil {
					errDuringForEach = err
					return
				}
				result[string(keyStr)] = val
			})
			if errDuringForEach != nil {
				return nil, errDuringForEach
			} else {
				return result, nil
			}
		}
	default:
		return nil, fmt.Errorf("Unsupported Lua type: %T", v)
	}
}