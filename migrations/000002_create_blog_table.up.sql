CREATE TABLE IF NOT EXISTS blog.post
(
    id         UUID      DEFAULT gen_random_uuid(),
    title      VARCHAR(1024)           NOT NULL,
    content    TEXT                    NOT NULL,
    published  BOOLEAN   DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted    BOOLEAN   DEFAULT FALSE NOT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL  NULL,
    CONSTRAINT pk_blog_post_id PRIMARY KEY (id),
    CONSTRAINT uq_blog_post_title UNIQUE (title)
);
