package main

import (
  "github.com/hanwen/go-fuse/fuse"
  "flag"
  "log"
  "os"
)

func main(){
  flag.Parse()
  if flag.NArg() != 2 {
    log.Fatalf("Usage: $0 mountpoint tarfile")
  }
  fp, err := os.OpenFile(flag.Arg(1), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
  if err != nil {
    log.Fatalf("Couldn't open tarfile: %v", err)
  }
  fs := NewWSOWFS(fp)
  state, _, err := fuse.MountFileSystem(flag.Arg(0), fuse.NewLoggingFileSystem(fs), nil)
  log.Printf("State: %v\tErr: %v", state, err)

  if err != nil {
    log.Fatal("Mount fail: %v\n", err)
  }
  state.Loop(true)
  fs.Unmount() // isn't getting called!
  fp.Close()

}
