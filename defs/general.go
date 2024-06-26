/*
 * Copyright (C) Zenkoo, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file general.go
 * @package defs
 * @author Dr.NP <conan.np@gmail.com>
 * @since 05/20/2024
 */

package defs

import "github.com/google/uuid"

const (
	SaltLength      = 16
	DefaultPageSize = 10

	TokenIssuer = "zenkoo-live"
)

var (
	EmptyUUID          = uuid.MustParse("00000000-0000-0000-0000-000000000000")
	PlatformMerchantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	PlatformExclusion  = uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")
)

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
