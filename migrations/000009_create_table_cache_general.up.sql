CREATE UNLOGGED TABLE IF NOT EXISTS cache.general
(
    id         uuid                                             NOT NULL,
    data       JSONB                                            NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()                          NOT NULL,
    expires_at TIMESTAMP DEFAULT (NOW() + INTERVAL '15 minute') NOT NULL,
    CONSTRAINT pk_cache_general_id PRIMARY KEY (id)
);
