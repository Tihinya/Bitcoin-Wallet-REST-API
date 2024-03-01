CREATE TABLE IF NOT EXISTS transactions (
    transaction_id character varying UNIQUE NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    amount double precision NOT NULL,
    spent BOOLEAN NOT NULL,
    created_at timestamp without time zone NOT NULL
);