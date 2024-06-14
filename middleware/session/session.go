/*
 * Copyright (C) LiangYu, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file session.go
 * @package session
 * @author Dr.NP <conan.np@gmail.com>
 * @since 06/14/2024
 */

package session

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zenkoo-live/svc.base/utils"
)

var (
	IDKey   string
	DataKey string
)

func New(config ...Config) fiber.Handler {
	cfg := configDefault(config...)
	IDKey = cfg.IDKey
	DataKey = cfg.DataKey

	return func(c *fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// No storage
		if cfg.Storage == nil {
			return c.Next()
		}

		sessionID := ""

		// Get session ID
		switch strings.ToLower(cfg.IDSource) {
		case "cookie":
			// Get session id from cookie
			sessionID = c.Cookies(cfg.IDKey)
		default:
			// Header
			v := c.Get(cfg.IDKey)
			if v == "" {
				v = c.Get("Authorization")
			}

			parts := strings.Split(v, " ")
			if len(parts) > 1 && parts[1] != "" {
				sessionID = parts[1]
			} else {
				sessionID = v
			}
		}

		if sessionID == "" {
			// Generate new
			sessionID = cfg.IDPrefix + uuid.NewString()
		}

		// Get value
		src, err := cfg.Storage.Get(c.Context(), sessionID).Bytes()
		if err != nil {
			src = nil
		}

		data := NewData(sessionID, src)
		if cfg.StrictAuth != "" {
			if data.Get(cfg.StrictAuth) == nil {
				resp := utils.WrapResponse(nil)
				resp.SetStatus(fiber.StatusUnauthorized)
				resp.SetCode(utils.CodeAuthFailed)
				resp.SetMessage(utils.MsgAuthFailed)

				return c.Format(resp)
			}
		}

		// Context
		c.Locals(cfg.IDKey, sessionID)
		c.Locals(cfg.DataKey, data)

		// Go next
		err = c.Next()
		if err != nil {
			return err
		}

		// Flush session
		src = data.Marshal()
		err = cfg.Storage.Set(c.Context(), sessionID, src, cfg.Expiration).Err()

		// Return session ID
		switch strings.ToLower(cfg.IDSource) {
		case "cookie":
			// Set header
			cookie := &fiber.Cookie{
				Name:    cfg.IDKey,
				Value:   sessionID,
				Expires: time.Now().Add(cfg.Expiration),
			}
			c.Cookie(cookie)
		default:
			// Set header
			c.Response().Header.Add(cfg.IDKey, sessionID)
		}

		return err
	}
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
