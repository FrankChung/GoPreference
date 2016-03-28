package pref

import (
	"concurrent"
	"encoding/gob"
	"log"
	"os"
	"reflect"
	"sync"
)

var (
	// basePath indicates the path of where the Preferences files is stored.
	basePath = "./"
	// prefLock is used for synchronization of creating a Preferences.
	prefLock = new(sync.Mutex)
	// prefMap keeps a map of Preferencess with their name as key.
	prefMap = make(map[string]*PreferencesImpl)
	// executor is a global executor for executing the go functions sequentially.
	executor = concurrent.SingleExecutor()
)

// PreferencesImpl is a basic struct for store/access data to/from memory and storage.
type PreferencesImpl struct {
	m            map[string]interface{}
	name         string
	observers    map[chan string]interface{}
	writeCh      chan map[string]interface{}
	diskLock     *sync.Mutex
	observerLock *sync.Mutex
	loadWg       *sync.WaitGroup
	*sync.Mutex
}

// EditorImpl is a modifier of Preferences.
type EditorImpl struct {
	modified map[string]interface{}
	pref     *PreferencesImpl
	cleared  bool
	*sync.Mutex
}

// InitBasePath should be called before NewPreferences to initialize the default storage path.
func InitBasePath(path string) {
	basePath = path
}

// NewPreferences gets or creates an instance of Preferences with a given name, remember to call gob.Register
// before writing/reading custom types to/from a Preferences.
func NewPreferences(name string) Preferences {
	prefLock.Lock()
	defer prefLock.Unlock()
	if _, exist := prefMap[name]; !exist {
		pref := &PreferencesImpl{
			m:            make(map[string]interface{}),
			name:         name,
			observers:    make(map[chan string]interface{}),
			writeCh:      make(chan map[string]interface{}, 10),
			diskLock:     &sync.Mutex{},
			observerLock: &sync.Mutex{},
			loadWg:       &sync.WaitGroup{},
			Mutex:        &sync.Mutex{}}
		pref.loadWg.Add(1)
		go pref.loadFromFile()
		prefMap[name] = pref
	}
	return prefMap[name]
}

func (p *PreferencesImpl) loadFromFile() {
	p.Lock()
	defer p.Unlock()
	path := basePath + p.name
	backupPath := path + "_bak"
	// Load backup file if exists.
	if _, err := os.Stat(backupPath); err == nil {
		os.Remove(path)
		os.Rename(backupPath, path)
	}
	if file, err := os.Open(basePath + p.name); err == nil {
		defer file.Close()
		dec := gob.NewDecoder(file)
		if err := dec.Decode(&p.m); err != nil {
			log.Printf("Error when decode the preference: %v", err)
		}
	} else if os.IsExist(err) {
		log.Printf("Error reading the preference file %s", basePath+p.name)
	}
	p.loadWg.Done()
}

// RegisterOnPreferenceChangeListener registers a listener for listening the changes of a preference.
func (p *PreferencesImpl) RegisterOnPreferenceChangeListener(observer OnPreferenceChangeListener) {
	p.observerLock.Lock()
	defer p.observerLock.Unlock()
	if observer != nil {
		p.observers[observer] = nil
	}
}

// UnregisterOnPreferenceChangeListener unregisters a listener, and caller needs to close the channel after unregister.
func (p *PreferencesImpl) UnregisterOnPreferenceChangeListener(observer OnPreferenceChangeListener) {
	p.observerLock.Lock()
	defer p.observerLock.Unlock()
	if observer != nil {
		delete(p.observers, observer)
	}
}

// Contains returns whether a key exists in this preference.
func (p *PreferencesImpl) Contains(key string) bool {
	p.Lock()
	defer p.Unlock()
	p.loadWg.Wait()
	_, exist := p.m[key]
	return exist
}

// GetBool returns the bool value from memory, and return default value if the key has not been set.
func (p *PreferencesImpl) GetBool(key string, defaultValue bool) bool {
	val, ok := p.GetObject(key, defaultValue).(bool)
	if !ok {
		return defaultValue
	}
	return val
}

// GetInt returns the int value from memory, and return default value if the key has not been set.
func (p *PreferencesImpl) GetInt(key string, defaultValue int) int {
	val, ok := p.GetObject(key, defaultValue).(int)
	if !ok {
		return defaultValue
	}
	return val
}

// GetInt32 returns the int32 value from memory, and return default value if the key has not been set.
func (p *PreferencesImpl) GetInt32(key string, defaultValue int32) int32 {
	val, ok := p.GetObject(key, defaultValue).(int32)
	if !ok {
		return defaultValue
	}
	return val
}

// GetInt64 returns the int64 value from memory, and return default value if the key has not been set.
func (p *PreferencesImpl) GetInt64(key string, defaultValue int64) int64 {
	val, ok := p.GetObject(key, defaultValue).(int64)
	if !ok {
		return defaultValue
	}
	return val
}

// GetUInt32 returns the uint32 value from memory, and return default value if the key has not been set.
func (p *PreferencesImpl) GetUInt32(key string, defaultValue uint32) uint32 {
	val, ok := p.GetObject(key, defaultValue).(uint32)
	if !ok {
		return defaultValue
	}
	return val
}

// GetUInt64 returns the uint64 value from memory, and return default value if the key has not been set.
func (p *PreferencesImpl) GetUInt64(key string, defaultValue uint64) uint64 {
	val, ok := p.GetObject(key, defaultValue).(uint64)
	if !ok {
		return defaultValue
	}
	return val
}

// GetFloat32 returns the float32 value from memory, and return default value if the key has not been set.
func (p *PreferencesImpl) GetFloat32(key string, defaultValue float32) float32 {
	val, ok := p.GetObject(key, defaultValue).(float32)
	if !ok {
		return defaultValue
	}
	return val
}

// GetFloat64 returns the float64 value from memory, and return default value if the key has not been set.
func (p *PreferencesImpl) GetFloat64(key string, defaultValue float64) float64 {
	val, ok := p.GetObject(key, defaultValue).(float64)
	if !ok {
		return defaultValue
	}
	return val
}

// GetByte returns the byte value from memory, and return default value if the key has not been set.
func (p *PreferencesImpl) GetByte(key string, defaultValue byte) byte {
	val, ok := p.GetObject(key, defaultValue).(byte)
	if !ok {
		return defaultValue
	}
	return val
}

// GetRune returns the rune value from memory, and return default value if the key has not been set.
func (p *PreferencesImpl) GetRune(key string, defaultValue rune) rune {
	val, ok := p.GetObject(key, defaultValue).(rune)
	if !ok {
		return defaultValue
	}
	return val
}

// GetString returns the string value from memory, and return default value if the key has not been set.
func (p *PreferencesImpl) GetString(key string, defaultValue string) string {
	val, ok := p.GetObject(key, defaultValue).(string)
	if !ok {
		return defaultValue
	}
	return val
}

// GetObject returns the object value from memory, and return default value if the key has not been set.
func (p *PreferencesImpl) GetObject(key string, defaultValue interface{}) interface{} {
	p.loadWg.Wait()
	p.Lock()
	defer p.Unlock()
	obj, exist := p.m[key]
	if !exist {
		return defaultValue
	}
	return obj
}

// Edit creates an editor to modify the value of Preferences.
func (p *PreferencesImpl) Edit() Editor {
	p.loadWg.Wait()
	return &EditorImpl{
		modified: make(map[string]interface{}),
		pref:     p,
		cleared:  false,
		Mutex:    &sync.Mutex{},
	}
}

// Put sets the modified object value in editor.
func (e *EditorImpl) Put(key string, value interface{}) Editor {
	e.Lock()
	defer e.Unlock()
	e.modified[key] = value
	return e
}

// Remove sets the nil value in editor.
func (e *EditorImpl) Remove(key string) Editor {
	e.Lock()
	defer e.Unlock()
	e.modified[key] = nil
	return e
}

// Clear the whole key-values in the preference.
func (e *EditorImpl) Clear() Editor {
	e.Lock()
	defer e.Unlock()
	e.cleared = true
	return e
}

func (p *PreferencesImpl) copyOfMapLocked() map[string]interface{} {
	dst := make(map[string]interface{})
	for k, v := range p.m {
		dst[k] = v
	}
	return dst
}

// Apply submits the changes to memory synchronously and submit the changes to disk later.
func (e *EditorImpl) Apply() {
	e.Lock()
	defer e.Unlock()
	e.pref.Lock()
	defer e.pref.Unlock()
	keys := e.commitToMemoryLocked()
	if len(keys) > 0 {
		executor.Execute(func() {
			e.pref.commitToDisk(e.pref.copyOfMapLocked())
		})
		e.pref.notifyObservers(keys)
	}
}

// Commit submits the changes to memory and disk synchronously.
func (e *EditorImpl) Commit() bool {
	e.Lock()
	defer e.Unlock()
	e.pref.Lock()
	defer e.pref.Unlock()
	var success = true
	keys := e.commitToMemoryLocked()
	if len(keys) > 0 {
		success = e.pref.commitToDisk(e.pref.m)
		e.pref.notifyObservers(keys)
	}
	return success
}

func (e *EditorImpl) commitToMemoryLocked() []string {
	// if clear flag is set, re-create a new modified map with all the keys in origin preference map,
	// and set the values of all keys to nil, then put the origin modified map to this new map.
	if e.cleared {
		newModified := e.pref.copyOfMapLocked()
		for k, _ := range newModified {
			newModified[k] = nil
		}
		for k, v := range e.modified {
			newModified[k] = v
		}
		e.modified = newModified
	}
	changedKeys := make([]string, 0)
	for k, v := range e.modified {
		old, exist := e.pref.m[k]
		// A nil value in modified map indicates the Preferences shall be removed.
		if v == nil {
			if exist {
				delete(e.pref.m, k)
				changedKeys = append(changedKeys, k)
			}
		} else if !reflect.DeepEqual(old, v) {
			e.pref.m[k] = v
			changedKeys = append(changedKeys, k)
		}
	}
	return changedKeys
}

// notifyObservers send the changed keys to all registered observers.
func (p *PreferencesImpl) notifyObservers(keys []string) {
	p.observerLock.Lock()
	p.observerLock.Unlock()
	for _, key := range keys {
		for ob, _ := range p.observers {
			select {
			case ob <- key:
			default:
			}
		}
	}
}

func (p *PreferencesImpl) commitToDisk(changedMap map[string]interface{}) bool {
	path := basePath + p.name
	backupPath := path + "_bak"
	p.diskLock.Lock()
	defer p.diskLock.Unlock()
	// Backup the normal file
	if _, err := os.Stat(path); err == nil {
		if _, err2 := os.Stat(backupPath); err2 == nil {
			os.Remove(path)
		} else {
			os.Rename(path, backupPath)
		}
	}
	if file, err := os.Create(path); err == nil {
		defer file.Close()
		enc := gob.NewEncoder(file)
		if err := enc.Encode(&changedMap); err != nil {
			log.Printf("Error when write preference: %v", err)
			// remove normal file if error
			os.Remove(path)
			return false
		} else {
			// remove backup file if success
			os.Remove(backupPath)
			return true
		}
	}
	return false
}
