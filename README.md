# kitcache

kitcache is an endpoint response caching framework and implementation to be used
with the [go-kit](https://github.com/go-kit/kit) framework.

## Usage

To cache a go-kit endpoint in with this framework:
1. implement a RequestFunnel for the endpoint's request value to create cache
keys
```go
type UserProfileRequest struct {
	UserId uint64
	RequestId string
}

type UserProfileResponse struct {}

type UserProfileCacheFunnel struct {}
func (f *UserProfileCacheFunnel) Hash(requestIface interface{}) ([]byte,
	error) {
	request := requestIface.(*UserProfileRequest)
	hashValue := fmt.Sprintf("%v", request.UserId)
	return []byte(hashValue), nil
}
...
funnel := &UserProfileCacheFunnel{}
```
2. implement a RequestMarshaler for the endpoint's request values or use the
provided [gob-based](https://golang.org/pkg/encoding/gob) marshaler
```go
func responseAllocator() interface{} {
	return &UserProfileResponse{}
}
...
marshaler := marshal.NewGobMarshaler(exampleStructAllocator)
```
3. implement a cache or use the built-in
[ccache-based](https://github.com/karlseguin/ccache) cache
```go
backingCache := cache.NewInMemoryCache(10, 10*time.Hour)
```
4. wrap the endpoint to be cached
```go
serviceEndpoint := BuildUserProfileEndpoint()
endpointCache, cacheCreateErr := kitcache.NewEndpointCache(serviceEndpoint,
	backingCache, marshaler, funnel)

// use this in the place of serviceEndpoint to use the cache
cachedEndpoint := endpointCache.Endpoint()
```
5. (optional) Invalidate cached items
```go
cacheLine := []byte()
cacheInvalidator := endpointCache.Invalidator()
cacheInvalidator.Invalidate(cacheLine)
```
