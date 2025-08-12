CREATE OR REPlACE FUNCTION delete_blog_post_deleted_at_timestamp()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.deleted = TRUE;
    NEW.deleted_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
