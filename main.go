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
	rootpath    = kingpin.Flag("path", "target path of proof existence").Default(getCurrentDirectory()).String()
	ignorePaths = kingpin.Flag("ignore", "ignore paths, default is empty").Strings()
	hostpoe     = kingpin.Flag("host", "host of poe").Default("http://0.0.0.0:9694").String()
	wp          = kingpin.Flag("wp", "wait period, using duration string, default '1m'").Default("1m").String()

	wg        = &sync.WaitGroup{}
	waitFiles = make(chan *FileInfo, 1000)
)

func main() {
	kingpin.Parse()
	logrus.SetLevel(logrus.InfoLevel)

	go readDir(*rootpath, *ignorePaths)
	go proof2POE()

	time.Sleep(time.Second)
	wg.Wait()
}
