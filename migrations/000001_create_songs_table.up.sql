CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    lyrics TEXT,
    release_date DATE,
    link VARCHAR(255),
    UNIQUE (group_name, title)
);

INSERT INTO songs (group_name, title, lyrics) VALUES
('Muse', 'Supermassive Black Hole', 
 'Ooh baby, don''t you know I suffer?
 Ooh baby, can you hear me moan?
 You caught me under false pretenses
 How long before you let me go?

 Ooh
 You set my soul alight
 Ooh
 You set my soul alight

 Glaciers melting in the dead of night (ooh)
 And the superstars sucked into the supermassive (you set my soul alight)
 Glaciers melting in the dead of night
 And the superstars sucked into the (you set my soul)
 (Into the supermassive)

 I thought I was a fool for no one
 Ooh baby, I''m a fool for you
 You''re the queen of the superficial
 And how long before you tell the truth?

 Ooh
 You set my soul alight
 Ooh
 You set my soul alight

 Glaciers melting in the dead of night (ooh)
 And the superstars sucked into the supermassive (you set my soul alight)
 Glaciers melting in the dead of night
 And the superstars sucked into the (you set my soul)
 (Into the supermassive)

 Supermassive black hole
 Supermassive black hole
 Supermassive black hole
 Supermassive black hole
 Glaciers melting in the dead of night
 And the superstars sucked into the supermassive
 Glaciers melting in the dead of night
 And the superstars sucked into the supermassive
 Glaciers melting in the dead of night (ooh)
 And the superstars sucked into the supermassive (you set my soul alight)
 Glaciers melting in the dead of night
 And the superstars sucked into the (you set my soul)
 (Into the supermassive)

 Supermassive black hole
 Supermassive black hole
 Supermassive black hole
 Supermassive black hole');


INSERT INTO songs (group_name, title) VALUES
('Radiohead', 'Creep'),
('Nirvana', 'Smells Like Teen Spirit'),
('The Beatles', 'Hey Jude'),
('Muse', 'Hysteria'),
('Coldplay', 'Fix You'),
('Linkin Park', 'Numb'),
('Linkin Park', 'In the End'),
('Linkin Park', 'Somewhere I Belong'),
('Linkin Park', 'Crawling'),
('Linkin Park', 'Breaking the Habit'),
('Coldplay', 'Yellow'),
('Green Day', 'Boulevard of Broken Dreams'),
('The Killers', 'Mr. Brightside'),
('Red Hot Chili Peppers', 'Californication'),
('Queen', 'Bohemian Rhapsody'),
('Adele', 'Rolling in the Deep'),
('Foo Fighters', 'Everlong'),
('Oasis', 'Wonderwall'),
('The White Stripes', 'Seven Nation Army'),
('Paramore', 'Misery Business'),
('Imagine Dragons', 'Radioactive'),
('Kings of Leon', 'Use Somebody'),
('The Lumineers', 'Ho Hey'),
('Linkin Park', 'Faint'),
('My Chemical Romance', 'Welcome to the Black Parade'),
('The Strokes', 'Last Nite'),
('Fall Out Boy', 'Sugar, We’re Goin Down'),
('Arctic Monkeys', 'Do I Wanna Know?'),
('Billie Eilish', 'Bad Guy'),
('Dua Lipa', 'Levitating'),
('Taylor Swift', 'Shake It Off'),
('Beyoncé', 'Single Ladies (Put a Ring on It)'),
('Ed Sheeran', 'Shape of You');
