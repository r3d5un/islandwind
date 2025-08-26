CREATE TABLE IF NOT EXISTS auth.refresh_token
(
    id             UUID    DEFAULT gen_random_uuid(),
    issuer         VARCHAR(512)         NOT NULL,
    expiration     TIMESTAMPTZ          NOT NULL,
    issued_at      TIMESTAMPTZ          NOT NULL,
    invalidated    BOOLEAN DEFAULT FALSE,
    invalidated_by UUID    DEFAULT NULL NULL,
    CONSTRAINT pk_refresh_token_id PRIMARY KEY (id),
    CONSTRAINT fk_refresh_token_invalidated_by FOREIGN KEY (invalidated_by)
        REFERENCES auth.refresh_token (id)
);
