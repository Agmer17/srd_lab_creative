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

CREATE TABLE categories (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    slug varchar(255) not null unique, 
    description TEXT
);

CREATE TABLE products (
    id            UUID                 PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(255)         NOT NULL,
    slug          VARCHAR(255)         NOT NULL UNIQUE,
    description   TEXT,
    price         DECIMAL(15, 2)       NOT NULL,
    status        product_status_enum  NOT NULL DEFAULT 'active',
    is_featured   BOOLEAN              NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMPTZ            NULL
);

CREATE TABLE product_images (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id  UUID        NOT NULL,
    image_url   TEXT        NOT NULL,
    is_primary  BOOLEAN     NOT NULL DEFAULT FALSE,
    sort_order  INT         NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ   NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_product
        FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE CASCADE
);

CREATE TABLE product_categories (
    product_id  UUID REFERENCES products(id)   ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, category_id)
);


