// Package qq 实现 QQ 互联 OAuth 2.0 登录的 goth Provider
package qq

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/markbates/goth"
	"golang.org/x/oauth2"
)

const (
	authURL     = "https://graph.qq.com/oauth2.0/authorize"
	tokenURL    = "https://graph.qq.com/oauth2.0/token" // #nosec G101 - QQ OAuth 标准接口地址
	openIDURL   = "https://graph.qq.com/oauth2.0/me"
	userInfoURL = "https://graph.qq.com/user/get_user_info"
)

// Provider 实现 goth.Provider 接口
type Provider struct {
	ClientKey    string
	Secret       string
	CallbackURL  string
	HTTPClient   *http.Client
	config       *oauth2.Config
	providerName string
}

// New 创建一个新的 QQ Provider
func New(clientKey, secret, callbackURL string, scopes ...string) *Provider {
	p := &Provider{
		ClientKey:    clientKey,
		Secret:       secret,
		CallbackURL:  callbackURL,
		providerName: "qq",
	}
	p.config = newConfig(p, scopes)
	return p
}

func newConfig(provider *Provider, scopes []string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     provider.ClientKey,
		ClientSecret: provider.Secret,
		RedirectURL:  provider.CallbackURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
		Scopes: []string{},
	}

	if len(scopes) > 0 {
		c.Scopes = append(c.Scopes, scopes...)
	} else {
		c.Scopes = []string{"get_user_info"}
	}

	return c
}

// Name 返回 provider 名称
func (p *Provider) Name() string {
	return p.providerName
}

// SetName 设置 provider 名称
func (p *Provider) SetName(name string) {
	p.providerName = name
}

// Client 返回 HTTP 客户端
func (p *Provider) Client() *http.Client {
	return goth.HTTPClientWithFallBack(p.HTTPClient)
}

// Debug 是否启用调试模式
func (p *Provider) Debug(_ bool) {}

// BeginAuth 开始 OAuth 认证流程
func (p *Provider) BeginAuth(state string) (goth.Session, error) {
	authURL := p.config.AuthCodeURL(state)
	session := &Session{
		AuthURL: authURL,
	}
	return session, nil
}

// UnmarshalSession 反序列化 session
func (p *Provider) UnmarshalSession(data string) (goth.Session, error) {
	s := &Session{}
	err := json.Unmarshal([]byte(data), s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// FetchUser 获取用户信息
func (p *Provider) FetchUser(session goth.Session) (goth.User, error) {
	s := session.(*Session)

	if s.AccessToken == "" {
		return goth.User{}, errors.New("缺少访问令牌")
	}

	if s.OpenID == "" {
		return goth.User{}, errors.New("缺少 OpenID")
	}

	user := goth.User{
		AccessToken: s.AccessToken,
		Provider:    p.Name(),
		ExpiresAt:   s.ExpiresAt,
		UserID:      s.OpenID,
	}

	// 获取用户详细信息
	userInfo, err := p.getUserInfo(s.AccessToken, s.OpenID)
	if err != nil {
		return user, err
	}

	user.NickName = userInfo.Nickname
	user.AvatarURL = userInfo.FigureURLQQ2
	if user.AvatarURL == "" {
		user.AvatarURL = userInfo.FigureURLQQ1
	}
	if user.AvatarURL == "" {
		user.AvatarURL = userInfo.FigureURL2
	}
	user.Name = userInfo.Nickname
	user.RawData = map[string]interface{}{
		"nickname":       userInfo.Nickname,
		"figureurl":      userInfo.FigureURL,
		"figureurl_1":    userInfo.FigureURL1,
		"figureurl_2":    userInfo.FigureURL2,
		"figureurl_qq_1": userInfo.FigureURLQQ1,
		"figureurl_qq_2": userInfo.FigureURLQQ2,
		"gender":         userInfo.Gender,
	}

	return user, nil
}

// qqUserInfo QQ 用户信息结构
type qqUserInfo struct {
	Ret          int    `json:"ret"`
	Msg          string `json:"msg"`
	Nickname     string `json:"nickname"`
	FigureURL    string `json:"figureurl"`
	FigureURL1   string `json:"figureurl_1"`
	FigureURL2   string `json:"figureurl_2"`
	FigureURLQQ1 string `json:"figureurl_qq_1"`
	FigureURLQQ2 string `json:"figureurl_qq_2"`
	Gender       string `json:"gender"`
}

// getUserInfo 获取用户详细信息
func (p *Provider) getUserInfo(accessToken, openID string) (*qqUserInfo, error) {
	params := url.Values{}
	params.Set("access_token", accessToken)
	params.Set("oauth_consumer_key", p.ClientKey)
	params.Set("openid", openID)

	reqURL := fmt.Sprintf("%s?%s", userInfoURL, params.Encode())

	resp, err := p.Client().Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo qqUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("解析用户信息失败: %v, body: %s", err, string(body))
	}

	if userInfo.Ret != 0 {
		return nil, fmt.Errorf("获取用户信息失败: %s", userInfo.Msg)
	}

	return &userInfo, nil
}

// RefreshToken 刷新 token（QQ 暂不支持）
func (p *Provider) RefreshToken(_ string) (*oauth2.Token, error) {
	return nil, errors.New("QQ 不支持刷新 token")
}

// RefreshTokenAvailable 返回是否支持刷新 token
func (p *Provider) RefreshTokenAvailable() bool {
	return false
}

// exchangeToken 使用授权码换取 access_token
func (p *Provider) exchangeToken(code string) (string, string, int, error) {
	params := url.Values{}
	params.Set("grant_type", "authorization_code")
	params.Set("client_id", p.ClientKey)
	params.Set("client_secret", p.Secret)
	params.Set("code", code)
	params.Set("redirect_uri", p.CallbackURL)

	reqURL := fmt.Sprintf("%s?%s", tokenURL, params.Encode())

	resp, err := p.Client().Get(reqURL)
	if err != nil {
		return "", "", 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", 0, err
	}

	// 尝试解析 JSON 格式
	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		Error        int    `json:"error,omitempty"`
		ErrorDesc    string `json:"error_description,omitempty"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		// 如果 JSON 解析失败，尝试解析 URL 编码格式
		values, parseErr := url.ParseQuery(string(body))
		if parseErr != nil {
			return "", "", 0, fmt.Errorf("解析 token 响应失败: %v, body: %s", err, string(body))
		}
		tokenResp.AccessToken = values.Get("access_token")
		tokenResp.RefreshToken = values.Get("refresh_token")
		if exp := values.Get("expires_in"); exp != "" {
			if _, err := fmt.Sscanf(exp, "%d", &tokenResp.ExpiresIn); err != nil {
				return "", "", 0, fmt.Errorf("解析 expires_in 失败: %v", err)
			}
		}
	}

	if tokenResp.Error != 0 {
		return "", "", 0, fmt.Errorf("获取 token 失败: %s", tokenResp.ErrorDesc)
	}

	if tokenResp.AccessToken == "" {
		return "", "", 0, errors.New("响应中缺少 access_token")
	}

	return tokenResp.AccessToken, tokenResp.RefreshToken, tokenResp.ExpiresIn, nil
}

// getOpenID 获取用户 OpenID
func (p *Provider) getOpenID(accessToken string) (string, error) {
	params := url.Values{}
	params.Set("access_token", accessToken)
	params.Set("fmt", "json") // 请求 JSON 格式

	reqURL := fmt.Sprintf("%s?%s", openIDURL, params.Encode())

	resp, err := p.Client().Get(reqURL)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	bodyStr := string(body)

	// 尝试解析 JSON 格式
	var openIDResp struct {
		ClientID  string `json:"client_id"`
		OpenID    string `json:"openid"`
		Error     int    `json:"error,omitempty"`
		ErrorDesc string `json:"error_description,omitempty"`
	}

	if err := json.Unmarshal(body, &openIDResp); err != nil {
		// 如果 JSON 解析失败，尝试解析 callback 格式: callback( {"client_id":"...","openid":"..."} );
		re := regexp.MustCompile(`callback\(\s*(\{.*\})\s*\)`)
		matches := re.FindStringSubmatch(bodyStr)
		if len(matches) < 2 {
			return "", fmt.Errorf("解析 OpenID 响应失败: %v, body: %s", err, bodyStr)
		}
		if err := json.Unmarshal([]byte(matches[1]), &openIDResp); err != nil {
			return "", fmt.Errorf("解析 OpenID JSON 失败: %v", err)
		}
	}

	if openIDResp.Error != 0 {
		return "", fmt.Errorf("获取 OpenID 失败: %s", openIDResp.ErrorDesc)
	}

	if openIDResp.OpenID == "" {
		return "", errors.New("响应中缺少 openid")
	}

	return openIDResp.OpenID, nil
}

// CompleteAuth 完成 OAuth 认证
func (p *Provider) CompleteAuth(session goth.Session, params goth.Params) (goth.User, error) {
	s := session.(*Session)

	code := params.Get("code")
	if code == "" {
		return goth.User{}, errors.New("缺少授权码")
	}

	// 使用授权码换取 access_token
	accessToken, refreshToken, expiresIn, err := p.exchangeToken(code)
	if err != nil {
		return goth.User{}, err
	}

	s.AccessToken = accessToken
	s.RefreshToken = refreshToken
	if expiresIn > 0 {
		s.ExpiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
	}

	// 获取 OpenID
	openID, err := p.getOpenID(accessToken)
	if err != nil {
		return goth.User{}, err
	}
	s.OpenID = openID

	return p.FetchUser(s)
}

// Session QQ OAuth Session
type Session struct {
	AuthURL      string    `json:"auth_url"`
	AccessToken  string    `json:"access_token"`  // #nosec G101 - OAuth 标准字段名
	RefreshToken string    `json:"refresh_token"` // #nosec G101 - OAuth 标准字段名
	ExpiresAt    time.Time `json:"expires_at"`
	OpenID       string    `json:"openid"`
}

// GetAuthURL 获取授权 URL
func (s *Session) GetAuthURL() (string, error) {
	if s.AuthURL == "" {
		return "", errors.New("缺少授权 URL")
	}
	return s.AuthURL, nil
}

// Authorize 处理授权回调
func (s *Session) Authorize(provider goth.Provider, params goth.Params) (string, error) {
	p := provider.(*Provider)
	user, err := p.CompleteAuth(s, params)
	if err != nil {
		return "", err
	}
	return user.AccessToken, nil
}

// Marshal 序列化 session
func (s *Session) Marshal() string {
	// #nosec G117 - 这是 OAuth Session 的标准序列化，字段名是 OAuth 标准命名
	b, _ := json.Marshal(s)
	return string(b)
}

func (s *Session) String() string {
	return s.Marshal()
}
