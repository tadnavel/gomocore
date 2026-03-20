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

import (
	"context"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type Requester interface {
	GetSubject() string
	GetTokenId() string
	GetRole() Role
}

type UserRequester struct {
	Sub string
	Tid string
}

func (u *UserRequester) GetSubject() string {
	return u.Sub
}

func (u *UserRequester) GetTokenId() string {
	return u.Tid
}

func (u *UserRequester) GetRole() Role {
	return RoleUser
}

type AdminRequester struct {
	Sub         string
	Tid         string
	Permissions map[string]struct{}
}

func (a *AdminRequester) GetSubject() string {
	return a.Sub
}

func (a *AdminRequester) GetTokenId() string {
	return a.Tid
}

func (a *AdminRequester) GetRole() Role {
	return RoleAdmin
}

func (a *AdminRequester) HasPermission(p string) bool {
	_, ok := a.Permissions[p]
	return ok
}

type requesterKeyType struct{}

var requesterKey = requesterKeyType{}

func GetRequester(ctx context.Context) Requester {
	if r, ok := ctx.Value(requesterKey).(Requester); ok {
		return r
	}
	return nil
}

func ContextWithRequester(ctx context.Context, r Requester) context.Context {
	return context.WithValue(ctx, requesterKey, r)
}
