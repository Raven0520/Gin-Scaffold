package util

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// GetAbsolutePath 获取绝对路径
func GetAbsolutePath() (p string, err error) {
	p, err = GetAbsolutePathByExecutable()
	if err != nil {
		return
	}
	if strings.Contains(p, GetSystemTempDir()) {
		return GetAbsolutePathByCaller(), nil
	}
	return
}

// GetSystemTempDir 获取系统临时目录
func GetSystemTempDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

// GetAbsolutePathByCaller go Run 获取执行路径
func GetAbsolutePathByCaller() (p string) {
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		p = path.Dir(filename)
	}
	return
}

// GetAbsolutePathByExecutable go Build 获取执行路径
func GetAbsolutePathByExecutable() (p string, err error) {
	ep, err := os.Executable()
	if err != nil {
		return
	}
	p, err = filepath.EvalSymlinks(filepath.Dir(ep))
	return
}
