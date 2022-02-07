// localcache is simple key/value cache store
package localcache

type Cache interface {
	Get(k string) (interface{}, error)
	Set(k string, v interface{}) error
}