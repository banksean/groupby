// Package groupby implements grouping for slices of interface{}.
// It is not intended to be performant or type-safe, just generic and easy to use.
//
// While Field and Func both assume each value in in is to be grouped by a single key value,
// FuncChan can be used in situations where a value in in may be grouped into multiple, or
// no groups.
package groupby

import "reflect"

// Field uses reflection to group structs by values of fieldName, which must be a valid field name on each struct in in.
func Field(in interface{}, fieldName string) map[interface{}][]interface{} {
	ret := map[interface{}][]interface{}{}
	inval := reflect.ValueOf(in)
	for i := 0; i < inval.Len(); i++ {
		iv := inval.Index(i)
		v := reflect.Indirect(iv)
		fv := v.FieldByName(fieldName)
		key := reflect.Indirect(fv).Interface()
		ret[key] = append(ret[key], iv.Interface())
	}
	return ret
}

// Func groups structs by values returned by keyFunc, which must return a type that may be used as a key in a map.
func Func(in interface{}, keyFunc func(i interface{}) interface{}) map[interface{}][]interface{} {
	ret := map[interface{}][]interface{}{}
	inval := reflect.ValueOf(in)
	for i := 0; i < inval.Len(); i++ {
		iv := inval.Index(i)
		v := reflect.Indirect(iv)
		key := keyFunc(v.Interface())
		ret[key] = append(ret[key], iv.Interface())
	}
	return ret
}

// FuncChan groups structs by values sent on a channel returned by keyFunc. It assumes keyFunc closes the chan it returns.
func FuncChan(in interface{}, keyFunc func(i interface{}) chan interface{}) map[interface{}][]interface{} {
	ret := map[interface{}][]interface{}{}
	inval := reflect.ValueOf(in)
	for i := 0; i < inval.Len(); i++ {
		iv := inval.Index(i)
		v := reflect.Indirect(iv)
		for key := range keyFunc(v.Interface()) {
			ret[key] = append(ret[key], iv.Interface())
		}
	}
	return ret
}
