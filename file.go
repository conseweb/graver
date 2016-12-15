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
	"strings"
	"fmt"
	"crypto/sha1"
	"os"
	"io"
	"path/filepath"
	"regexp"

	"github.com/Sirupsen/logrus"
)

type FileInfo struct {
	Name string
	Size int64
	Hash string
}

func (this *FileInfo) String() string {
	return strings.Join([]string{
		this.Name,
		this.Hash,
		fmt.Sprintf("%v", this.Size),
	}, ",")
}

func FileSha1(fpath string) string {
	file, err := os.Open(fpath)
	if err != nil {
		return ""
	}

	sha1N := sha1.New()
	if _, err := io.Copy(sha1N, file); err != nil {
		return ""
	}

	return fmt.Sprintf("%x", sha1N.Sum(nil))
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
			Name: info.Name(),
			Size: info.Size(),
			Hash: FileSha1(path),
		}

		logrus.Debugf("File: %s", fi)
		wg.Add(1)
		waitFiles <- fi
		return nil
	}); err != nil {
		logrus.Panic(err)
	}
}