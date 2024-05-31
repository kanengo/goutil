package utils

import (
	"io"
	"sync"
)

func CopyZeroAlloc(w io.Writer, r io.Reader) (int64, error) {
	if wt, ok := r.(io.WriterTo); ok {
		return wt.WriteTo(w)
	}
	if rt, ok := w.(io.ReaderFrom); ok {
		return rt.ReadFrom(r)
	}
	vBuf := copyBufPool.Get()
	buf := vBuf.([]byte)
	n, err := io.CopyBuffer(w, r, buf)
	copyBufPool.Put(vBuf)

	return n, err
}

var copyBufPool = sync.Pool{
	New: func() any {
		return make([]byte, 4096)
	},
}
