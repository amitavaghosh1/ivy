package src

import (
	"context"
	"crypto/aes"
	"crypto/md5"
	"net"
)

type MacEncrypter struct{}

func (MacEncrypter) encrypt(ctx context.Context, value string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	mac := interfaces[0].HardwareAddr.String()
	key := md5.Sum([]byte(mac))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	data := []byte(value)

	encrypted := make([]byte, len(data))
	block.Encrypt(encrypted, data)

	return string(encrypted), nil
}

func (MacEncrypter) decrypt(ctx context.Context, value string) (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	mac := interfaces[0].HardwareAddr.String()
	key := md5.Sum([]byte(mac))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	data := []byte(value)

	decrrypted := make([]byte, len(data))
	block.Decrypt(decrrypted, data)

	return string(decrrypted), nil
}
