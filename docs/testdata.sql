-- Вставка пользователей в таблицу users
INSERT INTO users (username, email, password, created_date) VALUES
('user1', 'user1@example.com', 'password1', '2024-08-10'),
('user2', 'user2@example.com', 'password2', '2024-08-10'),
('user3', 'user3@example.com', 'password3', '2024-08-10'),
('user4', 'user4@example.com', 'password4', '2024-08-10'),
('user5', 'user5@example.com', 'password5', '2024-08-10');

-- Вставка постов в таблицу posts
INSERT INTO posts (title, content, user_id, created_date) VALUES
('Post 1 Title', 'Post 1 Content', 1, '2024-08-10'),
('Post 2 Title', 'Post 2 Content', 2, '2024-08-10'),
('Post 3 Title', 'Post 3 Content', 3, '2024-08-10'),
('Post 4 Title', 'Post 4 Content', 4, '2024-08-10'),
('Post 5 Title', 'Post 5 Content', 5, '2024-08-10');

-- Вставка комментариев в таблицу comments
INSERT INTO comments (post_id, user_id, content, created_date) VALUES
(1, 2, 'Comment 1 on Post 1', '2024-08-10'),
(2, 3, 'Comment 2 on Post 2', '2024-08-10'),
(3, 4, 'Comment 3 on Post 3', '2024-08-10'),
(4, 5, 'Comment 4 on Post 4', '2024-08-10'),
(5, 1, 'Comment 5 on Post 5', '2024-08-10');

-- Вставка реакций на посты в таблицу post_reactions
INSERT INTO post_reactions (user_id, post_id, is_like) VALUES
(1, 2, 1),
(2, 3, 1),
(3, 4, 0),
(4, 5, 1),
(5, 1, 0);

-- Вставка реакций на комментарии в таблицу comment_reactions
INSERT INTO comment_reactions (user_id, comment_id, is_like) VALUES
(1, 1, 1),
(2, 2, 1),
(3, 3, 0),
(4, 4, 1),
(5, 5, 0);

-- Вставка категорий в таблицу categories
INSERT INTO categories (name) VALUES
('Technology'),
('Science'),
('Art'),
('Music'),
('Literature');

-- Вставка категорий для постов в таблицу post_categories
INSERT INTO post_categories (post_id, category_id) VALUES
(1, 1),
(2, 2),
(3, 3),
(4, 4),
(5, 5);