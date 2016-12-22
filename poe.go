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
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/conseweb/poe/api"
	"github.com/parnurzeal/gorequest"
)

const (
	// api url for document register
	api_doc_reg = "/poe/v1/documents"
	// api url for document proof status
	api_doc_proof_status = "/poe/v1/documents/%s/status"
	// api url for proof result check
	api_doc_proof_result = "/poe/v1/documents/%s/result"
)

func proof2POE(poes *DirPoes) {
	for {
		select {
		case fi := <-waitFiles:
			if fi.Name != fileName && poes.PutFile(fi) {
				docId := docRegister(*wp, fi)
				logrus.Infof("%s -> %s", docId, fi.Sha256)
				poes.Files[fi.Sha256].Id = docId
			}

			wg.Done()
		}
	}
}

func docRegister(wp string, fi *FileInfo) string {
	submitResp := new(api.DocumentSubmitResponse)
	_, _, errs := gorequest.New().Post(fmt.Sprintf("%s%s", *hostpoe1, api_doc_reg)).Type("json").Send(map[string]string{
		"proofWaitPeriod": wp,
		"rawDocument":     fi.Raw(),
		"metadata":        fi.Metadata(),
	}).EndStruct(submitResp)
	if len(errs) != 0 {
		logrus.Panic(errs[0])
	}

	return submitResp.DocumentID
}

func docProofStatus(docId string) *api.GetProofStatusResponse {
	statusResp := new(api.GetProofStatusResponse)
	_, _, errs := gorequest.New().Get(fmt.Sprintf("%s%s", *hostpoe2, fmt.Sprintf(api_doc_proof_status, docId))).EndStruct(statusResp)
	if len(errs) != 0 {
		panic(errs[0])
	}

	if statusResp.Status == "none" {
		panic(fmt.Errorf("document[%s] has not been accepted by poe", docId))
	}

	return statusResp
}

func docProofResult(docId string) *api.GetProofResponse {
	proofResp := new(api.GetProofResponse)
	_, _, errs := gorequest.New().Get(fmt.Sprintf("%s%s", *hostpoe2, fmt.Sprintf(api_doc_proof_result, docId))).EndStruct(proofResp)
	if len(errs) != 0 {
		logrus.Panic(errs[0])
	}

	if proofResp.Status == "none" || proofResp.Status == "invalid" {
		logrus.Panic(fmt.Errorf("document proof failure."))
	}

	return proofResp
}
