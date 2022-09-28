package ipassword

import (
	"errors"
)

const (
	ENCRYPT_ALGO_AES = "aes"
)

func Encrypt(origData, algo string) (string, error) {
	if algo == ENCRYPT_ALGO_AES {
		return AesEncrypt(origData)
	}

	return "", errors.New("no this encrypt algo!")
}

func Decrypt(origData, algo string) (string, error) {
	if algo == ENCRYPT_ALGO_AES {
		return AesDecrypt(origData)
	}

	return "", errors.New("no this encrypt algo!")
}
