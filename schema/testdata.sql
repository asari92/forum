-- Вставка пользователей в таблицу users
INSERT INTO users (username, email, password, created) VALUES
('user1', 'user1@example.com', 'password1', '2024-08-10 14:30:45'),
('user2', 'user2@example.com', 'password2', '2024-08-10 15:00:30'),
('user3', 'user3@example.com', 'password3', '2024-08-10 15:30:15'),
('user4', 'user4@example.com', 'password4', '2024-08-10 16:00:00'),
('user5', 'user5@example.com', 'password5', '2024-08-10 16:30:45');

-- Вставка постов в таблицу posts
INSERT INTO posts (title, content, user_id, created) VALUES
('Post 1 Title', 'Post 1 Content', 1, '2024-08-10 14:31:45'),
('Post 2 Title', 'Post 2 Content', 2, '2024-08-10 15:01:30'),
('Post 3 Title', 'Post 3 Content', 3, '2024-08-10 15:31:15'),
('Post 4 Title', 'Post 4 Content', 4, '2024-08-10 16:01:00'),
('Post 5 Title', 'Post 5 Content', 5, '2024-08-10 16:31:45');

-- Вставка комментариев в таблицу comments
INSERT INTO comments (post_id, user_id, content, created) VALUES
(1, 2, 'Comment 1 on Post 1', '2024-08-10 14:32:45'),
(2, 3, 'Comment 2 on Post 2', '2024-08-10 15:02:30'),
(3, 4, 'Comment 3 on Post 3', '2024-08-10 15:32:15'),
(4, 5, 'Comment 4 on Post 4', '2024-08-10 16:02:00'),
(5, 1, 'Comment 5 on Post 5', '2024-08-10 16:32:45');

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
('Other');
('Technology'),
('Science'),
('Art'),
('Music'),

-- Вставка категорий для постов в таблицу post_categories
INSERT INTO post_categories (post_id, category_id) VALUES
(1, 1),
(2, 2),
(3, 3),
(4, 4),
(5, 5);