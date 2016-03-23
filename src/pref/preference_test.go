package pref

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSharedPreference(t *testing.T) {
	assert.Empty(t, prefMap)
	GetSharedPreference("path1")
	assert.Equal(t, len(prefMap), 1)
	assert.Equal(t, prefMap["path1"].path, "path1")
	assert.Equal(t, len(prefMap["path1"].keyMap), 0)
}
