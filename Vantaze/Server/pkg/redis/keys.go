package redis

import "fmt"

const Namespace = "adminpanel:"

const (
	CacheKeyRefreshToken = Namespace + "refresh_token:"
	CacheKeyRateLimit    = Namespace + "ratelimit:"
	CacheKeyUserList     = Namespace + "users:list:"
	CacheKeyAdminList    = Namespace + "admins:list:"
	CacheKeyUserByID     = Namespace + "user:id:"
	CacheKeyAdminByID    = Namespace + "admin:id:"
	CacheKeyConfig       = Namespace + "config:"
	CacheKeyBlacklist    = Namespace + "blacklist:"
	CacheKeyUserByEmail  = Namespace + "user:email:"
)

func RefreshTokenKey(token string) string {
	return CacheKeyRefreshToken + token
}

func RateLimitKey(ip string) string {
	return CacheKeyRateLimit + ip
}

func UserListKey(search string, page, limit int) string {
	return fmt.Sprintf("%s%s:%d:%d", CacheKeyUserList, search, page, limit)
}

func AdminListKey(page, limit int) string {
	return fmt.Sprintf("%s%d:%d", CacheKeyAdminList, page, limit)
}

func UserByIDKey(id uint) string {
	return fmt.Sprintf("%s%d", CacheKeyUserByID, id)
}

func AdminByIDKey(id uint) string {
	return fmt.Sprintf("%s%d", CacheKeyAdminByID, id)
}

func ConfigKey(key string) string {
	return CacheKeyConfig + key
}

func BlacklistKey(token string) string {
	return CacheKeyBlacklist + token
}

func UserByEmailKey(email string) string {
	return CacheKeyUserByEmail + email
}

const (
	UserListPattern  = CacheKeyUserList + "*"
	AdminListPattern = CacheKeyAdminList + "*"
	UserByIDPattern  = CacheKeyUserByID + "*"
	BlacklistPattern = CacheKeyBlacklist + "*"
)