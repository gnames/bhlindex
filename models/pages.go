package models

import "path/filepath"

type Page struct {
	ID      string
	TitleID int
}

func IsPageFile(f string) bool {
	res, _ := filepath.Match("*_[0-9][0-9][0-9][0-9].txt", f)
	return res
}

func PageID(f string) string {
	extLen := len(filepath.Ext(f))
	idLen := len(f) - extLen
	return f[0:idLen]
}
