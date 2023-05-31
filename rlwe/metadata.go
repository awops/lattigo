package rlwe

import (
	"fmt"
	"io"

	"github.com/google/go-cmp/cmp"
)

type EncodingDomain int

const (
	SlotsDomain        = EncodingDomain(0)
	CoefficientsDomain = EncodingDomain(1)
)

// MetaData is a struct storing metadata.
type MetaData struct {
	Scale
	EncodingDomain EncodingDomain
	LogSlots       [2]int
	IsNTT          bool
	IsMontgomery   bool
}

// Equal returns true if two MetaData structs are identical.
func (m *MetaData) Equal(other *MetaData) (res bool) {
	res = cmp.Equal(&m.Scale, &other.Scale)
	res = res && m.EncodingDomain == other.EncodingDomain
	res = res && m.LogSlots == other.LogSlots
	res = res && m.IsNTT == other.IsNTT
	res = res && m.IsMontgomery == other.IsMontgomery
	return
}

// Slots returns the number of slots.
func (m *MetaData) Slots() [2]int {
	return [2]int{1 << m.LogSlots[0], 1 << m.LogSlots[1]}
}

// BinarySize returns the size in bytes that the object once marshalled into a binary form.
func (m *MetaData) BinarySize() int {
	return 5 + m.Scale.BinarySize()
}

// MarshalBinary encodes the object into a binary form on a newly allocated slice of bytes.
func (m *MetaData) MarshalBinary() (p []byte, err error) {
	p = make([]byte, m.BinarySize())
	_, err = m.Encode(p)
	return
}

// UnmarshalBinary decodes a slice of bytes generated by
// MarshalBinary or WriteTo on the object.
func (m *MetaData) UnmarshalBinary(p []byte) (err error) {
	_, err = m.Decode(p)
	return
}

// WriteTo writes the object on an io.Writer.
func (m *MetaData) WriteTo(w io.Writer) (int64, error) {
	if p, err := m.MarshalBinary(); err != nil {
		return 0, err
	} else {
		if n, err := w.Write(p); err != nil {
			return int64(n), err
		} else {
			return int64(n), nil
		}
	}
}

func (m *MetaData) ReadFrom(r io.Reader) (int64, error) {
	p := make([]byte, m.BinarySize())
	if n, err := r.Read(p); err != nil {
		return int64(n), nil
	} else {
		return int64(n), m.UnmarshalBinary(p)
	}
}

// Encode encodes the object into a binary form on a preallocated slice of bytes
// and returns the number of bytes written.
func (m *MetaData) Encode(p []byte) (n int, err error) {

	if len(p) < m.BinarySize() {
		return 0, fmt.Errorf("cannot Encode: len(p) is too small")
	}

	if n, err = m.Scale.Encode(p[n:]); err != nil {
		return 0, err
	}

	p[n] = uint8(m.EncodingDomain)
	n++

	p[n] = uint8(m.LogSlots[0])
	n++

	p[n] = uint8(m.LogSlots[1])
	n++

	if m.IsNTT {
		p[n] = 1
	}

	n++

	if m.IsMontgomery {
		p[n] = 1
	}

	n++

	return
}

// Decode decodes a slice of bytes generated by Encode
// on the object and returns the number of bytes read.
func (m *MetaData) Decode(p []byte) (n int, err error) {

	if len(p) < m.BinarySize() {
		return 0, fmt.Errorf("canoot Decode: len(p) is too small")
	}

	if n, err = m.Scale.Decode(p[n:]); err != nil {
		return
	}

	m.EncodingDomain = EncodingDomain(p[n])
	n++

	m.LogSlots[0] = int(int8(p[n]))
	n++

	m.LogSlots[1] = int(int8(p[n]))
	n++

	m.IsNTT = p[n] == 1
	n++

	m.IsMontgomery = p[n] == 1
	n++

	return
}
