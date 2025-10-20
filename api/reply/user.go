package reply

type TokenReply struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	Duration      int64  `json:"duration"`
	SrvCreateTime string `json:"srv_create_time"`
}

type UserInfoReply struct {
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	LoginName string `json:"login_name"`
	Verified  int    `json:"verified"`
	Avatar    string `json:"avatar"`
	Slogan    string `json:"slogan"`
	IsBlocked int    `json:"is_blocked"`
	CreatedAt string `json:"created_at"`
}

// PasswordResetApply 申请重置密码的响应
type PasswordResetApply struct {
	PasswordResetToken string `json:"password_reset_token"`
}

type UserAddress struct {
	ID              int64  `json:"id"`
	UserName        string `json:"user_name"`
	UserPhone       string `json:"user_phone"`
	MaskedUserName  string `json:"masked_user_name"`  // 用于前台展示的脱敏后的用户姓名
	MaskedUserPhone string `json:"masked_user_phone"` // 用于前台展示的脱敏后的用户手机号
	Default         int    `json:"default"`
	ProvinceName    string `json:"province_name"`
	CityName        string `json:"city_name"`
	RegionName      string `json:"region_name"`
	DetailAddress   string `json:"detail_address"`
	CreatedAt       string `json:"created_at"`
}
