/*
 * Copyright (C) Zenkoo, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file misc.go
 * @package utils
 * @author Dr.NP <conan.np@gmail.com>
 * @since 06/22/2024
 */

package utils

import "encoding/json"

type Additions struct {
	original map[string]any
}

func NewAction(data json.RawMessage) *Additions {
	a := &Additions{
		original: make(map[string]any),
	}

	if data != nil {
		json.Unmarshal(data, &a.original)
	}

	return a
}

func (a *Additions) Get(key string) any {
	return a.original[key]
}

func (a *Additions) Set(key string, value any) {
	a.original[key] = value
}

func (a *Additions) All() map[string]any {
	return a.original
}

func (a *Additions) Marshal() (json.RawMessage, error) {
	return json.Marshal(a.original)
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
