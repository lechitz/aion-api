CREATE TABLE IF NOT EXISTS aion_api.user_preferences (
    user_id          BIGINT PRIMARY KEY REFERENCES aion_api.users(user_id) ON DELETE CASCADE,
    theme_preset     VARCHAR(32)  NOT NULL DEFAULT 'midnight',
    theme_mode       VARCHAR(16)  NOT NULL DEFAULT 'system',
    compact_mode     BOOLEAN      NOT NULL DEFAULT false,
    reduced_motion   BOOLEAN      NOT NULL DEFAULT false,
    custom_overrides JSONB        NOT NULL DEFAULT '{}',
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT now()
);
