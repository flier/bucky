package core

import (
	"fmt"
	"reflect"
)

func Assert(target interface{}, format string, args ...interface{}) (ok bool) {
	if err, ok := target.(error); ok && err != nil {
		panic(fmt.Errorf("%s, %s", fmt.Errorf(format, args...), err))
	}

	v := reflect.ValueOf(target)

	if v.IsValid() {
		switch v.Type().Kind() {
		case reflect.Bool:
			ok = v.Bool()
		case reflect.Func, reflect.Interface, reflect.Ptr:
			ok = !v.IsNil()
		case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
			ok = v.Len() > 0
		case reflect.UnsafePointer:
			ok = v.Pointer() != 0
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			ok = v.Int() != 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			ok = v.Uint() != 0
		case reflect.Float32, reflect.Float64:
			ok = v.Float() != 0
		}
	}

	if !ok {
		panic(fmt.Errorf(format, args...))
	}

	return
}
