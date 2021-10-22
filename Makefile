## environment
GO_VERSION := "go1.1*.*"

# OS Detection
ifeq ($(OS),Windows_NT)
	# Left here for when/if we will support building on windows
	IS_WINDOWS:=true
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		IS_LINUX:=true
	endif
	ifeq ($(UNAME_S),Darwin) # Mac OS X
		IS_MAC_OS_X:=true
	endif
endif

## 
SHELL=/usr/bin/env bash

%.c: %.y
%.c: %.l


GO_VERSION ?= $(GO_VERSION)

# Note about difference between = and :=
# := means to not evaluate everytime it's expanded.
# = means to evaluate variable everytime it's expanded
# Reference: <http://www.gnu.org/software/make/manual/html_node/Flavors.html#Flavors>
BRANCH = "`git symbolic-ref HEAD | cut -b 12-`"
COMMITCMD = "`git rev-parse HEAD`"
COMMIT := $(shell echo $(COMMITCMD))
DATE := $(shell echo `date +%FT%T%z`)

VERSIONCMD = "`git describe --exact-match --tags $(git log -n1 --pretty='%h')`"

VERSION := $(shell echo $(VERSIONCMD))

ifeq ($(strip $(VERSION)),)

   BRANCHCMD := "`git describe --contains --all HEAD`-`git rev-parse HEAD`"
   VERSION = $(shell echo $(BRANCHCMD))

else

   TAGCMD := "`git describe --exact-match --tags $(git log -n1 --pretty='%h')`-`git rev-parse HEAD`"
   VERSION =  $(shell echo $(TAGCMD))

endif

# Later versions of git supports --count argument to rev-list subcommand
# but not the version on RHEL6 so for now we will just use wc-l
COMMIT_COUNT_CMD = "`git rev-list HEAD | wc -l`"
GIT_BRANCH_CMD = "`git symbolic-ref HEAD | awk -F / '{print $$3}'`"

COMMIT_COUNT := $(shell echo $(COMMIT_COUNT_CMD))
GIT_BRANCH   := $(shell echo $(GIT_BRANCH_CMD))

# if there are any changes not committed, modify the version
CHANGES := $(shell echo `git status --porcelain | wc -l`)
ifneq ($(strip $(CHANGES)), 0)
	VERSION := dirty-build-$(VERSION)
endif


# These are possible to predefine in Makefiles importing this file. in case
# the application needs to add more tags or linker flags
ifdef IS_LINUX
	CLDFLAGS += -L/usr/local/yay/lib
else ifdef IS_MAC_OS_X
	CLDFLAGS += -L /usr/local/lib/gcc/4.9 -L/usr/local/lib
endif

REMOVESYMBOL := -w -s
ifeq (true, $(DEBUG))
	REMOVESYMBOL = ""
endif


CLDFLAGS += -lfreetype -lbz2 -lz -lgomp -lpthread
LDFLAGS += -X $(LDFLAGSPREFIX)/version.version=$(VERSION) -X $(LDFLAGSPREFIX)/version.date=$(DATE) -X $(LDFLAGSPREFIX)/version.commit=$(COMMIT) $(REMOVESYMBOL)

#
BUILD_DIR := $(CURDIR)/build
PACKAGE_DIR := $(CURDIR)/package
BUILD := $(BUILD_DIR)/built
BUILD_TAGS := gm no_development
GOTOOLS_DIR := $(CURDIR)/../../../backend/go-tools

export GOPATH=$(GOTOOLS_DIR):$(CURDIR)/../../../../
export GOBIN := $(BUILD_DIR)/bin
export GOPROXY=https://goproxy.cn,direct

### build step
default: all

FIND_IGNORES= "!" -path "*/.git/*" $(shell git clean -ndX | perl -pe 's/^Would remove (.*)\n/ "!" -path ".\/$$1*" /')
SOURCES := $(shell find . -name '*.go' '!' -path './.*'  $(FIND_IGNORES)) 


.PHONY: all build clean test install pre-build dist dist-clean dist-tar build-dirs


check-go-version:
	@if ! go version | grep " $(GO_VERSION) " >/dev/null; then \
		printf "Wrong go version: "; \
		go version; \
		echo "Requires go version: $(GO_VERSION)"; \
		exit 2; \
	fi

pre-build: check-go-version

clean: build-clean dist-clean

build-clean:
	@echo NOTE: clean
	rm -rf $(CURDIR)/build      || true && \
	rm -rf $(GOPATH)/pkg  || true && \
	rm -f $(BUILD_DIR)/built || true \
	rm -f ${PACKAGE_DIR} || true

vet:
	go vet ./...

### build
all: clean test build

build-dirs:
	@for dir in $(BUILD_DIR)/{bin,lib,app,etc}; do \
		test -d $$dir || mkdir -p $$dir; \
	done


build:
	@make pre-build
	@make build-dirs
	go install -tags="$(BUILD_TAGS)" -ldflags "$(LDFLAGS)" github.com/wtifs/room-booking/cmd/...
	touch $(BUILD)
	mv $(CURDIR)/build/bin/cmd $(CURDIR)/build/bin/room-booking
	@echo -e "\033[32mbuild successfully\033[0m"


