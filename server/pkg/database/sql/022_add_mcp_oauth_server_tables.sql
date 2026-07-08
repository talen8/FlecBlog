-- MCP embedded OAuth authorization-server state.
-- Access tokens are signed JWTs and are intentionally not persisted.

CREATE TABLE IF NOT EXISTS mcp_oauth_clients (
    client_id VARCHAR(80) PRIMARY KEY,
    client_secret_hash VARCHAR(64) NOT NULL,
    redirect_uris TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_mcp_oauth_clients_last_used_at
    ON mcp_oauth_clients (last_used_at);

CREATE TABLE IF NOT EXISTS mcp_oauth_authorization_codes (
    code_hash VARCHAR(64) PRIMARY KEY,
    client_id VARCHAR(80) NOT NULL REFERENCES mcp_oauth_clients(client_id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    redirect_uri TEXT NOT NULL,
    resource TEXT NOT NULL,
    scopes TEXT NOT NULL,
    code_challenge VARCHAR(128) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_mcp_oauth_authorization_codes_expires_at
    ON mcp_oauth_authorization_codes (expires_at);

CREATE TABLE IF NOT EXISTS mcp_oauth_refresh_tokens (
    token_hash VARCHAR(64) PRIMARY KEY,
    client_id VARCHAR(80) NOT NULL REFERENCES mcp_oauth_clients(client_id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    scopes TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_mcp_oauth_refresh_tokens_expires_at
    ON mcp_oauth_refresh_tokens (expires_at);
CREATE INDEX IF NOT EXISTS idx_mcp_oauth_refresh_tokens_user_id
    ON mcp_oauth_refresh_tokens (user_id);
