CREATE TABLE IF NOT EXISTS posts(
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  user_id INTEGER NOT NULL,
  created TEXT NOT NULL,
  CONSTRAINT users_posts
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE No action
      ON UPDATE No action
);

CREATE TABLE IF NOT EXISTS users(
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  username TEXT NOT NULL,
  email TEXT NOT NULL,
  password TEXT NOT NULL,
  created TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS users_uc_email ON users(email);

INSERT INTO users (username, email, password, created) VALUES (
    'Alice Jones',
    'alice@example.com',
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
    '2022-01-01 10:00:00'
);