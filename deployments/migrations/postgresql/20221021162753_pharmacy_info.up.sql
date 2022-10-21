CREATE TABLE IF NOT EXISTS public.pharmacy_info
(
    uid uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    pharmacy_id uuid NOT NULL,
    day BIGINT,
    open_hour DOUBLE PRECISION,
    close_hour DOUBLE PRECISION,
    CONSTRAINT fk_pharmacy FOREIGN KEY(pharmacy_id) REFERENCES pharmacy(uid) ON DELETE CASCADE
);
