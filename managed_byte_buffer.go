package gobuf

import (
	"sync/atomic"
)

type ManagedByteBuffer struct {
	*ByteBuffer
	refCount int32
	buffers  *Buffers // pool
	// lock     sync.Mutex
}

func newManagedByteBuffer(buffers *Buffers) *ManagedByteBuffer {
	return &ManagedByteBuffer{
		ByteBuffer: NewByteBuffer([]byte{}), // underlying byte buffer object (initially nil bytes in general situation)
		refCount:   0,                       // refCount
		buffers:    buffers,                 // manager owned this bytebuffer object
	}
}

func (buf *ManagedByteBuffer) Retain() IByteBuffer {
	// buf.lock.Lock()
	// buf.refCount++
	// buf.lock.Unlock()
	atomic.AddInt32(&buf.refCount, 1)
	return buf
}

func (buf *ManagedByteBuffer) Release() {
	// buf.lock.Lock()
	// buf.refCount--
	// refCount := buf.refCount
	// buf.lock.Unlock()
	refCount := atomic.AddInt32(&buf.refCount, -1)
	if refCount == 0 {
		buf.Clear()
		buffers := buf.buffers
		if buffers != nil {
			buffers.pool.Put(buf)
		}
	}
}
