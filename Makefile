.phony: clean

PROJECT = ld33
SOURCES = $(wildcard src/*.go)
ASSETS  = $(wildcard src/resources/*)
VERSION = $(shell cat VARS | grep TDP_VERSION | sed s/export\ TDP_VERSION=//g)
REPLACE = s/9\.9\.9/$(VERSION)/g
UNAME   = $(shell uname)

ifeq ($(UNAME), Darwin)
	PLATFORM := osx
else ifeq ($(UNAME), MINGW32_NT-6.2)
	PLATFORM := win
else
	PLATFORM := nix
endif

include Makefile.$(PLATFORM)

clean:
	rm -rf build
