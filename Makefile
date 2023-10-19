GOCMD = go

all: build

build: format
	@mkdir -p bin
	@GO111MODULE=on $(GOCMD) build -o bin/node Node/node.go
	@GO111MODULE=on $(GOCMD) build -o bin/tracker Tracker/tracker.go
	@echo Compiling Done!

vet:
	GO111MODULE=on $(GOCMD) vet Node/node.go
	GO111MODULE=on $(GOCMD) vet Tracker/tracker.go

clean:
	@rm -fr ./bin
	@echo Cleaning Done!

format:
	@$(GOCMD) fmt Node/node.go
	@$(GOCMD) fmt Tracker/tracker.go
	@echo Formatting Done!
