ALTER TABLE oidc_authorization_codes ADD COLUMN code_challenge TEXT;
ALTER TABLE oidc_authorization_codes ADD COLUMN code_challenge_method_sha256 NUMERIC;
ALTER TABLE oidc_clients ADD COLUMN is_public BOOLEAN DEFAULT FALSE;