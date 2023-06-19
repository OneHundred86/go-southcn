package cache

import "encoding/json"

type Encoder interface {
	Encode(v interface{}) (string, error)
	Decode(str string, v interface{}) error
}

type JsonEncoder struct {
	Encoder
}

func (*JsonEncoder) Encode(v interface{}) (result string, err error) {
	bytes, err := json.Marshal(v)
	result = string(bytes)
	return
}

func (*JsonEncoder) Decode(str string, v interface{}) error {
	return json.Unmarshal([]byte(str), v)
}
