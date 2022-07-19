package collection

type Iterator[T any] interface {
	// Range calls f sequentially for each key and value present in the map.
	// If f returns false, range stops the iteration.
	Range(f func(val T) bool)
	//Slice 会返回一个新的符合条件的Slice
	Slice() Slice[T]
	//Map 根据输入值映射对应的输出值
	Map(f func(val T) T) Iterator[T]
	//Filter 过滤符合条件的元素
	Filter(f func(val T) bool) Iterator[T]
}

const (
	IterTypeNone = iota
	IterTypeFilter
	IterTypeMap
)
