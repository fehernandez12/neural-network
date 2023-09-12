package cache

type Cache interface {
	Get(key string) (string, error)
	Put(key string, value string) error
	Delete(key string) error
}

var impl Cache

func SetRepository(repository Cache) {
	impl = repository
}

func Get(key string) (string, error) {
	return impl.Get(key)
}

func Put(key string, value string) error {
	return impl.Put(key, value)
}

func Delete(key string) error {
	return impl.Delete(key)
}
