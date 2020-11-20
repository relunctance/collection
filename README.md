# collection
go 集合工具包 


## example:


### install 

```
go get -v github.com/relunctance/collection
```


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



```go
func main() {
    data := map[string]User{
        "hello1": User{1, 32, "hello1"},
        "hello2": User{2, 32, "hello2"},
    }
    ret := collection.New().Value(data).Field("Id").IntSlice() // []int{1,2}
}
```
