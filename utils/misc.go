/*
 * Copyright (C) Zenkoo, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file misc.go
 * @package utils
 * @author Dr.NP <conan.np@gmail.com>
 * @since 05/08/2024
 */

package utils

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const numericBytes = "0123456789"

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func RandomBytes(length int) []byte {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return b
}

func RandomNumericCode(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = numericBytes[rand.Intn(len(numericBytes))]
	}

	return string(b)
}

func CryptoPassword(original, salt string) string {
	h := sha256.New()
	h.Write([]byte(original + "@@" + salt))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func EnsureStatus(status, bit int64) bool {
	return (status & (1 << bit)) != 0
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
