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
	"github.com/conseweb/poe/api"
	"github.com/parnurzeal/gorequest"
	"github.com/Sirupsen/logrus"
)

const (
	// api url for document register
	api_doc_reg = "/poe/v1/documents"
	// api url for document proof status
	api_doc_proof_status = "/poe/v1/documents/%s/status"
	// api url for proof result check
	api_doc_proof_result = "/poe/v1/documents/result"
)

func proof2POE() {
	for {
		select {
		case fi := <- waitFiles:
			finfo := fi.String()
			docId := docRegister(*wp, finfo)

			logrus.Infof("%s -> %s", docId, finfo)
			wg.Done()
		}
	}
}

func docRegister(wp, data string) string {
	submitResp := new(api.DocumentSubmitResponse)
	_, _, errs := gorequest.New().Post(fmt.Sprintf("%s%s", *hostpoe, api_doc_reg)).Type("json").Send(map[string]string{
		"proofWaitPeriod": wp,
		"rawDocument":     data,
	}).EndStruct(submitResp)
	if len(errs) != 0 {
		panic(errs[0])
	}

	return submitResp.DocumentID
}

func docProofStatus(docId string) *api.GetProofStatusResponse {
	statusResp := new(api.GetProofStatusResponse)
	_, _, errs := gorequest.New().Get(fmt.Sprintf("%s%s", *hostpoe, fmt.Sprintf(api_doc_proof_status, docId))).EndStruct(statusResp)
	if len(errs) != 0 {
		panic(errs[0])
	}

	if statusResp.Status == "none" {
		panic(fmt.Errorf("document[%s] has not been accepted by poe", docId))
	}

	return statusResp
}

func docProofResult(data string) *api.GetProofResponse {
	proofResp := new(api.GetProofResponse)
	_, _, errs := gorequest.New().Post(fmt.Sprintf("%s%s", *hostpoe, api_doc_proof_result)).Type("json").Send(map[string]string{
		"rawDocument": data,
	}).EndStruct(proofResp)
	if len(errs) != 0 {
		panic(errs[0])
	}

	if proofResp.Status == "none" || proofResp.Status == "invalid" {
		panic(fmt.Errorf("document proof failure."))
	}

	return proofResp
}
