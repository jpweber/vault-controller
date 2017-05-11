// Copyright 2016 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/vault/api"
)

type tokenHandler struct {
	vaultAddr string
	namespace string
	name      string
}

func (h tokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat(tokenFile)
	if !os.IsNotExist(err) {
		log.Println("Token file already exists")
		w.WriteHeader(409)
		return
	}

	var swi api.SecretWrapInfo
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	r.Body.Close()

	err = json.Unmarshal(data, &swi)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	client.SetToken(swi.Token)
	client.SetAddress(h.vaultAddr)

	// Vault knows to unwrap the client token if the token to unwrap is empty.
	secret, err := client.Logical().Unwrap("")
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	f, err := os.Create(tokenFile)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(&secret)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	// save vault secrets to kube secrets
	// get service token
	// token, _ := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")

	// secMetadata := map[string]string{"name": "vault-test2"}
	// secData := map[string]string{"secret_id": b64.URLEncoding.EncodeToString([]byte(secret.Data["secret_id"].(string)))}
	// var payload = make(map[string]interface{})
	// payload["metadata"] = secMetadata
	// payload["data"] = secData
	// payload["type"] = "Opaque"
	// payload["apiVersion"] = ""
	// payload["kind"] = ""
	// jsonBody, _ := json.Marshal(payload)
	// log.Println(jsonBody)

	// secU := fmt.Sprintf("https://kubernetes/api/v1/namespaces/%s/secrets/%s", h.namespace, h.name)
	// r, _ = http.NewRequest("POST", secU, bytes.NewBuffer(jsonBody))
	// r.Header.Add("Authorization", "Bearer "+string(token))
	// r.Header.Add("Content-Type", "application/json")
	// htClient := &http.Client{}
	// resp, err := htClient.Do(r)
	// if err != nil {
	// 	log.Println("Error Saving to Kube Secrets", err)
	// }
	// log.Println(resp.Status)
	// end saving to secrets

	log.Printf("wrote %s", tokenFile)
	w.WriteHeader(200)
}
