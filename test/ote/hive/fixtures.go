package hive

import (
	"embed"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

//go:embed testdata/*
var testdataFS embed.FS

var (
	extractedDir string
	extractOnce  sync.Once
)

// FixturePath returns the absolute path to the testdata directory or a
// subdirectory within it. When running from the source tree the original
// files are referenced directly; when the source tree is absent (e.g. inside
// a container) the embedded testdata files are extracted to a temporary
// directory on disk first.
func FixturePath(elem ...string) string {
	_, file, _, ok := runtime.Caller(0)
	if ok {
		srcDir := filepath.Dir(file)
		candidate := filepath.Join(append([]string{srcDir}, elem...)...)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	extractOnce.Do(func() {
		tmp, err := os.MkdirTemp("", "hive-ote-testdata-")
		if err != nil {
			panic("failed to create temp dir for testdata: " + err.Error())
		}
		tdDir := filepath.Join(tmp, "testdata")
		if err := os.MkdirAll(tdDir, 0755); err != nil {
			panic("failed to create testdata subdir: " + err.Error())
		}
		entries, err := testdataFS.ReadDir("testdata")
		if err != nil {
			panic("failed to read embedded testdata: " + err.Error())
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			data, err := testdataFS.ReadFile("testdata/" + entry.Name())
			if err != nil {
				panic("failed to read embedded file " + entry.Name() + ": " + err.Error())
			}
			if err := os.WriteFile(filepath.Join(tdDir, entry.Name()), data, 0644); err != nil {
				panic("failed to write extracted file " + entry.Name() + ": " + err.Error())
			}
		}
		extractedDir = tmp
	})

	return filepath.Join(append([]string{extractedDir}, elem...)...)
}

// Asset reads an embedded testdata file and returns its contents.
// This replaces testdata.Asset from openshift-tests-private.
func Asset(name string) ([]byte, error) {
	return testdataFS.ReadFile(name)
}
