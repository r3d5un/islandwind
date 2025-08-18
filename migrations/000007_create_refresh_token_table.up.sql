CREATE TABLE IF NOT EXISTS auth.refresh_token
(
    id         UUID DEFAULT gen_random_uuid(),
    issuer     VARCHAR(512) NOT NULL,
    expiration TIMESTAMP    NOT NULL,
    issued_at  TIMESTAMP    NOT NULL,
    not_before TIMESTAmP    NOT NULL,
    CONSTRAINT pk_refresh_token_id PRIMARY KEY (id)
);
