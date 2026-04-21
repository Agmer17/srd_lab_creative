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

CREATE TABLE orders (
    id             UUID               PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID               NOT NULL REFERENCES users(id)    ON DELETE RESTRICT,
    product_id     UUID               NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    ordered_price  DECIMAL(15, 2)     NOT NULL,
    status         order_status_enum  NOT NULL DEFAULT 'pending',
    created_at     TIMESTAMPTZ          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMPTZ          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at     TIMESTAMPTZ          NULL  
);

CREATE TABLE payments (
    id                   UUID                  PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id             UUID                  NOT NULL UNIQUE REFERENCES orders(id) ON DELETE RESTRICT,
    method               VARCHAR(100),
    status               payment_status_enum   NOT NULL DEFAULT 'unpaid',
    amount               DECIMAL(15, 2)        NOT NULL,
    payment_gateway_ref  VARCHAR(255),
    paid_at              TIMESTAMPTZ             NULL,
    created_at           TIMESTAMPTZ             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at           TIMESTAMPTZ             NULL  
);

CREATE TABLE projects (
    id                 UUID                  PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id           UUID                  NOT NULL UNIQUE REFERENCES orders(id) ON DELETE RESTRICT,
    name               VARCHAR(255)          NOT NULL,
    description        TEXT,
    status             project_status_enum   NOT NULL DEFAULT 'in_progress',
    allowed_revision_count     INT                   NOT NULL DEFAULT 3,
    actual_start_date  TIMESTAMPTZ             NULL,
    end_date           TIMESTAMPTZ             NULL,
    created_at         TIMESTAMPTZ             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMPTZ             NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE project_members (
    id         UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID      NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id    UUID      NOT NULL REFERENCES users(id)    ON DELETE CASCADE,
    role_id    UUID      NOT NULL REFERENCES roles(id)    ON DELETE RESTRICT,
    is_owner   BOOLEAN   NOT NULL DEFAULT FALSE,
    joined_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at    TIMESTAMPTZ NULL,  
    UNIQUE (project_id, user_id)
);

CREATE TABLE progresses (
    id           UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id   UUID           NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    project_member_id UUID      REFERENCES project_members(id) ON DELETE SET NULL,
    title        VARCHAR(255)   NOT NULL,
    weight       DECIMAL(5, 2)  NOT NULL, 
    is_completed BOOLEAN        NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ      NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE revision_requests (
    id         UUID                  PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID                  NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title      VARCHAR(255)          NOT NULL,
    reason     TEXT                  NOT NULL,
    status     revision_status_enum  NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ             NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE chatrooms (
    id              UUID                PRIMARY KEY DEFAULT gen_random_uuid(),
    type            chatroom_type_enum  NOT NULL,
    project_id      UUID                REFERENCES projects(id) ON DELETE CASCADE,
    participant_key VARCHAR(73)         NULL, 

    created_at      TIMESTAMPTZ           NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_project_chatroom
        CHECK (
            (type = 'project' AND project_id IS NOT NULL AND participant_key IS NULL) OR
            (type = 'personal' AND project_id IS NULL AND participant_key IS NOT NULL)
        ),

    CONSTRAINT uq_project_chatroom  UNIQUE (project_id),
    CONSTRAINT uq_participant_key   UNIQUE (participant_key)
);

CREATE TABLE chatroom_participants (
    id          UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    chatroom_id UUID      NOT NULL REFERENCES chatrooms(id) ON DELETE CASCADE,
    user_id     UUID      NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
    joined_at   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at     TIMESTAMPTZ NULL,

    UNIQUE (chatroom_id, user_id)
);

CREATE TABLE chats (
    id          UUID      PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id     UUID      NOT NULL REFERENCES chatrooms(id) ON DELETE CASCADE,
    sender_id   UUID      REFERENCES users(id) ON DELETE SET NULL,
    text        TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMPTZ NULL
);

CREATE TABLE chat_medias (
    id          UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id     UUID          NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    file_name   VARCHAR(255)  NOT NULL,
    media_type  VARCHAR(100),
    size        BIGINT,
    is_one_time BOOLEAN       NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

