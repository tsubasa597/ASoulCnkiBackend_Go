package cache

import (
	"bytes"
	"encoding/gob"
)

// Serialize 以 gob 序列化 value
func Serialize(value interface{}) ([]byte, error) {
	var b bytes.Buffer

	encoder := gob.NewEncoder(&b)
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// Deserialize 反序列化
func Deserialize(byt []byte, ptr interface{}) (err error) {
	b := bytes.NewBuffer(byt)

	decoder := gob.NewDecoder(b)
	if err = decoder.Decode(ptr); err != nil {
		return err
	}

	return nil
}
