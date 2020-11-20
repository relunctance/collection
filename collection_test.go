package collection

import (
	"fmt"
	"testing"
)

type User struct {
	Id   int    `json:"id"`
	Age  int    `json:"age"`
	Name string `json:"name"`
}

func TestMapValue(t *testing.T) {

	data := map[string]User{
		"hello1": User{1, 32, "hello1"},
		"hello2": User{2, 32, "hello2"},
	}
	ret := New().Value(data).Field("Id").IntSlice() // []int{1,2}
	if len(ret) != 2 {
		t.Fatalf("should be == 1")
	}
	if ret[0] != 1 || ret[1] != 2 {
		t.Fatalf("should be == []int{1,2}")
	}
}

func TestInputNotSlice(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fatalf("err can not be nil , should be panic ")
		}
	}()
	data := "abc"
	New().Value(data).Field("Id").IntSlice() // 'id' is not exists
}
func TestIsNotStruct(t *testing.T) {
	data := []string{"a", "b"}
	ret := New().Value(data).Field("Id").IntSlice() // 'id' is not exists
	if len(ret) != 0 {
		t.Fatalf("should be == 0")
	}
}

func TestPtrIsNil(t *testing.T) {
	data := []*User{
		nil,
		&User{2, 32, "hello2"},
	}
	ret := New().Value(data).Field("Id").IntSlice() // 'id' is not exists
	if len(ret) != 1 {
		t.Fatalf("should be == 1")
	}
	if ret[0] != 2 {
		t.Fatalf("should be == 2")
	}
}

func TestFieldNotExists(t *testing.T) {
	data := getData()
	ret := New().Value(data).Field("id").IntSlice() // 'id' is not exists
	if len(ret) != 0 {
		t.Fatalf("length should be == 0")
	}
}

func getData() []User {
	return []User{
		User{1, 32, "hello1"},
		User{2, 33, "hello2"},
	}
}

func TestIntSlice(t *testing.T) {
	data := getData()
	ret := New().Value(data).Field("Id").IntSlice() // []int{1,2}
	if ret[0] != 1 || ret[1] != 2 {
		t.Fatalf("should be == []int{1,2}")
	}
}

func TestStringSlice(t *testing.T) {
	data := getData()
	ret := New().Value(data).Field("Name").StringSlice() // []int{1,2}
	if ret[0] != "hello1" || ret[1] != "hello2" {
		t.Fatalf("should be == []int{'hello1','hello2'}")
	}
}
