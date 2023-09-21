package lib

import (
	"os"
	"path"
	"sync"

	"github.com/mitchellh/go-homedir"
	"github.com/moonfdd/x264-go/libx264common"
)

var once sync.Once

func Init() {
	once.Do(func() {
		dir, err := homedir.Dir()
		if err != nil {
			panic(err)
		}
		libPath := path.Join(dir, ".lbbniu")
		if err = RestoreAssets(libPath, libName); err != nil {
			panic(err)
		}
		if err = os.Setenv("Path", os.Getenv("Path")+";./lib"); err != nil {
			panic(err)
		}
		libx264common.SetLibx264Path(path.Join(libPath, libName))
	})
}
