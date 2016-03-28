package main

import (
	// "encoding/gob"
	"fmt"
	"pref"
	"time"
)

type Orz struct {
	str string
}

// func (o *Orz) call() {
// 	o.str = "123"
// }

func (o Orz) call() {

}

type Qo interface {
	call()
}

func main() {
	orz := new(Orz)
	orz.str = "qq"
	fmt.Println("before ", orz)

	var cc interface{} = *orz
	cc.call()
	fmt.Println("assign ", cc, orz)

	// gob.Register(new(Test))
	// go TestChannel()
	// for i := 0; i < 10; i++ {
	// 	TestReadWrite(i)
	// }

	// // for i := 0; i < 3; i++ {
	// // 	TestReadWrite()
	// // }
	// time.Sleep(1000000)
}

type Test struct {
	Hello string
	Qq    int
}

func TestReadWrite(i int) {
	t := new(Test)
	t.Hello = "hello"
	t.Qq = 4243
	time.Sleep(10000)

	p := pref.NewPreferences("pref1")
	p.Edit().Put("key1", i).Put("key2", "hello").Put("key3", 0.5).Put("key4", true).Put("obj", t).Apply()
	fmt.Printf("Get key1=%d key2=%s key3=%f key4=%v obj=%v\n", p.GetInt("key1", 0), p.GetString("key2", ""), p.GetFloat64("key3", 0), p.GetBool("key4", false), p.GetObject("obj", nil))
	fmt.Printf("Get key5=%v key4=%s key7=%f key8=%v key9=%v\n", p.GetInt("key5", 0), p.GetString("key6", ""), p.GetFloat64("key7", 0), p.GetBool("key8", false), p.GetObject("key9", nil))
}

func TestChannel() {
	ch := make(chan string, 10)
	p := pref.NewPreferences("pref1")
	p.RegisterOnPreferenceChangeListener(ch)
	for i := 0; i < 15; i++ {
		key := <-ch
		fmt.Println("receive ", key)
	}
	p.UnregisterOnPreferenceChangeListener(ch)
	for i := 0; i < 15; i++ {
		key := <-ch
		fmt.Println("receive2 ", key)
	}
	close(ch)
	// for {
	// 	key := <-p.Channel
	// 	fmt.Println("receive ", key)
	// }
}
