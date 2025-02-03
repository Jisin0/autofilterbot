package fileid

import "encoding/binary"

const (
	// Word represents 4-byte sequence.
	// Values in TL are generally aligned to Word.
	Word = 4

	// If L <= 253, the serialization contains one byte with the value of L,
	// then L bytes of the string followed by 0 to 3 characters containing 0,
	// such that the overall length of the value be divisible by 4,
	// whereupon all of this is interpreted as a sequence of int(L/4)+1 32-bit numbers.
	maxSmallStringLength = 253
	// If L >= 254, the serialization contains byte 254, followed by 3 bytes with
	// the string length L, followed by L bytes of the string, further followed
	// by 0 to 3 null padding bytes.
	firstLongStringByte = 254
)

func encodeString(b []byte, v string) []byte {
	l := len(v)
	if l <= maxSmallStringLength {
		b = append(b, byte(l))
		b = append(b, v...)
		currentLen := l + 1
		b = append(b, make([]byte, nearestPaddedValueLength(currentLen)-currentLen)...)
		return b
	}

	b = append(b, firstLongStringByte, byte(l), byte(l>>8), byte(l>>16))
	b = append(b, v...)
	currentLen := l + 4
	b = append(b, make([]byte, nearestPaddedValueLength(currentLen)-currentLen)...)

	return b
}

func nearestPaddedValueLength(l int) int {
	n := Word * (l / Word)
	if n < l {
		n += Word
	}
	return n
}

// Buffer implements low level binary (de-)serialization for TL.
type Buffer struct {
	Buf []byte
}

// PutUint32 serializes unsigned 32-bit integer.
func (b *Buffer) PutUint32(v uint32) {
	t := make([]byte, Word)
	binary.LittleEndian.PutUint32(t, v)
	b.Buf = append(b.Buf, t...)
}

// PutString serializes bare string.
func (b *Buffer) PutString(s string) {
	b.Buf = encodeString(b.Buf, s)
}

// PutLong serializes v as signed integer.
func (b *Buffer) PutLong(v int64) {
	b.PutUint64(uint64(v))
}

// PutUint64 serializes v as unsigned 64-bit integer.
func (b *Buffer) PutUint64(v uint64) {
	t := make([]byte, Word*2)
	binary.LittleEndian.PutUint64(t, v)
	b.Buf = append(b.Buf, t...)
}

// PutBytes serializes bare byte string.
func (b *Buffer) PutBytes(v []byte) {
	b.Buf = encodeBytes(b.Buf, v)
}

// encodeBytes is same as encodeString, but for bytes.
func encodeBytes(b, v []byte) []byte {
	l := len(v)
	if l <= maxSmallStringLength {
		b = append(b, byte(l))
		b = append(b, v...)
		currentLen := l + 1
		b = append(b, make([]byte, nearestPaddedValueLength(currentLen)-currentLen)...)
		return b
	}

	b = append(b, firstLongStringByte, byte(l), byte(l>>8), byte(l>>16))
	b = append(b, v...)
	currentLen := l + 4
	b = append(b, make([]byte, nearestPaddedValueLength(currentLen)-currentLen)...)

	return b
}
