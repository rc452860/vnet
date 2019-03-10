package binaryx

import (
	"bytes"
	"encoding/binary"
)

func BEBytesToInt(data []byte) (ret int) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &ret)
	return
}

func BEBytesToInt32(data []byte) (ret int32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &ret)
	return
}

func BEBytesToUint32(data []byte) (ret uint32) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &ret)
	return
}

func LEBytesToUint64(data []byte) uint64 {
	if len(data) < 8 {
		return 0
	}
	return binary.LittleEndian.Uint64(data)
}

func LEUInt16ToBytes(data uint16) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, data)
	return buf
}
