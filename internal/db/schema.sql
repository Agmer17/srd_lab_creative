CREATE TABLE roles (
    id   UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL unique,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id               UUID               PRIMARY KEY DEFAULT gen_random_uuid(),
    global_role      global_role_enum   NOT NULL DEFAULT 'USER',
    full_name        VARCHAR(255)       NOT NULL,
    email            VARCHAR(255)       NOT NULL,
    phone_number     VARCHAR(50),
    profile_picture  TEXT,
    gender           user_gender not null default 'male',
    provider         VARCHAR(50)        NOT NULL,
    provider_user_id VARCHAR(255)       NOT NULL,
    created_at       TIMESTAMPTZ          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMPTZ          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at       TIMESTAMPTZ          NULL
);