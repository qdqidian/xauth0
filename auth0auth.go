package xauth0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/duke-git/lancet/v2/netutil"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	json_filter "qidian.xyz/scrap/json-filter"
)

func init() {
	caddy.RegisterModule(&HTTPAuth0Auth{})
}

type HTTPAuth0Auth struct {
	URL          string `json:"url,omitempty"`
	HeaderPrefix string `json:"header_prefix,omitempty"`

	logger *zap.Logger
}

// Provision sets up a.
func (a *HTTPAuth0Auth) Provision(ctx caddy.Context) error {
	a.logger = ctx.Logger()
	return nil
}

func (a *HTTPAuth0Auth) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.auth0",
		New: func() caddy.Module { return new(HTTPAuth0Auth) },
	}
}

func (a *HTTPAuth0Auth) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = r.Body.Close()
	if err != nil {
		return err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(b))

	schema := "http"
	if r.TLS != nil {
		schema = "https"
	}

	u, err := url.ParseRequestURI(schema + "://" + r.Host)
	if err != nil {
		return err
	}

	req := &Req{
		Schema: u.Scheme,
		Domain: u.Hostname(),
		Port:   u.Port(),
		Path:   r.URL.Path,
		Form:   r.Form,
		Body:   string(b),
	}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	request := &netutil.HttpRequest{
		RawURL:  a.URL,
		Method:  "POST",
		Headers: r.Header,
		Body:    reqBody,
	}
	request.Headers["Content-Type"] = []string{
		"application/json",
	}
	httpClient := netutil.NewHttpClient()
	resp, err := httpClient.SendRequest(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		result, _ := io.ReadAll(resp.Body)
		resultBody := string(result)
		a.logger.Warn("auth0 return status is not 200. body: ", zap.String("resp body: ", resultBody))
		return caddyhttp.Error(resp.StatusCode, fmt.Errorf(resultBody))
	}

	var res Res
	httpClient.DecodeResponse(resp, &res)

	repl := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	if res.User != nil {
		repl.Set(fmt.Sprintf("%sid", a.HeaderPrefix), res.User.Id)
		repl.Set(fmt.Sprintf("%sname", a.HeaderPrefix), res.User.Name)
		repl.Set(fmt.Sprintf("%scode", a.HeaderPrefix), res.User.Code)

		for k, v := range res.User.Meta {
			repl.Set(fmt.Sprintf("%s%s", a.HeaderPrefix, k), v)
		}
	}

	if res.AllowCrossDomain {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,POST")
		w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")
	}

	if res.FieldPermissionDefine == "" {
		return next.ServeHTTP(w, r)
	}

	writer := &Auth0Writer{body: bytes.NewBufferString(""), ResponseWriter: w}

	err = next.ServeHTTP(writer, r)
	// 字段级权限处理
	body := writer.body.String()
	// todo:: 状态码不为200，则直接跳过
	if body == "" || err != nil || writer.StatusCode != 200 {
		return err
	}

	resultBody := json_filter.Filter(res.FieldPermissionDefine, body)
	_, err = w.Write(resultBody)
	a.logger.Debug("json filter. ", zap.String("originBody: ", body), zap.String("define: ", res.FieldPermissionDefine), zap.String("result: ", string(resultBody)))
	return err
}
