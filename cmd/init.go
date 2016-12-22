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
package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/conseweb/graver/file"
	"github.com/conseweb/graver/utils"
)

const (
	metadataFileName = ".poe.json"
)

func StartInitCmd() {
	curdir := utils.GetCurDir()
	readOrcreatePoeResultFile(curdir)
}

func readOrcreatePoeResultFile(path string) *file.DirPoes {
	dirPoes := &file.DirPoes{
		Path:  path,
		Files: make(map[string]*file.PoeFile),
	}
	filePath := filepath.Join(path, metadataFileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		dirPoes.Save(filePath)
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
