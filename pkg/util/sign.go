package util

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io"

	"github.com/sirupsen/logrus"
)

func RsaCheckSign(data string, sign []byte, publicKeyStr string) error {
	block, _ := pem.Decode([]byte(publicKeyStr))
	w := md5.New()
	io.WriteString(w, data)
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"func": "RsaCheckSign",
		}).Errorf("ParsePKIXPublicKey error :%v", err)
		return err
	}
	md5Bytes := w.Sum(nil)
	return rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.MD5, md5Bytes, sign)
}

func RsaSign(data string, keyBytes string) string {
	w := md5.New()
	io.WriteString(w, data)
	md5_byte := w.Sum(nil)
	value := RsaSignWithMD5(md5_byte, []byte(keyBytes))
	encodeString := base64.StdEncoding.EncodeToString(value)
	return encodeString
}

func RsaSignWithMD5(data []byte, keyBytes []byte) []byte {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		logrus.Error("private key error")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		logrus.Error("ParsePKCS8PrivateKey err")
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.MD5, data)
	if err != nil {
		logrus.Error("Error from signing")
	}

	return signature
}

func RsaSignCheckWithMD5(data []byte, keyBytes []byte) {

}
