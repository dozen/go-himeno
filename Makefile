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

MIMAX=$($(SIZE)_MIMAX)
MJMAX=$($(SIZE)_MJMAX)
MKMAX=$($(SIZE)_MKMAX)

DIST=himeno_test.go
SRC=_$(DIST)

.PHONY: build
build:
	sed -E \
	-e "s/(const *MIMAX).*$$/\1 = $(MIMAX)/" \
	-e "s/(const *MJMAX).*$$/\1 = $(MJMAX)/" \
	-e "s/(const *MKMAX).*$$/\1 = $(MKMAX)/" \
	$(SRC) > $(DIST)
	make test

.PHONY: test
test:
	go test -bench . -benchmem

.PHONY: fmt
fmt:
	mv $(SRC) $(DIST)
	go fmt $(DIST)
	mv $(DIST) $(SRC)
