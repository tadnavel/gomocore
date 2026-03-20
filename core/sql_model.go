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

package core

import "time"

type SQLModel struct {
	ID        int64      `json:"-" gorm:"column:id;" db:"id"`
	FakeID    *UID       `json:"id" gorm:"-"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at;" db:"created_at"`
}

func NewSQLModel() SQLModel {
	now := time.Now().UTC()
	return SQLModel{
		CreatedAt: &now,
	}
}

func (s *SQLModel) Mask(objectType int, shardID uint32) {
	uid := NewUID(uint32(s.ID), objectType, shardID)
	s.FakeID = &uid
}

type SQLModelWithUpdatedAt struct {
	SQLModel
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;" db:"updated_at"`
}

func NewSQLModelWithUpdatedAt() SQLModelWithUpdatedAt {
	now := time.Now().UTC()
	return SQLModelWithUpdatedAt{
		SQLModel:  SQLModel{CreatedAt: &now},
		UpdatedAt: &now,
	}
}
