package cache

import (
	"github.com/karlseguin/ccache"
	"github.com/mrwinstead/kitcache"
	"time"
)

type inMemoryCache struct {
	cache     *ccache.Cache
	marshaler kitcache.ResponseMarshaler
	itemTTL   time.Duration
}

// NewInMemoryCache constructs a new in-cache cache
func NewInMemoryCache(maxItemCount uint32, itemTTL time.Duration,
) kitcache.Cache {
	config := ccache.Configure().MaxSize(int64(maxItemCount))

	created := &inMemoryCache{
		cache:   ccache.New(config),
		itemTTL: itemTTL,
	}
	return created
}

// Get fetches an item from the cache given the provided key. If the an item is
// not found under the provided value, then it will return
// kitcache.ErrorCacheMiss
func (i *inMemoryCache) Get(key []byte) (interface{}, error) {
	found := i.cache.Get(string(key))
	if nil != found && !found.Expired() {
		return found.Value(), nil
	}
	if nil == found {
		return nil, kitcache.ErrorCacheMiss
	}
	return nil, nil
}

// Put emplaces a key/value pair into the underlying cache implementation
func (i *inMemoryCache) Put(key []byte, value interface{}) error {
	i.cache.Set(string(key), value, i.itemTTL)
	return nil
}

// Invalidate deletes the key from the in-cache cache. It does not check
// whether the underlying implementation had the specified key
func (i *inMemoryCache) Invalidate(key []byte) error {
	i.cache.Delete(string(key))
	return nil
}
