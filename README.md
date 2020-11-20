# collection
go 集合工具包 


## example:

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
	ret := New().Value(data).Field("Id").IntSlice() // []int{1,2}
}
```
