package file

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
)

type DirPoes struct {
	Path  string              `json:"path"`
	Count int                 `json:"count"`
	Files map[string]*PoeFile `json:"files"`
}

type PoeFile struct {
	Fi      *FileInfo `json:"fi"`
	Id      string    `json:"id"`
	Proofed bool      `json:"proofed"`
}

func (dp *DirPoes) Save(path string) {
	if path == "" {
		return
	}

	dpBytes, err := json.MarshalIndent(dp, "", "    ")
	if err != nil {
		logrus.Panic(err)
	}

	if err := ioutil.WriteFile(path, dpBytes, 0666); err != nil {
		logrus.Panic(err)
	}
}

func (dp *DirPoes) PutFile(fi *FileInfo) bool {
	if _, ok := dp.Files[fi.Sha256]; ok {
		return false
	}

	dp.Files[fi.Sha256] = &PoeFile{
		Fi: fi,
	}
	dp.Count++
	return true
}
