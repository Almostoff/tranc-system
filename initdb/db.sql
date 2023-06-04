CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    balance INT
);

ALTER TABLE clients ALTER COLUMN balance SET DEFAULT 0;
INSERT INTO clients (id) VALUES (1);
INSERT INTO clients (id) VALUES (2);
INSERT INTO clients (id) VALUES (3);

CREATE TABLE invoice (
    u_id INTEGER NOT NULL REFERENCES clients (id),
    amount INT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
