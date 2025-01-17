CREATE TABLE IF NOT EXISTS Users(
    ID varchar(36) NOT NULL PRIMARY KEY,
    Name text NOT NULL,
    Password text NOT NULL,
    Email text NOT NULL,
    Age integer,
    DateRegistration timestamp NOT NULL default NOW()
    UNIQUE (email)
    );

CREATE TABLE IF NOT EXISTS Books(
    ID varchar(36) NOT NULL PRIMARY KEY,
    Name text NOT NULL,
    Description text,
    Author text NOT NULL,
    DateWriting timestamp
);