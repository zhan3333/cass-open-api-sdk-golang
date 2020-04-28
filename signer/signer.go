package signer

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
)

type Signer interface {
	// 签名
	Sign(src []byte, hash crypto.Hash) ([]byte, error)
	// 验证签名
	Verify(src []byte, sign []byte, hash crypto.Hash) error
}

type rsaClient struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func (client *rsaClient) Sign(src []byte, hash crypto.Hash) ([]byte, error) {
	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, client.PrivateKey, hash, hashed)
}

func (client *rsaClient) Verify(src []byte, sign []byte, hash crypto.Hash) error {
	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(client.PublicKey, hash, hashed, sign)
}

func NewSign(privateKey string) (Signer, error) {
	priKey, err := readPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	return &rsaClient{
		PrivateKey: priKey,
		PublicKey:  nil,
	}, nil
}

func New(privateKey, publicKey string) (Signer, error) {
	priKey, err := readPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	pubKey, err := readPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return &rsaClient{
		PrivateKey: priKey,
		PublicKey:  pubKey,
	}, nil
}

// 读取私钥对象
func readPrivateKey(key string) (*rsa.PrivateKey, error) {
	bytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	var privateKey *rsa.PrivateKey
	prkI, err := x509.ParsePKCS8PrivateKey(bytes)
	privateKey = prkI.(*rsa.PrivateKey)
	return privateKey, nil
}

// 读取公钥对象
// PKCS8格式单行key处理
func readPublicKey(key string) (*rsa.PublicKey, error) {
	bytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	var publicKey *rsa.PublicKey
	pubKI, err := x509.ParsePKIXPublicKey(bytes)
	publicKey = pubKI.(*rsa.PublicKey)
	return publicKey, nil
}
