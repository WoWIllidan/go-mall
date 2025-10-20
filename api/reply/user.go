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
