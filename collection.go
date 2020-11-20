package collection

import (
	"fmt"
	"reflect"
)

type Collection struct {
	data  interface{}
	field string
}

func NewWithValue(data interface{}) *Collection {
	s := New()
	s.Value(data)
	return s
}

// New().Value(data).Field("id").IntSlice()	// []int{1,2}
// New().Value([]*User)
// New().Value([]User)
func New() *Collection {
	return &Collection{}
}

func (s *Collection) Value(data interface{}) *Collection {
	if !IsSlice(data) {
		panic(fmt.Errorf("input value is not slice"))
	}
	s.data = data
	return s
}

func (s *Collection) Field(field string) *Collection {
	s.field = field
	return s
}

func (s *Collection) Data() interface{} {
	return s.data
}

func (s *Collection) getValue(v reflect.Value) interface{} {
	return getValueByKey(v, s.field)
}

// 类型转换
func (s *Collection) IntSlice() (ret []int) {
	slice := s.buildSlice()
	ret = make([]int, 0, len(slice))
	for _, v := range slice {
		ret = append(ret, v.(int))
	}
	return
}
func (s *Collection) StringSlice() (ret []string) {
	slice := s.buildSlice()
	ret = make([]string, 0, len(slice))
	for _, v := range slice {
		ret = append(ret, v.(string))
	}
	return

}

func (s *Collection) buildSlice() (ret []interface{}) {
	v := reflect.ValueOf(s.data)
	vk := v.Kind()
	l := v.Len()
	if l == 0 {
		return
	}
	ret = make([]interface{}, 0, l)
	switch vk {
	case reflect.Slice:
		for i := 0; i < l; i++ {
			item := v.Index(i)
			v := s.getValue(item)
			if v != nil {
				ret = append(ret, v) //
			}
		}
	}
	return
}

// New().Value(data).Field("id").IntSlice()	// []int{1,2}
// New().Value([]*User)
// New().Value([]User)
// New().Value([]string).Field("xxx");
func getValueByKey(item reflect.Value, key string) interface{} {
	k := item.Kind()
	if k == reflect.Ptr {
		if item.IsNil() {
			return nil // 取不到, 就不取了
		}
		item = item.Elem() //转换指针为结构
		k = item.Kind()
	}
	if k != reflect.Struct {
		return nil //不是结构体, 找不到
		// panic(fmt.Errorf("must be struct , you input type is : %s ", k.String()))
	}
	_, ok := item.Type().FieldByName(key)
	if !ok {
		return nil
	}
	if !item.FieldByName(key).CanInterface() {
		return nil
	}
	return item.FieldByName(key).Interface()
}

// 判断传入的元素是否是slice
// 支持检测 []string, []int , []int64 , []in32 , []interface{} , []float64 , []uint32 , []uint64 , ...
func IsSlice(v interface{}) bool {
	return reflect.ValueOf(v).Type().Kind().String() == "slice"
}
