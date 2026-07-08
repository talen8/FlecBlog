package oauthserver

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type signer struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
	kid        string
}

type accessTokenClaims struct {
	Scope        string `json:"scope"`
	MCPUserID    uint   `json:"mcp_user_id"`
	TokenVersion uint   `json:"token_version"`
	jwt.RegisteredClaims
}

type jwksDocument struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	Kid string `json:"kid"`
}

func loadSigner(path string) (*signer, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("stat MCP OAuth signing key: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("MCP OAuth signing key must not be a symlink")
	}
	if info.Mode().Perm()&0o077 != 0 {
		return nil, fmt.Errorf("MCP OAuth signing key must not be group/world accessible")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read MCP OAuth signing key: %w", err)
	}
	block, rest := pem.Decode(data)
	if block == nil || len(rest) != 0 {
		return nil, fmt.Errorf("MCP OAuth signing key must contain exactly one PEM block")
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse MCP OAuth signing key: %w", err)
	}
	privateKey, ok := parsed.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("MCP OAuth signing key must be Ed25519 PKCS#8")
	}
	publicKey := privateKey.Public().(ed25519.PublicKey)
	digest := sha256.Sum256(publicKey)
	return &signer{
		privateKey: privateKey,
		publicKey:  publicKey,
		kid:        hex.EncodeToString(digest[:8]),
	}, nil
}

func (s *signer) signAccessToken(cfg Config, userID, tokenVersion uint, scope string, now time.Time) (string, error) {
	claims := accessTokenClaims{
		Scope:        scope,
		MCPUserID:    userID,
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			Subject:   "user:" + strconv.FormatUint(uint64(userID), 10),
			Audience:  jwt.ClaimStrings{cfg.Audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(cfg.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        randomToken(16),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = s.kid
	return token.SignedString(s.privateKey)
}

func (s *signer) jwks() jwksDocument {
	return jwksDocument{Keys: []jwk{{
		Kty: "OKP",
		Crv: "Ed25519",
		X:   base64.RawURLEncoding.EncodeToString(s.publicKey),
		Use: "sig",
		Alg: "EdDSA",
		Kid: s.kid,
	}}}
}
