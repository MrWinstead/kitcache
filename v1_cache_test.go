package kitcache_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"

	"github.com/mrwinstead/kitcache"
	"github.com/mrwinstead/kitcache/cache"
	"github.com/mrwinstead/kitcache/marshal"
)

type exampleRequest struct {
	UserId  uint64
	Payload []byte
}

func exampleRequestAllocator() interface{} {
	return &exampleRequest{}
}

type exampleRequestFunnel struct{}

func (f *exampleRequestFunnel) Hash(requestIface interface{}) ([]byte, error) {
	request := requestIface.(*exampleRequest)

	hashValue := fmt.Sprintf("%v", request.UserId)

	return []byte(hashValue), nil
}

// buildLoopbackEndpoint is used to use the same structure for request and
// responses
func buildLoopbackEndpoint(err error) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return request, err
	}
}

func buildTestEndpointCache(t *testing.T) kitcache.EndpointCache {
	l := buildLoopbackEndpoint(nil)
	funnel := &exampleRequestFunnel{}
	marshaler := marshal.NewGobMarshaler(exampleRequestAllocator)
	backingCache := cache.NewInMemoryCache(10, 10*time.Hour)
	created, err := kitcache.NewEndpointCache(l, backingCache, marshaler,
		funnel)
	assert.Nil(t, err)
	return created
}

func TestNewEndpointCacheCreated(t *testing.T) {
	l := buildLoopbackEndpoint(nil)
	funnel := &exampleRequestFunnel{}
	marshaler := marshal.NewGobMarshaler(exampleRequestAllocator)
	backingCache := cache.NewInMemoryCache(10, 10*time.Hour)

	created, err := kitcache.NewEndpointCache(l, backingCache, marshaler,
		funnel)

	assert.Nil(t, err)
	assert.NotNil(t, created)
}

func TestNewEndpointCache_Endpoint_NoError(t *testing.T) {
	var loopbackEndpointError error
	created := buildTestEndpointCache(t)

	cachedEndpoint := created.Endpoint()

	request := &exampleRequest{
		UserId:  rand.Uint64(),
		Payload: []byte(randomdata.RandStringRunes(10)),
	}

	result, endpointError := cachedEndpoint(context.Background(), request)
	assert.Equal(t, loopbackEndpointError, endpointError)
	assert.Equal(t, request, result)

	assert.Equal(t, uint64(1), created.Misses())
}

func TestNewEndpointCache_Endpoint_CacheHit(t *testing.T) {
	var loopbackEndpointError error
	created := buildTestEndpointCache(t)

	cachedEndpoint := created.Endpoint()

	request := &exampleRequest{
		UserId:  rand.Uint64(),
		Payload: []byte(randomdata.RandStringRunes(10)),
	}

	result, endpointError := cachedEndpoint(context.Background(), request)
	assert.Equal(t, loopbackEndpointError, endpointError)
	assert.Equal(t, request, result)

	result2, endpointError2 := cachedEndpoint(context.Background(), request)
	assert.Equal(t, loopbackEndpointError, endpointError2)
	assert.Equal(t, request, result2)

	assert.Equal(t, uint64(1), created.Misses())
	assert.Equal(t, uint64(1), created.Hits())
}

func TestNewEndpointCache_ResetStatistics(t *testing.T) {
	var loopbackEndpointError error
	created := buildTestEndpointCache(t)

	cachedEndpoint := created.Endpoint()

	request := &exampleRequest{
		UserId:  rand.Uint64(),
		Payload: []byte(randomdata.RandStringRunes(10)),
	}

	do2Requests := func() {
		result, endpointError := cachedEndpoint(context.Background(), request)
		assert.Equal(t, loopbackEndpointError, endpointError)
		assert.Equal(t, request, result)

		result2, endpointError2 := cachedEndpoint(context.Background(), request)
		assert.Equal(t, loopbackEndpointError, endpointError2)
		assert.Equal(t, request, result2)
	}

	do2Requests()
	assert.Equal(t, uint64(1), created.Misses())
	assert.Equal(t, uint64(1), created.Hits())

	created.ResetStatistics()
	assert.Equal(t, uint64(0), created.Misses())
	assert.Equal(t, uint64(0), created.Hits())

	// the cache will be warm, so there'll be only hits from here on out

	do2Requests()
	assert.Equal(t, uint64(0), created.Misses())
	assert.Equal(t, uint64(2), created.Hits())

	do2Requests()
	assert.Equal(t, uint64(0), created.Misses())
	assert.Equal(t, uint64(4), created.Hits())
}

func TestNewEndpointCache_Invalidator(t *testing.T) {
	var loopbackEndpointError error
	l := buildLoopbackEndpoint(nil)
	funnel := &exampleRequestFunnel{}
	marshaler := marshal.NewGobMarshaler(exampleRequestAllocator)
	backingCache := cache.NewInMemoryCache(10, 10*time.Hour)

	created, err := kitcache.NewEndpointCache(l, backingCache, marshaler,
		funnel)

	assert.Nil(t, err)
	assert.NotNil(t, created)

	cachedEndpoint := created.Endpoint()

	request := &exampleRequest{
		UserId:  rand.Uint64(),
		Payload: []byte(randomdata.RandStringRunes(10)),
	}

	result, endpointError := cachedEndpoint(context.Background(), request)
	assert.Equal(t, loopbackEndpointError, endpointError)
	assert.Equal(t, request, result)

	result2, endpointError2 := cachedEndpoint(context.Background(), request)
	assert.Equal(t, loopbackEndpointError, endpointError2)
	assert.Equal(t, request, result2)

	assert.Equal(t, uint64(1), created.Misses())
	assert.Equal(t, uint64(1), created.Hits())

	requestKey, funnelErr := funnel.Hash(request)
	assert.Nil(t, funnelErr)

	i := created.Invalidator()
	invalidationErr := i.Invalidate(requestKey)
	assert.Nil(t, invalidationErr)

	result3, endpointError3 := cachedEndpoint(context.Background(), request)
	assert.Equal(t, loopbackEndpointError, endpointError3)
	assert.Equal(t, request, result3)

	assert.Equal(t, uint64(2), created.Misses())
	assert.Equal(t, uint64(1), created.Hits())
}
