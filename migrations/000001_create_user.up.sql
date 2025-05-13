CREATE TABLE Users
(
    guid     UUID PRIMARY KEY,
    username varchar(30) UNIQUE NOT NULL,
    password char(60)
);