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
