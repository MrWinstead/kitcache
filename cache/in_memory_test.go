package cache_test

import (
	"testing"
	"time"

	"github.com/pallinder/go-randomdata"
	"github.com/stretchr/testify/assert"

	"github.com/mrwinstead/kitcache"
	"github.com/mrwinstead/kitcache/cache"
)

func TestNewInMemoryCache(t *testing.T) {
	created := cache.NewInMemoryCache(10, 10*time.Hour)
	assert.NotNil(t, created)
}

func TestInMemoryCache_Put_Get(t *testing.T) {
	created := cache.NewInMemoryCache(10, 10*time.Hour)
	key := []byte(randomdata.RandStringRunes(10))
	value := []byte(randomdata.RandStringRunes(10))

	putErr := created.Put(key, value)
	assert.Nil(t, putErr)

	fetched, fetchErr := created.Get(key)
	assert.Nil(t, fetchErr)

	assert.Equal(t, value, fetched)
}

func TestInMemoryCache_Get_Put_Get(t *testing.T) {
	created := cache.NewInMemoryCache(10, 10*time.Hour)
	key := []byte(randomdata.RandStringRunes(10))
	value := []byte(randomdata.RandStringRunes(10))

	missingValue, shouldBeCacheMiss := created.Get(key)
	assert.Equal(t, kitcache.ErrorCacheMiss, shouldBeCacheMiss)
	assert.Nil(t, missingValue)

	putErr := created.Put(key, value)
	assert.Nil(t, putErr)

	fetched, fetchErr := created.Get(key)
	assert.Nil(t, fetchErr)
	assert.Equal(t, value, fetched)
}

func TestInMemoryCache_Put_Get_Invalidate_Get(t *testing.T) {
	created := cache.NewInMemoryCache(10, 10*time.Hour)
	key := []byte(randomdata.RandStringRunes(10))
	value := []byte(randomdata.RandStringRunes(10))

	putErr := created.Put(key, value)
	assert.Nil(t, putErr)

	fetched, fetchErr := created.Get(key)
	assert.Nil(t, fetchErr)
	assert.Equal(t, value, fetched)

	invalidationErr := created.Invalidate(key)
	assert.Nil(t, invalidationErr)

	invalidatedValue, shouldBeCacheMiss := created.Get(key)
	assert.Equal(t, kitcache.ErrorCacheMiss, shouldBeCacheMiss)
	assert.Nil(t, invalidatedValue)
}
