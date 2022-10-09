CREATE TABLE Product (
    UID          BYTES(16)           NOT NULL,
    ProductID    BYTES(16)           NOT NULL,
    Name         STRING(MAX)         NOT NULL,
    Price        FLOAT64             NOT NULL,
    CreatedTime  TIMESTAMP           NOT NULL
) PRIMARY KEY(UID, ProductID),
  INTERLEAVE IN PARENT Pharmacy ON DELETE CASCADE;
