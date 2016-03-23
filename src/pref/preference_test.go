package pref

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPreference(t *testing.T) {
	assert.Equal(t, len(prefMap), 0)
	GetPreference("path1")
	assert.Equal(t, len(prefMap), 1)
	assert.Equal(t, prefMap["path1"].path, "path1")
}
