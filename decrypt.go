// SPDX-License-Identifier: MIT

package sic

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// ErrIncorrectPassword means that the credentials are incorrect
var ErrIncorrectPassword = errors.New("incorrect credentials")

// Decrypt decrypts a SafeInCloud database by a given file (e.g. os.Open)
// and a password
func Decrypt(file io.Reader, password string) ([]byte, error) {
	data := bufio.NewReader(file)
	var magic uint16
	// Decrypt the FD
	if err := binary.Read(data, binary.LittleEndian, &magic); err != nil {
		return nil, fmt.Errorf("could not read magic: %w", err)
	}
	if _, err := data.ReadByte(); err != nil {
		return nil, fmt.Errorf("could not read sver: %w", err)
	}
	salt, err := readByteArray(data)
	if err != nil {
		return nil, fmt.Errorf("could not read salt: %w", err)
	}
	nonce, err := readByteArray(data)
	if err != nil {
		return nil, fmt.Errorf("could not read nonce: %w", err)
	}
	pwd, err := pbkdf2.Key(sha1.New, password, salt, 10000, 32)
	if err != nil {
		return nil, fmt.Errorf("could not derive key: %w", err)
	}
	_, err = readByteArray(data) // Idk what this is; salt but not necessary?!
	if err != nil {
		return nil, fmt.Errorf("could not read salt: %w", err)
	}
	block, err := readByteArray(data)
	if err != nil {
		return nil, fmt.Errorf("could not read subfd: %w", err)
	}
	if err := decryptAES(pwd, nonce, &block); err != nil {
		return nil, fmt.Errorf("could not decrypt aes: %w", err)
	}
	fd := bufio.NewReader(bytes.NewBuffer(block))
	encFile, err := io.ReadAll(data)
	if err != nil {
		return nil, fmt.Errorf("could not read remaining encrypted content: %w", err)
	}
	innerNonce, err := readByteArray(fd)
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, ErrIncorrectPassword
	} else if err != nil {
		return nil, fmt.Errorf("could not read nonce: %w", err)
	}
	innerKey, err := readByteArray(fd)
	if err != nil {
		return nil, ErrIncorrectPassword
	}
	if _, err = readByteArray(fd); err != nil {
		return nil, err
	}
	if err := decryptAES(innerKey, innerNonce, &encFile); err != nil {
		return nil, fmt.Errorf("could not decrypt aes: %w", err)
	}
	zReader, err := zlib.NewReader(bytes.NewReader(encFile))
	if err != nil {
		return nil, err
	}
	defer func() { _ = zReader.Close() }()
	return io.ReadAll(zReader)
}

func decryptAES(pwd, nonce []byte, content *[]byte) error {
	block, err := aes.NewCipher(pwd)
	if err != nil {
		return fmt.Errorf("could not create cipher: %w", err)
	}
	cipher.NewCBCDecrypter(block, nonce).CryptBlocks(*content, *content)
	return nil
}

// readByteArray reads a byte array with the given size in the next byte
func readByteArray(data *bufio.Reader) ([]byte, error) {
	size, err := data.ReadByte()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, size)
	if err = binary.Read(data, binary.LittleEndian, &buf); err != nil {
		return nil, err
	}
	return buf, nil
}
