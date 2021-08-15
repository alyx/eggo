package main

import (
	"github.com/go-acme/lego/v4/certificate"
	"gorm.io/gorm"
)

type CertificateStore struct {
	gorm.Model
	Domain            string
	CertURL           string
	CertStableURL     string
	PrivateKey        string
	Certificate       string
	IssuerCertificate string
	CSR               string
}

func CertificateStoreFromResource(f *certificate.Resource) *CertificateStore {
	c := &CertificateStore{
		Domain:            f.Domain,
		CertURL:           f.CertURL,
		CertStableURL:     f.CertStableURL,
		PrivateKey:        string(f.PrivateKey),
		Certificate:       string(f.Certificate),
		IssuerCertificate: string(f.IssuerCertificate),
		CSR:               string(f.CSR),
	}
	return c
}

func (c *CertificateStore) ToResource() *certificate.Resource {
	var f *certificate.Resource

	f.Domain = c.Domain
	f.CertURL = c.CertURL
	f.CertStableURL = c.CertStableURL
	f.PrivateKey = []byte(c.PrivateKey)
	f.Certificate = []byte(c.Certificate)
	f.IssuerCertificate = []byte(c.IssuerCertificate)
	f.CSR = []byte(c.CSR)

	return f
}
