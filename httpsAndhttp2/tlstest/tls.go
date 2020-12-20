package tlstest

import (
	"path/filepath"
	"runtime"
)

var basePath string

func init() {
	// 读取当前文件路径
	_, currentFile, _, _ := runtime.Caller(0)
	// 获取当前文件的文件夹路径
	basePath = filepath.Dir(currentFile)
}

func Path(rel string) string {
	// 判定输入的是不是绝对路径
	if filepath.IsAbs(rel) {
		return rel
	}
	return filepath.Join(basePath, rel)
}
