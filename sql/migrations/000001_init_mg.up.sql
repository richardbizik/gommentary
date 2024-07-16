CREATE TABLE subject (
    id text NOT NULL PRIMARY KEY,
    name text,
    created timestamp NOT NULL
);

CREATE TABLE comment (
    id text NOT NULL PRIMARY KEY,
    position int NOT NULL,
    subject text NOT NULL,
    author text,
    date timestamp NOT NULL,
    text text NOT NULL,
    FOREIGN KEY (subject) REFERENCES subject (id)
);

CREATE TABLE comment_edit (
    comment text NOT NULL,
    date timestamp NOT NULL,
    old_text text NOT NULL,
    FOREIGN KEY (comment) REFERENCES comment (id)
);
