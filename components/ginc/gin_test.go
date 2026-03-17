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
