-- MDB generated migration file
-- rambler up

DROP TABLE IF EXISTS blogs;
CREATE TABLE blogs (
  id   BIGSERIAL PRIMARY KEY,
  name VARCHAR(30) UNIQUE  NOT NULL,
  url  VARCHAR(100) UNIQUE NOT NULL
);

DROP TABLE IF EXISTS blog_posts;
CREATE TABLE blog_posts (
  id         BIGSERIAL PRIMARY KEY,
  blog_id    BIGINT REFERENCES blogs (id)                  NOT NULL,
  wp_id      BIGINT                                        NOT NULL,
  title      TEXT                                          NOT NULL,
  content    TEXT                                          NOT NULL,
  posted_at  TIMESTAMP WITHOUT TIME ZONE                   NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now_utc()    NOT NULL
);

CREATE INDEX IF NOT EXISTS blog_posts_created_at_idx
  ON blog_posts
  USING BTREE (created_at);

CREATE INDEX IF NOT EXISTS blog_posts_blog_id_posted_at_idx
  ON blog_posts
  USING BTREE (blog_id, posted_at);

CREATE INDEX IF NOT EXISTS blog_posts_blog_id_wp_id_idx
  ON blog_posts
  USING BTREE (blog_id, wp_id);

insert into blogs (name, url) values
  ('laitman-ru', 'https://laitman.ru'),
  ('laitman-com', 'http://laitman.com/blog'),
  ('laitman-es', 'http://laitman.es'),
  ('laitman-co-il', 'http://laitman.co.il');

-- rambler down

DROP INDEX IF EXISTS blog_posts_created_at_idx;
DROP INDEX IF EXISTS blog_posts_blog_id_posted_at_idx;
DROP INDEX IF EXISTS blog_posts_blog_id_wp_id_idx;

DROP TABLE IF EXISTS blogs;
DROP TABLE IF EXISTS blog_posts;
