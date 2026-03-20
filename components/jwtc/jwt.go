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
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	sctx "github.com/tadnavel/gomocore"
	"github.com/tadnavel/gomocore/components/loggerc"
)

// private key
// openssl genrsa -out private.pem 2048

// public key
// openssl rsa -in private.pem -pubout -out public.pem

const (
	defaultAccessTokenExpireSeconds  = 10 * 60     // 10 minutes
	defaultRefreshTokenExpireSeconds = 1 * 60 * 60 // 1 hours
	minAccessTokenLifeTime           = 60
	minRefreshTokenLifeTime          = 60 * 60
)

var (
	ErrPrivateKeyNotValid           = errors.New("private key not valid")
	ErrPublicKeyNotValid            = errors.New("public key not valid")
	ErrAccessTokenLifeTimeTooShort  = fmt.Errorf("access token life time too short(min=%d)", minAccessTokenLifeTime)
	ErrRefreshTokenLifeTimeTooShort = fmt.Errorf("refresh token life time too short(min=%d)", minRefreshTokenLifeTime)
	ErrRefreshTokenLifeTimeNotValid = errors.New("refresh token life time not valid")
)

type jwtx struct {
	id                          string
	privateKeyPath              string
	publicKeyPath               string
	expireAccessTokenInSeconds  int
	expireRefreshTokenInSeconds int
	privateKey                  *rsa.PrivateKey
	publicKey                   *rsa.PublicKey
	logger                      loggerc.Logger
}

type CustomClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

// Return JWT component
func NewJWTComponent(id string) *jwtx {
	return &jwtx{
		id: id,
	}
}

func (j *jwtx) InitFlags() {
	flag.StringVar(
		&j.privateKeyPath,
		"jwt-private-key",
		"",
		"path to RSA private key PEM file (only required on auth service)",
	)
	flag.StringVar(
		&j.publicKeyPath,
		"jwt-public-key",
		"",
		"path to RSA public key PEM file",
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

	// public key
	if j.publicKeyPath == "" {
		return ErrPublicKeyNotValid
	}
	pubBytes, err := os.ReadFile(j.publicKeyPath)
	if err != nil {
		return fmt.Errorf("read public key: %w", err)
	}
	pubBlock, _ := pem.Decode(pubBytes)
	if pubBlock == nil {
		return ErrPublicKeyNotValid
	}
	pubInterface, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		return fmt.Errorf("parse public key: %w", err)
	}
	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return ErrPublicKeyNotValid
	}
	j.publicKey = pub

	// private key only need in auth service
	if j.privateKeyPath != "" {
		privBytes, err := os.ReadFile(j.privateKeyPath)
		if err != nil {
			return fmt.Errorf("read private key: %w", err)
		}
		privBlock, _ := pem.Decode(privBytes)
		if privBlock == nil {
			return ErrPrivateKeyNotValid
		}
		priv, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
		if err != nil {
			return fmt.Errorf("parse private key: %w", err)
		}
		j.privateKey = priv
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
	sub, id, role string,
	exp int,
) (string, int, error) {
	if j.privateKey == nil {
		return "", 0, errors.New("private key not loaded, this service cannot issue tokens")
	}

	j.logger.With("id", id).With("subject", sub).Debug("generating token")
	now := time.Now().UTC()

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(exp))),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        id,
		},
		Role: role,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := t.SignedString(j.privateKey)
	if err != nil {
		j.logger.Errorf("failed to sign token: id=%s, err=%v", id, err)
		return "", 0, err
	}

	return tokenStr, exp, nil
}

func (j *jwtx) IssueAccessToken(
	ctx context.Context,
	id string,
	sub string,
	role string,
) (token string, tokenLifeTime int, err error) {
	tokenStr, exp, err := j.generateToken(ctx, sub, id, role, j.expireAccessTokenInSeconds)

	if err != nil {
		return "", 0, err
	}

	return tokenStr, exp, nil
}

func (j *jwtx) IssueRefreshToken(
	ctx context.Context,
	id string,
	sub string,
	role string,
) (token string, tokenLifeTime int, err error) {
	tokenStr, exp, err := j.generateToken(ctx, sub, id, role, j.expireRefreshTokenInSeconds)

	if err != nil {
		return "", 0, err
	}

	return tokenStr, exp, nil
}

func (j *jwtx) ParseToken(
	ctx context.Context,
	tokenString string,
) (*CustomClaims, error) {
	var cc CustomClaims

	token, err := jwt.ParseWithClaims(tokenString, &cc, func(token *jwt.Token) (any, error) {
		// only accept RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if token == nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return &cc, nil
}
