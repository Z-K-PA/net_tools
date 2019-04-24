package util

import (
	"io/ioutil"
	"os"
)

func ReadFile2Buffer(fileName string) ([]byte, error) {
	inFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	buffers, err := ioutil.ReadAll(inFile)
	return buffers, err
}

type UnmarshalFunc func(buf []byte, v interface{}) error

func UnMarshalFile2Object(unmarshalFunc UnmarshalFunc, fileName string, val interface{}) error {
	buffers, err := ReadFile2Buffer(fileName)
	if err != nil {
		return err
	}
	err = unmarshalFunc(buffers, val)
	return err
}
