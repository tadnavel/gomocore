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

package gormc

import (
	"log"
	"os"
	"testing"

	sctx "github.com/tadnavel/gomocore"
)

var testServiceCtx sctx.ServiceContext

func TestMain(m *testing.M) {
	_ = os.Setenv("DB_DRIVER", "sqlite")
	_ = os.Setenv("DB_DSN", ":memory:")

	testServiceCtx = sctx.NewServiceContext(
		sctx.WithName("test"),
		sctx.WithComponent(NewGormDB("gorm", "")),
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

func TestGormDB_Connect_Success(t *testing.T) {
	gormComp := testServiceCtx.MustGet("gorm").(*gormDB)

	db := gormComp.db
	if db == nil {
		t.Fatal("gorm db should not be nil")
	}
}

type TestUser struct {
	ID   uint
	Name string
}

func TestGormDB_InsertFind_Success(t *testing.T) {
	gormComp := testServiceCtx.MustGet("gorm").(*gormDB)
	db := gormComp.GetDB()

	err := db.AutoMigrate(&TestUser{})
	if err != nil {
		t.Fatal(err)
	}

	user := TestUser{Name: "dat"}

	if err := db.Create(&user).Error; err != nil {
		t.Fatal(err)
	}

	var found TestUser
	if err := db.First(&found, "name = ?", "dat").Error; err != nil {
		t.Fatal(err)
	}
}

func Test_gormDB_Activate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		id     string
		prefix string
		// Named input parameters for target function.
		serviceCtx sctx.ServiceContext
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gdb := NewGormDB(tt.id, tt.prefix)
			gotErr := gdb.Activate(tt.serviceCtx)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Activate() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Activate() succeeded unexpectedly")
			}
		})
	}
}
