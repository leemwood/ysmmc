-- Track successful public model-detail previews and support hot-list queries.
ALTER TABLE models ADD COLUMN IF NOT EXISTS views INTEGER NOT NULL DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_models_hot ON models (views DESC, downloads DESC, created_at DESC);
