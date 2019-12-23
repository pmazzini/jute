// Autogenerated jute compiler
// @generated from '/home/bbennett/src/jute/testdata/zookeeper.jute'

package txn // github.com/go-zookeeper/zk/internal/txn

import (
	"fmt"

	"github.com/go-zookeeper/jute/lib/go/jute"
	"github.com/go-zookeeper/zk/internal/data"
)

type SetACLTxn struct {
	Path    *string     // path
	Acl     []*data.ACL // acl
	Version int32       // version
}

func (r *SetACLTxn) Read(dec jute.Decoder) (err error) {
	var size int
	if err = dec.ReadStart(); err != nil {
		return err
	}
	r.Path, err = dec.ReadUstring()
	if err != nil {
		return err
	}
	size, err = dec.ReadVectorStart()
	if err != nil {
		return err
	}
	if size < 0 {
		r.Acl = nil
	} else {
		r.Acl = make([]*data.ACL, size)
		for i := 0; i < size; i++ {
			if err = dec.ReadRecord(r.Acl[i]); err != nil {
				return err
			}
		}
	}
	if err = dec.ReadVectorEnd(); err != nil {
		return err
	}
	r.Version, err = dec.ReadInt()
	if err != nil {
		return err
	}
	if err = dec.ReadEnd(); err != nil {
		return err
	}
	return nil
}

func (r *SetACLTxn) Write(enc jute.Encoder) error {
	if err := enc.WriteStart(); err != nil {
		return err
	}
	if err := enc.WriteUstring(r.Path); err != nil {
		return err
	}
	if err := enc.WriteVectorStart(len(r.Acl), r.Acl == nil); err != nil {
		return err
	}
	for _, v := range r.Acl {
		if err := enc.WriteRecord(v); err != nil {
			return err
		}
	}
	if err := enc.WriteVectorEnd(); err != nil {
		return err
	}
	if err := enc.WriteInt(r.Version); err != nil {
		return err
	}
	if err := enc.WriteEnd(); err != nil {
		return err
	}
	return nil
}

func (r *SetACLTxn) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SetACLTxn(%+v)", *r)
}
