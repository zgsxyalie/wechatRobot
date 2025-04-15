package helper

import "os"

func Exists(p string) bool {
	if _, err := os.Stat(p); err != nil {
		return !os.IsNotExist(err)
	}
	return true

}
