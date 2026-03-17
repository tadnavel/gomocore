// -----------------------------------------------------------------------------
// Copyright (C) 2026 tadnavel
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
