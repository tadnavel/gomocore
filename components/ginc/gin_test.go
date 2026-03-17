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

package ginc

import (
	"log"
	"os"
	"testing"

	sctx "github.com/tadnavel/gomocore"
)

var testServiceCtx sctx.ServiceContext

const ID = "gin"

func TestMain(m *testing.M) {
	testServiceCtx = sctx.NewServiceContext(
		sctx.WithName("test"),
		sctx.WithComponent(NewGin(ID)),
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

func TestGinRun(t *testing.T) {
	ginComp := testServiceCtx.MustGet(ID).(*ginEngine)

	ginComp.Run()

	if ginComp.GetRouter() == nil {
		t.Fatal("router should not be nil")
	}
}
