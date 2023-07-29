package connection

import (
    "path/filepath"
)

func isExecutable(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".exe"
}