package single_file

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type testTmpData struct {
	Name string
}

func (customTTD *testTmpData) GetFilename() string {
	return "test_tmp_data.tmp"
}

func (customTTD *testTmpData) Store() error {
	return Store(customTTD)
}

func (customTTD *testTmpData) Load() error {
	return Load(customTTD)
}

func TestStore(t *testing.T) {
	currentAssert := assert.New(t)

	tmpData := new(testTmpData)
	defer func() {
		_ = os.Remove(tmpData.GetFilename())
	}()
	tmpData.Name = "test"

	currentAssert.Nil(tmpData.Store())
}

func TestLoad(t *testing.T) {
	currentAssert := assert.New(t)

	tmpData := new(testTmpData)
	defer func() {
		_ = os.Remove(tmpData.GetFilename())
	}()
	tmpData.Name = "test"

	currentAssert.Nil(tmpData.Store())

	otherTmpData := new(testTmpData)

	currentAssert.Nil(otherTmpData.Load())
	currentAssert.Equal(tmpData.Name, otherTmpData.Name)
}
