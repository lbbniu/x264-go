package lib

// go:generate go install github.com/go-bindata/go-bindata@latest
//go:generate go-bindata -tags=windows -nomemcopy -o=libx264_windows.go -pkg=lib libx264-164.dll
//go:generate go-bindata -tags=darwin -nomemcopy -o=libx264_darwin.go -pkg=lib libx264.164.dylib
//go:generate go-bindata -tags=linux -nomemcopy -o=libx264_linux.go -pkg=lib libx264.so.164
