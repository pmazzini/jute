// Autogenerated jute compiler
// @generated from 'testdata/zookeeper.jute'

package txn // github.com/go-zookeeper/zk/internal/txn

import (
	"fmt"

	"github.com/go-zookeeper/jute/lib/go/jute"
)

type DeleteTxn struct {
	Path string // path
}

func (r *DeleteTxn) GetPath() string {
	if r != nil {
		return r.Path
	}
	return ""
}

func (r *DeleteTxn) Read(dec jute.Decoder) (err error) {
	if err = dec.ReadStart(); err != nil {
		return err
	}
	r.Path, err = dec.ReadString()
	if err != nil {
		return err
	}
	if err = dec.ReadEnd(); err != nil {
		return err
	}
	return nil
}

func (r *DeleteTxn) Write(enc jute.Encoder) error {
	if err := enc.WriteStart(); err != nil {
		return err
	}
	if err := enc.WriteString(r.Path); err != nil {
		return err
	}
	if err := enc.WriteEnd(); err != nil {
		return err
	}
	return nil
}

func (r *DeleteTxn) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf("DeleteTxn(%+v)", *r)
}
