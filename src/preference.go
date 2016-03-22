package main

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"sync"
)

var prefLock = new(sync.Mutex)

type observer interface {
	OnChanged(key string)
}

type Preference struct {
	prefMap   map[string]interface{}
	mutex     *sync.Mutex
	observers []observer
}

type editor struct {
	modified map[string]interface{}
}

var instance *Preference

func GetPrefernce() *Preference {
	if instance == nil {
		prefLock.Lock()
		if instance == nil {
			instance = new(Preference)
			instance.prefMap = make(map[string]interface{})
			instance.mutex = new(sync.Mutex)
			bytes, _ := ioutil.ReadFile("sdcard/test.txt")
			json.Unmarshal(bytes, &instance.prefMap)
			instance.observers = make([]observer, 0)
		}
		prefLock.Unlock()
	}
	return instance
}

func (p *Preference) GetInt(key string) int {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	obj, _ := p.prefMap[key]
	value := reflect.ValueOf(obj)
	if value.Kind() == reflect.Int {
		return int(value.Int())
	}
	return 0
}

func (p *Preference) RegisterObserver(o observer) {
	p.observers = append(p.observers, o)
}

func (p *Preference) UnregisterObserver(o observer) {

}

func (p *Preference) Edit() *editor {
	editor := new(editor)
	editor.modified = make(map[string]interface{})
	return editor
}

func (e *editor) PutInt(key string, value int) *editor {
	e.modified[key] = value
	return e
}

func (e *editor) Remove(key string) *editor {
	e.modified[key] = nil
	return e
}

func (e *editor) Apply() {
	e.commitToMemory()
	go e.commitToDisk()
}

func (e *editor) Commit() {
	e.commitToMemory()
	e.commitToDisk()
}

func (e *editor) commitToMemory() {
	instance.mutex.Lock()
	for k, v := range e.modified {
		if v == nil {
			delete(instance.prefMap, k)
		} else {
			instance.prefMap[k] = v
		}
		for _, observer := range instance.observers {
			observer.OnChanged(k)
		}
	}
	instance.mutex.Unlock()
}

func (e *editor) commitToDisk() {
	bytes, _ := json.Marshal(instance.prefMap)
	ioutil.WriteFile("sdcard/test.txt", bytes, 0644)
}
