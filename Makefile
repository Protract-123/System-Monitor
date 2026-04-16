BINARY    := system-monitor
LDFLAGS   := -X=runtime.godebugDefault=asyncpreemptoff=1
CGO_CXX   := -std=c++17

export CGO_CXXFLAGS := $(CGO_CXX)

.PHONY: run build install clean

run:
	go run -ldflags="$(LDFLAGS)" .

build:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) .

install:
	go install -ldflags="$(LDFLAGS)" .

clean:
	rm -f $(BINARY)
