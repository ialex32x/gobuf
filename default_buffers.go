package gobuf

import (
	"sync"
)

type Buffers struct {
	pool *sync.Pool
}

var DefaultBuffers = NewBuffers()

func NewBuffers() *Buffers {
	buffers := &Buffers{}
	buffers.pool = &sync.Pool{
		New: func() interface{} {
			return newManagedByteBuffer(buffers)
		},
	}
	return buffers
}

func (buffers *Buffers) Alloc(capacity int) *ManagedByteBuffer {
	buf := buffers.pool.Get().(*ManagedByteBuffer)
	buf.refCount = 1
	buf.ensureCapacity(capacity)
	return buf
}
