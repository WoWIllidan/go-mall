package request

type UserRegister struct {
	LoginName       string `json:"login_name" binding:"required,e164|email"` // 验证登录名必须为手机号或者邮箱地址
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
	Nickname        string `json:"nickname" binding:"max=30"`
	Slogan          string `json:"slogan" binding:"max=30"`
	Avatar          string `json:"avatar" binding:"max=100"`
}

// UserLogin 用户登录请求,需要同时验证和绑定Body和Header中的数据
// 使用Gin绑定RequestBoy和Header https://github.com/gin-gonic/gin/issues/2309#issuecomment-2020168668
type UserLogin struct {
	Body struct {
		LoginName string `json:"login_name" binding:"required,e164|email"`
		Password  string `json:"password" binding:"required,min=8"`
	}
	Header struct {
		Platform string `json:"platform" header:"platform" binding:"required,oneof=H5 APP"`
	}
}

type UserInfoUpdate struct {
	Nickname string `json:"nickname" binding:"max=30"`
	Slogan   string `json:"slogan" binding:"max=30"`
	Avatar   string `json:"avatar" binding:"max=100"`
}

type PasswordResetApply struct {
	LoginName string `json:"login_name" binding:"required,e164|email"` // 验证登录名必须为手机号或者邮箱地址
}

type PasswordReset struct {
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
	Token           string `json:"password_reset_token" binding:"required"`
	Code            string `json:"password_reset_code" binding:"required"`
}
