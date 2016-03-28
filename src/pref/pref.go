package pref

type OnPreferenceChangeListener chan string

type Preferences interface {
	Contains(string) bool
	GetBool(string, bool) bool
	GetInt(string, int) int
	GetInt32(string, int32) int32
	GetInt64(string, int64) int64
	GetUInt32(string, uint32) uint32
	GetUInt64(string, uint64) uint64
	GetFloat32(string, float32) float32
	GetFloat64(string, float64) float64
	GetByte(string, byte) byte
	GetRune(string, rune) rune
	GetString(string, string) string
	GetObject(string, interface{}) interface{}
	RegisterOnPreferenceChangeListener(OnPreferenceChangeListener)
	UnregisterOnPreferenceChangeListener(OnPreferenceChangeListener)

	Edit() Editor
}

type Editor interface {
	Apply()
	Commit() bool
	Clear() Editor
	Remove(string) Editor
	Put(string, interface{}) Editor
}
