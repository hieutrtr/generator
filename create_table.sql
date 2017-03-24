# psql -U username -d myDataBase -a -f myInsertFile
CREATE TABLE users (
    user_id   serial,
    name      varchar,
    age       smallint,
    friends         int,
    salary   money,
    ipv4        inet,
    metadata         jsonb,
    CONSTRAINT user_id PRIMARY KEY(user_id)
);
