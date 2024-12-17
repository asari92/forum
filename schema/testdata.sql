-- Вставка пользователей в таблицу users
INSERT INTO users (username, email, password, role, created) VALUES
('admin', 'admin@mail.com', '12345678', 'admin', '2024-01-01 15:00:00'),
('moder', 'moder@mail.com', '12345678', 'moderator', '2024-01-01 15:00:00'),
('MusicLover1', 'musiclover1@example.com', 'password1', 'user', '2024-08-01 15:00:00'),
('JazzFanatic', 'jazzfan@example.com', 'password2', 'user', '2024-08-01 15:00:00'),
('VinylCollector', 'vinylcollector@example.com', 'password3', 'user', '2024-08-01 15:00:00'),
('ClassicalMaestro', 'classicalmaestro@example.com', 'password4', 'user', '2024-08-01 15:00:00'),
('PopStarFan', 'popstarfan@example.com', 'password5', 'user', '2024-08-01 15:00:00'),
('ElectronicVibes', 'electronicvibes@example.com', 'password6', 'user', '2024-08-01 15:00:00'),
('HipHopHead', 'hiphophead@example.com', 'password7', 'user', '2024-08-01 15:00:00'),
('FolkMusicLover', 'folkmusiclover@example.com', 'password8', 'user', '2024-08-01 15:00:00'),
('BluesBrother', 'bluesbrother@example.com', 'password9', 'user', '2024-08-01 15:00:00'),
('CountryRoads', 'countryroads@example.com', 'password10', 'user', '2024-08-01 15:00:00');

-- Вставка постов в таблицу posts
INSERT INTO posts (title, content, user_id, created, is_approved) VALUES
('The Evolution of Rock Music', 'Explore the history of rock music from its roots in the early 20th century blues and country music to the diverse genres we see today. Understand how rock music has been a reflection of social change and technological advances.', 1, '2024-08-01 16:00:00', 1),
('The Impact of The Beatles', 'Analyze the profound impact The Beatles had on the music world, from their innovative recording techniques to their role in the socio-political movements of the 1960s. Discover how they influenced countless musicians and genres that followed.', 1, '2024-08-01 16:30:00', 1),
('The Roots of Jazz and Its Evolution', 'Delve into the origins of jazz in New Orleans and trace its evolution through the big band era to modern jazz. Explore key figures like Louis Armstrong, Duke Ellington, and Miles Davis.', 2, '2024-08-02 16:30:00', 1),
('The Rise of Electronic Music', 'From the avant-garde experiments of the 1970s to today''s dance music explosion, this post explores how electronic music has grown and the technology that fuels it.', 5, '2024-08-04 16:30:00', 1),
('Hip-Hop: More Than Music', 'Examine hip-hop as a cultural movement that encompasses not only music but also dance, art, and social activism. Understand the stories of its pioneers and the genres it has inspired.', 6, '2024-08-04 17:30:00', 1),
('The Folk Music Revival', 'Look at the 1960s folk revival in America with artists like Bob Dylan and Joan Baez who used their music for political protest and cultural expression.', 7, '2024-08-04 17:40:00', 1),
('The Global Impact of Reggae Music', 'Explore how reggae music from Jamaica became a global phenomenon influencing music, fashion, and politics worldwide.', 8, '2024-08-05 17:40:00', 1),
('The Power of Classical Music in Movies', 'Explore how classical music enhances movie storytelling, from dramatic opera to subtle symphonies that intensify the emotional landscape of films.', 3, '2024-08-06 17:40:00', 1),
('The Symbiosis of Jazz and Classical Music', 'Discover how jazz and classical music have borrowed from each other, resulting in rich, hybrid musical forms that challenge and expand the boundaries of both genres.', 2, '2024-08-07 17:40:00', 1),
('Innovations in Modern Rock', 'Discuss how modern rock bands have incorporated technology and diverse musical influences to redefine what rock music can be in the 21st century.', 1, '2024-08-08 11:40:00', 1),
('The Enduring Popularity of K-Pop', 'Analyze the factors behind K-Pop''s worldwide popularity, from its catchy melodies and high-production music videos to the intensive training of its stars.', 4, '2024-08-09 11:40:00', 1),
('The Evolution of Jazz', 'Jazz music has undergone numerous transformations since its inception in the early 20th century. This post explores the key periods and changes in jazz, highlighting major artists and their influential works.', 2, '2024-08-10 15:02:00', 1),
('Exploring the Genres of Electronic Music', 'Electronic music encompasses a broad range of styles and genres. This post examines the various sub-genres of electronic music, their characteristics, and their cultural impacts.', 5, '2024-08-10 15:04:00', 1),
('History of Hip Hop', 'Hip hop started in the streets of New York in the 1970s and has grown to become a global phenomenon. This post delves into its origins, cultural significance, and evolution over the decades.', 6, '2024-08-10 15:06:00', 1),
('Folk Music Around the World', 'Folk music tells the stories of people and traditions. This post explores different folk music traditions from various cultures around the world.', 7, '2024-08-10 15:08:00', 1),
('Blues Music and Its Influence', 'Blues music has deeply influenced many other music genres. This post discusses the origins of the blues and its impact on the development of rock and jazz music.', 8, '2024-08-10 15:10:00', 1),
('Classical Music Composers of the 19th Century', 'The 19th century was a pivotal era for classical music. This post profiles some of the most significant composers from this period and their contributions to music.', 3, '2024-08-10 15:12:00', 1),
('The Impact of Music Festivals on Popular Music', 'Music festivals have become a vital part of music culture. This post examines how they influence popular music trends and artists’ careers.', 4, '2024-08-10 15:14:00', 1),
('The Cultural Impact of The Rolling Stones', 'The Rolling Stones played a crucial role in the development of rock music. This post explores their cultural impact and how they shaped music and society.', 1, '2024-08-10 15:16:00', 1),
('The Revival of Vinyl Records', 'Vinyl records are experiencing a resurgence in popularity. This post explores why music lovers are returning to this old format and its impact on the music industry.', 1, '2024-08-10 15:18:00', 1);

-- Вставка комментариев в таблицу comments
INSERT INTO comments (post_id, user_id, content, created) VALUES
(1, 2, 'Totally agree with your points!', '2024-08-10 15:00:00'),
(3, 5, 'Jazz is such a vibrant and expressive genre!', '2024-08-10 12:40:00'),
(4, 2, 'Cool post!', '2024-08-11 15:00:10'),
(5, 1, 'Classical music soothes the soul like nothing else.', '2024-08-10 12:50:00'),
(7, 3, 'Pop music really captures the spirit of its time!', '2024-08-10 13:50:00');

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
('Others'), ('Rock'), ('Jazz'), ('Classical'), ('Pop'),
('Electronic'), ('Hip-Hop'), ('Folk'), ('Blues'), ('Country');

-- Вставка категорий для постов в таблицу post_categories
INSERT INTO post_categories (category_id, post_id) VALUES
(2, 1), -- Rock music
(2, 2), -- The Beatles (Rock)
(3, 3), -- Jazz
(5, 4), -- Electronic
(6, 5), -- Hip-Hop
(7, 6), -- Folk
(1, 7), -- Others
(4, 8), -- Classical
(4, 9), -- Classical
(3, 9), -- Jazz
(2, 10), -- Rock
(5, 11), -- Pop
(1, 11), -- Others
(3, 12), -- Jazz
(5, 13), -- Electronic
(6, 14), -- Hip-Hop
(7, 15), -- Folk
(9, 16), -- Blues
(4, 17), -- Classical
(5, 18), -- Pop
(2, 19), -- Rock
(2, 20); -- Rock