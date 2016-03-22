package main

import (
	"fmt"
)

type test struct{}

func (o *test) OnChanged(key string) {
	fmt.Println("on changed ", key)
}

func main() {
	fmt.Println("hello")
	p := GetPrefernce("qq")
	// t := new(test)
	p.Edit().PutString("key", "3.5").Commit()
	fmt.Println("test get ", p.GetString("key"))
	p.Edit().PutString("key2", "3").Commit()
	fmt.Println("test get2 ", p.GetString("key2"))

	// p = GetPrefernce("orz")
	// // p.Edit().PutInt("ccc", 6).Commit()
	// fmt.Println("test get ", p.GetInt("ccc"))
	// // p.Edit().PutInt("ooo", 7).Commit()
	// fmt.Println("test get2 ", p.GetInt("ooo"))
}
