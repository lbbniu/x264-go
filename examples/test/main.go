package main

import (
	"fmt"

	"github.com/moonfdd/x264-go/lib"
	"github.com/moonfdd/x264-go/libx264"
)

func main() {
	lib.Init()
	fmt.Println(libx264.X264_POINTVER)
	p := new(libx264.X264ParamT)
	p.X264ParamDefault()
	p.X264ParamDefaultPreset("veryfast", "zerolatency")
}
