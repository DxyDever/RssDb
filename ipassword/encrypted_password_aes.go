package ipassword

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func AesDecrypt(org_data string) (string, error) {
	src, err := base64.StdEncoding.DecodeString(org_data)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed, err:%v", err)
	}
	key := "e87swCcx2mi/Rmeq"
	result, err := aesDecrypt(src, []byte(key))
	if err != nil {
		return "", fmt.Errorf("aes decode failed, err:%v", err)
	}
	return string(result), nil
}

func AesEncrypt(origData string) (string, error) {
	key := "e87swCcx2mi/Rmeq"
	keyslice := []byte(key)
	block, err := aes.NewCipher(keyslice)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	originDataSlice := []byte(origData)
	originDataSlice = PKCS7Padding(originDataSlice, blockSize)

	iv := make([]byte, 16)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(originDataSlice))
	blockMode.CryptBlocks(crypted, originDataSlice)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func aesDecrypt(crypted, key []byte) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, 16)
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}
