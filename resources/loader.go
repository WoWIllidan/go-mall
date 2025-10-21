package resources

import (
	"bytes"
	"embed"
	"io"
)

//go:embed *
var f embed.FS

// LoadResourceFile 把资源文件 filePath 加载到内存, 以 io.Reader 形式返回
func LoadResourceFile(filePath string) (io.Reader, error) {
	_bytes, err := f.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(_bytes), nil
}
