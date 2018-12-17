package kitcache

import (
	"errors"

	"github.com/go-kit/kit/endpoint"
)

var (
	// ErrorCacheMiss is returned if the item could not be found in the cache
	ErrorCacheMiss = errors.New("cached item not available")
)

// EndpointCache is the basic unit of the kitcache library. It will cache
// responses serviced by the wrapped endpoint
type EndpointCache interface {
	// Invalidator fetches a goroutine-safe Invalidator to remove items from
	// the cache
	Invalidator() Invalidator

	// Endpoint will generate an endpoint which will utilize the cache and call
	// the wrapped endpoint on cache misses
	Endpoint() endpoint.Endpoint

	// Hits returns how many times the cache successfully satisfied a request
	// from the cache only
	Hits() uint64

	// Misses returns how many times a request was processed without having been
	// previously cached
	Misses() uint64

	//ResetStatistics reset the hits & misses counters
	ResetStatistics()
}

// RequestFunnel facilitates converting a request into a key used in the cache
type RequestFunnel interface {
	// Hash should take a request provided to the cached endpoint and produces
	// a key under which the item should be cached
	Hash(interface{}) ([]byte, error)
}

// ResponseMarshaler will marshal a request object to and from the cache
type ResponseMarshaler interface {

	// Marshal saves the provided empty interface to a byte array and can
	// provide an error if unsuccessful
	Marshal(interface{}) ([]byte, error)

	// Unmarshal loads the provided bytes into an object viable to be served in
	// place of the cache forwarding a request to the wrapped Endpoint.
	Unmarshal([]byte) (interface{}, error)
}

// Invalidator is used to remove cached items
type Invalidator interface {
	// Invalidate will remove the provided cache key from the underlying cache.
	// It returns the error encountered while removing the a cached object
	Invalidate([]byte) error
}

// Cache is the interface the backing cache implementation must satisfy to be
// used by kitcache
type Cache interface {
	// Get attempts to fetch a cached item by its key. If not found,
	// nil and ErrorCacheMiss will be returned
	Get([]byte) (interface{}, error)

	// Put takes a cache key and item to cache and emplaces it under the key
	// provided. If the item exists, then it is overwritten
	Put([]byte, interface{}) error

	// Invalidate removes an item from the cache.
	Invalidate([]byte) error
}

// V1ConstructorOption is used by the cache implementation for optional parameters
type V1ConstructorOption func(i *v1Cache) error
