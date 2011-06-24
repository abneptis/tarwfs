package main

import (
  . "log"
  "io"
  "os"
  "github.com/hanwen/go-fuse/fuse"
  "path"
  "sync"
  "syscall"
  "archive/tar"
  "time"
)

type WSOWFS struct {
  fuse.DefaultFileSystem
  Files       map[string]*os.FileInfo
  IsDir   map[string]bool
  lock    *sync.Mutex 
  w       *tar.Writer
}

func NewWSOWFS(w io.Writer)(wfs *WSOWFS){
  wfs = &WSOWFS{
    Files: map[string]*os.FileInfo{
      "": &os.FileInfo{Mode: 0755 | syscall.S_IFDIR},
    },
    IsDir: map[string]bool{"":true},
    lock: &sync.Mutex{},
    w: tar.NewWriter(w),
  }
  
  return
} 

func (self *WSOWFS)GetAttr(name string)(fi *os.FileInfo, eno fuse.Status){
  Printf("GetAttr:start:%s", name)
  self.lock.Lock()
  fi, exists := self.Files[name]
  self.lock.Unlock() 
  if exists {
    eno = fuse.OK
  } else {
    eno = fuse.ENOENT
  }
  Printf("GetAttr:finish:%s:%v\t%+v", name, eno, fi)

  return
}

func (self *WSOWFS) OpenDir(name string) (c chan fuse.DirEntry, eno fuse.Status) {
  Printf("GetAttr:start:%s", name)
  self.lock.Lock()
  is_dir, exists := self.IsDir[name]
  self.lock.Unlock()
  if exists && is_dir {
    c =  make(chan fuse.DirEntry,16)
    eno = fuse.OK
    go func(){
      self.lock.Lock()
      // send ourselves first.
      dent := fuse.DirEntry{Name: ".", Mode: self.Files[name].Mode}
      Printf("Opendir(%s) -> %s {%+v}", name, ".", dent)
      self.lock.Unlock() 
      c <- dent
      self.lock.Lock()
      for k,v := range(self.Files){
        self.lock.Unlock() 
        k_d, k_n := path.Split(k)
        if path.Join(name,k_n) == path.Join(k_d, k_n) && path.Join(name,k_n) != name {
          dent := fuse.DirEntry{Name: k_n, Mode: v.Mode}
          Printf("Opendir(%s) ->  %s|%s", name, k_d, k_n)
          c <- dent
        } else {
          Printf("%s is not the parent of %s|%s", name, k_d, k_n)
        }
        self.lock.Lock()
      }

      self.lock.Unlock() 
      close(c)

    }()
  } else {
    if exists { 
      eno = fuse.ENOTDIR
    } else {
      eno = fuse.ENOENT
    }
  }
  return
}

func (self *WSOWFS)Mkdir(name string, mode uint32)( eno fuse.Status){
  Printf("Open:Mkdir:%s:%x", name,mode)
  self.lock.Lock()
  _, exists := self.IsDir[name]
  if exists {
    self.lock.Unlock()
    eno = fuse.EPERM
  } else {
    self.Files[name] = &os.FileInfo{
      Mode: mode | syscall.S_IFDIR,
    }
    self.IsDir[name] = true
    self.lock.Unlock()
    self.lock.Lock()
    now := time.Seconds()
    err := self.w.WriteHeader(&tar.Header{
       Typeflag: tar.TypeDir,
       Name: name,
       Mode: int64(mode),
       Size: 0,
       Ctime: now,
       Mtime: now, 
       Atime: now,
      })
    self.lock.Unlock()
    if err == nil {
      eno = fuse.OK
    } else {
      eno = fuse.EIO
    }
  }
  Printf("Open:Mkdir:%s:%x:%x:done\t%s", name,mode, eno)
  return
}

func (self *WSOWFS)Create(name string, flags, mode uint32) (file fuse.File, eno fuse.Status){
  Printf("Create:%s", name)
  self.lock.Lock()
  _, exists := self.IsDir[name]
  if ! exists {
    self.Files[name] = &os.FileInfo{ Mode: mode | syscall.S_IFREG, }
    self.IsDir[name] = false
    self.lock.Unlock()
    start := time.Seconds()
    file = NewFile(func(r io.Reader, rlen int64)(err os.Error){
      self.lock.Lock()
      err = self.w.WriteHeader(&tar.Header{
       Typeflag: tar.TypeReg,
       Name: name,
       Mode: int64(mode),
       Size: rlen,
       Ctime: start,
       Mtime: time.Seconds(),
       Atime: time.Seconds(),
      })
      if err == nil {
        _, err = io.Copyn(self.w, r, rlen)
      }
      self.lock.Unlock()
      return
    })
    eno = fuse.OK
  } else {
    self.lock.Unlock()
    eno = fuse.EINVAL
  }
  return
}



func (self *WSOWFS)Unmount(){
  // In case someone else is finishing up still, wait for them to unmount the lock.
  self.lock.Lock()
  self.w.Close()
  self.lock.Unlock()
}


