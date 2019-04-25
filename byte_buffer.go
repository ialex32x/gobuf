package gobuf

import (
	"encoding/binary"
	"log"
	"math"
)

type IByteBuffer interface {
	Retain() IByteBuffer
	Release()
	Clear()
	AddSize(size int) int
	Size() int
	Index() int
	Print()
	Tag() int
	SetTag(tag int) int

	Capacity() int
	Readable() int
	Bytes() []byte
	SharedBytes(start, end int) []byte
	Seek(index int)
	PushState()
	PopState(restore bool)

	ReadFloat32() (val float32)
	ReadFloat64() (val float64)
	ReadUint64BE() (val uint64)
	ReadUint32BE() (val uint32)
	ReadInt32BE() (val int32)
	ReadUint16BE() (val uint16)
	ReadUint64LE() (val uint64)
	ReadUint32LE() (val uint32)
	ReadInt32LE() (val int32)
	ReadUint16LE() (val uint16)
	ReadUint8() (val uint8)

	WriteUint8(v uint8)
	WriteUint16BE(v uint16)
	WriteUint32BE(v uint32)
	WriteInt32BE(v int32)
	WriteUint64BE(v uint64)
	WriteUint16LE(v uint16)
	WriteUint32LE(v uint32)
	WriteInt32LE(v int32)
	WriteUint64LE(v uint64)
	WriteFloat32(v float32)
	WriteFloat64(v float64)
	WriteBytes(data []byte)
}

type ByteBufferState struct {
	index int // reader position
	size  int // writer position
}

type ByteBuffer struct {
	tag    int
	data   []byte
	state  ByteBufferState
	states []ByteBufferState
}

// NewByteBuffer 创建无内存池管理的 ByteBuffer 实例
func NewByteBuffer(bytes []byte, tag int) *ByteBuffer {
	return &ByteBuffer{
		tag:  tag,
		data: bytes, // Underlying Bytes
		state: ByteBufferState{
			index: 0,          // Position for Reader
			size:  len(bytes), // Size
		},
	}
}

func (buf *ByteBuffer) Tag() int {
	return buf.tag
}

func (buf *ByteBuffer) SetTag(tag int) int {
	old := buf.tag
	buf.tag = tag
	return old
}

func (buf *ByteBuffer) Print() {
	log.Printf("buffer index=%v size=%v bytes=%v",
		buf.state.index,
		buf.state.size,
		buf.SharedBytes(buf.state.index, buf.state.size))
}

// 不能中途改容量 (产生slice分裂后会产生拷贝, 与原数据分离)
func (buf *ByteBuffer) ensureCapacity(size int) {
	if size > 0 {
		delta := size - len(buf.data)
		if delta > 0 {
			resized := make([]byte, size)
			copy(resized, buf.data)
			buf.data = resized
		}
	}
}

// PushState 存储当前状态
func (buf *ByteBuffer) PushState() {
	buf.states = append(buf.states, buf.state)
}

// PopState 移除存储的状态, 如果 restore=true 则当前状态赋值为该状态
func (buf *ByteBuffer) PopState(restore bool) {
	size := len(buf.states)
	if restore {
		buf.state = buf.states[size-1]
	}
	buf.states = buf.states[:size-1]
}

// Retain 增加引用计数
func (buf *ByteBuffer) Retain() IByteBuffer {
	return buf
}

// Release 减少引用计数
func (buf *ByteBuffer) Release() {
}

// Clear 清空 (不覆盖数据)
func (buf *ByteBuffer) Clear() {
	buf.state.size = 0
	buf.state.index = 0
}

// Readable 是否还有数据可读
func (buf *ByteBuffer) Readable() int {
	return buf.state.size - buf.state.index
}

// Capacity 底层容量
func (buf *ByteBuffer) Capacity() int {
	return len(buf.data)
}

// AddSize 直接增加大小 (需要通过其他方式填充数据, 比如 SharedBytes)
func (buf *ByteBuffer) AddSize(size int) int {
	buf.state.size += size
	return buf.state.size
}

// Index 当前读索引
func (buf *ByteBuffer) Index() int {
	return buf.state.index
}

// Size 内容长度 (当前写索引)
func (buf *ByteBuffer) Size() int {
	return buf.state.size
}

// SharedBytes 底层数据的指定切片 (不会改变状态)
func (buf *ByteBuffer) SharedBytes(start, end int) []byte {
	return buf.data[start:end]
}

func (buf *ByteBuffer) Bytes() []byte {
	return buf.data[buf.state.index:buf.state.size]
}

// Seek 指定度索引位置
func (buf *ByteBuffer) Seek(index int) {
	if index < buf.state.size {
		buf.state.index = index
	}
}

func (buf *ByteBuffer) ReadFloat32() (val float32) {
	val = math.Float32frombits(buf.ReadUint32LE())
	return
}

func (buf *ByteBuffer) ReadFloat64() (val float64) {
	val = math.Float64frombits(buf.ReadUint64LE())
	return
}

func (buf *ByteBuffer) ReadUint64BE() (val uint64) {
	val = binary.BigEndian.Uint64(buf.data[buf.state.index : buf.state.index+8])
	buf.state.index += 8
	return
}

func (buf *ByteBuffer) ReadUint32BE() (val uint32) {
	val = binary.BigEndian.Uint32(buf.data[buf.state.index : buf.state.index+4])
	buf.state.index += 4
	return
}

func (buf *ByteBuffer) ReadInt32BE() (val int32) {
	val = int32(binary.BigEndian.Uint32(buf.data[buf.state.index : buf.state.index+4]))
	buf.state.index += 4
	return
}

func (buf *ByteBuffer) ReadUint16BE() (val uint16) {
	val = binary.BigEndian.Uint16(buf.data[buf.state.index : buf.state.index+2])
	buf.state.index += 2
	return
}

func (buf *ByteBuffer) ReadUint64LE() (val uint64) {
	val = binary.LittleEndian.Uint64(buf.data[buf.state.index : buf.state.index+8])
	buf.state.index += 8
	return
}

func (buf *ByteBuffer) ReadUint32LE() (val uint32) {
	val = binary.LittleEndian.Uint32(buf.data[buf.state.index : buf.state.index+4])
	buf.state.index += 4
	return
}

func (buf *ByteBuffer) ReadInt32LE() (val int32) {
	val = int32(binary.LittleEndian.Uint32(buf.data[buf.state.index : buf.state.index+4]))
	buf.state.index += 4
	return
}

func (buf *ByteBuffer) ReadUint16LE() (val uint16) {
	val = binary.LittleEndian.Uint16(buf.data[buf.state.index : buf.state.index+2])
	buf.state.index += 2
	return
}

func (buf *ByteBuffer) ReadUint8() (val uint8) {
	val = uint8(buf.data[buf.state.index])
	buf.state.index += 1
	return
}

func (buf *ByteBuffer) WriteUint8(v uint8) {
	buf.data[buf.state.size] = byte(v)
	buf.state.size += 1
}

func (buf *ByteBuffer) WriteUint16BE(v uint16) {
	binary.BigEndian.PutUint16(buf.data[buf.state.size:buf.state.size+2], v)
	buf.state.size += 2
}

func (buf *ByteBuffer) WriteUint32BE(v uint32) {
	binary.BigEndian.PutUint32(buf.data[buf.state.size:buf.state.size+4], v)
	buf.state.size += 4
}

func (buf *ByteBuffer) WriteInt32BE(v int32) {
	binary.BigEndian.PutUint32(buf.data[buf.state.size:buf.state.size+4], uint32(v))
	buf.state.size += 4
}

func (buf *ByteBuffer) WriteUint64BE(v uint64) {
	binary.BigEndian.PutUint64(buf.data[buf.state.size:buf.state.size+8], v)
	buf.state.size += 8
}

func (buf *ByteBuffer) WriteUint16LE(v uint16) {
	binary.LittleEndian.PutUint16(buf.data[buf.state.size:buf.state.size+2], v)
	buf.state.size += 2
}

func (buf *ByteBuffer) WriteUint32LE(v uint32) {
	binary.LittleEndian.PutUint32(buf.data[buf.state.size:buf.state.size+4], v)
	buf.state.size += 4
}

func (buf *ByteBuffer) WriteInt32LE(v int32) {
	binary.LittleEndian.PutUint32(buf.data[buf.state.size:buf.state.size+4], uint32(v))
	buf.state.size += 4
}

func (buf *ByteBuffer) WriteUint64LE(v uint64) {
	binary.LittleEndian.PutUint64(buf.data[buf.state.size:buf.state.size+8], v)
	buf.state.size += 8
}

func (buf *ByteBuffer) WriteFloat32(v float32) {
	buf.WriteUint32LE(math.Float32bits(v))
}

func (buf *ByteBuffer) WriteFloat64(v float64) {
	buf.WriteUint64LE(math.Float64bits(v))
}

// WriteBytes 写入指定数据
func (buf *ByteBuffer) WriteBytes(data []byte) {
	size := len(data)
	copy(buf.data[buf.state.size:buf.state.size+size], data)
	buf.state.size += size
}
