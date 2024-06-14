/*
 * Copyright (C) LiangYu, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file data.go
 * @package session
 * @author Dr.NP <conan.np@gmail.com>
 * @since 06/14/2024
 */

package session

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vmihailenco/msgpack/v5"
)

type Data struct {
	id string
	d  map[string]interface{}
}

func (d *Data) Get(key string) interface{} {
	if d.d != nil {
		return d.d[key]
	}

	return nil
}

func (d *Data) Set(key string, value interface{}) {
	if d.d == nil {
		d.d = make(map[string]interface{})
	}

	d.d[key] = value
}

func (d *Data) Marshal() []byte {
	b, err := msgpack.Marshal(d.d)
	if err != nil {
		return nil
	}

	return b
}

func (d *Data) ID() string {
	return d.id
}

func (d *Data) Purge() {
	d.d = nil
}

func (d *Data) Remove(c *fiber.Ctx) {
	if Storage != nil {
		Storage.Del(c.Context(), d.id)
	}

	d.d = nil
	if IDSource == "cookie" {
		c.ClearCookie(IDKey)
	}

	c.Locals(IDKey, nil)
}

func NewData(id string, src []byte) *Data {
	ret := &Data{
		id: id,
	}
	if src != nil {
		d := make(map[string]interface{})
		err := msgpack.Unmarshal(src, &d)
		if err == nil {
			ret.d = d
		}
	}

	return ret
}

func FromFiber(c *fiber.Ctx) *Data {
	data, ok := c.Locals(DataKey).(*Data)
	if ok {
		return data
	}

	return &Data{}
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
