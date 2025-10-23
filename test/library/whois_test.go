package library

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/WoWBytePaladin/go-mall/common/util/httptool"
	"github.com/WoWBytePaladin/go-mall/library"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)
	// 把框架的httptool使用的http client 换成gock拦截的client
	httptool.SetUTHttpClient(client)
	os.Exit(m.Run())
}

func TestWhoisLib_GetHostIpDetail(t *testing.T) {
	defer gock.Off()
	gock.New("https://ipwho.is").
		MatchHeader("User-Agent", "curl/7.77.0").Get("").
		Reply(200).
		BodyString("{\"ip\":\"127.126.113.220\",\"success\":true}")

	ipDetail, err := library.NewWhoisLib(context.TODO()).GetHostIpDetail()
	assert.Nil(t, err)
	assert.Equal(t, "127.126.113.220", ipDetail.Ip)
}
