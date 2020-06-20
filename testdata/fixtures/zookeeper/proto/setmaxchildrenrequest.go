// Autogenerated jute compiler
// @generated from '/home/pmazzini/repos/jute/testdata/zookeeper.jute'

package proto // github.com/go-zookeeper/zk/internal/proto

import (
	"fmt"

	"github.com/go-zookeeper/jute/lib/go/jute"
)

type SetMaxChildrenRequest struct {
	Path *string // path
	Max  int32   // max
}

func (r *SetMaxChildrenRequest) GetPath() string {
	if r != nil && r.Path != nil {
		return *r.Path
	}
	return ""
}

func (r *SetMaxChildrenRequest) GetMax() int32 {
	if r != nil {
		return r.Max
	}
	return 0
}

func (r *SetMaxChildrenRequest) Read(dec jute.Decoder) (err error) {
	if err = dec.ReadStart(); err != nil {
		return err
	}
	r.Path, err = dec.ReadString()
	if err != nil {
		return err
	}
	r.Max, err = dec.ReadInt()
	if err != nil {
		return err
	}
	if err = dec.ReadEnd(); err != nil {
		return err
	}
	return nil
}

func (r *SetMaxChildrenRequest) Write(enc jute.Encoder) error {
	if err := enc.WriteStart(); err != nil {
		return err
	}
	if err := enc.WriteString(r.Path); err != nil {
		return err
	}
	if err := enc.WriteInt(r.Max); err != nil {
		return err
	}
	if err := enc.WriteEnd(); err != nil {
		return err
	}
	return nil
}

func (r *SetMaxChildrenRequest) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SetMaxChildrenRequest(%+v)", *r)
}
