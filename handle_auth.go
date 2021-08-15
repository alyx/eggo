package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

type ExternalAccountBindingJSON struct {
	Success bool   `json:"success"`
	KID     string `json:"eab_kid"`
	HMAC    string `json:"eab_hmac_key"`
}

func handleZSSLAuth(client *lego.Client, apiKey string) (*registration.Resource, error) {

	resp, err := http.Post("https://api.zerossl.com/acme/eab-credentials?access_key="+apiKey, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var eab *ExternalAccountBindingJSON
	err = json.Unmarshal(data, &eab)
	if err != nil {
		return nil, err
	}

	reg, err := client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
		TermsOfServiceAgreed: true,
		Kid:                  eab.KID,
		HmacEncoded:          eab.HMAC,
	})
	if err != nil {
		return nil, err
	}

	return reg, nil
}
