CREATE TABLE app_config_variables
(
    key           VARCHAR(100) NOT NULL PRIMARY KEY,
    value         TEXT NOT NULL,
    type          VARCHAR(20) NOT NULL,
    is_public     BOOLEAN DEFAULT FALSE NOT NULL,
    is_internal   BOOLEAN DEFAULT FALSE NOT NULL,
    default_value TEXT
);

CREATE TABLE user_groups
(
    id            UUID NOT NULL PRIMARY KEY,
    created_at    TIMESTAMPTZ,
    friendly_name VARCHAR(255) NOT NULL,
    name          VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE users
(
    id         UUID NOT NULL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    username   VARCHAR(255) NOT NULL UNIQUE,
    email      VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(100),
    last_name  VARCHAR(100),
    is_admin   BOOLEAN DEFAULT FALSE NOT NULL
);

CREATE TABLE audit_logs
(
    id         UUID NOT NULL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    event      VARCHAR(100) NOT NULL,
    ip_address INET NOT NULL,
    data       JSONB NOT NULL,
    user_id    UUID REFERENCES users ON DELETE SET NULL,
    user_agent TEXT,
    country    VARCHAR(100),
    city       VARCHAR(100)
);

CREATE TABLE custom_claims
(
    id            UUID NOT NULL PRIMARY KEY,
    created_at    TIMESTAMPTZ,
    key           VARCHAR(255) NOT NULL,
    value         TEXT NOT NULL,
    user_id       UUID REFERENCES users ON DELETE CASCADE,
    user_group_id UUID REFERENCES user_groups ON DELETE CASCADE,
    CONSTRAINT custom_claims_unique UNIQUE (key, user_id, user_group_id),
    CHECK (user_id IS NOT NULL OR user_group_id IS NOT NULL)
);

CREATE TABLE oidc_authorization_codes
(
    id                           UUID NOT NULL PRIMARY KEY,
    created_at                   TIMESTAMPTZ,
    code                         VARCHAR(255) NOT NULL UNIQUE,
    scope                        TEXT NOT NULL,
    nonce                        VARCHAR(255),
    expires_at                   TIMESTAMPTZ NOT NULL,
    user_id                      UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    client_id                    UUID NOT NULL,
    code_challenge               VARCHAR(255),
    code_challenge_method_sha256 BOOLEAN
);

CREATE TABLE oidc_clients
(
    id            UUID NOT NULL PRIMARY KEY,
    created_at    TIMESTAMPTZ,
    name          VARCHAR(255),
    secret        TEXT,
    callback_urls JSONB,
    image_type    VARCHAR(10),
    created_by_id UUID REFERENCES users ON DELETE SET NULL,
    is_public     BOOLEAN DEFAULT FALSE
);

CREATE TABLE one_time_access_tokens
(
    id         UUID NOT NULL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    token      VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    user_id    UUID NOT NULL REFERENCES users ON DELETE CASCADE
);

CREATE TABLE user_authorized_oidc_clients
(
    scope     VARCHAR(255),
    user_id   UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    client_id UUID NOT NULL REFERENCES oidc_clients ON DELETE CASCADE,
    PRIMARY KEY (user_id, client_id)
);

CREATE TABLE user_groups_users
(
    user_id       UUID NOT NULL REFERENCES users ON DELETE CASCADE,
    user_group_id UUID NOT NULL REFERENCES user_groups ON DELETE CASCADE,
    PRIMARY KEY (user_id, user_group_id)
);

CREATE TABLE webauthn_credentials
(
    id               UUID NOT NULL PRIMARY KEY,
    created_at       TIMESTAMPTZ,
    name             VARCHAR(255) NOT NULL,
    credential_id    BYTEA NOT NULL UNIQUE,
    public_key       BYTEA NOT NULL,
    attestation_type VARCHAR(20) NOT NULL,
    transport        JSONB NOT NULL,
    user_id          UUID REFERENCES users ON DELETE CASCADE,
    backup_eligible  BOOLEAN DEFAULT FALSE NOT NULL,
    backup_state     BOOLEAN DEFAULT FALSE NOT NULL
);

CREATE TABLE webauthn_sessions
(
    id                UUID NOT NULL PRIMARY KEY,
    created_at        TIMESTAMPTZ,
    challenge         VARCHAR(255) NOT NULL UNIQUE,
    expires_at        TIMESTAMPTZ NOT NULL,
    user_verification VARCHAR(255) NOT NULL
);