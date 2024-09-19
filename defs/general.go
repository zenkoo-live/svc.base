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

type RuntimePlatform int32

const (
	SaltLength      = 16
	DefaultPageSize = 10

	TokenIssuer = "zenkoo-live"

	Unknown RuntimePlatform = 0
	// Window PC 客户端
	Window RuntimePlatform = 1
	// IOS 客户端
	IOS RuntimePlatform = 2
	// Android 客户端
	Android RuntimePlatform = 4
	// IPad ios pad 客户端
	IPad RuntimePlatform = 8
	// APad android pad 客户端
	APad RuntimePlatform = 16
	// Web 浏览器web端
	Web RuntimePlatform = 32
)

var (
	EmptyUUID          = uuid.MustParse("00000000-0000-0000-0000-000000000000")
	PlatformMerchantID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	PlatformExclusion  = uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")
)

func (r RuntimePlatform) name() string {
	switch r {
	case Unknown:
		return "unknown"
	case Window:
		return "window"
	case IOS:
		return "iOS"
	case Android:
		return "android"
	case IPad:
		return "iPad"
	case APad:
		return "Android Pad"
	case Web:
		return "web"
	}
	return "unknown"
}

func FromPlatform(plat string) RuntimePlatform {
	switch plat {
	case "iOS", "ios", "IOS":
		return IOS
	case "Android", "android", "ANDROID":
		return Android
	case "ipad", "IPAD", "iPad":
		return IPad
	case "apad", "APAD", "aPad", "android pad", "Android Pad":
		return APad
	case "Window", "window", "WINDOW":
		return Window
	case "web", "Web", "WEB":
		return Web
	}
	return Unknown
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
