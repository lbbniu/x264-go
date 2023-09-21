package lib

import (
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/moonfdd/x264-go/libx264common"
)

func Init() {
	dir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	libPath := path.Join(dir, ".lbbniu")
	libFile := path.Join(libPath, libName)
	if _, err = os.Stat(libFile); err == nil {
		return
	}
	if err = RestoreAssets(libPath, libName); err != nil {
		panic(err)
	}
	if err = os.Setenv("Path", os.Getenv("Path")+";./lib"); err != nil {
		panic(err)
	}
	libx264common.SetLibx264Path(libFile)
}
