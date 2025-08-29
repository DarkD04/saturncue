
NAME := SaturnCue
SRC := main.go
BUILDDIR := build

ifeq ($(OS), Windows_NT)
	export PATH := winutils:$(PATH)
	NAME := $(NAME).exe
endif
EXE := $(BUILDDIR)/$(NAME)

$(MAKEDIRECTORY):
	mkdir -p $(MAKEDIRECTORY)

.PHONY: all build run clean

all: build 

build: $(MAKEDIRECTORY)
	@echo Building $(SRC)...
	@go build -o $(EXE) $(SRC)
	
run: build
	@$(EXE) $(ARGS)
	
clean:
	@rm -rf $(BUILDDIR)