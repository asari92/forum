CREATE TABLE posts(
   id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
   title TEXT NOT NULL,
   content TEXT NOT NULL,
   user_id INTEGER NOT NULL,
   created_date TEXT NOT NULL,
   CONSTRAINT id_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE users(
   id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
   username TEXT NOT NULL,
   email TEXT NOT NULL,
   password TEXT NOT NULL,
   created_date INTEGER NOT NULL
);

   CREATE UNIQUE INDEX users_uc_email ON users(email);
   
CREATE TABLE comments(
   id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
   post_id INTEGER NOT NULL,
   user_id INTEGER NOT NULL,
   content TEXT NOT NULL,
   created_date TEXT NOT NULL,
   CONSTRAINT id_post_id FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
   CONSTRAINT id_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE post_reactions(
   user_id INTEGER NOT NULL,
   post_id INTEGER NOT NULL,
   is_like BOOLEAN NOT NULL,
   PRIMARY KEY(user_id, post_id),
   CONSTRAINT id_user_id FOREIGN KEY (user_id) REFERENCES users (id),
   CONSTRAINT id_post_id FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
);

CREATE TABLE comment_reactions(
   user_id INTEGER NOT NULL,
   comment_id INTEGER NOT NULL,
   is_like BOOLEAN NOT NULL,
   PRIMARY KEY(user_id, comment_id),
   CONSTRAINT id_user_id FOREIGN KEY (user_id) REFERENCES users (id),
   CONSTRAINT id_comment_id FOREIGN KEY (comment_id) REFERENCES comments (id) ON DELETE CASCADE
);

CREATE TABLE categories(
   id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 
   name TEXT NOT NULL,
   CONSTRAINT category_name_uk UNIQUE(name) 
);

CREATE TABLE post_categories(
   category_id INTEGER NOT NULL,
   post_id INTEGER NOT NULL,
   PRIMARY KEY(post_id, category_id),
   CONSTRAINT id_post_id FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
   CONSTRAINT id_category_id FOREIGN KEY (category_id) REFERENCES categories (id)
);

-- Включить проверку внешних ключей
PRAGMA foreign_keys = ON;