package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/conseweb/graver/utils"
	"github.com/conseweb/poe/api"
	"github.com/parnurzeal/gorequest"
)

var (
	verifyHost string
)

func StartVerifyCmd(host string) {
	verifyHost = host

	curdir := utils.GetCurDir()
	poes := readOrcreatePoeResultFile(curdir)
	for key, file := range poes.Files {
		if file.Proofed {
			continue
		}

		rsp := docProofResult(file.Id)
		if rsp.Status == "valid" {
			poes.Files[key].Proofed = true
		}
	}

	poes.Save(filepath.Join(curdir, metadataFileName))
}

func docProofResult(docId string) *api.GetProofResponse {
	proofResp := new(api.GetProofResponse)
	_, _, errs := gorequest.New().Get(fmt.Sprintf("%s/poe/v1/documents/%s/result", verifyHost, docId)).EndStruct(proofResp)
	if len(errs) != 0 {
		logrus.Panic(errs[0])
	}

	if proofResp.Status == "none" || proofResp.Status == "invalid" {
		logrus.Panic(fmt.Errorf("document proof failure."))
	}

	return proofResp
}
