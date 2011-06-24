package tarwfs

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/abneptis/onclose_writer"
	"io"
	"os"
	"bytes"
)

type file struct {
	fuse.DefaultFile
	w io.WriteCloser
}

func newFile(f func(r io.Reader, l int64) (err os.Error)) *file {
	return &file{w: onclose_writer.New(bytes.NewBuffer(nil), f)}
}

func (self *file) Write(wi *fuse.WriteIn, ib []byte) (written uint32, code fuse.Status) {
	//log.Printf("File:write:start:%+v", wi)
	//log.Printf("File:write:start:%d", len(ib))
	// we need the result to fit into a uint32.
	if len(ib) > 2147483647 {
		ib = ib[0:2147483647]
	}
	n, err := self.w.Write(ib)
	written = uint32(n)
	if err == nil {
		code = fuse.OK
	} else {
		code = fuse.EIO
	}
	//log.Printf("File:write:end:%+v:%s\t%d", wi, code, written)
	return
}

func (self *file) Release() {
	//log.Printf("File:release")
	self.w.Close()
	//log.Printf("File:release:complete")
}
