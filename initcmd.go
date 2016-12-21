/*
Copyright Mojing Inc. 2016 All Rights Reserved.
Written by mint.zhao.chiu@gmail.com. github.com: https://www.github.com/mintzhao

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
)

const (
	fileName = ".poe.json"
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

	filePath := filepath.Join(path, fileName)
	if err := ioutil.WriteFile(filePath, dpBytes, 0666); err != nil {
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

func readOrcreatePoeResultFile(path string) *DirPoes {
	dirPoes := &DirPoes{
		Path:  path,
		Files: make(map[string]*PoeFile),
	}
	filePath := filepath.Join(path, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		dpBytes, err := json.MarshalIndent(dirPoes, "", "   ")
		if err != nil {
			logrus.Panic(err)
		}

		ioutil.WriteFile(filePath, dpBytes, 0666)
		return dirPoes
	}

	ibytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Panic(err)
	}

	if err := json.Unmarshal(ibytes, dirPoes); err != nil {
		logrus.Panic(err)
	}

	if path != dirPoes.Path {
		logrus.Panic("please don't copy file from other dir, just init it")
	}

	return dirPoes
}
