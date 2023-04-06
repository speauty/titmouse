// Package single_file
/**
 * 独立文件数据存储
 * @todo 未针对文件名称进行存在性检测
 */
package single_file

import (
	"bytes"
	"encoding/gob"
	"os"
)

type SingleFileInterface interface {
	GetFilename() string
	Store() error
	Load() error
}

func Store(singleFile SingleFileInterface) error {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(singleFile); err != nil {
		return err
	}
	return os.WriteFile(singleFile.GetFilename(), buffer.Bytes(), 0666)
}

func Load(singleFile SingleFileInterface) error {
	raw, err := os.ReadFile(singleFile.GetFilename())
	if err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(raw)).Decode(singleFile)
}
