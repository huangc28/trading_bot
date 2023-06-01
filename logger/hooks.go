package logger

import "os"

type FileHook struct {
	file *os.File
}
