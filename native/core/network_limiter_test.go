package core

import (
	"net/http"
	"testing"
)

func TestCatalogResponseBlockReason(t *testing.T) {
	cases := []struct {
		name   string
		header http.Header
		body   string
		want   string
	}{
		{
			name: "hard block",
			body: "<title>Attention Required! | Cloudflare</title>Sorry, you have been blocked",
			want: "Cloudflare 拒绝了当前客户端，当前请求无法通过",
		},
		{
			name: "unable to access",
			body: "cloudflare you are unable to access this site",
			want: "Cloudflare 拒绝了当前客户端，当前请求无法通过",
		},
		{
			name:   "challenge header",
			header: http.Header{"Cf-Mitigated": []string{"challenge"}},
			body:   "ok",
			want:   "站点要求浏览器验证，当前请求无法通过",
		},
		{
			name: "playlist",
			body: "#EXTM3U\n#EXT-X-VERSION:3\n",
			want: "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			header := tc.header
			if header == nil {
				header = http.Header{}
			}
			got := catalogResponseBlockReason(&http.Response{Header: header, StatusCode: http.StatusForbidden}, []byte(tc.body))
			if got != tc.want {
				t.Fatalf("got %q", got)
			}
		})
	}
}
