// Autogenerated jute compiler
// @generated from 'testdata/zookeeper.jute'

package persistence // github.com/go-zookeeper/zk/internal/server/persistence

import (
	"fmt"

	"github.com/go-zookeeper/jute/lib/go/jute"
)

type FileHeader struct {
	Magic   int32 // magic
	Version int32 // version
	Dbid    int64 // dbid
}

func (r *FileHeader) GetMagic() int32 {
	if r != nil {
		return r.Magic
	}
	return 0
}

func (r *FileHeader) GetVersion() int32 {
	if r != nil {
		return r.Version
	}
	return 0
}

func (r *FileHeader) GetDbid() int64 {
	if r != nil {
		return r.Dbid
	}
	return 0
}

func (r *FileHeader) Read(dec jute.Decoder) (err error) {
	if err = dec.ReadStart(); err != nil {
		return err
	}
	r.Magic, err = dec.ReadInt()
	if err != nil {
		return err
	}
	r.Version, err = dec.ReadInt()
	if err != nil {
		return err
	}
	r.Dbid, err = dec.ReadLong()
	if err != nil {
		return err
	}
	if err = dec.ReadEnd(); err != nil {
		return err
	}
	return nil
}

func (r *FileHeader) Write(enc jute.Encoder) error {
	if err := enc.WriteStart(); err != nil {
		return err
	}
	if err := enc.WriteInt(r.Magic); err != nil {
		return err
	}
	if err := enc.WriteInt(r.Version); err != nil {
		return err
	}
	if err := enc.WriteLong(r.Dbid); err != nil {
		return err
	}
	if err := enc.WriteEnd(); err != nil {
		return err
	}
	return nil
}

func (r *FileHeader) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf("FileHeader(%+v)", *r)
}
