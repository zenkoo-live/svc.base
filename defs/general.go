/*
 * Copyright (C) LiangYu, Inc - All Rights Reserved
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

var (
	EmptyUUID          = uuid.MustParse("00000000-0000-0000-0000-000000000000")
	PlatformMerchantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

	SaltLength = 16
)

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
