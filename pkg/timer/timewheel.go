package timer

import (
	"runtime"
	"sync/atomic"
	"time"
)

const (
	WheelLength    = 64
	TimeNearShift  = 8
	TimeNear       = 1 << TimeNearShift
	TimeNearMask   = TimeNear - 1
	TimeLevelShift = 6
	TimeLevel      = 1 << TimeLevelShift
	TimeLevelMask  = TimeLevel - 1

	TimingTypeTimeout = 0
	TimingTypeTick    = 1
)

type TimeWheel[T any] struct {
	interval time.Duration
	near     []timeList[T]
	t        [][]timeList[T]
	time     uint32
	cb       func(T)
	locker   *spinLock
}
type spinLock struct {
	status int32
}

func (sl *spinLock) Lock() {
	for atomic.CompareAndSwapInt32(&sl.status, 0, 1) {
		runtime.Gosched()
	}
}

func (sl *spinLock) Unlock() {
	atomic.StoreInt32(&sl.status, 0)
}

type timeList[T any] struct {
	head *timeNode[T]
	tail *timeNode[T]
}

type timeNode[T any] struct {
	expire uint32
	timing uint8
	data   T
	next   *timeNode[T]
}

func newTimeList[T any]() timeList[T] {
	var zero T
	h := &timeNode[T]{
		next:   nil,
		data:   zero,
		timing: 0,
	}
	return timeList[T]{
		head: h,
		tail: h,
	}
}

func (tl *timeList[T]) addNode(node *timeNode[T]) {
	if tl.head == tl.tail {
		tl.head.next = node
	} else {
		tl.tail.next = node
	}
	tl.tail = node
	node.next = nil
}

func (tl *timeList[T]) clearList() *timeNode[T] {
	node := tl.head.next
	tl.head.next = nil
	tl.tail = tl.head

	return node
}

func NewTimeWheel[T any](interval time.Duration, cb func(T)) *TimeWheel[T] {
	tw := new(TimeWheel[T])
	tw.interval = interval
	tw.cb = cb
	tw.locker = &spinLock{}
	tw.near = make([]timeList[T], TimeNear)
	for i := range tw.near {
		tw.near[i] = newTimeList[T]()
	}
	tw.t = make([][]timeList[T], 4)
	for i := range tw.t {
		tw.t[i] = make([]timeList[T], TimeLevel)
		for k := range tw.t[i] {
			tw.t[i][k] = newTimeList[T]()
		}
	}
	go func() {
		ticker := time.NewTicker(tw.interval)
		for {
			select {
			case <-ticker.C:
				tw.tick()
			}
		}
	}()
	return tw
}

func (tw *TimeWheel[T]) Timeout(timeout time.Duration, data T) {
	if timeout <= 0 {
		tw.cb(data)
		return
	}
	node := &timeNode[T]{
		expire: tw.time + uint32(timeout/tw.interval),
		next:   nil,
		data:   data,
		timing: 0,
	}
	tw.locker.Lock()
	defer tw.locker.Unlock()
	tw.addNode(node)
}

func (tw *TimeWheel[T]) addNode(node *timeNode[T]) {
	t := node.expire
	ct := tw.time
	if (t | TimeNearMask) == (ct | TimeNearMask) {
		tw.near[t&TimeNearMask].addNode(node)
	} else {
		mask := uint32(TimeNear << TimeLevelShift)
		i := uint32(0)
		for ; i < 3; i++ {
			if (t | (mask - 1)) == (ct | (mask - 1)) {
				break
			}
			mask <<= TimeLevelShift
		}
		tw.t[i][(t>>(TimeNearShift+i*TimeLevelShift))&TimeLevelMask].addNode(node)
	}
}

func (tw *TimeWheel[T]) moveList(level, bucket int) {
	//fmt.Println("movelist", level, bucket, tw.time)
	current := tw.t[level][bucket].clearList()
	for current != nil {
		node := current
		current = current.next
		node.next = nil
		tw.addNode(node)
	}
}

func (tw *TimeWheel[T]) execute() {
	idx := tw.time & TimeNearMask
	for tw.near[idx].head.next != nil {
		current := tw.near[idx].clearList()
		tw.locker.Unlock()
		for current != nil {
			tw.cb(current.data)
			current = current.next
		}
		tw.locker.Lock()
	}
}

func (tw *TimeWheel[T]) shift() {
	tw.time += 1
	ct := tw.time
	//fmt.Println("ct", ct)
	if ct == 0 {
		tw.moveList(3, 0)
	} else {
		mask := uint32(TimeNear)
		t := ct >> TimeNearShift
		i := 0
		for ct&(mask-1) == 0 {
			idx := t & TimeLevelMask
			if idx != 0 {
				tw.moveList(i, int(idx))
				break

			}
			mask <<= TimeLevelShift
			t >>= TimeLevelShift
			i += 1
		}
	}
}

func (tw *TimeWheel[T]) tick() {
	tw.locker.Lock()
	defer tw.locker.Unlock()
	tw.shift()
	tw.execute()
}
