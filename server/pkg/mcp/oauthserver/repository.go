package oauthserver

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrStateLimitReached = errors.New("MCP OAuth state limit reached")

type OAuthClient struct {
	ClientID         string     `gorm:"column:client_id;primaryKey"`
	ClientSecretHash string     `gorm:"column:client_secret_hash"`
	RedirectURIs     string     `gorm:"column:redirect_uris"`
	CreatedAt        time.Time  `gorm:"column:created_at"`
	LastUsedAt       *time.Time `gorm:"column:last_used_at"`
}

func (OAuthClient) TableName() string { return "mcp_oauth_clients" }

type AuthorizationCode struct {
	CodeHash      string    `gorm:"column:code_hash;primaryKey"`
	ClientID      string    `gorm:"column:client_id"`
	UserID        uint      `gorm:"column:user_id"`
	RedirectURI   string    `gorm:"column:redirect_uri"`
	Resource      string    `gorm:"column:resource"`
	Scopes        string    `gorm:"column:scopes"`
	CodeChallenge string    `gorm:"column:code_challenge"`
	ExpiresAt     time.Time `gorm:"column:expires_at"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

func (AuthorizationCode) TableName() string { return "mcp_oauth_authorization_codes" }

type RefreshToken struct {
	TokenHash    string    `gorm:"column:token_hash;primaryKey"`
	ClientID     string    `gorm:"column:client_id"`
	UserID       uint      `gorm:"column:user_id"`
	TokenVersion *uint     `gorm:"column:token_version"`
	Scopes       string    `gorm:"column:scopes"`
	ExpiresAt    time.Time `gorm:"column:expires_at"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (RefreshToken) TableName() string { return "mcp_oauth_refresh_tokens" }

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CleanupExpired(now time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("expires_at <= ?", now).Delete(&AuthorizationCode{}).Error; err != nil {
			return err
		}
		return tx.Where("expires_at <= ?", now).Delete(&RefreshToken{}).Error
	})
}

func (r *Repository) CreateClientBounded(client *OAuthClient, maxClients int, inactiveBefore, unusedBefore time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("LOCK TABLE mcp_oauth_clients IN SHARE ROW EXCLUSIVE MODE").Error; err != nil {
			return err
		}
		var count int64
		if err := tx.Model(&OAuthClient{}).Count(&count).Error; err != nil {
			return err
		}
		if count >= int64(maxClients) {
			var stale OAuthClient
			result := tx.Where(
				"(last_used_at IS NULL AND created_at <= ?) OR (last_used_at IS NOT NULL AND last_used_at <= ?)",
				unusedBefore, inactiveBefore,
			).
				Order("COALESCE(last_used_at, created_at) ASC").
				Limit(1).Find(&stale)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected != 1 {
				return ErrStateLimitReached
			}
			if err := tx.Delete(&stale).Error; err != nil {
				return err
			}
		}
		return tx.Create(client).Error
	})
}

func (r *Repository) GetClient(clientID string) (*OAuthClient, error) {
	var client OAuthClient
	if err := r.db.Where("client_id = ?", clientID).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *Repository) TouchClient(clientID string, now time.Time) error {
	return r.db.Model(&OAuthClient{}).Where("client_id = ?", clientID).Update("last_used_at", now).Error
}

func (r *Repository) StoreAuthorizationCodeBounded(code *AuthorizationCode, maxCodes int, now time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("expires_at <= ?", now).Delete(&AuthorizationCode{}).Error; err != nil {
			return err
		}
		if err := tx.Exec("LOCK TABLE mcp_oauth_authorization_codes IN SHARE ROW EXCLUSIVE MODE").Error; err != nil {
			return err
		}
		var count int64
		if err := tx.Model(&AuthorizationCode{}).Count(&count).Error; err != nil {
			return err
		}
		if count >= int64(maxCodes) {
			return ErrStateLimitReached
		}
		return tx.Create(code).Error
	})
}

func (r *Repository) ConsumeAuthorizationCodeBound(
	codeHash, clientID, redirectURI, resource, codeChallenge string,
	now time.Time,
) (*AuthorizationCode, error) {
	var result AuthorizationCode
	err := r.db.Transaction(func(tx *gorm.DB) error {
		query := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where(
				"code_hash = ? AND client_id = ? AND redirect_uri = ? AND resource = ? AND code_challenge = ? AND expires_at > ?",
				codeHash, clientID, redirectURI, resource, codeChallenge, now,
			).
			Limit(1).
			Find(&result)
		if query.Error != nil {
			return query.Error
		}
		if query.RowsAffected != 1 {
			return gorm.ErrRecordNotFound
		}
		return tx.Delete(&result).Error
	})
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *Repository) StoreRefreshTokenBounded(token *RefreshToken, maxTokens int, now time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("expires_at <= ?", now).Delete(&RefreshToken{}).Error; err != nil {
			return err
		}
		if err := tx.Exec("LOCK TABLE mcp_oauth_refresh_tokens IN SHARE ROW EXCLUSIVE MODE").Error; err != nil {
			return err
		}
		var count int64
		if err := tx.Model(&RefreshToken{}).Count(&count).Error; err != nil {
			return err
		}
		if count >= int64(maxTokens) {
			return ErrStateLimitReached
		}
		return tx.Create(token).Error
	})
}

func (r *Repository) GetRefreshToken(tokenHash string) (*RefreshToken, error) {
	var result RefreshToken
	query := r.db.Where("token_hash = ?", tokenHash).Limit(1).Find(&result)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected != 1 {
		return nil, gorm.ErrRecordNotFound
	}
	return &result, nil
}

func (r *Repository) RotateRefreshToken(oldHash string, expectedClientID string, expectedUserID uint, replacement *RefreshToken) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var current RefreshToken
		query := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("token_hash = ? AND client_id = ? AND user_id = ?", oldHash, expectedClientID, expectedUserID).
			Limit(1).Find(&current)
		if query.Error != nil {
			return query.Error
		}
		if query.RowsAffected != 1 {
			return gorm.ErrRecordNotFound
		}
		if err := tx.Delete(&current).Error; err != nil {
			return err
		}
		return tx.Create(replacement).Error
	})
}
