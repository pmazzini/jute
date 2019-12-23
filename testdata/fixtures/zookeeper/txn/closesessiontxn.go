// Autogenerated jute compiler
// @generated from '/home/bbennett/src/jute/testdata/zookeeper.jute'

package txn // github.com/go-zookeeper/zk/internal/txn

import (
	"fmt"

	"github.com/go-zookeeper/jute/lib/go/jute"
)

type CloseSessionTxn struct {
	Paths2Delete []*string // paths2Delete
}

func (r *CloseSessionTxn) Read(dec jute.Decoder) (err error) {
	var size int
	if err = dec.ReadStart(); err != nil {
		return err
	}
	size, err = dec.ReadVectorStart()
	if err != nil {
		return err
	}
	if size < 0 {
		r.Paths2Delete = nil
	} else {
		r.Paths2Delete = make([]*string, size)
		for i := 0; i < size; i++ {
			r.Paths2Delete[i], err = dec.ReadUstring()
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

func (r *CloseSessionTxn) Write(enc jute.Encoder) error {
	if err := enc.WriteStart(); err != nil {
		return err
	}
	if err := enc.WriteVectorStart(len(r.Paths2Delete), r.Paths2Delete == nil); err != nil {
		return err
	}
	for _, v := range r.Paths2Delete {
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

func (r *CloseSessionTxn) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf("CloseSessionTxn(%+v)", *r)
}
