BINARY    := system-monitor
APP_NAME  := System-Monitor.app
VERSION   ?= 0.0.0-dev
LDFLAGS   := -X=runtime.godebugDefault=asyncpreemptoff=1
CGO_CXX   := -std=c++17

export CGO_CXXFLAGS := $(CGO_CXX)

.PHONY: run build install app clean

run:
	go run -ldflags="$(LDFLAGS)" .

build:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) .

install:
	go install -ldflags="$(LDFLAGS)" .

app: build
	rm -rf $(APP_NAME)
	mkdir -p $(APP_NAME)/Contents/MacOS
	cp .github/Info.plist $(APP_NAME)/Contents/Info.plist
	cp $(BINARY) $(APP_NAME)/Contents/MacOS/$(BINARY)
	plutil -replace CFBundleShortVersionString -string "$(VERSION)" $(APP_NAME)/Contents/Info.plist
	plutil -replace CFBundleVersion -string "$(VERSION)" $(APP_NAME)/Contents/Info.plist
	macdeployqt $(APP_NAME)

clean:
	rm -rf $(BINARY) $(APP_NAME)
