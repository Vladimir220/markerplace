BEGIN;

CREATE TABLE IF NOT EXISTS users (
    login VARCHAR(50) PRIMARY KEY,
    group_name VARCHAR(50),
    password VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS announcements (
    id SERIAL PRIMARY KEY,

    title VARCHAR(200) NOT NULL,
    body TEXT NOT NULL,
    pic_link TEXT NOT NULL,
    author_login VARCHAR(50) NOT NULL,
    price INT,
    date TIMESTAMP

    FOREIGN KEY (author_login) REFERENCES users(login) ON UPDATE CASCADE
);

COMMIT;