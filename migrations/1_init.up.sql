CREATE TABLE IF NOT EXISTS Users(
    ID varchar(36) not null primary key,
    Name text not null,
    Password text not null,
    Email text not null,
    Age integer,
    DateRegistration timestamp not null default now(),
    unique (email)
);

CREATE TABLE IF NOT EXISTS Books(
    ID varchar(36) not null primary key,
    Name text not null,
    Description text,
    Author text not null,
    DateWriting timestamp,
    Deleted boolean not null default false
);