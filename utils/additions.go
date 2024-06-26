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
	original map[string]string
}

func NewAdditions(data json.RawMessage) *Additions {
	a := &Additions{
		original: make(map[string]string),
	}

	if data != nil {
		json.Unmarshal(data, &a.original)
	}

	return a
}

func (a *Additions) Get(key string) string {
	return a.original[key]
}

func (a *Additions) Set(key, value string) {
	if key != "" {
		if value != "" {
			a.original[key] = value
		} else {
			delete(a.original, key)
		}
	}
}

func (a *Additions) All() map[string]string {
	return a.original
}

func (a *Additions) Marshal() json.RawMessage {
	o, _ := json.Marshal(a.original)

	return o
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
