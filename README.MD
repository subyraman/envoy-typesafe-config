compiling validate

`protoc -I proto $(find proto/validate -iname "*.proto") --go_out=:$GOPATH/src`