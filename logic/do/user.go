package do

import "time"

type SessionInfo struct {
	UserId       int64  `json:"user_id"`
	Platform     string `json:"platform"` //平台 app,h5,wx，
	SessionId    string `json:"session_id"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenInfo struct {
	AccessToken   string    `json:"access_token"`
	RefreshToken  string    `json:"refresh_token"`
	Duration      int64     `json:"duration"`
	SrvCreateTime time.Time `json:"srv_create_time"`
}

type UserBaseInfo struct {
	ID        int64     `json:"id"`
	Nickname  string    `json:"nickname"`
	LoginName string    `json:"login_name"`
	Verified  int       `json:"verified"`
	Avatar    string    `json:"avatar"`
	Slogan    string    `json:"slogan"`
	IsBlocked int       `json:"is_blocked"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TokenVerify struct {
	Approved  bool   // 验证结果
	UserId    int64  // 用户ID
	SessionId string // SessionId 可以用于存储一些与登录相关的东西, 用户不重新登录不会变
}
