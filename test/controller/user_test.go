package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WoWBytePaladin/go-mall/api/reply"
	"github.com/WoWBytePaladin/go-mall/api/request"
	"github.com/WoWBytePaladin/go-mall/api/router"
	"github.com/WoWBytePaladin/go-mall/logic/appservice"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLoginUser(t *testing.T) {
	Convey("Given right login name and password", t, func() {
		loginName := "yourName@go-mall.com"
		password := "12Qa@6783Wxf3~!45"

		Convey("When use them to Login through API /user/login", func() {
			var s *appservice.UserAppSvc
			gomonkey.ApplyMethod(s, "UserLogin", func(_ *appservice.UserAppSvc, _ *request.UserLogin) (*reply.TokenReply, error) {
				LoginReply := &reply.TokenReply{
					AccessToken:   "70624d19b6644b0bbf8169f51fb5a91f132edebc",
					RefreshToken:  "d16e22fef5cb7f6c69355c9a3c6ce8d1d3b37a84",
					Duration:      7200,
					SrvCreateTime: "2025-02-01 15:34:35",
				}
				return LoginReply, nil
			})

			var b bytes.Buffer
			json.NewEncoder(&b).Encode(map[string]string{"login_name": loginName, "password": password})
			req := httptest.NewRequest(http.MethodPost, "/user/login", &b)
			req.Header.Set("platform", "H5")
			gin.SetMode(gin.ReleaseMode) // 不让它在控制台里输出路由信息
			g := gin.New()
			router.RegisterRoutes(g)
			// mock一个响应记录器
			w := httptest.NewRecorder()
			// 让server端处理mock请求并记录返回的响应内容
			g.ServeHTTP(w, req)

			Convey("Then the user will login successfully", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				// 检验响应内容是否复合预期
				var resp map[string]interface{}
				json.Unmarshal([]byte(w.Body.String()), &resp)
				respData := resp["data"].(map[string]interface{})
				So(respData["access_token"], ShouldNotBeEmpty)
			})
		})
	})
}
