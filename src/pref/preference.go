package pref

import (
	"encoding/gob"
	"os"
	"sync"
)

const channelSize = 4

var (
	basePath = "./"
	prefLock = new(sync.Mutex)
	prefMap  = make(map[string]*Preference)
)

// Preference is a basic struct for store/access data to/from memory and storage.
type Preference struct {
	keyMap  map[string]interface{}
	name    string
	Channel chan string
	*sync.Mutex
}

// Editor is a modifier of Preference.
type Editor struct {
	modified map[string]interface{}
	pref     *Preference
}

// InitBasePath should be called before GetPreference to initialize the default storage path.
func InitBasePath(path string) {
	basePath = path
}

// RegisterCustomType registers a custom type for serialize and de-serialize the data, must be called at start.
func RegisterCustomType(value interface{}) {
	gob.Register(value)
}

// GetPreference gets or creates an instance of Preference with a given name.
func GetPreference(name string) *Preference {
	if _, exist := prefMap[name]; !exist {
		prefLock.Lock()
		if _, exist := prefMap[name]; !exist {
			prefMap[name] = &Preference{
				keyMap:  make(map[string]interface{}),
				name:    name,
				Channel: make(chan string, channelSize),
				Mutex:   &sync.Mutex{}}
			if file, err := os.Open(basePath + name); err == nil {
				dec := gob.NewDecoder(file)
				dec.Decode(&prefMap[name].keyMap)
				file.Close()
			}
		}
		prefLock.Unlock()
	}
	return prefMap[name]
}

// GetInt returns the int value from memory.
func (p *Preference) GetInt(key string) int {
	return p.GetIntOrDefault(key, 0)
}

// GetIntOrDefault returns the int value from memory, and return default value if the key has not been set.
func (p *Preference) GetIntOrDefault(key string, defaultValue int) int {
	p.Lock()
	defer p.Unlock()
	obj, exist := p.keyMap[key]
	if !exist {
		return defaultValue
	}
	return obj.(int)
}

// GetFloat returns the float value from memory.
func (p *Preference) GetFloat(key string) float64 {
	return p.GetFloatOrDefault(key, 0)
}

// GetFloatOrDefault returns the float value from memory, and return default value if the key has not been set.
func (p *Preference) GetFloatOrDefault(key string, defaultValue float64) float64 {
	p.Lock()
	defer p.Unlock()
	obj, exist := p.keyMap[key]
	if !exist {
		return defaultValue
	}
	return obj.(float64)
}

// GetBool returns the bool value from memory.
func (p *Preference) GetBool(key string) bool {
	return p.GetBoolOrDefault(key, false)
}

// GetBoolOrDefault returns the bool value from memory, and return default value if the key has not been set.
func (p *Preference) GetBoolOrDefault(key string, defaultValue bool) bool {
	p.Lock()
	defer p.Unlock()
	obj, exist := p.keyMap[key]
	if !exist {
		return defaultValue
	}
	return obj.(bool)
}

// GetString returns the string value from memory.
func (p *Preference) GetString(key string) string {
	return p.GetStringOrDefault(key, "")
}

// GetStringOrDefault returns the string value from memory, and return default value if the key has not been set.
func (p *Preference) GetStringOrDefault(key string, defaultValue string) string {
	p.Lock()
	defer p.Unlock()
	obj, exist := p.keyMap[key]
	if !exist {
		return defaultValue
	}
	return obj.(string)
}

// GetObject returns the object value from memory.
func (p *Preference) GetObject(key string) interface{} {
	return p.GetObjectOrDefault(key, nil)
}

// GetObjectOrDefault returns the object value from memory, and return default value if the key has not been set.
func (p *Preference) GetObjectOrDefault(key string, defaultValue interface{}) interface{} {
	p.Lock()
	defer p.Unlock()
	obj, exist := p.keyMap[key]
	if !exist {
		return defaultValue
	}
	return obj
}

// Edit creates an editor to modify the value of preference.
func (p *Preference) Edit() *Editor {
	return &Editor{modified: make(map[string]interface{}), pref: p}
}

// PutInt sets the modified int value in editor.
func (e *Editor) PutInt(key string, value int) *Editor {
	e.modified[key] = value
	return e
}

// PutFloat sets the modified float value in editor.
func (e *Editor) PutFloat(key string, value float64) *Editor {
	e.modified[key] = value
	return e
}

// PutBool sets the modified bool value in editor.
func (e *Editor) PutBool(key string, value bool) *Editor {
	e.modified[key] = value
	return e
}

// PutString sets the modified string value in editor.
func (e *Editor) PutString(key string, value string) *Editor {
	e.modified[key] = value
	return e
}

func (e *Editor) PutObject(key string, value interface{}) *Editor {
	e.modified[key] = value
	return e
}

// Remove sets the nil value in editor.
func (e *Editor) Remove(key string) *Editor {
	e.modified[key] = nil
	return e
}

// Apply submits the changes to memory synchronously and submit the changes to disk later.
func (e *Editor) Apply() {
	e.pref.Lock()
	e.commitToMemory()
	e.pref.Unlock()
	go e.commitToDisk()
}

// Commit submits the changes to memory and disk synchronously.
func (e *Editor) Commit() {
	e.pref.Lock()
	e.commitToMemory()
	e.commitToDisk()
	e.pref.Unlock()
}

func (e *Editor) commitToMemory() {
	for k, v := range e.modified {
		if v == nil {
			// A nil value in modified map indicates the preference shall be removed.
			delete(e.pref.keyMap, k)
		} else {
			e.pref.keyMap[k] = v
		}
		// Try push changed keys to channel and not block the go routine
		select {
		case e.pref.Channel <- k:
		default:
		}
	}
}

func (e *Editor) commitToDisk() {
	if file, err := os.Create(basePath + e.pref.name); err == nil {
		enc := gob.NewEncoder(file)
		enc.Encode(&e.pref.keyMap)
		file.Close()
	}
}
