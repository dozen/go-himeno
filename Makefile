SSMALL_MIMAX=33
SSMALL_MJMAX=33
SSMALL_MKMAX=65

SMALL_MIMAX=65
SMALL_MJMAX=65
SMALL_MKMAX=129

MIDDLE_MIMAX=129
MIDDLE_MJMAX=129
MIDDLE_MKMAX=257

LARGE_MIMAX=257
LARGE_MJMAX=257
LARGE_MKMAX=513

ELARGE_MIMAX=513
ELARGE_MJMAX=513
ELARGE_MKMAX=1025

VAR_SIZE=$(if $(SIZE),$(SIZE),LARGE)

MIMAX=$($(VAR_SIZE)_MIMAX)
MJMAX=$($(VAR_SIZE)_MJMAX)
MKMAX=$($(VAR_SIZE)_MKMAX)

DIST=himeno.go
SRC=_$(DIST)

.PHONY: build
build:
	sed -E \
	-e "s/(const *MIMAX).*$$/\1 = $(MIMAX)/" \
	-e "s/(const *MJMAX).*$$/\1 = $(MJMAX)/" \
	-e "s/(const *MKMAX).*$$/\1 = $(MKMAX)/" \
	$(SRC) > $(DIST)
	go build
	#rm $(DIST)

.PHONY: mng
mng:
	go build cmd/himeno-mng/himeno-mng.go

.PHONY: test
test:
	sed -E \
	-e "s/(const *MIMAX).*$$/\1 = $(MIMAX)/" \
	-e "s/(const *MJMAX).*$$/\1 = $(MJMAX)/" \
	-e "s/(const *MKMAX).*$$/\1 = $(MKMAX)/" \
	$(SRC) > $(DIST)
	go test ./...
	rm $(DIST)

.PHONY: fmt
fmt:
	mv $(SRC) $(DIST)
	go fmt ./...
	mv $(DIST) $(SRC)

.PHONY: proto
proto:
	protoc -Imanager/proto/ manager/proto/proto.proto --go_out=plugins=grpc:manager/proto
