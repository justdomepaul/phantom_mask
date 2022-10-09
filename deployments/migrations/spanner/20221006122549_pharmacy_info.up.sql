CREATE TABLE PharmacyInfo (
    UID          BYTES(16)           NOT NULL,
    Day          INT64,
    OpenHour     FLOAT64,
    CloseHour    FLOAT64
) PRIMARY KEY(UID, Day ASC),
  INTERLEAVE IN PARENT Pharmacy ON DELETE CASCADE;
