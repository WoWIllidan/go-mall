package reply

type TokenReply struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	Duration      int64  `json:"duration"`
	SrvCreateTime string `json:"srv_create_time"`
}

// PasswordResetApply 申请重置密码的响应
type PasswordResetApply struct {
	PasswordResetToken string `json:"password_reset_token"`
}
