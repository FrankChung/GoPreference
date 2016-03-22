package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello")
	p := GetPrefernce()
	p.Edit().PutInt("key", 3).Commit()
	fmt.Println("test get ", p.GetInt("key"))
}
