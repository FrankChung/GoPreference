package main

import (
	"fmt"
	"pref"
	"time"
)

func main() {

	pref.RegisterCustomType(new(Test))
	go TestChannel()
	for i := 0; i < 3; i++ {
		TestReadWrite()
	}

	// for i := 0; i < 3; i++ {
	// 	TestReadWrite()
	// }
	time.Sleep(100)
}

type Test struct {
	Hello string
	Qq    int
}

func TestReadWrite() {
	t := new(Test)
	t.Hello = "hello"
	t.Qq = 4243

	p := pref.GetPreference("pref1")
	p.Edit().PutInt("key1", 3).PutString("key2", "hello").PutFloat("key3", 0.5).PutBool("key4", true).PutObject("obj", t).Commit()
	fmt.Printf("Get key1=%d key2=%s key3=%f key4=%v obj=%v\n", p.GetInt("key1"), p.GetString("key2"), p.GetFloat("key3"), p.GetBool("key4"), p.GetObject("obj"))
	fmt.Printf("Get key5=%v key4=%s key7=%f key8=%v key9=%v\n", p.GetIntOrDefault("key5", 6), p.GetString("key6"), p.GetFloat("key7"), p.GetBool("key8"), p.GetObject("key9"))
}

func TestChannel() {

	p := pref.GetPreference("pref1")
	for {
		key := <-p.Channel
		fmt.Println("receive ", key)
	}
}
