CREATE EXTENSION pgcrypto;

CREATE TABLE public.transaction (
    id UUID PRIMARY KEY default (gen_random_uuid()),
    type VARCHAR(100) NOT NULL,
    amount INT NOT NULL,
    comment VARCHAR(500)
);