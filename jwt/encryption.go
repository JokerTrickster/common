package jwt

/*
	AES 암호화 및 복호화 로직
*/

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"

	_error "github.com/JokerTrickster/common/error"
)

// AES encryption logic

func AesEncrypt(ctx context.Context, byteToEncrypt []byte, keyString string) (string, error) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		return "", _error.CreateError(ctx, string(_error.ErrInternalServer), _error.Trace(), "failed to decode AES key", string(_error.ErrFromClient))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", _error.CreateError(ctx, string(_error.ErrInternalServer), _error.Trace(), "failed to create cipher", string(_error.ErrFromInternal))
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", _error.CreateError(ctx, string(_error.ErrInternalServer), _error.Trace(), "failed to create GCM", string(_error.ErrFromInternal))
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", _error.CreateError(ctx, string(_error.ErrInternalServer), _error.Trace(), "failed to create nonce", string(_error.ErrFromInternal))
	}

	return base64.StdEncoding.EncodeToString(aesGCM.Seal(nonce, nonce, byteToEncrypt, nil)), nil
}
