// Autogenerated jute compiler
// @generated from '/home/pmazzini/repos/jute/testdata/zookeeper.jute'

package proto // github.com/go-zookeeper/zk/internal/proto

import (
	"fmt"

	"github.com/go-zookeeper/jute/lib/go/jute"
	"github.com/go-zookeeper/zk/internal/data"
)

type SetDataResponse struct {
	Stat *data.Stat // stat
}

func (r *SetDataResponse) GetStat() *data.Stat {
	if r != nil && r.Stat != nil {
		return r.Stat
	}
	return nil
}

func (r *SetDataResponse) Read(dec jute.Decoder) (err error) {
	if err = dec.ReadStart(); err != nil {
		return err
	}
	if err = dec.ReadRecord(r.Stat); err != nil {
		return err
	}
	if err = dec.ReadEnd(); err != nil {
		return err
	}
	return nil
}

func (r *SetDataResponse) Write(enc jute.Encoder) error {
	if err := enc.WriteStart(); err != nil {
		return err
	}
	if err := enc.WriteRecord(r.Stat); err != nil {
		return err
	}
	if err := enc.WriteEnd(); err != nil {
		return err
	}
	return nil
}

func (r *SetDataResponse) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SetDataResponse(%+v)", *r)
}
