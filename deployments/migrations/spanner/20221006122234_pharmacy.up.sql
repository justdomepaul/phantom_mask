CREATE TABLE Pharmacy (
    UID          BYTES(16)           NOT NULL,
    Name         STRING(MAX)         NOT NULL,
    CashBalance  FLOAT64             NOT NULL,
    CreatedTime  TIMESTAMP           NOT NULL
) PRIMARY KEY(UID);
