ALTER TABLE oidc_authorization_codes DROP COLUMN code_challenge;
ALTER TABLE oidc_authorization_codes DROP COLUMN code_challenge_method_sha256;
ALTER TABLE oidc_clients DROP COLUMN is_public;