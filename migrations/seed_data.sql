-- Тестовые данные для групп
INSERT INTO groups (name) VALUES 
('Muse'),
('Radiohead'),
('Nirvana'),
('The Beatles'),
('Coldplay'),
('Linkin Park'),
('Queen'),
('Pink Floyd'),
('The Rolling Stones'),
('Led Zeppelin'),
('Green Day'),
('Foo Fighters'),
('U2'),
('Metallica'),
('Arctic Monkeys'),
('Imagine Dragons'),
('Red Hot Chili Peppers'),
('The Killers'),
('Evanescence'),
('Paramore')
ON CONFLICT (name) DO NOTHING;

-- Тестовые данные для песен
INSERT INTO songs (group_id, group_name, title, lyrics, release_date, link) VALUES
(1, 'Muse', 'Supermassive Black Hole', '', '2006-05-09', NULL),
(1, 'Muse', 'Hysteria', '', '2003-12-01', NULL),
(2, 'Radiohead', 'Creep', '', '1992-09-21', NULL),
(2, 'Radiohead', 'Karma Police', '', '1997-08-25', NULL),
(3, 'Nirvana', 'Smells Like Teen Spirit', '', '1991-09-10', NULL),
(3, 'Nirvana', 'Come as You Are', '', '1991-03-02', NULL),
(4, 'The Beatles', 'Hey Jude', '', '1968-08-26', NULL),
(4, 'The Beatles', 'Let It Be', '', '1970-03-06', NULL),
(5, 'Coldplay', 'Fix You', '', '2005-09-05', NULL),
(5, 'Coldplay', 'Viva La Vida', '', '2008-05-25', NULL),
(6, 'Linkin Park', 'Numb', '', '2003-03-25', NULL),
(6, 'Linkin Park', 'In the End', '', '2000-10-24', NULL),
(7, 'Queen', 'Bohemian Rhapsody', '', '1975-10-31', NULL),
(7, 'Queen', 'We Will Rock You', '', '1977-10-07', NULL),
(8, 'Pink Floyd', 'Comfortably Numb', '', '1979-11-30', NULL),
(8, 'Pink Floyd', 'Wish You Were Here', '', '1975-09-12', NULL),
(9, 'The Rolling Stones', 'Paint It Black', '', '1966-05-06', NULL),
(9, 'The Rolling Stones', 'Angie', '', '1973-08-20', NULL),
(10, 'Led Zeppelin', 'Stairway to Heaven', '', '1971-11-08', NULL),
(10, 'Led Zeppelin', 'Kashmir', '', '1975-02-24', NULL)
ON CONFLICT (title) DO NOTHING;

