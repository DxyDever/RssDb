package ipassword

import (
	"testing"
)

func TestEncryptedPwdAes(t *testing.T) {
	key := "pwd123"
	encryptedKey := "YOivt1fCjWab5B9cRKKEaA=="

	{
		t.Log("====== test encrypt ======")
		output, err := Encrypt(key, ENCRYPT_ALGO_AES)
		if err != nil {
			t.Fatalf("encrypt error:%v\n", err)
		} else if output != encryptedKey {
			t.Fatalf("encrypt error. expected:%s,result:%s\n", encryptedKey, output)
		} else {
			t.Log("encrypt successfully!")
		}
	}

	{
		t.Log("====== test decrypt ======")
		output, err := Decrypt(encryptedKey, ENCRYPT_ALGO_AES)
		if err != nil {
			t.Fatalf("decrypt error:%v\n", err)
		} else if output != key {
			t.Fatalf("decrypt error. expected:%s,result:%s\n", key, output)
		} else {
			t.Log("decrypt successfully!")
		}
	}
}
