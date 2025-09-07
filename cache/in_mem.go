package cache

type InMemoryCache struct {
	store map[string]string
}

func NewInMemoryCacheRepository() *InMemoryCache {
	return &InMemoryCache{
		store: make(map[string]string),
	}
}

func (c *InMemoryCache) Get(key string) (string, error) {
	val, exists := c.store[key]
	if !exists {
		return "", nil
	}
	return val, nil
}

func (c *InMemoryCache) Put(key string, val string) error {
	c.store[key] = val
	return nil
}

func (c *InMemoryCache) Delete(key string) error {
	delete(c.store, key)
	return nil
}
