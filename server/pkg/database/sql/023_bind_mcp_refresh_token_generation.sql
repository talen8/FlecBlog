-- Bind embedded MCP refresh tokens to the local user token revocation generation.
-- NULL is intentional: legacy rows fail closed and cannot be silently promoted.
ALTER TABLE mcp_oauth_refresh_tokens
    ADD COLUMN IF NOT EXISTS token_version BIGINT NULL;
