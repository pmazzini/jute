// Autogenerated jute compiler
// @generated from '/home/bbennett/src/jute/testdata/zookeeper.jute'

package proto // github.com/go-zookeeper/zk/internal/proto

import (
	"fmt"

	"github.com/go-zookeeper/jute/lib/go/jute"
)

type GetChildrenResponse struct {
	Children []string // children
}

func (r *GetChildrenResponse) Read(dec jute.Decoder) (err error) {
	var size int
	if err = dec.ReadStart(); err != nil {
		return err
	}
	size, err = dec.ReadVectorStart()
	if err != nil {
		return err
	}
	if size < 0 {
		r.Children = nil
	} else {
		r.Children = make([]string, size)
		for i := 0; i < size; i++ {
			r.Children[i], err = dec.ReadUstring()
			if err != nil {
				return err
			}
		}
	}
	if err = dec.ReadVectorEnd(); err != nil {
		return err
	}
	if err = dec.ReadEnd(); err != nil {
		return err
	}
	return nil
}

func (r *GetChildrenResponse) Write(enc jute.Encoder) error {
	if err := enc.WriteStart(); err != nil {
		return err
	}
	if err := enc.WriteVectorStart(len(r.Children), r.Children == nil); err != nil {
		return err
	}
	for _, v := range r.Children {
		if err := enc.WriteUstring(v); err != nil {
			return err
		}
	}
	if err := enc.WriteVectorEnd(); err != nil {
		return err
	}
	if err := enc.WriteEnd(); err != nil {
		return err
	}
	return nil
}

func (r *GetChildrenResponse) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetChildrenResponse(%+v)", *r)
}