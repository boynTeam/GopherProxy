package proxy

import (
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/BoynChan/GopherProxy/internal/loadbalance"
	"github.com/BoynChan/GopherProxy/internal/urls"
	"github.com/spf13/viper"
)

// Author:Boyn
// Date:2020/8/31

// TODO(Boyn): make it configurable
var transport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second, //连接超时
		KeepAlive: 30 * time.Second, //长连接超时时间
	}).DialContext,
	MaxIdleConns:          100,              //最大空闲连接
	IdleConnTimeout:       90 * time.Second, //空闲超时时间
	TLSHandshakeTimeout:   10 * time.Second, //tls握手超时时间
	ExpectContinueTimeout: 1 * time.Second,  //100-continue状态码超时时间
}

func NewHttpProxyHandler(urlSli []string, lbType loadbalance.Type) (*httputil.ReverseProxy, error) {

	registerAddr := viper.GetString("ZookeeperAddr")
	dyUrls, err := urls.NewDynamicUrls(urlSli, lbType, registerAddr)
	if err != nil {
		return nil, err
	}
	director := func(req *http.Request) {
		nextAddr, err := dyUrls.GetNext(req.RemoteAddr)
		if err != nil {
			log.Fatal("get next addr fail")
		}
		target, err := url.Parse(nextAddr)
		if err != nil {
			log.Fatal(err)
		}
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}
		req.Header.Set("X-Real-Ip", req.RemoteAddr)
	}

	//更改内容
	modifyFunc := func(resp *http.Response) error {
		if resp.StatusCode != 200 {
			//获取内容
			oldPayload, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			//追加内容
			newPayload := []byte("StatusCode error:" + string(oldPayload))
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(newPayload))
			resp.ContentLength = int64(len(newPayload))
			resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(newPayload)), 10))
		}
		return nil
	}

	//错误回调 ：关闭real_server时测试，错误回调
	//范围：transport.RoundTrip发生的错误、以及ModifyResponse发生的错误
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "ErrorHandler error:"+err.Error(), 500)
	}

	return &httputil.ReverseProxy{Director: director, Transport: transport, ModifyResponse: modifyFunc, ErrorHandler: errFunc}, nil

}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
