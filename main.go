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
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	wg        = &sync.WaitGroup{}
	waitFiles = make(chan *FileInfo, 1000)

	initCmd   = kingpin.Command("init", "initial poe file")

	submitCmd = kingpin.Command("submit", "submit files to poe platform")
	ignorePaths = submitCmd.Flag("ignore", "ignore paths, default is empty").Strings()
	hostpoe1     = submitCmd.Flag("host", "host of poe").Default("http://0.0.0.0:9694").String()
	wp          = submitCmd.Flag("wp", "wait period, using duration string, default '1m'").Default("1m").String()

	verifyCmd = kingpin.Command("verify", "verify files which has been submited to the poe platform")
	hostpoe2 = verifyCmd.Flag("host", "host of poe").Default("http://0.0.0.0:9694").String()
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	curDir := getCurrentDirectory()
	switch kingpin.Parse() {
	case initCmd.FullCommand():
		readOrcreatePoeResultFile(curDir)
	case submitCmd.FullCommand():
		poes := readOrcreatePoeResultFile(curDir)
		go readDir(curDir, *ignorePaths)
		go proof2POE(poes)

		time.Sleep(time.Second)
		wg.Wait()

		poes.Save(curDir)
	case verifyCmd.FullCommand():
		poes := readOrcreatePoeResultFile(curDir)
		for key, file := range poes.Files {
			if file.Proofed {
				continue
			}

			rsp := docProofResult(file.Id)
			if rsp.Status == "valid" {
				poes.Files[key].Proofed = true
			}
		}

		poes.Save(curDir)
	}
}
