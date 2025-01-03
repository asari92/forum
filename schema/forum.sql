-- Включить проверку внешних ключей
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS categories(
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
  name TEXT NOT NULL,
  CONSTRAINT category_name_uk UNIQUE(name)
);


CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    user_id INTEGER,          -- Владелец поста, который получает уведомление
    post_id INTEGER,          -- ID поста
    action_type TEXT,         -- 'like', 'dislike', или 'comment'
    comment_id INTEGER,
    trigger_user_id INTEGER,  -- ID пользователя, который вызвал уведомление
    created TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (trigger_user_id) REFERENCES users(id) ON DELETE CASCADE

);


CREATE TABLE IF NOT EXISTS comment_reactions(
  user_id INTEGER NOT NULL,
  comment_id INTEGER NOT NULL,
  is_like BOOLEAN NOT NULL,
  PRIMARY KEY(user_id, comment_id),
  CONSTRAINT comments_comment_reactions
    FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE Cascade
      ON UPDATE No action,
  CONSTRAINT users_comment_reactions
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE No action
      ON UPDATE No action
);

CREATE TABLE IF NOT EXISTS post_images (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  post_id INTEGER NOT NULL,
  image_url TEXT NOT NULL,
  FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE Cascade
);

CREATE TABLE IF NOT EXISTS comments(
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  post_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  content TEXT NOT NULL,
  created TEXT NOT NULL,
  CONSTRAINT users_comments
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE No action
      ON UPDATE No action,
  CONSTRAINT posts_comments
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE Cascade
      ON UPDATE No action
);

CREATE INDEX IF NOT EXISTS comments_idx_post_id ON comments(post_id);


CREATE TABLE IF NOT EXISTS post_categories(
  category_id INTEGER NOT NULL,
  post_id INTEGER NOT NULL,
  PRIMARY KEY(post_id, category_id),
  CONSTRAINT categories_post_categories
    FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE Cascade
      ON UPDATE No action,
  CONSTRAINT posts_post_categories
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE Cascade
      ON UPDATE No action
);

CREATE TABLE IF NOT EXISTS post_reactions(
  user_id INTEGER NOT NULL,
  post_id INTEGER NOT NULL,
  is_like BOOLEAN NOT NULL,
  PRIMARY KEY(user_id, post_id),
  CONSTRAINT posts_post_reactions
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE Cascade
      ON UPDATE No action,
  CONSTRAINT users_post_reactions
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE No action
      ON UPDATE No action
);

CREATE TABLE IF NOT EXISTS posts(
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  user_id INTEGER NOT NULL,
  created TEXT NOT NULL,
  is_approved BOOLEAN DEFAULT FALSE, -- для модерации
  CONSTRAINT users_posts
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE No action
      ON UPDATE No action
);

CREATE TABLE IF NOT EXISTS users(
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  username TEXT NOT NULL,
  email TEXT NOT NULL,
  password TEXT NOT NULL,
  role TEXT NOT NULL,
  created TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS moderation_requests (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  user_id INTEGER NOT NULL,
  created TEXT NOT NULL, 
  reason TEXT NOT NULL, 
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE

);

CREATE UNIQUE INDEX IF NOT EXISTS users_uc_email ON users(email);
CREATE UNIQUE INDEX IF NOT EXISTS users_uc_username ON users(username);


CREATE TABLE IF NOT EXISTS reports(
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  post_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  reason TEXT NOT NULL,
  created TEXT NOT NULL,
  CONSTRAINT posts_comments
  FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
    ON UPDATE No action
  CONSTRAINT users_posts
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE No action
      ON UPDATE No action
);


