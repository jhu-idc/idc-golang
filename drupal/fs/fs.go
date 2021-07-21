// Deprecated package of functions used to discover test resources
//
// Instead of using this package, consider using go:embed instead
package fs

import (
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

// Searches the file system for the named file, and answers the first path matching the file name.
//
// The `name` or optional `searchdirs` should not contain any path separators.  If `searchdirs` is supplied, the
// returned file will have at least one of the `searchdirs` as an ancestor.
//
// This function allows for an IDE to discover test resources while allowing for IDC test framework (the one invoked by
// `make test`) to discover those same resources without hard coding paths.
func FindExpectedJson(t *testing.T, name string, searchdirs ...string) string {
	// the resolved json file, including its path relative to the working directory.
	var expectedJsonFile string

	if strings.Contains(name, string(os.PathSeparator)) {
		logger.Panic().Msgf("Supplied file name '%s' must not contain path separator '%s'", name, string(os.PathSeparator))
	}

	if searchdirs != nil && len(searchdirs) > 0 {
		for _, dir := range searchdirs {
			if strings.Contains(dir, string(os.PathSeparator)) {
				logger.Panic().Msgf("Supplied search directory '%s' must not contain path separator '%s'", dir, string(os.PathSeparator))
			}
		}
	}

	var basedirs = []string{}

	if searchdirs != nil && len(searchdirs) > 0 {
		filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				if pathContains(path, searchdirs) {
					basedirs = append(basedirs, path)
					return filepath.SkipDir
				} else {
					logger.Printf("skipping dir: %s", info.Name())
				}
			}
			return nil
		})
	} else {
		basedirs = append(basedirs, ".")
	}

	for _, basedir := range basedirs {
		filepath.Walk(basedir, func(path string, info os.FileInfo, err error) error {
			require.Nil(t, err, "Unexpected error when searching for '%s': %s", name, err)

			if info.IsDir() {
				logger.Printf("searching dir: %s", info.Name())
			}

			// Resolve the file
			if info.Name() == name {
				expectedJsonFile = path
				return errors.New(fmt.Sprintf("Found file %s", expectedJsonFile))
			}
			return nil
		})
	}

	if expectedJsonFile == "" {
		logger.Panic().Msgf("Could not locate file '%v'", name)
	}
	return expectedJsonFile
}

func pathContains(path string, candidates []string) bool {
	for _, pathelement := range strings.Split(path, string(os.PathSeparator)) {
		for _, candidate := range candidates {
			if pathelement == candidate {
				return true
			}
		}
	}

	return false
}
