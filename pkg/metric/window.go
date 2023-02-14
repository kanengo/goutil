package metric

import "fmt"

type Bucket struct {
	Points []float64
	Count  int64
	next   *Bucket
}

func (b *Bucket) Append(val float64) {
	b.Points = append(b.Points, val)
	b.Count++
}

func (b *Bucket) Add(offset int, val float64) {
	if offset < 0 || offset >= len(b.Points) {
		return
	}
	b.Points[offset] += val
	b.Count++
}

func (b *Bucket) Reset() {
	b.Points = b.Points[:0]
	b.Count = 0
}

func (b *Bucket) Next() *Bucket {
	return b.next
}

type BucketIterator struct {
	pos    int
	curPos int
	cur    *Bucket
}

func (i *BucketIterator) Next() bool {
	return i.pos != i.curPos
}

func (i *BucketIterator) Bucket() Bucket {
	if !i.Next() {
		panic(fmt.Errorf("metric bucket iteration out of range pos:%d curPos %d", i.pos, i.curPos))
	}
	bucket := *i.cur
	i.curPos++
	i.cur = i.cur.Next()
	return bucket
}

type Window struct {
	window []Bucket
	size   int
}

type WindowOpts struct {
	Size int
}

func NewWindow(opts WindowOpts) *Window {
	buckets := make([]Bucket, opts.Size)
	for offset := range buckets {
		buckets[offset] = Bucket{
			Points: make([]float64, 0),
			Count:  0,
			next:   nil,
		}
		nextOffset := offset + 1
		if nextOffset == opts.Size {
			nextOffset = 0
		}
		buckets[offset].next = &buckets[nextOffset]
	}
	return &Window{window: buckets, size: opts.Size}
}

func (w *Window) ResetBucket(offset int) {
	w.window[offset].Reset()
}

func (w *Window) ResetBuckets(offsets []int) {
	for _, offset := range offsets {
		w.ResetBucket(offset)
	}
}

func (w *Window) ResetWindow() {
	for offset := range w.window {
		w.ResetBucket(offset)
	}
}

func (w *Window) Append(offset int, val float64) {
	w.window[offset].Append(val)
}

func (w *Window) Add(offset int, val float64) {
	if w.window[offset].Count == 0 {
		w.window[offset].Append(val)
		return
	}

	w.window[offset].Add(0, val)
}

func (w *Window) Bucket(offset int) Bucket {
	return w.window[offset]
}

func (w *Window) Size() int {
	return w.size
}

func (w *Window) Iterator(offset int, count int) BucketIterator {
	return BucketIterator{
		pos:    count,
		curPos: 0,
		cur:    &w.window[offset],
	}
}
