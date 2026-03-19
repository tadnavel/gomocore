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

package svce

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	sctx "github.com/tadnavel/gomocore"
)

var testServiceCtx sctx.ServiceContext

func TestMain(m *testing.M) {
	testServiceCtx = sctx.NewServiceContext(
		sctx.WithName("test"),
	)

	if err := testServiceCtx.Load(); err != nil {
		log.Fatalln(err)
	}

	code := m.Run()

	if err := testServiceCtx.Stop(); err != nil {
		log.Fatalln(err)
	}

	os.Exit(code)
}

func TestServiceErrorWrap(t *testing.T) {
	errDB := errors.New("can't connect to database")
	svcErr := NewServiceError(
		http.StatusNotFound,
		"USR_01",
		"user not found",
		errDB,
		WithField("user_id", "abc123"),
	)

	wrapped := fmt.Errorf("handler: %w", svcErr)

	var target *ServiceError
	if !errors.As(wrapped, &target) {
		t.Fatal("expected to unwrap ServiceError, got nil")
	}

	if !errors.Is(wrapped, errDB) {
		t.Fatal("expected errors.Is to find ErrNotFound through wrap chain")
	}

	if target.HTTPStatus != http.StatusNotFound {
		t.Errorf("HTTPStatus: got %d, want %d", target.HTTPStatus, http.StatusNotFound)
	}
	if target.Message != "user not found" {
		t.Errorf("Message: got %s, want %s", target.Message, "user not found")
	}
	if target.Fields["user_id"] != "abc123" {
		t.Errorf("Fields[user_id]: got %v, want abc123", target.Fields["user_id"])
	}

	logger := testServiceCtx.Logger("error-logger")
	logger.WithFields(target.GetFields()).Info(target)

	logger.Info("--------------------------------------")

	t.Logf("Error()     : %s", wrapped.Error())
	t.Logf("HTTPStatus  : %d", target.HTTPStatus)
	t.Logf("Code        : %s", target.Code)
	t.Logf("Message     : %s", target.Message)
	t.Logf("Fields      : %v", target.Fields)
	t.Logf("errors.Is(ErrNotFound): %v", errors.Is(wrapped, errDB))
}
