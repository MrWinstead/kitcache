package marshal_test

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"

	"github.com/mrwinstead/kitcache/marshal"
)

type exampleStruct struct {
	Data string
}

func exampleStructAllocator() interface{} {
	return &exampleStruct{}
}

func TestNewGobMarshaler(t *testing.T) {
	created := marshal.NewGobMarshaler(exampleStructAllocator)
	assert.NotNil(t, created)
}

func TestGobMarshaler_Marshal(t *testing.T) {
	created := marshal.NewGobMarshaler(exampleStructAllocator)
	assert.NotNil(t, created)

	source := &exampleStruct{randomdata.RandStringRunes(10)}

	serialized, marshalEerr := created.Marshal(source)
	assert.Nil(t, marshalEerr)
	assert.NotEmpty(t, serialized)

	decoder := gob.NewDecoder(bytes.NewBuffer(serialized))
	destination := &exampleStruct{}
	decodeErr := decoder.Decode(destination)
	assert.Nil(t, decodeErr)

	assert.Equal(t, source.Data, destination.Data)
}

func TestGobMarshaler_Unmarshal(t *testing.T) {
	created := marshal.NewGobMarshaler(exampleStructAllocator)
	assert.NotNil(t, created)

	original := &exampleStruct{randomdata.RandStringRunes(10)}
	serializedOriginal := bytes.Buffer{}

	serializationErr := gob.NewEncoder(&serializedOriginal).Encode(original)
	assert.Nil(t, serializationErr)

	unmarshaledIface, unmarshalErr := created.Unmarshal(
		serializedOriginal.Bytes())
	assert.Nil(t, unmarshalErr)

	unmarshaled := unmarshaledIface.(*exampleStruct)
	assert.Equal(t, original.Data, unmarshaled.Data)
}
