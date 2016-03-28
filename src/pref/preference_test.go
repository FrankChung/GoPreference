package pref

import (
	"github.com/stretchr/testify/suite"
	"os"
	"sync"
	"testing"
)

const PrefName = "pref_unit"

var pref *PreferencesImpl
var editor *EditorImpl

type TestSuite struct {
	suite.Suite
}

func (suite *TestSuite) SetupTest() {
	basePath = "./"
	prefMap = make(map[string]*PreferencesImpl)
	pref = &PreferencesImpl{
		m:            make(map[string]interface{}),
		name:         PrefName,
		observers:    make(map[chan string]interface{}),
		writeCh:      make(chan map[string]interface{}),
		diskLock:     &sync.Mutex{},
		observerLock: &sync.Mutex{},
		loadWg:       &sync.WaitGroup{},
		Mutex:        &sync.Mutex{},
	}
	editor = &EditorImpl{
		modified: make(map[string]interface{}),
		pref:     pref,
		cleared:  false,
		Mutex:    &sync.Mutex{},
	}
}

func (suite *TestSuite) TearDownTest() {
	os.Remove(basePath + PrefName)
	os.Remove(basePath + PrefName + "_bak")
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestInitBasePath() {
	suite.Equal(basePath, "./")
	InitBasePath("test/path")
	suite.Equal(basePath, "test/path")
}

func (suite *TestSuite) TestGetSharedPreferenceEmpty() {
	suite.Empty(prefMap)
	NewPreferences("name1")
	suite.Len(prefMap, 1)
	suite.Equal(prefMap["name1"].name, "name1")
	suite.Empty(prefMap["name1"].m)
	suite.Len(prefMap["name1"].observers, 0)
}

func (suite *TestSuite) TestGetSharedPreferenceExist() {
	suite.Empty(prefMap)
	prefMap[PrefName] = pref
	NewPreferences(PrefName)
	suite.Len(prefMap, 1)
	suite.Equal(prefMap[PrefName].name, PrefName)
	suite.Empty(prefMap[PrefName].m)
	suite.Len(prefMap[PrefName].observers, 0)
}

func (suite *TestSuite) TestRegisterObserver() {
	o := make(chan string)
	suite.Len(pref.observers, 0)
	pref.RegisterOnPreferenceChangeListener(o)
	suite.Len(pref.observers, 1)
	pref.RegisterOnPreferenceChangeListener(o)
	suite.Len(pref.observers, 1)
	pref.UnregisterOnPreferenceChangeListener(o)
	suite.Len(pref.observers, 0)
	pref.UnregisterOnPreferenceChangeListener(o)
	suite.Len(pref.observers, 0)
}

func (suite *TestSuite) TestGetNotFound() {
	var boolean = true
	var i int = 4
	var i32 int32 = 4
	var i64 int64 = 4
	var u32 uint32 = 4
	var u64 uint64 = 4
	var f32 float32 = 4.5
	var f64 float64 = 4.5
	var b byte = 0x61
	var r rune = 'a'
	var str string = "a"
	var obj = struct {
		field string
	}{field: "str"}
	suite.Equal(pref.GetBool("key1", boolean), boolean)
	suite.Equal(pref.GetInt("key2", i), i)
	suite.Equal(pref.GetInt32("key3", i32), i32)
	suite.Equal(pref.GetInt64("key4", i64), i64)
	suite.Equal(pref.GetUInt32("key5", u32), u32)
	suite.Equal(pref.GetUInt64("key6", u64), u64)
	suite.Equal(pref.GetFloat32("key7", f32), f32)
	suite.Equal(pref.GetFloat64("key8", f64), f64)
	suite.Equal(pref.GetByte("key9", b), b)
	suite.Equal(pref.GetRune("key10", r), r)
	suite.Equal(pref.GetString("key11", str), str)
	suite.Equal(pref.GetObject("key12", obj), obj)
}

func (suite *TestSuite) TestGet() {
	var boolean = true
	var i int = 4
	var i32 int32 = 4
	var i64 int64 = 4
	var u32 uint32 = 4
	var u64 uint64 = 4
	var f32 float32 = 4.5
	var f64 float64 = 4.5
	var b byte = 0x61
	var r rune = 'a'
	var str string = "a"
	var obj = struct {
		field string
	}{field: "str"}
	pref.m["key1"] = boolean
	pref.m["key2"] = i
	pref.m["key3"] = i32
	pref.m["key4"] = i64
	pref.m["key5"] = u32
	pref.m["key6"] = u64
	pref.m["key7"] = f32
	pref.m["key8"] = f64
	pref.m["key9"] = b
	pref.m["key10"] = r
	pref.m["key11"] = str
	pref.m["key12"] = obj
	suite.Equal(pref.GetBool("key1", false), boolean)
	suite.Equal(pref.GetInt("key2", 0), i)
	suite.Equal(pref.GetInt32("key3", 0), i32)
	suite.Equal(pref.GetInt64("key4", 0), i64)
	suite.Equal(pref.GetUInt32("key5", 0), u32)
	suite.Equal(pref.GetUInt64("key6", 0), u64)
	suite.Equal(pref.GetFloat32("key7", 0), f32)
	suite.Equal(pref.GetFloat64("key8", 0), f64)
	suite.Equal(pref.GetByte("key9", 0), b)
	suite.Equal(pref.GetRune("key10", ' '), r)
	suite.Equal(pref.GetString("key11", ""), str)
	suite.Equal(pref.GetObject("key12", nil), obj)
}

func (suite *TestSuite) TestContains() {
	pref.m["key"] = 3
	suite.False(pref.Contains("other"))
	suite.True(pref.Contains("key"))
}

func (suite *TestSuite) TestEdit() {
	editor := pref.Edit().(*EditorImpl)
	suite.Empty(editor.modified)
	suite.Equal(editor.pref, pref)
}

func (suite *TestSuite) TestPutInt() {
	suite.Empty(editor.modified)
	editor.Put("key", 3)
	suite.Len(editor.modified, 1)
	suite.Equal(editor.modified["key"], 3)
}

func (suite *TestSuite) TestPutFloat() {
	suite.Empty(editor.modified)
	editor.Put("key", 3.5)
	suite.Len(editor.modified, 1)
	suite.Equal(editor.modified["key"], 3.5)
}

func (suite *TestSuite) TestPutBool() {
	suite.Empty(editor.modified)
	editor.Put("key", true)
	suite.Len(editor.modified, 1)
	suite.Equal(editor.modified["key"], true)
}

func (suite *TestSuite) TestPutString() {
	suite.Empty(editor.modified)
	editor.Put("key", "hello")
	suite.Len(editor.modified, 1)
	suite.Equal(editor.modified["key"], "hello")
}

func (suite *TestSuite) TestPutObject() {
	obj := struct {
		field string
	}{field: "str"}
	suite.Empty(editor.modified)
	editor.Put("key", obj)
	suite.Len(editor.modified, 1)
	suite.Equal(editor.modified["key"], obj)
}

func (suite *TestSuite) TestRemove() {
	suite.Empty(editor.modified)
	editor.Remove("key")
	suite.Len(editor.modified, 1)
	suite.Nil(editor.modified["key"])
}

func (suite *TestSuite) TestClear() {
	suite.False(editor.cleared)
	editor.Clear()
	suite.True(editor.cleared)
}

func (suite *TestSuite) TestApply() {
	suite.Empty(pref.m)
	editor.Put("key", 5).Apply()
	suite.Equal(pref.m["key"], 5)
	editor.Put("key", 15).Apply()
	suite.Equal(pref.m["key"], 15)
	editor.Remove("key").Apply()
	suite.Empty(pref.m)
}

func (suite *TestSuite) TestCommit() {
	suite.Empty(pref.m)
	editor.Put("key", 5).Commit()
	suite.Equal(pref.m["key"], 5)
	editor.Put("key", 15).Commit()
	suite.Equal(pref.m["key"], 15)
	editor.Remove("key").Commit()
	suite.Empty(pref.m)
}

func (suite *TestSuite) TestObserver() {
	ch := make(chan string, 4)
	pref.RegisterOnPreferenceChangeListener(ch)
	editor.Put("key1", 1).Put("key2", true).Put("key3", 2.5).Put("key4", "value").Commit()
	s := make(map[string]interface{})
	s[<-ch] = nil
	s[<-ch] = nil
	s[<-ch] = nil
	s[<-ch] = nil
	suite.Contains(s, "key1")
	suite.Contains(s, "key2")
	suite.Contains(s, "key3")
	suite.Contains(s, "key4")
	pref.UnregisterOnPreferenceChangeListener(ch)
	close(ch)
}

func (suite *TestSuite) TestObserverWithNoChanges() {
	ch := make(chan string, 4)
	editor.Put("key1", 1).Put("key2", true).Put("key3", 2.5).Put("key4", "value").Commit()
	pref.RegisterOnPreferenceChangeListener(ch)
	editor.Put("key1", 1).Put("key2", true).Put("key3", 0).Put("key4", "value").Commit()
	suite.Equal("key3", <-ch)
	pref.UnregisterOnPreferenceChangeListener(ch)
	close(ch)
}

func (suite *TestSuite) TestApplyConcurrency() {
	ch := make(chan bool, 2)
	go func() {
		for i := 0; i < 100; i++ {
			pref.Edit().Put("key1", i).Apply()
		}
		ch <- true
	}()
	go func() {
		for i := 100; i < 200; i++ {
			pref.Edit().Put("key2", i).Apply()
		}
		ch <- true
	}()
	<-ch
	<-ch
	suite.Equal(pref.GetInt("key1", 0), 99)
	suite.Equal(pref.GetInt("key2", 0), 199)
}

func (suite *TestSuite) TestCommitConcurrency() {
	ch := make(chan bool, 2)
	go func() {
		for i := 0; i < 100; i++ {
			pref.Edit().Put("key1", i).Commit()
		}
		ch <- true
	}()
	go func() {
		for i := 100; i < 200; i++ {
			pref.Edit().Put("key2", i).Commit()
		}
		ch <- true
	}()
	<-ch
	<-ch
	suite.Equal(pref.GetInt("key1", 0), 99)
	suite.Equal(pref.GetInt("key2", 0), 199)
}

func (suite *TestSuite) TestApplyConcurrency2() {
	ch := make(chan bool, 2)
	go func() {
		for i := 0; i < 100; i++ {
			pref.Edit().Put("key1", i).Put("key3", i).Apply()
		}
		ch <- true
	}()
	go func() {
		for i := 100; i < 200; i++ {
			pref.Edit().Put("key2", i).Put("key3", i).Apply()
		}
		ch <- true
	}()
	<-ch
	<-ch
	suite.Equal(pref.GetInt("key1", 0), 99)
	suite.Equal(pref.GetInt("key2", 0), 199)
	val3 := pref.GetInt("key3", 0)
	suite.True(val3 == 99 || val3 == 199)
}

func (suite *TestSuite) TestCommitConcurrency2() {
	ch := make(chan bool, 2)
	go func() {
		for i := 0; i < 100; i++ {
			pref.Edit().Put("key1", i).Put("key3", i).Commit()
		}
		ch <- true
	}()
	go func() {
		for i := 100; i < 200; i++ {
			pref.Edit().Put("key2", i).Put("key3", i).Commit()
		}
		ch <- true
	}()
	<-ch
	<-ch
	suite.Equal(pref.GetInt("key1", 0), 99)
	suite.Equal(pref.GetInt("key2", 0), 199)
	val3 := pref.GetInt("key3", 0)
	suite.True(val3 == 99 || val3 == 199)
}

func (suite *TestSuite) TestReadWritePrefFile() {
	_, err := os.Open(basePath + PrefName)
	suite.True(os.IsNotExist(err))
	editor.Put("key", "value").Commit()
	_, err = os.Open(basePath + PrefName)
	suite.Nil(err)
	_, err = os.Open(basePath + PrefName + "_bak")
	suite.True(os.IsNotExist(err))

	pref.loadWg.Add(1)
	go pref.loadFromFile()
	pref.loadWg.Wait()
	suite.Len(pref.m, 1)
	suite.Equal(pref.m["key"], "value")
}

func (suite *TestSuite) TestReadWriteBackupFile() {
	type Stranger struct {
	}
	editor.Put("key", "value").Commit()
	_, err := os.Open(basePath + PrefName)
	suite.Nil(err)
	editor.Put("key", Stranger{}).Commit()
	_, err = os.Open(basePath + PrefName)
	suite.True(os.IsNotExist(err))
	_, err = os.Open(basePath + PrefName + "_bak")
	suite.Nil(err)

	pref.loadWg.Add(1)
	go pref.loadFromFile()
	pref.loadWg.Wait()
	suite.Len(pref.m, 1)
	suite.Equal(pref.m["key"], "value")
}
