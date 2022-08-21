package collection

type Iterator[T any] interface {
	// Range calls f sequentially for each key and value present in the map.
	// If f returns false, range stops the iteration.
	Range(f func(val T) bool)
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
