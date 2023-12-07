GOCMD = /home/core/.asdf/installs/golang/1.21.2/go/bin/go

all: build

build: format
	@mkdir -p bin
	@GO111MODULE=on $(GOCMD) build -o bin/node Node/node.go Node/menu.go
	@GO111MODULE=on $(GOCMD) build Protocols/centralProtocol.go
	@GO111MODULE=on $(GOCMD) build Protocols/taxiProtocol.go
	@GO111MODULE=on $(GOCMD) build util/util.go
	@GO111MODULE=on $(GOCMD) build -o bin/tracker Tracker/tracker.go
	@echo Compiling Done!

vet:
	GO111MODULE=on $(GOCMD) vet Node/node.go
	GO111MODULE=on $(GOCMD) vet Tracker/tracker.go

clean:
	@rm -fr ./bin
	@echo Cleaning Done!

lint:
	golint Tracker/tracker.go
	@echo ---
	golint Node/node.go
	@echo Linting Done!
  
format:
	@$(GOCMD) fmt Node/node.go
	@$(GOCMD) fmt Tracker/tracker.go
	@echo Formatting Done!
