package collection

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"sync"
)

var collection *Collection

type Collection struct {
	data       interface{}
	field      string
	uniqueFlag bool
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
	if collection == nil { // 单例模式
		return &Collection{} // 默认unique = false
	}
	collection.Init()
	return collection
}

func (s *Collection) Init() {
	s.data = nil
	s.field = ""
	s.uniqueFlag = false
	return
}

func (s *Collection) Value(data interface{}) *Collection {
	vk := reflect.ValueOf(data).Kind()
	if vk != reflect.Slice && vk != reflect.Map {
		panic("must Slice or Map")
	}
	s.uniqueFlag = false // 每次重新设置值时, 默认不进行去重
	s.data = data
	return s
}

func (s *Collection) Field(field string) *Collection {
	s.field = field
	return s
}

// 去重
func (s *Collection) Unique() *Collection {
	s.uniqueFlag = true
	return s
}

func (s *Collection) Data() interface{} {
	return s.data
}

// 返回根据Field取到的元素数组
func (s *Collection) Slice() []interface{} {
	return s.buildSlice()
}

func (s *Collection) getValue(v reflect.Value) interface{} {
	return getValueByKey(v, s.field)
}

// 类型转换
func (s *Collection) IntSlice() (ret []int) {
	slice := s.buildSlice()
	ret = make([]int, 0, len(slice))
	for _, v := range slice {
		if vint, err := interface2Int(v); err == nil {
			ret = append(ret, vint)
		}
	}
	if s.uniqueFlag {
		ret = SliceIntUnique(ret)
	}
	return
}

func (s *Collection) StringSlice() (ret []string) {
	slice := s.buildSlice()
	ret = make([]string, 0, len(slice))
	for _, v := range slice {
		vstr := interface2String(v)
		ret = append(ret, vstr)
	}
	if s.uniqueFlag {
		ret = SliceStringUnique(ret)
	}
	return
}

// data_trun_key();
func (s *Collection) IntMapSlice() (ret map[int][]interface{}) {
	mapdata := s.dataTrunMulti()
	ret = make(map[int][]interface{}, len(mapdata))
	for k, items := range mapdata {
		if vint, err := interface2Int(k); err == nil {
			ret[vint] = items
		}
	}
	return
}
func interface2Int(v interface{}) (int, error) {
	switch val := v.(type) {
	case int:
		return val, nil
	case string:
		return strconv.Atoi(val)
	}
	return 0, fmt.Errorf("can not ParseInt success")
}

// data_trun_key();
func (s *Collection) StringMapSlice() (ret map[string][]interface{}) {
	mapdata := s.dataTrunMulti()
	ret = make(map[string][]interface{}, len(mapdata))
	for k, items := range mapdata {
		ret[interface2String(k)] = items
	}
	return
}

func interface2String(v interface{}) string {
	var key string
	switch val := v.(type) {
	case int:
		key = strconv.Itoa(val)
	case string:
		key = val
	default:
		key = fmt.Sprintf("%v", val)
	}
	return key
}

// data_trun_key();
func (s *Collection) IntMap() (ret map[int]interface{}) {
	mapdata := s.dataTrunValue()
	ret = make(map[int]interface{}, len(mapdata))
	for k, items := range mapdata {
		if vint, err := interface2Int(k); err == nil {
			ret[vint] = items
		}
	}
	return
}

// data_trun_key();
func (s *Collection) StringMap() (ret map[string]interface{}) {
	mapdata := s.dataTrunValue()
	ret = make(map[string]interface{}, len(mapdata))
	for k, items := range mapdata {
		ret[interface2String(k)] = items
	}
	return
}

// 仅用于定义类型
// valKey 为Field("field") 取到的对应item的value值
// itemInterfaceValue 为对应item
type F func(valKey interface{}, itemInterfaceValue interface{})

func (s *Collection) commonDataTrun(f F) {
	v := reflect.ValueOf(s.data)
	vk := v.Kind()
	l := v.Len()
	if l == 0 {
		return
	}

	switch vk {
	case reflect.Slice:
		for i := 0; i < l; i++ {
			item := v.Index(i)
			if item.CanInterface() {
				valKey := s.getValue(item)
				f(valKey, item.Interface()) // 利用map是指针传递这个特性
			}
		}

	case reflect.Map:
		mapkeys := v.MapKeys()
		for _, idx := range mapkeys {
			item := v.MapIndex(idx)
			if item.CanInterface() {
				valKey := s.getValue(item)
				f(valKey, item.Interface())
			}
		}
	}
	return

}
func (s *Collection) dataTrunValue() (res map[interface{}]interface{}) {
	lock := new(sync.Mutex)
	res = make(map[interface{}]interface{}, 10)
	var f F = func(valKey interface{}, itemInterfaceValue interface{}) {
		if valKey == nil {
			return
		}
		lock.Lock()
		res[valKey] = itemInterfaceValue
		lock.Unlock()
	}
	s.commonDataTrun(f)
	return
}

func (s *Collection) dataTrunMulti() (res map[interface{}][]interface{}) {
	lock := new(sync.Mutex)
	res = make(map[interface{}][]interface{}, 10)
	var f F = func(valKey interface{}, itemInterfaceValue interface{}) {
		if valKey == nil {
			return
		}
		lock.Lock()
		_, ok := res[valKey]
		if !ok {
			res[valKey] = make([]interface{}, 0, 10) //减少分配内存
		}
		res[valKey] = append(res[valKey], itemInterfaceValue)
		lock.Unlock()
	}

	s.commonDataTrun(f)
	return
}

//commonBuild
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
			val := s.getValue(item)
			if val != nil {
				ret = append(ret, val)
			}
		}

	case reflect.Map:
		mapkeys := v.MapKeys()
		for _, idx := range mapkeys {
			item := v.MapIndex(idx)
			val := s.getValue(item)
			if val != nil {
				ret = append(ret, val)
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

func SliceIntUnique(slice []int) []int {
	sort.Ints(slice)
	i := 0
	var j int
	for {
		if i >= len(slice)-1 {
			break
		}
		for j = i + 1; j < len(slice) && slice[i] == slice[j]; j++ {
		}
		slice = append(slice[:i+1], slice[j:]...)
		i++
	}
	return slice
}

func SliceStringUnique(slice []string) []string {
	sort.Strings(slice)
	i := 0
	var j int
	for {
		if i >= len(slice)-1 {
			break
		}
		for j = i + 1; j < len(slice) && slice[i] == slice[j]; j++ {
		}
		slice = append(slice[:i+1], slice[j:]...)
		i++
	}
	return slice

}
