package crypt

import (
	"bytes"
	"encoding/gob"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

func Decode(c []byte) (credValues credentials.Value, err error) {
	buf := bytes.NewBuffer(c)
	err = gob.NewDecoder(buf).Decode(&credValues)
	return
}

func Encode(credValues credentials.Value) (encodedCredValues []byte, err error) {
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(credValues)
	encodedCredValues = buf.Bytes()
	return
}
