package projectpath

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// RootProject repository root project
	RootProject = filepath.Join(filepath.Dir(b), "../..")
)
