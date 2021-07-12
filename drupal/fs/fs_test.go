package fs

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)
import "github.com/stretchr/testify/assert"

func Test_FindExpectedJsonCurrentDir(t *testing.T) {
	path := FindExpectedJson(t, "fs.go")
	assert.NotEmpty(t, path)
	assert.Equal(t, path, "fs.go")
}

func Test_FindExpectedJsonSubDir(t *testing.T) {
	d, err := os.MkdirTemp(".", "")
	assert.Nil(t, err, "error creating temporary directory for test")
	f, err := os.CreateTemp(d, "")
	assert.Nil(t, err, "error creating temporary file for test")
	info, err := f.Stat()
	assert.Nil(t, err, "error performing stat on temporary file")

	path := FindExpectedJson(t, info.Name())

	assert.Equal(t, strings.TrimPrefix(f.Name(), "./"), path)
	assert.Nil(t, os.RemoveAll(d))
}

func Test_FindExpectedJsonSubSubDir(t *testing.T) {
	d1, err := os.MkdirTemp(".", "")
	assert.Nil(t, err, "error creating temporary directory for test")
	d2, err := os.MkdirTemp(d1, "")
	assert.Nil(t, err, "error creating temporary directory for test")
	f, err := os.CreateTemp(d2, "")
	assert.Nil(t, err, "error creating temporary file for test")
	info, err := f.Stat()
	assert.Nil(t, err, "error performing stat on temporary file")

	path := FindExpectedJson(t, info.Name())
	log.Printf("Found file %s: %s", info.Name(), path)
	assert.Equal(t, strings.TrimPrefix(f.Name(), "./"), path)
	assert.Nil(t, os.RemoveAll(d1))
}

func Test_FindExpectedJsonSubSubDirWithBasedir(t *testing.T) {
	d1, err := os.MkdirTemp(".", "")
	assert.Nil(t, err, "error creating temporary directory for test")
	d2, err := os.MkdirTemp(d1, "")
	assert.Nil(t, err, "error creating temporary directory for test")
	d3, err := os.MkdirTemp(d2, "")
	assert.Nil(t, err, "error creating temporary directory for test")

	_, err = os.MkdirTemp(d1, "")
	assert.Nil(t, err, "error creating temporary directory for test")
	_, err = os.MkdirTemp(d2, "")
	assert.Nil(t, err, "error creating temporary directory for test")
	_, err = os.MkdirTemp(d3, "")

	f, err := os.CreateTemp(d3, "")
	assert.Nil(t, err, "error creating temporary file for test")
	fileInfo, err := f.Stat()
	assert.Nil(t, err, "error performing stat on temporary file")
	dirInfo, err := os.Stat(d3)
	assert.Nil(t, err, "error performing stat on temporary file")

	path := FindExpectedJson(t, fileInfo.Name(), dirInfo.Name())
	log.Printf("Found file %s: %s", fileInfo.Name(), path)
	assert.Equal(t, strings.TrimPrefix(f.Name(), "./"), path)
	assert.Nil(t, os.RemoveAll(d1))
}

func Test_FindExpectedJsonPathElement(t *testing.T) {
	assert.Panics(t, func() {
		FindExpectedJson(t, filepath.Join("foo", "bar"))
	})
}

func Test_FindExpectedJsonDirPathElement(t *testing.T) {
	assert.Panics(t, func() {
		FindExpectedJson(t, "moo", filepath.Join("foo", "bar"))
	})
}
