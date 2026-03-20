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
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	defaultPasswordCost = 12
	maxBcryptBytes      = 72
)

var ErrPasswordTooLong = errors.New("password too long")

type Hasher struct {
	cost int
}

func NewHasher(cost int) *Hasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = defaultPasswordCost
	}
	return &Hasher{cost: cost}
}

func (r *Hasher) HashPassword(salt, password string) (string, error) {
	spStr := fmt.Sprintf("%s.%s", salt, password)
	if len(spStr) > maxBcryptBytes {
		return "", ErrPasswordTooLong
	}

	h, err := bcrypt.GenerateFromPassword([]byte(spStr), r.cost)
	if err != nil {
		return "", err
	}

	return string(h), nil
}

func (r *Hasher) CompareHashPassword(hashedPassword, salt, password string) bool {
	spStr := fmt.Sprintf("%s.%s", salt, password)
	if len(spStr) > maxBcryptBytes {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(spStr)) == nil
}
