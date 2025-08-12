DROP TRIGGER IF EXISTS trigger_blog_post_timestamp_on_insert ON blog.post;
DROP TRIGGER IF EXISTS trigger_blog_post_timestamp_on_update ON blog.post;

CREATE TRIGGER trigger_blog_post_timestamp_on_insert
    BEFORE INSERT ON blog.post
    FOR EACH ROW
EXECUTE PROCEDURE update_blog_post_updated_at_timestamp();

CREATE TRIGGER trigger_blog_post_timestamp_on_update
    BEFORE UPDATE ON blog.post
    FOR EACH ROW
EXECUTE PROCEDURE update_blog_post_updated_at_timestamp();
