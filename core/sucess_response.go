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

type successResponse struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
	Paging  any  `json:"paging,omitempty"`
	Extra   any  `json:"extra,omitempty"`
}

func SuccessResponse(data, paging, extra any) *successResponse {
	return &successResponse{
		Success: true,
		Data:    data,
		Paging:  paging,
		Extra:   extra,
	}
}

func ResponseData(data any) *successResponse {
	return SuccessResponse(data, nil, nil)
}

func ResponseWithPaging(data any, paging *Paging) *successResponse {
	return SuccessResponse(data, paging, nil)
}
