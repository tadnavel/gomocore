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
