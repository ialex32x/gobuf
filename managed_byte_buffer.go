package gobuf

import (
	"sync/atomic"
)

type ManagedByteBuffer struct {
	*ByteBuffer
	refCount int32
	buffers  *Buffers // pool
}

func newManagedByteBuffer(buffers *Buffers) *ManagedByteBuffer {
	return &ManagedByteBuffer{
		ByteBuffer: NewByteBuffer([]byte{}, 0), // underlying byte buffer object (initially nil bytes in general situation)
		refCount:   0,                          // refCount
		buffers:    buffers,                    // manager owned this bytebuffer object
	}
}

func (buf *ManagedByteBuffer) Retain() IByteBuffer {
	atomic.AddInt32(&buf.refCount, 1)
	return buf
}

func (buf *ManagedByteBuffer) Release() {
	refCount := atomic.AddInt32(&buf.refCount, -1)
	if refCount == 0 {
		buf.Clear()
		buffers := buf.buffers
		if buffers != nil {
			buffers.pool.Put(buf)
			// atomic.AddInt32(&buffers.alive, -1)
			// fmt.Printf("buffers release %v\n", buffers.alive)
		}
	}
}
