package main

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"sync"
)

var filePath = "./"
var prefLock = new(sync.Mutex)
var prefMap = make(map[string]*Preference)

type Preference struct {
	keyMap  map[string]interface{}
	name    string
	channel chan string
	*sync.Mutex
}

type editor struct {
	modified map[string]interface{}
	pref     *Preference
}

func InitPath(path string) {
	filePath = path
}

func GetPreference(name string) *Preference {
	if prefMap[name] == nil {
		prefLock.Lock()
		if prefMap[name] == nil {
			prefMap[name] = &Preference{keyMap: make(map[string]interface{}), name: name, channel: make(chan string, 4), Mutex: &sync.Mutex{}}
			bytes, _ := ioutil.ReadFile(filePath + name + ".json")
			json.Unmarshal(bytes, &prefMap[name].keyMap)
		}
		prefLock.Unlock()
	}
	return prefMap[name]
}

func (p *Preference) GetInt(key string) int {
	p.Lock()
	defer p.Unlock()
	obj, _ := p.keyMap[key]
	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Int:
		return int(value.Int())
	case reflect.Float64:
		return int(value.Float())
	default:
		return 0
	}
}

func (p *Preference) GetFloat(key string) float64 {
	p.Lock()
	defer p.Unlock()
	obj, _ := p.keyMap[key]
	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Int:
		return float64(value.Int())
	case reflect.Float64:
		return value.Float()
	default:
		return 0
	}
}

func (p *Preference) GetBool(key string) bool {
	p.Lock()
	defer p.Unlock()
	obj, _ := p.keyMap[key]
	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Bool:
		return value.Bool()
	default:
		return false
	}
}

func (p *Preference) GetString(key string) string {
	p.Lock()
	defer p.Unlock()
	obj, _ := p.keyMap[key]
	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.String:
		return value.String()
	default:
		return ""
	}
}

func (p *Preference) Edit() *editor {
	editor := &editor{modified: make(map[string]interface{}), pref: p}
	return editor
}

func (e *editor) PutInt(key string, value int) *editor {
	e.modified[key] = value
	return e
}

func (e *editor) PutFloat(key string, value float64) *editor {
	e.modified[key] = value
	return e
}

func (e *editor) PutBool(key string, value bool) *editor {
	e.modified[key] = value
	return e
}

func (e *editor) PutString(key string, value string) *editor {
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
	e.pref.Lock()
	defer e.pref.Unlock()
	for k, v := range e.modified {
		if v == nil {
			delete(e.pref.keyMap, k)
		} else {
			e.pref.keyMap[k] = v
		}
		go func(key string) {
			// try push changed keys to channel and not block the go routine
			select {
			case e.pref.channel <- key:
			default:
			}
		}(k)
	}
}

func (e *editor) commitToDisk() {
	bytes, _ := json.Marshal(e.pref.keyMap)
	ioutil.WriteFile(filePath+e.pref.name+".json", bytes, 0644)
}
