DROP TRIGGER IF EXISTS trigger_blog_post_timestamp_on_delete ON blog.post;

CREATE TRIGGER trigger_blog_post_timestamp_on_delete
    BEFORE DELETE ON blog.post
    FOR EACH ROW
EXECUTE PROCEDURE delete_blog_post_deleted_at_timestamp();
