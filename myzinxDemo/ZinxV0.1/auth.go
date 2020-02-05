package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	APP_ID = "bf796c1d7081462a49042c0a71ed9b143"
	KEY    = "8bf76c1d7081462a9042c0a71ed9b142"
)

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		values := req.URL.Query()
		data := url.QueryEscape(values.Encode()) // Encode() 实现了排序
		fmt.Fprintf(res, "步骤1：URL 参数规范化。\n\t%s\n\n", data)

		t := time.Now().UnixNano() / int64(time.Millisecond)
		toSign := fmt.Sprintf("%s&%s&%d&%s", req.Method, url.QueryEscape(req.URL.Path), t, data)
		fmt.Fprintf(res, "步骤2：构造用于计算签名的字符串。\n\t%s\n\n", toSign)

		hasher := hmac.New(sha256.New, []byte(KEY))
		hasher.Write([]byte(toSign))
		result := base64.StdEncoding.EncodeToString(hasher.Sum(nil))
		fmt.Fprintf(res, "步骤3：计算签名。\n\t%s\n\n", result)

		fmt.Fprintf(
			res,
			"步骤4：将 Authorization 填写到 openApi 请求消息的 HTTP 头中。\n\tAuthorization: algorithm=HMAC-SHA256,appId=%s,time=%d,sign=%s\n\n",
			APP_ID,
			t,
			result,
		)
	})

	if err := http.ListenAndServe(":9000", nil); err != nil {
		panic(err)
	}
}
