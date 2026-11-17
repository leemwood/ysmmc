-- 模型文件上传归属。原有已发布文件仍可读取，新提交的文件必须先上传。
CREATE TABLE IF NOT EXISTS model_uploads (
    path VARCHAR(500) PRIMARY KEY,
    user_id UUID NOT NULL,
    size BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_model_uploads_user_id ON model_uploads(user_id);
