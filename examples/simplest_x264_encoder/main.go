// https://github.com/leixiaohua1020/simplest_encoder/blob/master/simplest_x264_encoder/simplest_x264_encoder.cpp
package main

import (
	"fmt"
	"os"
	"path"
	"unsafe"

	"github.com/moonfdd/ffmpeg-go/ffcommon"
	"github.com/moonfdd/x264-go/lib"
	"github.com/moonfdd/x264-go/libx264"
)

func main0() ffcommon.FInt {

	var ret ffcommon.FInt
	var ySize ffcommon.FInt
	var i, j ffcommon.FInt

	fpSrc, _ := os.Open("./resources/cuc_ieschool_640x360_yuv420p.yuv")
	fpDstFile := "./out/cuc_ieschool_640x360_yuv420p.h264"
	_ = os.MkdirAll(path.Dir(fpDstFile), os.ModeDir)
	fpDst, _ := os.Create(fpDstFile)

	//Encode 50 frame
	//if set 0, encode all frame
	var frameNum ffcommon.FInt = 0
	var csp ffcommon.FInt = libx264.X264_CSP_I420
	var width, height ffcommon.FInt = 640, 360

	var iNal ffcommon.FInt = 0
	var pNals *libx264.X264NalT
	var pHandle *libx264.X264T
	picIn := new(libx264.X264PictureT)
	picOut := new(libx264.X264PictureT)
	pParam := new(libx264.X264ParamT)

	//Check
	if fpSrc == nil || fpDst == nil {
		fmt.Printf("Error open files.\n")
		return -1
	}

	pParam.X264ParamDefault()
	pParam.IWidth = width
	pParam.IHeight = height
	/*
		//Param
		pParam->i_log_level  = X264_LOG_DEBUG;
		pParam->i_threads  = X264_SYNC_LOOKAHEAD_AUTO;
		pParam->i_frame_total = 0;
		pParam->i_keyint_max = 10;
		pParam->i_bframe  = 5;
		pParam->b_open_gop  = 0;
		pParam->i_bframe_pyramid = 0;
		pParam->rc.i_qp_constant=0;
		pParam->rc.i_qp_max=0;
		pParam->rc.i_qp_min=0;
		pParam->i_bframe_adaptive = X264_B_ADAPT_TRELLIS;
		pParam->i_fps_den  = 1;
		pParam->i_fps_num  = 25;
		pParam->i_timebase_den = pParam->i_fps_num;
		pParam->i_timebase_num = pParam->i_fps_den;
	*/
	pParam.ICsp = csp
	pParam.X264ParamApplyProfile(libx264.X264ProfileNames[5])

	pHandle = pParam.X264EncoderOpen164()

	picOut.X264PictureInit()
	picIn.X264PictureAlloc(csp, pParam.IWidth, pParam.IHeight)

	//ret = x264_encoder_headers(pHandle, &pNals, &iNal);

	ySize = pParam.IWidth * pParam.IHeight
	//detect frame number
	fi, _ := fpSrc.Stat()
	switch csp {
	case libx264.X264_CSP_I444:
		frameNum = int32(fi.Size()) / (ySize * 3)
	case libx264.X264_CSP_I420:
		frameNum = int32(fi.Size()) / (ySize * 3 / 2)
	default:
		fmt.Printf("Colorspace Not Support.\n")
		return -1
	}

	//Loop to Encode
	for i = 0; i < frameNum; i++ {
		switch csp {
		case libx264.X264_CSP_I444:
			fpSrc.Read(ffcommon.ByteSliceFromByteP(picIn.Img.Plane[0], int(ySize))) //Y
			fpSrc.Read(ffcommon.ByteSliceFromByteP(picIn.Img.Plane[1], int(ySize))) //U
			fpSrc.Read(ffcommon.ByteSliceFromByteP(picIn.Img.Plane[2], int(ySize))) //V

		case libx264.X264_CSP_I420:
			fpSrc.Read(ffcommon.ByteSliceFromByteP(picIn.Img.Plane[0], int(ySize)))   //Y
			fpSrc.Read(ffcommon.ByteSliceFromByteP(picIn.Img.Plane[1], int(ySize/4))) //U
			fpSrc.Read(ffcommon.ByteSliceFromByteP(picIn.Img.Plane[2], int(ySize/4))) //V

		default:
			fmt.Printf("Colorspace Not Support.\n")
			return -1
		}
		picIn.IPts = int64(i)

		ret = pHandle.X264EncoderEncode(&pNals, &iNal, picIn, picOut)
		if ret < 0 {
			fmt.Printf("Error.\n")
			return -1
		}

		fmt.Printf("Succeed encode frame: %5d\n", i)

		for j = 0; j < iNal; j++ {
			a := unsafe.Sizeof(libx264.X264NalT{})
			pNal := (*libx264.X264NalT)(unsafe.Pointer(uintptr(unsafe.Pointer(pNals)) + uintptr(a*uintptr(j))))
			fpDst.Write(ffcommon.ByteSliceFromByteP(pNal.PPayload, int(pNal.IPayload)))
		}
	}
	i = 0
	//flush encoder
	for {
		ret = pHandle.X264EncoderEncode(&pNals, &iNal, nil, picOut)
		if ret == 0 {
			break
		}
		fmt.Printf("Flush 1 frame.\n")
		for j = 0; j < iNal; j++ {
			a := unsafe.Sizeof(libx264.X264NalT{})
			pNal := (*libx264.X264NalT)(unsafe.Pointer(uintptr(unsafe.Pointer(pNals)) + uintptr(a*uintptr(j))))
			fpDst.Write(ffcommon.ByteSliceFromByteP(pNal.PPayload, int(pNal.IPayload)))
		}
		i++
	}
	picIn.X264PictureClean()
	pHandle.X264EncoderClose()
	pHandle = nil

	fpSrc.Close()
	fpDst.Close()

	fmt.Printf("\nffplay %s\n", fpDstFile)
	return 0
}

func main() {
	lib.Init()
	fmt.Println(libx264.X264_POINTVER)
	main0()
}
