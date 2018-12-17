package marshal

import (
	"bytes"
	"encoding/gob"

	"github.com/mrwinstead/kitcache"
)

// GobMarshaler implements the kitcache.ResponseMarshaler interface using the
// build-int golang encoding/gob library
type GobMarshaler struct {
	allocator UnmarshalObjectAllocator
}

// UnmarshalObjectAllocator will be called each time Unmarshal is called in
// order to get an object which this marshaler will produce with the Unmarshal
// call. It should return a pointer to the struct to be returned so as to avoid
// unexpected behavior surrounding struct copying.
//
// This function is needed in order to have a struct allocated which the
// underlying serialization library (encoding/gob) will populate.
type UnmarshalObjectAllocator func() interface{}

// NewGobMarshaler constructs a GobMarshaler
func NewGobMarshaler(allocator UnmarshalObjectAllocator,
) kitcache.ResponseMarshaler {
	created := &GobMarshaler{
		allocator: allocator,
	}
	return created
}

// Marshal encodes the provided value using a gob encoding
func (g *GobMarshaler) Marshal(unserialized interface{}) ([]byte, error) {
	output := &bytes.Buffer{}
	enc := gob.NewEncoder(output)

	err := enc.Encode(unserialized)
	if nil != err {
		return nil, err
	}

	return output.Bytes(), nil
}

// Unmarshal allocates a value by calling the UnmarshalObjectAllocator and
// decodes the provided buffer into the blank value
func (g *GobMarshaler) Unmarshal(source []byte) (interface{}, error) {
	sourceBuffer := bytes.NewBuffer(source)
	dec := gob.NewDecoder(sourceBuffer)

	destination := g.allocator()
	decodeErr := dec.Decode(destination)
	if nil != decodeErr {
		return nil, decodeErr
	}

	return destination, nil
}
