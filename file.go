package main

import (
  "github.com/hanwen/go-fuse/fuse"
  "github.com/abneptis/onclose_writer"
  "io"
  . "log"
  "os"
  "bytes"
)

type File struct {
  fuse.DefaultFile
  w io.WriteCloser
}

func NewFile(f func(r io.Reader, l int64)(err os.Error))(*File){
  return &File{w: onclose_writer.New(bytes.NewBuffer(nil), f)}
}

func (self *File)Write(wi *fuse.WriteIn, ib []byte) (written uint32, code fuse.Status){
  Printf("File:write:start:%+v", wi)
  Printf("File:write:start:%d", len(ib))
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
  Printf("File:write:end:%+v:%s\t%d", wi, code, written)
  return
}

func (self *File)Release(){
  Printf("File:release")
  self.w.Close()
  Printf("File:release:complete")
}