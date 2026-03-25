CREATE TABLE tasks (
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title       VARCHAR(100) NOT NULL,
    description TEXT,
    done        BOOLEAN NOT NULL DEFAULT FALSE
);
