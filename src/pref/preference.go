package pref

import (
	"encoding/gob"
	"github.com/deckarep/golang-set"
	"os"
	"sync"
)

var (
	// basePath indicates the path of where the Preferences files is stored.
	basePath = "./"
	// prefLock is used for synchronization of creating a Preferences.
	prefLock = new(sync.Mutex)
	// prefMap keeps a map of Preferencess with their name as key.
	prefMap = make(map[string]*Preferences)
)

// Preferences is a basic struct for store/access data to/from memory and storage.
type Preferences struct {
	keyMap    map[string]interface{}
	name      string
	observers mapset.Set
	*sync.Mutex
}

// Editor is a modifier of Preferences.
type Editor struct {
	modified map[string]interface{}
	pref     *Preferences
}

// InitBasePath should be called before GetPreferences to initialize the default storage path.
func InitBasePath(path string) {
	basePath = path
}

// RegisterCustomType registers a custom type for serialize and de-serialize the data, must be called at start.
func RegisterCustomType(value interface{}) {
	gob.Register(value)
}

// GetPreferences gets or creates an instance of Preferences with a given name.
func GetPreferences(name string) *Preferences {
	if _, exist := prefMap[name]; !exist {
		prefLock.Lock()
		if _, exist := prefMap[name]; !exist {
			prefMap[name] = &Preferences{
				keyMap:    make(map[string]interface{}),
				name:      name,
				observers: mapset.NewSet(),
				Mutex:     &sync.Mutex{}}
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

// RegisterObserver registers a channel for listening the changes of a preference.
func (p *Preferences) RegisterObserver(observer chan string) {
	if observer != nil {
		p.observers.Add(observer)
	}
}

// UnregisterObserver unregisters a channel, and caller needs to close the channel after unregister.
func (p *Preferences) UnregisterObserver(observer chan string) {
	if observer != nil {
		p.observers.Remove(observer)
	}
}

// GetInt returns the int value from memory.
func (p *Preferences) GetInt(key string) int {
	return p.GetIntOrDefault(key, 0)
}

// GetIntOrDefault returns the int value from memory, and return default value if the key has not been set.
func (p *Preferences) GetIntOrDefault(key string, defaultValue int) int {
	p.Lock()
	defer p.Unlock()
	obj, exist := p.keyMap[key]
	if !exist {
		return defaultValue
	}
	return obj.(int)
}

// GetFloat returns the float value from memory.
func (p *Preferences) GetFloat(key string) float64 {
	return p.GetFloatOrDefault(key, 0)
}

// GetFloatOrDefault returns the float value from memory, and return default value if the key has not been set.
func (p *Preferences) GetFloatOrDefault(key string, defaultValue float64) float64 {
	p.Lock()
	defer p.Unlock()
	obj, exist := p.keyMap[key]
	if !exist {
		return defaultValue
	}
	return obj.(float64)
}

// GetBool returns the bool value from memory.
func (p *Preferences) GetBool(key string) bool {
	return p.GetBoolOrDefault(key, false)
}

// GetBoolOrDefault returns the bool value from memory, and return default value if the key has not been set.
func (p *Preferences) GetBoolOrDefault(key string, defaultValue bool) bool {
	p.Lock()
	defer p.Unlock()
	obj, exist := p.keyMap[key]
	if !exist {
		return defaultValue
	}
	return obj.(bool)
}

// GetString returns the string value from memory.
func (p *Preferences) GetString(key string) string {
	return p.GetStringOrDefault(key, "")
}

// GetStringOrDefault returns the string value from memory, and return default value if the key has not been set.
func (p *Preferences) GetStringOrDefault(key string, defaultValue string) string {
	p.Lock()
	defer p.Unlock()
	obj, exist := p.keyMap[key]
	if !exist {
		return defaultValue
	}
	return obj.(string)
}

// GetObject returns the object value from memory.
func (p *Preferences) GetObject(key string) interface{} {
	return p.GetObjectOrDefault(key, nil)
}

// GetObjectOrDefault returns the object value from memory, and return default value if the key has not been set.
func (p *Preferences) GetObjectOrDefault(key string, defaultValue interface{}) interface{} {
	p.Lock()
	defer p.Unlock()
	obj, exist := p.keyMap[key]
	if !exist {
		return defaultValue
	}
	return obj
}

// Edit creates an editor to modify the value of Preferences.
func (p *Preferences) Edit() *Editor {
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

// PutObject sets the modified object value in editor.
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
			// A nil value in modified map indicates the Preferences shall be removed.
			delete(e.pref.keyMap, k)
		} else {
			e.pref.keyMap[k] = v
		}
		e.pref.tryNotify(k)
	}
}

func (p *Preferences) tryNotify(key string) {
	// Recover if caller close the channel before unregister it.
	defer func() { recover() }()
	for ob := range p.observers.Iter() {
		select {
		case ob.(chan string) <- key:
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
