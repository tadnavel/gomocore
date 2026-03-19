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
	"log"
	"os"
	"testing"

	sctx "github.com/tadnavel/gomocore"
)

var testServiceCtx sctx.ServiceContext

const ID = "jwt"

func TestMain(m *testing.M) {
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

	os.Exit(code)
}

func TestIssueToken(t *testing.T) {
	jwtComp := testServiceCtx.MustGet(ID).(*jwtx)
	if _, _, err := jwtComp.IssueAccessToken(context.Background(), "tokenID", "userID"); err != nil {
		t.Fatalf("unexpected error(issue access token): %v", err)
	}

	if _, _, err := jwtComp.IssueRefreshToken(context.Background(), "tokenID", "userID"); err != nil {
		t.Fatalf("unexpected error(issue refresh token): %v", err)
	}
}

func TestParseToken(t *testing.T) {
	jwtComp := testServiceCtx.MustGet(ID).(*jwtx)

	tokenString, _, err := jwtComp.IssueAccessToken(context.Background(), "tokenID", "userID")
	if err != nil {
		t.Fatalf("unexpected error(issue access token): %v", err)
	}

	claims, err := jwtComp.ParseToken(context.Background(), tokenString)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if claims.Subject != "userID" {
		t.Fatalf("expected sub=userID, get sub=%s", claims.Subject)
	}

	if claims.ID != "tokenID" {
		t.Fatalf("expected sub=tokenID, get sub=%s", claims.Subject)
	}
}
