package util

import (
	"bytes"
	"encoding/binary"
)

var byteOrder = binary.BigEndian

func UInt64ToBytes(num uint64) (ret []byte, err error) {
	var buffer bytes.Buffer
	err = binary.Write(&buffer, byteOrder, num)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func BytesToUInt64(in []byte) (ret uint64, err error) {
	return byteOrder.Uint64(in), nil
}
