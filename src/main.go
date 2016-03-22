package main

import (
	"fmt"
	"time"
)

func main() {

	for i := 0; i < 30; i++ {
		TestReadWrite()
	}
	go TestChannel()
	time.Sleep(1)

}

func TestReadWrite() {
	p := GetPreference("pref1")
	p.Edit().PutInt("key1", 3).PutString("key2", "hello").PutFloat("key3", 0.5).PutBool("key4", true).Commit()
	fmt.Printf("Get key1=%d key2=%s key3=%f key4=%v\n", p.GetInt("key1"), p.GetString("key2"), p.GetFloat("key3"), p.GetBool("key4"))
	fmt.Printf("Get key5=%v key4=%s key7=%f key8=%v\n", p.GetInt("key5"), p.GetString("key6"), p.GetFloat("key7"), p.GetBool("key8"))
}

func TestChannel() {

	// key := <-channel
	// fmt.Println("test get channel ", key)
	// p := GetPreference("pref1")
	// for {
	// 	key := <-p.channel
	// 	fmt.Println("receive ", key)
	// }
}
