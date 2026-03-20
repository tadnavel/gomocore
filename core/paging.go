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

type Paging struct {
	// offset-based
	Page  int   `json:"page,omitempty"  form:"page"`
	Limit int   `json:"limit"           form:"limit"`
	Total int64 `json:"total,omitempty" form:"-"`

	// cursor-based
	Cursor     string `json:"cursor,omitempty"      form:"cursor"`
	NextCursor string `json:"next_cursor,omitempty" form:"-"`
}

func (p *Paging) Process() {
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.Limit > 200 {
		p.Limit = 200
	}

	if p.Cursor == "" && p.Page < 1 {
		p.Page = 1
	}
}

func (p *Paging) IsOffsetBased() bool {
	return p.Cursor == ""
}

func (p *Paging) Offset() int {
	if p.Page < 1 {
		return 0
	}
	return (p.Page - 1) * p.Limit
}
