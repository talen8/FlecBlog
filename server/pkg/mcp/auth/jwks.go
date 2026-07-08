package auth

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	maxJWKSResponseBytes         = 1 << 20
	jwksCacheTTL                 = 5 * time.Minute
	jwksHTTPTimeout              = 10 * time.Second
	jwksForcedRefreshMinInterval = 30 * time.Second
)

var allowedJWTAlgorithms = map[string]struct{}{
	"RS256": {}, "RS384": {}, "RS512": {},
	"PS256": {}, "PS384": {}, "PS512": {},
	"ES256": {}, "ES384": {}, "ES512": {},
	"EdDSA": {},
}

type jwksDocument struct {
	Keys []rawJWK `json:"keys"`
}

type rawJWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

type parsedJWK struct {
	kid string
	alg string
	use string
	key any
}

type jwksCache struct {
	url    string
	client *http.Client

	mu                sync.RWMutex
	refreshMu         sync.Mutex
	keys              []parsedJWK
	expiresAt         time.Time
	lastForcedRefresh time.Time
}

func newJWKSCache(rawURL string) *jwksCache {
	transport := &http.Transport{}
	if defaultTransport, ok := http.DefaultTransport.(*http.Transport); ok {
		transport = defaultTransport.Clone()
	}
	transport.Proxy = nil
	return &jwksCache{
		url: rawURL,
		client: &http.Client{
			Transport: transport,
			Timeout:   jwksHTTPTimeout,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (c *jwksCache) keyForToken(ctx context.Context, token *jwt.Token) (any, error) {
	alg := token.Method.Alg()
	if _, ok := allowedJWTAlgorithms[alg]; !ok {
		return nil, fmt.Errorf("unsupported JWT alg %q", alg)
	}

	kid, _ := token.Header["kid"].(string)
	if key, ok := c.lookup(kid, alg, false); ok {
		return key, nil
	}
	forceRefresh := c.cacheFresh()
	if err := c.refresh(ctx, forceRefresh); err != nil {
		return nil, err
	}
	if key, ok := c.lookup(kid, alg, true); ok {
		return key, nil
	}
	if kid == "" {
		return nil, fmt.Errorf("JWT kid is missing and no unique compatible JWK exists")
	}
	return nil, fmt.Errorf("no compatible JWK found for kid %q", kid)
}

func (c *jwksCache) lookup(kid, alg string, afterRefresh bool) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if !afterRefresh && (len(c.keys) == 0 || time.Now().After(c.expiresAt)) {
		return nil, false
	}

	var candidates []parsedJWK
	for _, candidate := range c.keys {
		if candidate.use != "" && candidate.use != "sig" {
			continue
		}
		if kid != "" && candidate.kid != kid {
			continue
		}
		if candidate.alg != "" && candidate.alg != alg {
			continue
		}
		if !jwkKeySupportsAlg(candidate.key, alg) {
			continue
		}
		candidates = append(candidates, candidate)
	}
	if len(candidates) != 1 {
		return nil, false
	}
	return candidates[0].key, true
}

func (c *jwksCache) cacheFresh() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.keys) > 0 && time.Now().Before(c.expiresAt)
}

func (c *jwksCache) refresh(ctx context.Context, force bool) error {
	c.refreshMu.Lock()
	defer c.refreshMu.Unlock()

	now := time.Now()
	c.mu.RLock()
	fresh := len(c.keys) > 0 && now.Before(c.expiresAt)
	lastForcedRefresh := c.lastForcedRefresh
	c.mu.RUnlock()
	if !force && fresh {
		return nil
	}
	if force && fresh && !lastForcedRefresh.IsZero() && now.Sub(lastForcedRefresh) < jwksForcedRefreshMinInterval {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return fmt.Errorf("create JWKS request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("fetch JWKS: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("fetch JWKS: HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxJWKSResponseBytes+1))
	if err != nil {
		return fmt.Errorf("read JWKS: %w", err)
	}
	if len(body) > maxJWKSResponseBytes {
		return fmt.Errorf("JWKS response exceeds size limit")
	}

	var document jwksDocument
	if err := json.Unmarshal(body, &document); err != nil {
		return fmt.Errorf("decode JWKS: %w", err)
	}
	if len(document.Keys) == 0 {
		return fmt.Errorf("JWKS contains no keys")
	}

	parsed := make([]parsedJWK, 0, len(document.Keys))
	for _, raw := range document.Keys {
		key, err := parseJWK(raw)
		if err != nil {
			continue
		}
		parsed = append(parsed, parsedJWK{kid: raw.Kid, alg: raw.Alg, use: raw.Use, key: key})
	}
	if len(parsed) == 0 {
		return fmt.Errorf("JWKS contains no supported public keys")
	}

	c.mu.Lock()
	c.keys = parsed
	c.expiresAt = time.Now().Add(jwksCacheTTL)
	if force {
		c.lastForcedRefresh = time.Now()
	}
	c.mu.Unlock()
	return nil
}

func parseJWK(raw rawJWK) (any, error) {
	switch raw.Kty {
	case "RSA":
		return parseRSAJWK(raw)
	case "EC":
		return parseECJWK(raw)
	case "OKP":
		return parseOKPJWK(raw)
	default:
		return nil, fmt.Errorf("unsupported JWK kty %q", raw.Kty)
	}
}

func parseRSAJWK(raw rawJWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(raw.N)
	if err != nil || len(nBytes) == 0 {
		return nil, fmt.Errorf("invalid RSA modulus")
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(raw.E)
	if err != nil || len(eBytes) == 0 || len(eBytes) > 4 {
		return nil, fmt.Errorf("invalid RSA exponent")
	}
	exponent := 0
	for _, b := range eBytes {
		exponent = exponent<<8 | int(b)
	}
	if exponent < 3 || exponent%2 == 0 {
		return nil, fmt.Errorf("invalid RSA exponent")
	}
	modulus := new(big.Int).SetBytes(nBytes)
	if modulus.BitLen() < 2048 {
		return nil, fmt.Errorf("RSA modulus is smaller than 2048 bits")
	}
	return &rsa.PublicKey{N: modulus, E: exponent}, nil
}

func parseECJWK(raw rawJWK) (*ecdsa.PublicKey, error) {
	var curve elliptic.Curve
	switch raw.Crv {
	case "P-256":
		curve = elliptic.P256()
	case "P-384":
		curve = elliptic.P384()
	case "P-521":
		curve = elliptic.P521()
	default:
		return nil, fmt.Errorf("unsupported EC curve %q", raw.Crv)
	}
	xBytes, err := base64.RawURLEncoding.DecodeString(raw.X)
	if err != nil || len(xBytes) == 0 {
		return nil, fmt.Errorf("invalid EC x coordinate")
	}
	yBytes, err := base64.RawURLEncoding.DecodeString(raw.Y)
	if err != nil || len(yBytes) == 0 {
		return nil, fmt.Errorf("invalid EC y coordinate")
	}
	x := new(big.Int).SetBytes(xBytes)
	y := new(big.Int).SetBytes(yBytes)
	if !curve.IsOnCurve(x, y) {
		return nil, fmt.Errorf("EC point is not on curve")
	}
	return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}, nil
}

func parseOKPJWK(raw rawJWK) (ed25519.PublicKey, error) {
	if raw.Crv != "Ed25519" {
		return nil, fmt.Errorf("unsupported OKP curve %q", raw.Crv)
	}
	xBytes, err := base64.RawURLEncoding.DecodeString(raw.X)
	if err != nil || len(xBytes) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid Ed25519 public key")
	}
	return ed25519.PublicKey(xBytes), nil
}

func jwkKeySupportsAlg(key any, alg string) bool {
	switch key.(type) {
	case *rsa.PublicKey:
		return alg == "RS256" || alg == "RS384" || alg == "RS512" ||
			alg == "PS256" || alg == "PS384" || alg == "PS512"
	case *ecdsa.PublicKey:
		return alg == "ES256" || alg == "ES384" || alg == "ES512"
	case ed25519.PublicKey:
		return alg == "EdDSA"
	default:
		return false
	}
}
