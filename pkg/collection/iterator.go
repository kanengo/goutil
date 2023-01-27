package collection

type Ranger[T any] interface {
	// Range calls f sequentially for each key and value present in the map.
	// If f returns false, range stops the iteration.
	Range(f func(val T) bool)
}

type Iterator[T any] interface {
	Ranger[T]
	//Slice 会返回一个新的符合条件的Slice
	Slice() *Slice[T]
	//Map 根据输入值映射对应的输出值
	Map(f func(val T) T) Iterator[T]
	//Filter 过滤符合条件的元素
	Filter(f func(val T) bool) Iterator[T]
	//Reduce 类似js的reduce
	Reduce(f func(previousValue any, val T) any, initialValue any) any
}

type Pairs[K comparable, V any] struct {
	k K
	v V
}

func (p *Pairs[K, V]) Key() K {
	return p.k
}

func (p *Pairs[K, V]) Value() V {
	return p.v
}

func (p *Pairs[K, V]) SetKey(k K) {
	p.k = k
}

func (p *Pairs[K, V]) SetValue(v V) {
	p.v = v
}

const (
	IterTypeNone = iota
	IterTypeFilter
	IterTypeMap
)

type iterTypeFunc[T any] struct {
	typ     int
	mapF    func(val T) T
	filterF func(val T) bool
	next    *iterTypeFunc[T]
}

type defaultIterator[T any] struct {
	Ranger[T]
	next *iterTypeFunc[T]
}

func (di *defaultIterator[T]) Range(fn func(val T) bool) {
	di.Ranger.Range(func(val T) bool {
		iterF := di.next
		for iterF != nil {
			switch iterF.typ {
			case IterTypeFilter:
				if !iterF.filterF(val) {
					goto NEXT
				}
			case IterTypeMap:
				val = iterF.mapF(val)
			}
			iterF = iterF.next
		}
		if !fn(val) {
			return false
		}
	NEXT:
		return true
	})
}

func (di *defaultIterator[T]) Slice() *Slice[T] {
	ret := NewSlice[T]()
	di.Range(func(val T) bool {
		ret.Append(val)
		return true
	})
	return ret
}

func (di *defaultIterator[T]) insertFunc(itf *iterTypeFunc[T]) {
	if di.next == nil {
		di.next = itf
	} else {
		di.next.next = itf
	}
}

func (di *defaultIterator[T]) Map(f func(val T) T) Iterator[T] {
	itf := &iterTypeFunc[T]{
		typ:  IterTypeMap,
		mapF: f,
	}

	di.insertFunc(itf)

	return di
}

func (di *defaultIterator[T]) Filter(f func(val T) bool) Iterator[T] {
	itf := &iterTypeFunc[T]{
		typ:     IterTypeFilter,
		filterF: f,
	}
	di.insertFunc(itf)
	return di
}

func (di *defaultIterator[T]) Reduce(f func(previousValue any, val T) any, initialValue any) any {
	ret := initialValue
	di.Range(func(val T) bool {
		ret = f(ret, val)
		return true
	})
	return ret
}
