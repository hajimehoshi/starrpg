include $(GOROOT)/src/Make.inc

EXEC:=starrpg
TARG:=hajimehoshi/starrpg
GOFILES:=\
	dummy_storage.go \
	file_cache.go \
	handler.go \
	map_storage.go \

include $(GOROOT)/src/Make.pkg

main.$O: $(INSTALLFILES) main.go
	$(QUOTED_GOBIN)/$(GC) -o $@ main.go

$(EXEC): main.$O
	$(QUOTED_GOBIN)/$(LD) -o $@ main.$O
