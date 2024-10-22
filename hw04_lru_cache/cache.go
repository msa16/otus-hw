package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache // Remove me after realization.

	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// структура для хранения в очереди
type pair struct {
	// ключ. нужен для удаления элемента из словаря при удалении элемента из очереди
	key  Key
	data interface{}
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	elem, ok := cache.items[key]
	if ok {
		// если элемент присутствует в словаре, то обновить его значение и переместить элемент в начало очереди;
		elem.Value.(*pair).data = value
		cache.queue.MoveToFront(elem)
	} else {
		// если элемента нет в словаре, то добавить в словарь и в начало очереди
		// (при этом, если размер очереди больше ёмкости кэша, то необходимо удалить последний элемент из очереди и его значение из словаря);
		if cache.queue.Len() == cache.capacity {
			back := cache.queue.Back()
			cache.queue.Remove(back)
			delete(cache.items, back.Value.(*pair).key)
		}
		cache.items[key] = cache.queue.PushFront(&pair{key: key, data: value})
	}
	// возвращаемое значение - флаг, присутствовал ли элемент в кэше.
	return ok
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	elem, ok := cache.items[key]
	if ok {
		// если элемент присутствует в словаре, то переместить элемент в начало очереди и вернуть его значение и true;
		cache.queue.MoveToFront(elem)
		return elem.Value.(*pair).data, true
	}
	// если элемента нет в словаре, то вернуть nil и false
	return nil, false
}

func (cache *lruCache) Clear() {
	// чтобы не удалять по одному, пересоздаем словарь и очередь
	cache.items = make(map[Key]*ListItem, cache.capacity)
	cache.queue = NewList()
}
