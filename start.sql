-- Użyj bazy danych db
\c db

CREATE TABLE IF NOT EXISTS userdata (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    nickname VARCHAR(64),
    password VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS post (
    id SERIAL PRIMARY KEY,
    text VARCHAR(1500) NOT NULL
);

CREATE TABLE IF NOT EXISTS tag (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS tagposts (
    tag INT,
    posts JSON NOT NULL,
    FOREIGN KEY (tag) REFERENCES tag(id)
    );

CREATE TABLE IF NOT EXISTS profileposts (
    userid INT,
    posts JSON NOT NULL,
    FOREIGN KEY (userid) REFERENCES userdata(id)
);

-- Tabela "user"
INSERT INTO "userdata" (email, nickname, password) VALUES 
    ('1234', 'szop', '1234');

-- Tabela "post"
INSERT INTO post (text) VALUES 
    ('Jestę na mainie!!!'),
    ('Widzę!!!'),
    ('Dzisiaj byłem na spacerze po parku i zobaczyłem piękne kwiaty.'),
    ('Zaczynam nowy kurs programowania. Jestem bardzo podekscytowany.'),
    ('Właśnie wróciłem z podróży do Hiszpanii. Zobaczyłem tam wiele wspaniałych miejsc.'),
    ('Obejrzałem wczoraj nowy film na Netflixie. Naprawdę mi się podobał.'),
    ('Dziś zacząłem nową książkę. Mam nadzieję, że będzie równie dobra jak poprzednia.'),
    ('Spotkałem wczoraj starego przyjaciela. Było miło porozmawiać z nim po latach.'),
    ('Dzisiaj świętujemy urodziny mojego brata. Będzie dużo jedzenia i zabawy.'),
    ('W końcu udało mi się zrealizować jeden z moich życiowych celów. Jestem bardzo dumny z siebie.'),
    ('Planuję wkrótce zacząć nowy projekt. Mam wiele pomysłów, którymi chcę się podzielić.'),
    ('Dzisiaj jest piękna pogoda. Idealna na piknik w parku z przyjaciółmi.');

-- Tabela "tag"
INSERT INTO tag (name) VALUES 
    ('main'),
    ('programowanie'),
    ('technologia'),
    ('zdrowie'),
    ('sport'),
    ('moda'),
    ('kuchnia'),
    ('podróże'),
    ('film'),
    ('muzyka'),
    ('motoryzacja');

-- Tabela "tagposts"
INSERT INTO tagposts (tag, posts) VALUES 
    (1, '[1, 2, 3]'),
    (2, '[4, 5, 6]'),
    (3, '[7, 8, 9]'),
    (4, '[10, 1, 2]'),
    (5, '[3, 4, 5]'),
    (6, '[6, 7, 8]'),
    (7, '[9, 10, 1]'),
    (8, '[2, 3, 4]'),
    (9, '[5, 6, 7]'),
    (10, '[8, 9, 10]');

-- Tabela "profileposts"
INSERT INTO profileposts (userid, posts) VALUES 
    (1, '[1, 2, 3]');
