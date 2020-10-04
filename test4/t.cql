CREATE TABLE demo.verysimple (
    id int PRIMARY KEY,
    message text,
    WITH CLUSTERING ORDER BY (id ASC)
);
