CREATE INDEX IF NOT EXISTS idx_blog_post_filter
    ON blog.post (id, title, published, deleted);
