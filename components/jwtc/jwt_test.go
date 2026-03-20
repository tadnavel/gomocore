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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"testing"

	sctx "github.com/tadnavel/gomocore"
)

var testServiceCtx sctx.ServiceContext

const (
	ID              = "jwt"
	testPrivKeyPath = "testdata/private.pem"
	testPubKeyPath  = "testdata/public.pem"
)

func TestMain(m *testing.M) {
	if err := generateTestKeys(); err != nil {
		log.Fatalf("generate test keys: %v", err)
	}

	os.Args = append(os.Args,
		"-jwt-private-key="+testPrivKeyPath,
		"-jwt-public-key="+testPubKeyPath,
	)

	testServiceCtx = sctx.NewServiceContext(
		sctx.WithName("test"),
		sctx.WithComponent(NewJWTComponent(ID)),
	)

	if err := testServiceCtx.Load(); err != nil {
		log.Fatalln(err)
	}

	code := m.Run()

	if err := testServiceCtx.Stop(); err != nil {
		log.Fatalln(err.Error())
	}

	_ = os.RemoveAll("testdata")

	os.Exit(code)
}

func TestIssueToken(t *testing.T) {
	jwtComp := testServiceCtx.MustGet(ID).(*jwtx)

	if _, _, err := jwtComp.IssueAccessToken(context.Background(), "tokenID", "userID", "user"); err != nil {
		t.Fatalf("unexpected error(issue access token): %v", err)
	}

	if _, _, err := jwtComp.IssueRefreshToken(context.Background(), "tokenID", "userID", "user"); err != nil {
		t.Fatalf("unexpected error(issue refresh token): %v", err)
	}
}

func TestParseToken(t *testing.T) {
	jwtComp := testServiceCtx.MustGet(ID).(*jwtx)

	tokenString, _, err := jwtComp.IssueAccessToken(context.Background(), "tokenID", "userID", "user")
	if err != nil {
		t.Fatalf("unexpected error(issue access token): %v", err)
	}

	claims, err := jwtComp.ParseToken(context.Background(), tokenString)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if claims.Subject != "userID" {
		t.Fatalf("expected sub=userID, got sub=%s", claims.Subject)
	}

	if claims.ID != "tokenID" {
		t.Fatalf("expected jti=tokenID, got jti=%s", claims.ID)
	}

	if claims.Role != "user" {
		t.Fatalf("expected role=user, got role=%s", claims.Role)
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	jwtComp := testServiceCtx.MustGet(ID).(*jwtx)

	_, err := jwtComp.ParseToken(context.Background(), "invalid.token.string")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}

func TestIssueToken_WithoutPrivateKey(t *testing.T) {
	// simulate service chỉ có public key (non-auth service)
	j := &jwtx{
		publicKey:                   testServiceCtx.MustGet(ID).(*jwtx).publicKey,
		expireAccessTokenInSeconds:  defaultAccessTokenExpireSeconds,
		expireRefreshTokenInSeconds: defaultRefreshTokenExpireSeconds,
	}

	_, _, err := j.IssueAccessToken(context.Background(), "tokenID", "userID", "user")
	if err == nil {
		t.Fatal("expected error when private key not loaded, got nil")
	}
}

func generateTestKeys() error {
	if err := os.MkdirAll("testdata", 0700); err != nil {
		return err
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// private key
	privFile, err := os.Create(testPrivKeyPath)
	if err != nil {
		return err
	}
	defer privFile.Close()
	if err := pem.Encode(privFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	}); err != nil {
		return err
	}

	// public key
	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return err
	}
	pubFile, err := os.Create(testPubKeyPath)
	if err != nil {
		return err
	}
	defer pubFile.Close()
	return pem.Encode(pubFile, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})
}
