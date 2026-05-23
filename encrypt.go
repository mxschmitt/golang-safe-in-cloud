// SPDX-License-Identifier: MIT

package sic

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
)

// Encrypt encrypts decrypted SafeInCloud database XML using the given
// password, producing a payload that round-trips through Decrypt.
func Encrypt(raw []byte, password string) ([]byte, error) {
	out := &bytes.Buffer{}
	if err := binary.Write(out, binary.LittleEndian, uint16(0x0505)); err != nil {
		return nil, fmt.Errorf("could not write magic: %w", err)
	}
	if err := out.WriteByte(0x01); err != nil {
		return nil, fmt.Errorf("could not write sver: %w", err)
	}
	salt, err := randomBytes(16)
	if err != nil {
		return nil, fmt.Errorf("could not generate salt: %w", err)
	}
	if err := writeByteArray(out, salt); err != nil {
		return nil, fmt.Errorf("could not write salt: %w", err)
	}
	outerNonce, err := randomBytes(16)
	if err != nil {
		return nil, fmt.Errorf("could not generate nonce: %w", err)
	}
	if err := writeByteArray(out, outerNonce); err != nil {
		return nil, fmt.Errorf("could not write nonce: %w", err)
	}
	outerKey, err := pbkdf2.Key(sha1.New, password, salt, 10000, 32)
	if err != nil {
		return nil, fmt.Errorf("could not derive key: %w", err)
	}
	if err := writeByteArray(out, nil); err != nil {
		return nil, fmt.Errorf("could not write mystery salt: %w", err)
	}

	innerNonce, err := randomBytes(16)
	if err != nil {
		return nil, fmt.Errorf("could not generate inner nonce: %w", err)
	}
	innerKey, err := randomBytes(32)
	if err != nil {
		return nil, fmt.Errorf("could not generate inner key: %w", err)
	}
	fd := &bytes.Buffer{}
	if err := writeByteArray(fd, innerNonce); err != nil {
		return nil, fmt.Errorf("could not write inner nonce: %w", err)
	}
	if err := writeByteArray(fd, innerKey); err != nil {
		return nil, fmt.Errorf("could not write inner key: %w", err)
	}
	if err := writeByteArray(fd, nil); err != nil {
		return nil, fmt.Errorf("could not write inner trailer: %w", err)
	}
	fdCipher := pkcs7Pad(fd.Bytes(), aes.BlockSize)
	if err := encryptAES(outerKey, outerNonce, &fdCipher); err != nil {
		return nil, fmt.Errorf("could not encrypt fd: %w", err)
	}
	if err := writeByteArray(out, fdCipher); err != nil {
		return nil, fmt.Errorf("could not write fd: %w", err)
	}

	zBuf := &bytes.Buffer{}
	zWriter := zlib.NewWriter(zBuf)
	if _, err := zWriter.Write(raw); err != nil {
		return nil, fmt.Errorf("could not zlib-compress: %w", err)
	}
	if err := zWriter.Close(); err != nil {
		return nil, fmt.Errorf("could not close zlib writer: %w", err)
	}
	body := pkcs7Pad(zBuf.Bytes(), aes.BlockSize)
	if err := encryptAES(innerKey, innerNonce, &body); err != nil {
		return nil, fmt.Errorf("could not encrypt body: %w", err)
	}
	if _, err := out.Write(body); err != nil {
		return nil, fmt.Errorf("could not write body: %w", err)
	}
	return out.Bytes(), nil
}

func encryptAES(key, nonce []byte, content *[]byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("could not create cipher: %w", err)
	}
	cipher.NewCBCEncrypter(block, nonce).CryptBlocks(*content, *content)
	return nil
}

func writeByteArray(w io.Writer, b []byte) error {
	if len(b) > 255 {
		return fmt.Errorf("byte array exceeds 255 bytes: %d", len(b))
	}
	if _, err := w.Write([]byte{byte(len(b))}); err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}
	_, err := w.Write(b)
	return err
}

func pkcs7Pad(b []byte, blockSize int) []byte {
	pad := blockSize - len(b)%blockSize
	out := make([]byte, len(b)+pad)
	copy(out, b)
	for i := len(b); i < len(out); i++ {
		out[i] = byte(pad)
	}
	return out
}

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}
