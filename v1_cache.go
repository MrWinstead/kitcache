package kitcache

import (
	"context"
	"sync/atomic"

	"github.com/go-kit/kit/endpoint"
)

var (
	defaultV1CacheOptions = []V1ConstructorOption{
		V1WithMaximumItems(DefaultMaxCachedUnits),
	}
)

type v1Cache struct {
	cache           Cache
	funnel          RequestFunnel
	marshaler       ResponseMarshaler
	maximumItems    uint64
	wrappedEndpoint endpoint.Endpoint

	hits   uint64
	misses uint64
}

type v1CacheInvalidator struct {
	cache Cache
}

// V1WithMaximumItems sets the maximum items which the cache will hold
func V1WithMaximumItems(maximumItems uint64) V1ConstructorOption {
	return func(i *v1Cache) error {
		i.maximumItems = maximumItems
		return nil
	}
}

// NewEndpointCache wraps an endpoint with a cache
func NewEndpointCache(wrapped endpoint.Endpoint, cache Cache,
	marshaler ResponseMarshaler, funnel RequestFunnel,
	opts ...V1ConstructorOption) (
	EndpointCache, error) {
	created := &v1Cache{
		cache:           cache,
		funnel:          funnel,
		marshaler:       marshaler,
		wrappedEndpoint: wrapped,
	}

	allOpts := append(defaultV1CacheOptions, opts...)
	for optIdx := range allOpts {
		optErr := allOpts[optIdx](created)
		if nil != optErr {
			return nil, optErr
		}
	}

	return created, nil
}

func (c *v1Cache) Invalidator() Invalidator {
	i := &v1CacheInvalidator{
		cache: c.cache,
	}
	return i
}

func (c *v1Cache) Endpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		requestKey, funnelErr := c.funnel.Hash(request)
		if nil != funnelErr {
			return nil, funnelErr
		}

		response, cacheFetchErr := c.cache.Get(requestKey)
		if nil != cacheFetchErr && ErrorCacheMiss != cacheFetchErr {
			return nil, cacheFetchErr
		}

		if ErrorCacheMiss == cacheFetchErr {
			atomic.AddUint64(&c.misses, 1)

			var endpointErr error
			response, endpointErr = c.wrappedEndpoint(ctx, request)
			if nil != endpointErr {
				return nil, endpointErr
			}

			cachePutErr := c.cache.Put(requestKey, response)
			if nil != cachePutErr {
				return nil, cachePutErr
			}
		} else {
			atomic.AddUint64(&c.hits, 1)
		}

		return response, nil
	}
}

func (c *v1Cache) Hits() uint64 {
	return atomic.LoadUint64(&c.hits)
}

func (c *v1Cache) Misses() uint64 {
	return atomic.LoadUint64(&c.misses)
}

func (c *v1Cache) ResetStatistics() {
	atomic.SwapUint64(&c.hits, 0)
	atomic.SwapUint64(&c.misses, 0)
}

func (i *v1CacheInvalidator) Invalidate(key []byte) error {
	return i.cache.Invalidate(key)
}
