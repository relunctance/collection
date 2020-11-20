# collection
go 集合工具包  同PHP  array_collumn()

## install 

```
go get -v github.com/relunctance/collection
```


## example Slice:


``` go

type User struct {
	Id   int    `json:"id"`
	Age  int    `json:"age"`
	Name string `json:"name"`
}

func main() {
	data := []User{
		User{1, 32, "hello1"},
		User{2, 33, "hello2"},
	}
	ret := collection.New().Value(data).Field("Id").IntSlice() // []int{1,2}
}

```
## example map:

```go
func main() {
    data := map[string]User{
        "hello1": User{1, 32, "hello1"},
        "hello2": User{2, 32, "hello2"},
    }
    ret := collection.New().Value(data).Field("Id").IntSlice() // []int{1,2}
}
```
## example3:

```go
	data := []User{
		User{1, 32, "hello1"},
		User{2, 32, "hello2"},
		User{3, 32, "hello2"},
		User{4, 32, "hello3"},
	}
	ret := New().Value(data).Field("Name").StringMap()
	/*
		return:  map[string][]interface{} {
			"hello1": interface{User{1, 32, "hello1"}},
			"hello2": interface{User{3, 32, "hello2"}},	// 其中User{2, 32, "hello2"}会被覆盖掉
			"hello3": interface{User{4, 32, "hello3"}},
		}
	*/
```
## example4:
```go
	data := []User{
		User{1, 32, "hello1"},
		User{2, 32, "hello2"},
		User{3, 32, "hello2"},
		User{4, 32, "hello3"},
	}
	ret := New().Value(data).Field("Name").StringMapSlice()
	/*
		return:  map[string][]interface{} {
			"hello1": []interface{User{1, 32, "hello1"}},
			"hello2": []interface{User{2, 32, "hello2"} , User{3, 32, "hello2"}  },
			"hello3": []interface{User{4, 32, "hello3"}},
		}
	*/
```
