package util

// 存放加解密相关的工具函数

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

// AesEncrypt AES加密  ｜ key长度为 16 字节才能加密成功
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unPadding 次
	unPadding := int(origData[length-1])
	if unPadding < 1 || unPadding > 32 {
		unPadding = 0
	}
	return origData[:(length - unPadding)]
}

// RsaSignPKCS1v15 对消息的散列值进行数字签名
func RsaSignPKCS1v15(msg, privateKey []byte, hashType crypto.Hash) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key decode error")
	}
	pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("parse private key error")
	}
	key, ok := pri.(*rsa.PrivateKey)
	if ok == false {
		return nil, errors.New("private key format error")
	}
	sign, err := rsa.SignPKCS1v15(cryptoRand.Reader, key, hashType, msg)
	if err != nil {
		return nil, errors.New("sign error")
	}
	return sign, nil
}

// SHA256HashString 对字符串消息进行 sha256 哈希
func SHA256HashString(stringMessage string) string {
	message := []byte(stringMessage) //字符串转化字节数组
	//创建一个基于SHA256算法的hash.Hash接口的对象
	hash := sha256.New() //sha-256加密
	//hash := sha512.New() //SHA-512加密
	//输入数据
	hash.Write(message)
	//计算哈希值
	bytes := hash.Sum(nil)
	//将字符串编码为16进制格式,返回字符串
	hashCode := hex.EncodeToString(bytes)
	//返回哈希值
	return hashCode
}

func SHA256HashBytes(stringMessage string) []byte {
	message := []byte(stringMessage)
	hash := sha256.New()
	hash.Write(message)
	bytes := hash.Sum(nil)
	return bytes
}
