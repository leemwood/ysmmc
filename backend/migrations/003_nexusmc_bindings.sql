-- 003: NexusMC OAuth 绑定表
CREATE TABLE IF NOT EXISTS nexus_mc_bindings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    sub VARCHAR(64) NOT NULL,
    uid BIGINT NOT NULL DEFAULT 0,
    username VARCHAR(100),
    slug VARCHAR(100),
    avatar TEXT,
    nexus_role VARCHAR(50),
    email VARCHAR(255),
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    scope VARCHAR(255),
    access_token TEXT,
    refresh_token TEXT,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_nexus_mc_bindings_user_id ON nexus_mc_bindings(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_nexus_mc_bindings_sub ON nexus_mc_bindings(sub);
