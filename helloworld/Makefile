# all builds the necessary parts of the two wasm files and the
# parigot system call library. because it is go, it doesn't
# build things it doesn't need to.
all:  hello.p.wasm greeting.p.wasm

#run unit test
.PHONY: test
test:
	GOROOT=/home/parigot/deps/go1.20.4 go1.20.4 test github.com/iansmith/parigot-example/hello-world/greeting

# first WASM file, compiled with GOOS and GOARCH set
hello.p.wasm: main.go
	GOROOT=/home/parigot/deps/go1.21 GOOS=wasip1 GOARCH=wasm go1.21 build -o build/hello.p.wasm ./main.go

# second WASM file, compiled with GOOS and GOARCH set
greeting.p.wasm: greeting/main.go
	GOROOT=/home/parigot/deps/go1.21 GOOS=wasip1 GOARCH=wasm go1.21 build -o build/greeting.p.wasm ./greeting/main.go

# this needs to be called to regenerate the generated files in
# g/ that are derived from the .proto files.  You only need this
# if you change the .proto schema.
.PHONY: generate
generate:
	buf generate

# clean
.PHONY:clean
clean:
	rm -rf build/*
	rm -rf g/*
	

PARIGOT_VERSION=54bda20b

# this needs to be called once when you open the dev container
# to initialize the tools. you will also need to call it if you
# do a "go get github.com/iansmith/parigot" or similar because the version
# affects the tools (runner and protoc-gen-parigot) and
# the lib, syscall.so.
.PHONY:tools
tools:
	go get github.com/iansmith/parigot@${PARIGOT_VERSION}
	go get github.com/iansmith/parigot/command/runner@${PARIGOT_VERSION}
	go get github.com/iansmith/parigot/command/protoc-gen-parigot@${PARIGOT_VERSION}
	go install github.com/iansmith/parigot/command/runner
	go install github.com/iansmith/parigot/command/protoc-gen-parigot
	rm -f /home/parigot/tools/bin/runner
	ln -s /home/parigot/go/bin/runner /home/parigot/tools/bin/runner
	rm -f /home/parigot/tools/bin/protoc-gen-parigot
	ln -s /home/parigot/go/bin/protoc-gen-parigot /home/parigot/tools/bin/protoc-gen-parigot
	go get github.com/iansmith/parigot/api/plugin/syscall/main@${PARIGOT_VERSION}
	go build -o build/syscall.so -buildmode=plugin github.com/iansmith/parigot/api/plugin/syscall/main

	