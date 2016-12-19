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
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Sirupsen/logrus"
)

type FileInfo struct {
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	LastModified int64  `json:"lastModified"`
	Type         string `json:"type"`
	Desc         string `json:"desc"`
	Sha256       string `json:"sha256"`
}

func (fi *FileInfo) Metadata() string {
	fiBytes, err := json.Marshal(fi)
	if err != nil {
		logrus.Panic(err)
	}

	return string(fiBytes)
}

func (fi *FileInfo) Raw() string {
	return fmt.Sprintf("%s:%d", fi.Sha256, fi.Size)
}

func FileSha256(fpath string) string {
	file, err := os.Open(fpath)
	if err != nil {
		logrus.Panic(err)
	}

	shaN := sha256.New()
	if _, err := io.Copy(shaN, file); err != nil {
		logrus.Panic(err)
	}

	return fmt.Sprintf("%x", shaN.Sum(nil))
}

func readDir(dir string, ignores []string) {
	if dir == "" {
		dir = getCurrentDirectory()
	}
	dir, _ = filepath.Abs(dir)

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		logrus.Debugf("file path: %s", path)

		if info == nil || err != nil {
			return err
		}

		for _, ignore := range ignores {
			match, _ := regexp.MatchString(ignore, path)
			if match {
				logrus.Debugf("file<%v> match ignore<%v>", path, ignore)
				return nil
			}
		}

		if info.IsDir() {
			return nil
		}

		fi := &FileInfo{
			Name:         info.Name(),
			Size:         info.Size(),
			LastModified: info.ModTime().UnixNano(),
			Type:         "",
			Desc:         "",
			Sha256:       FileSha256(path),
		}

		logrus.Debugf("File: %s", fi)
		wg.Add(1)
		waitFiles <- fi
		return nil
	}); err != nil {
		logrus.Panic(err)
	}
}
