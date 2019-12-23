// Autogenerated jute compiler
// @generated from '/home/bbennett/src/jute/testdata/test.jute'

package test // com/github/gozookeeper/jute/test

import (
	"fmt"

	"github.com/go-zookeeper/jute/lib/go/jute"
)

type Container struct {
	V []string        // v
	M map[int32]int32 // m
	B *Basic          // b
}

func (r *Container) Read(dec jute.Decoder) (err error) {
	var size int
	if err = dec.ReadStart(); err != nil {
		return err
	}
	size, err = dec.ReadVectorStart()
	if err != nil {
		return err
	}
	r.V = make([]string, size)
	for i := 0; i < size; i++ {
		r.V[i], err = dec.ReadUstring()
		if err != nil {
			return err
		}
	}
	if err = dec.ReadVectorEnd(); err != nil {
		return err
	}
	size, err = dec.ReadMapStart()
	if err != nil {
		return err
	}
	r.M = make(map[int32]int32)
	var k1 int32
	var v1 int32
	for i := 0; i < size; i++ {
		k1, err = dec.ReadInt()
		if err != nil {
			return err
		}
		v1, err = dec.ReadInt()
		if err != nil {
			return err
		}
		r.M[k1] = v1
	}
	if err = dec.ReadMapEnd(); err != nil {
		return err
	}
	if err = dec.ReadRecord(r.B); err != nil {
		return err
	}
	if err = dec.ReadEnd(); err != nil {
		return err
	}
	return nil
}

func (r *Container) Write(enc jute.Encoder) error {
	if err := enc.WriteStart(); err != nil {
		return err
	}
	if err := enc.WriteVectorStart(len(r.V)); err != nil {
		return err
	}
	for _, v := range r.V {
		if err := enc.WriteUstring(v); err != nil {
			return err
		}
	}
	if err := enc.WriteVectorEnd(); err != nil {
		return err
	}
	if err := enc.WriteMapStart(len(r.M)); err != nil {
		return err
	}
	for k, v := range r.M {
		if err := enc.WriteInt(k); err != nil {
			return err
		}
		if err := enc.WriteInt(v); err != nil {
			return err
		}
	}
	if err := enc.WriteMapEnd(); err != nil {
		return err
	}
	if err := enc.WriteRecord(r.B); err != nil {
		return err
	}
	if err := enc.WriteEnd(); err != nil {
		return err
	}
	return nil
}

func (r *Container) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Container(%+v)", *r)
}
