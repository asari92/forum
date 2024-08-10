CREATE TABLE categories(
id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, name TEXT NOT NULL,
  CONSTRAINT category_name_uk UNIQUE(name)
);

 

CREATE TABLE comment_reactions(
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

 

CREATE TABLE comments(
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  post_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  content TEXT NOT NULL,
  created_date TEXT NOT NULL,
  CONSTRAINT users_comments
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE No action
      ON UPDATE No action,
  CONSTRAINT posts_comments
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE Cascade
      ON UPDATE No action
);

 

CREATE TABLE post_categories(
  category_id INTEGER NOT NULL,
  post_id INTEGER NOT NULL,
  PRIMARY KEY(post_id, category_id),
  CONSTRAINT categories_post_categories
    FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE No action
      ON UPDATE No action,
  CONSTRAINT posts_post_categories
    FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE Cascade
      ON UPDATE No action
);

 

CREATE TABLE post_reactions(
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

 

CREATE TABLE posts(
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  user_id INTEGER NOT NULL,
  created_date TEXT NOT NULL,
  CONSTRAINT users_posts
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE No action
      ON UPDATE No action
);

 

CREATE TABLE users(
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  username TEXT NOT NULL,
  email TEXT NOT NULL,
  password TEXT NOT NULL,
  created_date INTEGER NOT NULL
);

  CREATE UNIQUE INDEX users_uc_email ON users(email);
  
 

CREATE TABLE sessions(
  token VARCHAR NOT NULL,
  data BLOB NOT NULL,
  expiry TEXT NOT NULL,
  PRIMARY KEY(token)
);

  CREATE INDEX sessions_expiry_idx ON sessions(expiry);
  