package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/conseweb/graver/file"
	"github.com/conseweb/graver/utils"
	"github.com/conseweb/poe/api"
	"github.com/parnurzeal/gorequest"
)

var (
	wg        = &sync.WaitGroup{}
	waitFiles = make(chan *file.FileInfo, 1000)
	submitHost string
	waitperiod string
)

func StartSubmitCmd(ignorePaths []string, host, wp string) {
	submitHost = host
	waitperiod = wp

	curdir := utils.GetCurDir()
	poes := readOrcreatePoeResultFile(curdir)
	go readDir(curdir, ignorePaths)
	go proof2POE(poes)

	time.Sleep(time.Second)
	wg.Wait()

	poes.Save(filepath.Join(curdir, metadataFileName))
}

func readDir(dir string, ignores []string) {
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

		fi := &file.FileInfo{
			Name:         info.Name(),
			Size:         info.Size(),
			LastModified: info.ModTime().UnixNano(),
			Type:         "",
			Desc:         "",
			Sha256:       file.FileSha256(path),
		}

		logrus.Debugf("File: %s", fi)
		waitFiles <- fi
		wg.Add(1)
		return nil
	}); err != nil {
		logrus.Panic(err)
	}
}


func proof2POE(poes *file.DirPoes) {
	for {
		select {
		case fi := <-waitFiles:
			if fi.Name != metadataFileName && poes.PutFile(fi) {
				docId := docRegister(fi)
				logrus.Infof("%s -> %s", docId, fi.Sha256)
				poes.Files[fi.Sha256].Id = docId
			}

			wg.Done()
		}
	}
}

func docRegister(fi *file.FileInfo) string {
	submitResp := new(api.DocumentSubmitResponse)
	_, _, errs := gorequest.New().Post(fmt.Sprintf("%s/poe/v1/documents", submitHost)).Type("json").Send(map[string]string{
		"proofWaitPeriod": waitperiod,
		"rawDocument":     fi.Raw(),
		"metadata":        fi.Metadata(),
	}).EndStruct(submitResp)
	if len(errs) != 0 {
		logrus.Panic(errs[0])
	}

	return submitResp.DocumentID
}
