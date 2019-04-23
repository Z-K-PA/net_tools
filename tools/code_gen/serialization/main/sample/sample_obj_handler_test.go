package sample

import (
	"bytes"
	"testing"
)

func TestNewSampleHandler(t *testing.T) {
	buf := make([]byte, 1024)
	//序列化
	var sample Sample2

	flag := float64(-0.0003)

	sample.Id = 1002010
	sample.Sample1List = make([]Sample1, 2)
	sample.Sample1List[0].Field1 = []byte{0, 'a', 'z', 0}
	sample.Sample1List[0].Field2 = "abc"
	sample.Sample1List[0].Field3 = float64(flag)

	writer := NewWriteSampleHandler(buf)
	err := writer.WriteSample2(sample)
	if err != nil {
		t.Errorf("marshal error:%+v", err)
		return
	}

	newData := writer.Data()
	t.Logf("new buf len:%+v cap:%+v", len(newData), cap(newData))

	reader, err := NewReadSampleHandler(newData)
	if err != nil {
		t.Errorf("new reader error:%+v", err)
		return
	} else {
		newSample, err := reader.ReadSample2()
		if err != nil {
			t.Errorf("unmarshal error:%+v", err)
		}
		if newSample.Id != 1002010 {
			t.Fail()
		}
		if len(newSample.Sample1List) != 2 {
			t.Fail()
		}
		if !bytes.Equal(newSample.Sample1List[0].Field1, []byte{0, 'a', 'z', 0}) {
			t.Fail()
		}
		if newSample.Sample1List[0].Field2 != "abc" {
			t.Fail()
		}
		if newSample.Sample1List[0].Field3 != flag {
			t.Fail()
		}
		if newSample.Sample1List[1].Field1 == nil {
			t.Fail()
		}
		if len(newSample.Sample1List[1].Field1) != 0 {
			t.Fail()
		}
		if newSample.Sample1List[1].Field2 != "" {
			t.Fail()
		}
		if newSample.Sample1List[1].Field3 != float64(0) {
			t.Fail()
		}
	}

}
