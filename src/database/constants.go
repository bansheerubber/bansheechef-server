package database

import (
	"os"
	"path/filepath"
)

var HOME = os.Getenv("HOME")
var LOCAL = filepath.Join(HOME, ".config", "bansheechef")
var LOCAL_STORAGE = filepath.Join(LOCAL, "storage")
var LOCAL_IMAGES = filepath.Join(LOCAL_STORAGE, "images")
