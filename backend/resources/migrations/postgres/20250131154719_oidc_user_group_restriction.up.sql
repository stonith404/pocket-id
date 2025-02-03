CREATE TABLE oidc_clients_allowed_user_groups
(
    user_group_id  UUID NOT NULL REFERENCES user_groups ON DELETE CASCADE,
    oidc_client_id UUID NOT NULL REFERENCES oidc_clients ON DELETE CASCADE,
    PRIMARY KEY (oidc_client_id, user_group_id)
);


