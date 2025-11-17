package plugins

import (
	"fmt"
	"reflect"

	"github.com/yuin/gopher-lua"
)

func goToLua(L *lua.LState, v any) (lua.LValue, error) {
    if v == nil {
        return lua.LNil, nil
    }

    switch x := v.(type) {
	    case lua.LValue:
	        return x, nil
	    case string:
	        return lua.LString(x), nil
	    case bool:
	        return lua.LBool(x), nil
    }

    rv := reflect.ValueOf(v)
    rt := rv.Type()

    switch rv.Kind() {
	    case reflect.String:
	        return lua.LString(rv.String()), nil
									
	    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	        return lua.LNumber(rv.Int()), nil
	
	    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	        return lua.LNumber(rv.Uint()), nil
	
	    case reflect.Float32, reflect.Float64:
	        return lua.LNumber(rv.Float()), nil
	
	    case reflect.Slice, reflect.Array:
	        t := L.NewTable()
	        for i := 0; i < rv.Len(); i++ {
	            elem := rv.Index(i).Interface()
	            val, err := goToLua(L, elem)
	            if err != nil {
	                return nil, err
	            }
	            t.RawSetInt(i+1, val)
	        }
	        return t, nil
	
	    case reflect.Map:
	        if rt.Key().Kind() != reflect.String {
	            return nil, fmt.Errorf("unsupported map key type: %s", rt.Key())
	        }
	        t := L.NewTable()
	        for _, key := range rv.MapKeys() {
	            k := key.String()
	            valIface := rv.MapIndex(key).Interface()
	            val, err := goToLua(L, valIface)
	            if err != nil {
	                return nil, err
	            }
	            t.RawSetString(k, val)
	        }
	        return t, nil
	
	    case reflect.Struct:
	        t := L.NewTable()
	        for i := 0; i < rt.NumField(); i++ {
	            f := rt.Field(i)
	            if f.PkgPath != "" { // unexported
	                continue
	            }
	            name := f.Name
	            // if tag := f.Tag.Get("lua"); tag != "" { name = tag }
	
	            fv := rv.Field(i).Interface()
	            val, err := goToLua(L, fv)
	            if err != nil {
	                return nil, fmt.Errorf("field %s: %w", name, err)
	            }
	            t.RawSetString(name, val)
	        }
	        return t, nil
	
	    case reflect.Pointer:
	        if rv.IsNil() {
	            return lua.LNil, nil
	        }
	        return goToLua(L, rv.Elem().Interface())
    }

    return nil, fmt.Errorf("Unsupported type: %T", v)
}


func luaToGo(L *lua.LState, v lua.LValue, backup lua.LValue, expected reflect.Type) (any, error) {
	if expected == nil {
        return nil, fmt.Errorf("luaToGo called with nil expected type")
    }

    switch expected.Kind() {
    case reflect.String:
        if s, ok := v.(lua.LString); ok {
            return string(s), nil
        }
        return nil, fmt.Errorf("expected string, got %s", v.Type().String())

    case reflect.Bool:
        if b, ok := v.(lua.LBool); ok {
            return bool(b), nil
        }
        return nil, fmt.Errorf("expected bool, got %s", v.Type().String())
        
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        ln, ok := v.(lua.LNumber)
        if !ok {
            return nil, fmt.Errorf("expected number, got %s", v.Type().String())
        }
        i64 := int64(ln)
        // convert to the concrete int type
        return reflect.ValueOf(i64).Convert(expected).Interface(), nil

    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
        ln, ok := v.(lua.LNumber)
        if !ok {
            return nil, fmt.Errorf("expected number, got %s", v.Type().String())
        }
        u64 := uint64(ln)
        return reflect.ValueOf(u64).Convert(expected).Interface(), nil

    case reflect.Float32, reflect.Float64:
        ln, ok := v.(lua.LNumber)
        if !ok {
            return nil, fmt.Errorf("expected number, got %s", v.Type().String())
        }
        f64 := float64(ln)
        return reflect.ValueOf(f64).Convert(expected).Interface(), nil


    case reflect.Map:
        // e.g. map[string]any or map[string]Something
        tbl, ok := v.(*lua.LTable)
        if !ok {
            return nil, fmt.Errorf("expected map (table), got %s", v.Type().String())
        }
        if expected.Key().Kind() != reflect.String {
            return nil, fmt.Errorf("only map[string]... supported, got key %v", expected.Key())
        }

        result := reflect.MakeMapWithSize(expected, tbl.Len())
        elemType := expected.Elem()

        tbl.ForEach(func(k, lv lua.LValue) {
            ks, ok := k.(lua.LString)
            if !ok {
                return
            }
            ev, err := luaToGo(L, lv, lua.LNil, elemType)
            if err != nil {
                return
            }
            result.SetMapIndex(
                reflect.ValueOf(string(ks)),
                reflect.ValueOf(ev).Convert(elemType),
            )
        })

        return result.Interface(), nil

    case reflect.Slice:
        tbl, ok := v.(*lua.LTable)
        if !ok {
            return nil, fmt.Errorf("expected slice (table), got %s", v.Type().String())
        }
        elemType := expected.Elem()
        length := tbl.Len()
        slice := reflect.MakeSlice(expected, length, length)
        for i := 1; i <= length; i++ {
            lv := tbl.RawGetInt(i)
            ev, err := luaToGo(L, lv, lua.LNil, elemType)
            if err != nil {
                return nil, fmt.Errorf("slice element %d: %w", i, err)
            }
            slice.Index(i-1).Set(reflect.ValueOf(ev).Convert(elemType))
        }
        return slice.Interface(), nil

    case reflect.Struct:
        tbl, ok := v.(*lua.LTable)
        if !ok {
            return nil, fmt.Errorf("expected struct (table), got %s", v.Type().String())
        }
        res := reflect.New(expected).Elem()
        for i := 0; i < expected.NumField(); i++ {
            field := expected.Field(i)
            fieldName := field.Name // TODO: tag or naming convention
            lv := tbl.RawGetString(fieldName)
            if lv == lua.LNil {
                continue
            }
            fv, err := luaToGo(L, lv, lua.LNil, field.Type)
            if err != nil {
                return nil, fmt.Errorf("field %s: %w", fieldName, err)
            }
            res.Field(i).Set(reflect.ValueOf(fv).Convert(field.Type))
        }
        return res.Interface(), nil

    case reflect.Pointer:
        val, err := luaToGo(L, v, backup, expected.Elem())
        if err != nil {
            return nil, err
        }
        ptr := reflect.New(expected.Elem())
        ptr.Elem().Set(reflect.ValueOf(val).Convert(expected.Elem()))
        return ptr.Interface(), nil

    default:
        return nil, fmt.Errorf("unsupported expected type: %v", expected)
    }
}
