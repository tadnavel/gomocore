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

package jwtc

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	sctx "github.com/tadnavel/gomocore"
	"github.com/tadnavel/gomocore/components/loggerc"
)

const (
	defaultSecret                    = "very-important-please-change-it!" // 32 bytes
	defaultAccessTokenExpireSeconds  = 10 * 60                            // 10 minutes
	defaultRefreshTokenExpireSeconds = 1 * 60 * 60                        // 1 hours
	minSecretLength                  = 32
	minAccessTokenLifeTime           = 60
	minRefreshTokenLifeTime          = 60 * 60
)

var (
	ErrSecretKeyNotValid            = errors.New("secret key must be in 32 bytes")
	ErrAccessTokenLifeTimeTooShort  = fmt.Errorf("access token life time too short(min=%d)", minAccessTokenLifeTime)
	ErrRefreshTokenLifeTimeTooShort = fmt.Errorf("refresh token life time too short(min=%d)", minRefreshTokenLifeTime)
	ErrRefreshTokenLifeTimeNotValid = errors.New("refresh token life time not valid")
)

type jwtx struct {
	id                          string
	secret                      string
	expireAccessTokenInSeconds  int
	expireRefreshTokenInSeconds int
	logger                      loggerc.Logger
}

// Return JWT component
func NewJWTComponent(id string) *jwtx {
	return &jwtx{
		id: id,
	}
}

func (j *jwtx) InitFlags() {
	flag.StringVar(
		&j.secret,
		"jwt-secret",
		defaultSecret,
		"secret key to sign JWT",
	)
	flag.IntVar(
		&j.expireAccessTokenInSeconds,
		"jwt-access-exp",
		defaultAccessTokenExpireSeconds,
		"access token life time in second",
	)
	flag.IntVar(
		&j.expireRefreshTokenInSeconds,
		"jwt-refresh-exp",
		defaultRefreshTokenExpireSeconds,
		"refresh token life time in second",
	)
}

func (j *jwtx) ID() string {
	return j.id
}

func (j *jwtx) Activate(serviceCtx sctx.ServiceContext) error {
	j.logger = serviceCtx.Logger(j.id)

	j.logger.Info("activating...")
	if len(j.secret) < minSecretLength {
		return ErrSecretKeyNotValid
	}

	if j.secret == defaultSecret {
		j.logger.Warn("using default secret key, please change it!")
	}

	if j.expireAccessTokenInSeconds < minAccessTokenLifeTime {
		return ErrAccessTokenLifeTimeTooShort
	}

	if j.expireRefreshTokenInSeconds < minRefreshTokenLifeTime {
		return ErrRefreshTokenLifeTimeTooShort
	}

	if j.expireRefreshTokenInSeconds < j.expireAccessTokenInSeconds {
		return ErrRefreshTokenLifeTimeNotValid
	}

	j.logger.Info("activated")

	return nil
}

func (j *jwtx) Stop() error {
	return nil
}

func (j *jwtx) generateToken(
	ctx context.Context,
	sub string,
	id string,
	exp int,
) (token string, tokenLifeTime int, err error) {
	j.logger.With("id", id).With("subject", sub).Debug("generating token")
	now := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(exp))),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        id,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenSignedStr, err := t.SignedString([]byte(j.secret))

	if err != nil {
		j.logger.Errorf("failed to sign token: id=%s, err=%w", id, err)
		return "", 0, err
	}

	return tokenSignedStr, exp, nil
}

func (j *jwtx) IssueAccessToken(
	ctx context.Context,
	id string,
	sub string,
) (token string, tokenLifeTime int, err error) {
	tokenStr, exp, err := j.generateToken(ctx, sub, id, j.expireAccessTokenInSeconds)

	if err != nil {
		return "", 0, err
	}

	return tokenStr, exp, nil
}

func (j *jwtx) IssueRefreshToken(
	ctx context.Context,
	id string,
	sub string,
) (token string, tokenLifeTime int, err error) {
	tokenStr, exp, err := j.generateToken(ctx, sub, id, j.expireRefreshTokenInSeconds)

	if err != nil {
		return "", 0, err
	}

	return tokenStr, exp, nil
}

func (j *jwtx) ParseToken(
	ctx context.Context,
	tokenString string,
) (*jwt.RegisteredClaims, error) {
	var rc jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(tokenString, &rc, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if token == nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return &rc, nil
}
