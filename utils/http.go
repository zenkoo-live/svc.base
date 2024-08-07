/*
 * Copyright (C) Zenkoo, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file http.go
 * @package utils
 * @author Dr.NP <conan.np@gmail.com>
 * @since 05/08/2024
 */

package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	CodeOK = 0
	MsgOK  = "OK"

	CodeBodyParseFailed = 9999400000
	MsgBodyParseFailed  = "Parse HTTP body failed"
	CodeValidateFailed  = 9999400001
	MsgValidateFailed   = "Validate HTTP request failed"
	CodeAuthFailed      = 9999401001
	MsgAuthFailed       = "Auth failed"
	CodeStorageFailed   = 9999500001
	MsgStorageFailed    = "Storage failed"

	CodeGeneralFailed = 9999999999
	MsgGeneralFailed  = "General failed"
)

type Envelope struct {
	Code      int         `json:"code"`
	Status    int         `json:"status"`
	Timestamp time.Time   `json:"timestamp"`
	Message   string      `json:"message"`
	RequestId string      `json:"request_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func WrapResponse(data interface{}) *Envelope {
	e := &Envelope{
		Code:      CodeOK,
		Status:    fiber.StatusOK,
		Timestamp: time.Now(),
		Message:   MsgOK,
		Data:      data,
	}

	return e
}

func (e *Envelope) SetStatus(status int) *Envelope {
	e.Status = status

	return e
}

func (e *Envelope) SetCode(code int) *Envelope {
	e.Code = code

	return e
}

func (e *Envelope) SetMessage(msg string) *Envelope {
	e.Message = msg

	return e
}

func (e *Envelope) SetData(data interface{}) *Envelope {
	e.Data = data

	return e
}
func (e *Envelope) SetRequestId(id string) *Envelope {
	e.RequestId = id

	return e
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
