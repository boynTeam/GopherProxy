package proxy

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
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
	"github.com/BoynChan/GopherProxy/pkg"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
)

// Author:Boyn
// Date:2020/8/31

// TODO(Boyn): make it configurable
var transport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second, //连接超时
		KeepAlive: 30 * time.Second, //长连接超时时间
	}).DialContext,
	TLSClientConfig: func() *tls.Config {
		pool := x509.NewCertPool()
		caCertPath := viper.GetString("Http.Cert")
		caCrt, _ := ioutil.ReadFile(caCertPath)
		pool.AppendCertsFromPEM(caCrt)
		return &tls.Config{RootCAs: pool, InsecureSkipVerify: true}
	}(),
	MaxIdleConns:          100,              //最大空闲连接
	IdleConnTimeout:       90 * time.Second, //空闲超时时间
	TLSHandshakeTimeout:   10 * time.Second, //tls握手超时时间
	ExpectContinueTimeout: 1 * time.Second,  //100-continue状态码超时时间
}

func NewHttpProxyHandler(urlSli []string, lbType loadbalance.Type) (gin.HandlerFunc, error) {

	registerAddr := viper.GetString("ZookeeperAddr")
	dyUrls, err := urls.NewDynamicUrls(urlSli, lbType, registerAddr)
	if err != nil {
		return nil, err
	}
	director := getDirector(dyUrls)

	//更改内容
	modifyFunc := getModifyFunc()

	//错误回调 ：关闭real_server时测试，错误回调
	//范围：transport.RoundTrip发生的错误、以及ModifyResponse发生的错误
	errFunc := getErrorfunc()

	http2.ConfigureTransport(transport)
	reverseProxy := &httputil.ReverseProxy{
		Director:       director,
		Transport:      transport,
		ModifyResponse: modifyFunc,
		ErrorHandler:   errFunc}

	return gin.WrapH(reverseProxy), nil

}

func getErrorfunc() func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "ErrorHandler error:"+err.Error(), 500)
	}
}

func getModifyFunc() func(resp *http.Response) error {
	return func(resp *http.Response) error {
		if strings.Contains(resp.Header.Get("Connection"), "Upgrade") {
			return nil
		}
		if resp.StatusCode != 200 {
			//获取内容
			oldPayload, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			message := pkg.NewMessageBuilder().Code(pkg.DownstreamErrorCode).Message(string(oldPayload)).Build()
			newPayload, _ := json.Marshal(message)
			//追加内容
			resp.StatusCode = http.StatusOK
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(newPayload))
			resp.ContentLength = int64(len(newPayload))
			resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(newPayload)), 10))
		}
		return nil
	}
}

func getDirector(dyUrls *urls.DynamicUrls) func(req *http.Request) {
	return func(req *http.Request) {
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
