include $(GOROOT)/src/Make.inc

TARG=tarwfs

DEPS=\
	onclose_writer/\
	go-fuse/fuse/\

GOFILES=\
	file.go\
	main.go\
        mounted.go\

include $(GOROOT)/src/Make.cmd

