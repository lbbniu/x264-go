package lib

//go:generate go-bindata -tags=windows -nomemcopy -o=libx264_windows.go -pkg=lib libx264-164.dll
//go:generate go-bindata -tags=darwin -nomemcopy -o=libx264_darwin.go -pkg=lib libx264.164.dylib
