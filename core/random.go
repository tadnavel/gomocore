// -----------------------------------------------------------------------------
// Copyright (C) 2026 tadnavel
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// -----------------------------------------------------------------------------

package core

import (
	"crypto/rand"
	"math/big"
)

const (
	Digits           = "0123456789"
	Letters          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LettersAndDigits = Letters + Digits
	HexCharacters    = "0123456789abcdef"
)

func GenOTP(length int) (string, error) {
	return randSequence(length, Digits)
}

func GenSalt(length int) (string, error) {
	if length <= 0 {
		length = 32
	}
	return randSequence(length, LettersAndDigits)
}

func RandSequence(n int, charset string) (string, error) {
	if n <= 0 {
		return "", nil
	}

	result := make([]byte, n)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := range n {
		num, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}

func randSequence(n int, charset string) (string, error) {
	return RandSequence(n, charset)
}
