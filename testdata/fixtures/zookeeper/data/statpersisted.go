// Autogenerated jute compiler
// @generated from '/home/bbennett/src/jute/testdata/zookeeper.jute'

package data // github.com/go-zookeeper/zk/internal/data

import (
	"fmt"

	"github.com/go-zookeeper/jute/lib/go/jute"
)

type StatPersisted struct {
	Czxid          int64 // czxid
	Mzxid          int64 // mzxid
	Ctime          int64 // ctime
	Mtime          int64 // mtime
	Version        int32 // version
	Cversion       int32 // cversion
	Aversion       int32 // aversion
	EphemeralOwner int64 // ephemeralOwner
	Pzxid          int64 // pzxid
}

func (r *StatPersisted) Read(dec jute.Decoder) (err error) {
	if err = dec.ReadStart(); err != nil {
		return err
	}
	r.Czxid, err = dec.ReadLong()
	if err != nil {
		return err
	}
	r.Mzxid, err = dec.ReadLong()
	if err != nil {
		return err
	}
	r.Ctime, err = dec.ReadLong()
	if err != nil {
		return err
	}
	r.Mtime, err = dec.ReadLong()
	if err != nil {
		return err
	}
	r.Version, err = dec.ReadInt()
	if err != nil {
		return err
	}
	r.Cversion, err = dec.ReadInt()
	if err != nil {
		return err
	}
	r.Aversion, err = dec.ReadInt()
	if err != nil {
		return err
	}
	r.EphemeralOwner, err = dec.ReadLong()
	if err != nil {
		return err
	}
	r.Pzxid, err = dec.ReadLong()
	if err != nil {
		return err
	}
	if err = dec.ReadEnd(); err != nil {
		return err
	}
	return nil
}

func (r *StatPersisted) Write(enc jute.Encoder) error {
	if err := enc.WriteStart(); err != nil {
		return err
	}
	if err := enc.WriteLong(r.Czxid); err != nil {
		return err
	}
	if err := enc.WriteLong(r.Mzxid); err != nil {
		return err
	}
	if err := enc.WriteLong(r.Ctime); err != nil {
		return err
	}
	if err := enc.WriteLong(r.Mtime); err != nil {
		return err
	}
	if err := enc.WriteInt(r.Version); err != nil {
		return err
	}
	if err := enc.WriteInt(r.Cversion); err != nil {
		return err
	}
	if err := enc.WriteInt(r.Aversion); err != nil {
		return err
	}
	if err := enc.WriteLong(r.EphemeralOwner); err != nil {
		return err
	}
	if err := enc.WriteLong(r.Pzxid); err != nil {
		return err
	}
	if err := enc.WriteEnd(); err != nil {
		return err
	}
	return nil
}

func (r *StatPersisted) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf("StatPersisted(%+v)", *r)
}