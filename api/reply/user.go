package reply

type TokenReply struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	Duration      int64  `json:"duration"`
	SrvCreateTime string `json:"srv_create_time"`
}
