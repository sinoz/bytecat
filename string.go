package bytecat

import (
	"errors"
	"io"
	"strings"
)

// emptyString is a String with an empty byte array.
var emptyString = &String{bytes: make([]byte, 0)}

// String is an immutable rope-like sequence of bytes that allows
// for fast concatenation without the copying overhead.
type String struct {
	bytes []byte
}

// Builder builds up a String. The internal byte array may grow
// as more bytes are written to this buffer.
type Builder struct {
	bytes []byte
	index int
}

// Iterator iterates through a String.
type Iterator struct {
	bytes *String
	index int
}

// NewDefaultBuilder produces a Builder with a default capacity of 16 bytes.
func NewDefaultBuilder() *Builder {
	return NewBuilder(16)
}

// NewBuilder constructs a new Builder with the specified initial capacity.
func NewBuilder(initialCapacity int) *Builder {
	return &Builder{bytes: make([]byte, initialCapacity)}
}

// EmptyString returns an empty String without any bytes in it.
func EmptyString() *String {
	return emptyString
}

// StringOf constructs a String with the given values.
func StringOf(values ...byte) *String {
	return StringWrap(values)
}

// StringWrap wraps the given byte array into a String instance.
func StringWrap(bytes []byte) *String {
	return &String{bytes: bytes}
}

// Drop drops the specified amount of bytes from this String. If more
// bytes are requested to be dropped than there are available, an empty
// String is returned.
func (s *String) Drop(amount int) *String {
	if s.IsEmpty() || s.Length() < amount {
		return emptyString
	}

	return StringWrap(s.bytes[amount:])
}

// Take takes a slice of this String with the specified amount of bytes.
// If more bytes are requested than there are available, only the available
// amount is provided in the returned String.
func (s *String) Take(amount int) *String {
	if s.IsEmpty() {
		return emptyString
	}

	if s.Length() < amount {
		amount = s.Length()
	}

	return StringWrap(s.bytes[:amount])
}

// Concat concatenates this String with the given String, creating a
// new String with both values concatenated to each other in a rope-like
// fashion.
func (s *String) Concat(other *String) *String {
	if s.IsEmpty() && other.IsEmpty() {
		return emptyString
	}

	if s.IsEmpty() && !other.IsEmpty() {
		return other
	}

	if !s.IsEmpty() && other.IsEmpty() {
		return s
	}

	return StringWrap(append(s.bytes, other.bytes...))
}

// IsEmpty returns whether this String is empty or not.
func (s *String) IsEmpty() bool {
	return s.Length() == 0
}

// Length returns the complete length of this String, in bytes.
func (s *String) Length() int {
	return len(s.bytes)
}

// ByteAt reads a byte value at the specified index.
func (s *String) ByteAt(index int) (byte, error) {
	if index < 0 || index > len(s.bytes) {
		return 0, errors.New("index out of bounds")
	}

	return s.bytes[index], nil
}

// ToByteArray exposes the contents of this byte String as a byte array.
func (s *String) ToByteArray() []byte {
	return s.bytes
}

// Iterator exposes a Iterator that allows the client caller to
// iterate through the String as if it were a regular buffer.
func (s *String) Iterator() *Iterator {
	return &Iterator{bytes: s}
}

// Append appends the given series of bytes to the tail of the
// String, resulting in a byte String.
func (s *String) Append(values ...byte) *String {
	return s.Concat(StringWrap(values))
}

// Prepend prepends the given series of bytes to the front of
// the String, resulting in a new byte String.
func (s *String) Prepend(values ...byte) *String {
	return StringWrap(values).Concat(s)
}

// ReadByte reads a single byte at the current index and advances the index
// after reading it. May return an error if the index is out of bounds.
func (i *Iterator) ReadByte() (byte, error) {
	value, err := i.bytes.ByteAt(i.index)
	i.index++
	return value, err
}

// GetByte reads a single byte at the current index. May return an error
// if the index is out of bounds.
func (i *Iterator) GetByte(index int) (byte, error) {
	return i.bytes.ByteAt(index)
}

// ReadBool reads a single byte at the current index, advances the index
// after reading it and returns it as a boolean value. May return an error
// if the index is out of bounds.
func (i *Iterator) ReadBool() (bool, error) {
	value, err := i.ReadByte()
	if err != nil {
		return false, err
	}

	return value == 1, nil
}

// ReadUInt16 reads an unsigned 16-bit integer at the current index and
// advances the index after reading it. May return an error if the index
// is out of bounds.
func (i *Iterator) ReadUInt16() (uint16, error) {
	if err := i.testAvailableBytes(2); err != nil {
		return 0, err
	}

	v1, _ := i.bytes.ByteAt(i.index)
	v2, _ := i.bytes.ByteAt(i.index + 1)

	i.index += 2
	return uint16(v1)<<8 | uint16(v2), nil
}

// ReadUInt24 reads an unsigned 24-bit integer at the current index and
// advances the index after reading it. May return an error if the index
// is out of bounds.
func (i *Iterator) ReadUInt24() (uint32, error) {
	if err := i.testAvailableBytes(3); err != nil {
		return 0, err
	}

	v1, _ := i.bytes.ByteAt(i.index)
	v2, _ := i.bytes.ByteAt(i.index + 1)
	v3, _ := i.bytes.ByteAt(i.index + 2)

	i.index += 3
	return uint32(v1)<<16 | uint32(v2)<<8 | uint32(v3), nil
}

// ReadUInt32 reads an unsigned 32-bit integer at the current index and
// advances the index after reading it. May return an error if the index
// is out of bounds.
func (i *Iterator) ReadUInt32() (uint32, error) {
	if err := i.testAvailableBytes(4); err != nil {
		return 0, err
	}

	v1, _ := i.bytes.ByteAt(i.index)
	v2, _ := i.bytes.ByteAt(i.index + 1)
	v3, _ := i.bytes.ByteAt(i.index + 2)
	v4, _ := i.bytes.ByteAt(i.index + 3)

	i.index += 4
	return uint32(v1)<<24 | uint32(v2)<<16 | uint32(v3)<<8 | uint32(v4), nil
}

// ReadUInt64 reads an unsigned 64-bit integer at the current index and
// advances the index after reading it. May return an error if the index
// is out of bounds.
func (i *Iterator) ReadUInt64() (uint64, error) {
	if err := i.testAvailableBytes(8); err != nil {
		return 0, err
	}

	v1, _ := i.bytes.ByteAt(i.index)
	v2, _ := i.bytes.ByteAt(i.index + 1)
	v3, _ := i.bytes.ByteAt(i.index + 2)
	v4, _ := i.bytes.ByteAt(i.index + 3)
	v5, _ := i.bytes.ByteAt(i.index + 4)
	v6, _ := i.bytes.ByteAt(i.index + 5)
	v7, _ := i.bytes.ByteAt(i.index + 6)
	v8, _ := i.bytes.ByteAt(i.index + 7)

	i.index += 8
	return uint64(v1)<<56 | uint64(v2)<<48 | uint64(v3)<<40 | uint64(v4)<<32 | uint64(v5)<<24 | uint64(v6)<<16 | uint64(v7)<<8 | uint64(v8), nil
}

// ReadString reads a series of characters and outputs a string value.
func (i *Iterator) ReadString(length int) (string, error) {
	var bldr strings.Builder

	for j := 0; j < length && i.IsReadable(); j++ {
		charByte, err := i.ReadByte()
		if err != nil {
			return "", err
		}

		bldr.WriteByte(charByte)
	}

	return bldr.String(), nil
}

// SkipBytes skips the specified amount of bytes within the iterator.
// Does not exceed the amount of bytes that are within the byte String.
func (i *Iterator) SkipBytes(amount int) {
	i.index += amount
	if i.index >= i.bytes.Length() {
		i.index = i.bytes.Length() - 1
	}
}

// Index returns the current reader index within the iterator.
func (i *Iterator) Index() int {
	return i.index
}

// SetIndex sets the current reader index to the specified value.
func (i *Iterator) SetIndex(v int) {
	i.index = v
}

// ReadableBytes returns the amount of bytes the iterator has left to read.
func (i *Iterator) ReadableBytes() int {
	return i.bytes.Length() - i.index
}

// CanRead returns whether the iterator has at least the specified amount
// of bytes left to read.
func (i *Iterator) CanRead(amount int) bool {
	return (i.index + amount) <= i.bytes.Length()
}

// IsReadable returns whether the iterator has any bytes left to read.
func (i *Iterator) IsReadable() bool {
	return i.index < i.bytes.Length()
}

// IsEmpty returns whether this Iterator has any bytes left to read.
func (i *Iterator) IsEmpty() bool {
	return !i.CanRead(1)
}

// testAvailableBytes tests whether the specified amount of bytes are
// available for reading.
func (i *Iterator) testAvailableBytes(amount int) error {
	if (i.index + amount) > i.bytes.Length() {
		return errors.New("index out of bounds")
	}

	return nil
}

// Read reads the next len(p) bytes from the iterator or until the
// iterator is drained. The return value n is the number of bytes read.
// If the iterator has no data to return, err is io.EOF (unless len(p)
// is zero); otherwise it is nil.
func (i *Iterator) Read(p []byte) (n int, err error) {
	if i.IsEmpty() {
		return 0, io.EOF
	}

	n = copy(p, i.bytes.bytes[i.index:])
	i.index += n

	return n, nil
}

// Write appends the contents of p to the Builder, growing the Builder as
// needed. The return value n is the length of p; err is always nil. If the
// buffer becomes too large, Write will panic with ErrTooLarge.
func (b *Builder) Write(p []byte) (n int, err error) {
	b.EnsureWritable(len(p))
	amtBytesWritten := copy(b.bytes[b.index:], p)
	b.index += amtBytesWritten
	return amtBytesWritten, nil
}

// Index returns the current writer index within the builder.
func (b *Builder) Index() int {
	return b.index
}

// SetIndex sets the current writer index to the specified value.
func (b *Builder) SetIndex(v int) {
	b.index = v
}

// InspectAt gets a single byte at the current index. May return an error
// if the index is out of bounds.
func (b *Builder) InspectAt(index int) (byte, error) {
	if index < 0 || index > len(b.bytes) {
		return 0, errors.New("index out of bounds")
	}

	return b.bytes[index], nil
}

// Capacity returns this builder's current capacity.
func (b *Builder) Capacity() int {
	return len(b.bytes)
}

// IsWritable returns whether a single byte can be written now to the builder.
func (b *Builder) IsWritable() bool {
	return b.CanWrite(1)
}

// CanWrite returns whether the specified amount of bytes can be written
// to this builder without the act of regrowth.
func (b *Builder) CanWrite(amount int) bool {
	return (b.Capacity() - b.index) >= amount
}

// EnsureWritable ensures that the specified amount of bytes can be written.
func (b *Builder) EnsureWritable(amount int) {
	for !b.CanWrite(amount) {
		trail := make([]byte, b.Capacity())
		b.bytes = append(b.bytes, trail...)
	}
}

// WriteByte writes the given byte value as a single byte to the builder.
func (b *Builder) WriteByte(value byte) *Builder {
	b.EnsureWritable(1)

	b.bytes[b.index] = value
	b.index++

	return b
}

// SetByte sets the specified byte value at the specified index within
// the builder. May return an error if the given index is out of bounds.
func (b *Builder) SetByte(index int, value byte) error {
	if index < 0 || index >= len(b.bytes) {
		return errors.New("index out of bounds")
	}

	b.bytes[index] = value

	return nil
}

// WriteInt16 writes the given int16 value as a 16-bit integer to the builder.
func (b *Builder) WriteInt16(value int16) *Builder {
	b.EnsureWritable(2)

	b.bytes[b.index] = byte(value >> 8)
	b.bytes[b.index+1] = byte(value)

	b.index += 2

	return b
}

// WriteInt32 writes the given int32 value as a 32-bit integer to the builder.
func (b *Builder) WriteInt32(value int32) *Builder {
	b.EnsureWritable(4)

	b.bytes[b.index] = byte(value >> 24)
	b.bytes[b.index+1] = byte(value >> 16)
	b.bytes[b.index+2] = byte(value >> 8)
	b.bytes[b.index+3] = byte(value)

	b.index += 4

	return b
}

// WriteInt64 writes the given int64 value as a 64-bit integer to the builder.
func (b *Builder) WriteInt64(value int64) *Builder {
	b.EnsureWritable(8)

	b.bytes[b.index] = byte(value >> 56)
	b.bytes[b.index+1] = byte(value >> 48)
	b.bytes[b.index+2] = byte(value >> 40)
	b.bytes[b.index+3] = byte(value >> 32)
	b.bytes[b.index+4] = byte(value >> 24)
	b.bytes[b.index+5] = byte(value >> 16)
	b.bytes[b.index+6] = byte(value >> 8)
	b.bytes[b.index+7] = byte(value)

	b.index += 8

	return b
}

// WriteString writes the given string value to the builder without any
// sort of terminator value. It is up to the established user protocol to
// determine what the actual length is.
func (b *Builder) WriteString(value string) *Builder {
	b.EnsureWritable(len(value))

	for character := range value {
		b.WriteByte(byte(character))
	}

	return b
}

// Build constructs a String out of the written bytes.
func (b *Builder) Build() *String {
	return StringWrap(b.bytes[:b.index])
}
