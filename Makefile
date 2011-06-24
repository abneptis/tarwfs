include $(GOROOT)/src/Make.inc

TARG=github.com/abneptis/tarwfs

DEPS=\
	onclose_writer/\
	go-fuse/fuse/\

GOFILES=\
	file.go\
        mounted.go\


include $(GOROOT)/src/Make.pkg
