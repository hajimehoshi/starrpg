include $(GOROOT)/src/Make.inc

TARG:=starrpg
GOFILES:=\
	file_cache.go \
	file_storage.go \
	handler.go \

include $(GOROOT)/src/Make.pkg

main.$O: $(INSTALLFILES) main.go
	$(QUOTED_GOBIN)/$(GC) -o $@ main.go

$(TARG): main.$O
	$(QUOTED_GOBIN)/$(LD) -o $@ main.$O
