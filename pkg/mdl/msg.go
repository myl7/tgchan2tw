// Copyright 2021-2022 myl7
// SPDX-License-Identifier: Apache-2.0

package mdl

type Msg struct {
	ID        string
	Body      string
	ImageUrls []string
	FwdFrom   string
	ReplyTo   string
}
