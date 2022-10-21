CREATE TABLE IF NOT EXISTS public.pharmacy
(
    uid uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    name VARCHAR(256) NOT NULL,
    Cash_balance DOUBLE PRECISION NOT NULL,
    created_time timestamp NOT NULL DEFAULT now()
);