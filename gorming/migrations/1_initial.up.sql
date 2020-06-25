CREATE TABLE IF NOT EXISTS persons (
    name text PRIMARY KEY,
    surname text
);

CREATE TABLE IF NOT EXISTS pets (
    name text,
    kind text,
    owner_name text REFERENCES persons (name) ON DELETE CASCADE
);