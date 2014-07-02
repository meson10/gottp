package utils

import "encoding/json"

type Q map[string]interface{}

type StreamReader interface {
	Close() error
	Read([]byte) (int, error)
}

func Encoder(data interface{}) []byte {
	b, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return b
}

func DecodeStream(s StreamReader, f interface{}) {
	decoder := json.NewDecoder(s)
	err := decoder.Decode(f)
	if err != nil {
		panic(err)
	}
}

func Decoder(b []byte, f interface{}) {
	err := json.Unmarshal(b, f)
	if err != nil {
		panic(err)
	}
}

func Convert(input interface{}, output interface{}) {
	Decoder(Encoder(input), output)
}
