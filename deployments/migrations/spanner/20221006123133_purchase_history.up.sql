CREATE TABLE PurchaseHistory (
    UID               BYTES(16)           NOT NULL,
    PharmacyUID       BYTES(16)           NOT NULL,
    ProductID         BYTES(16)           NOT NULL,
    TransactionAmount FLOAT64             NOT NULL,
    TransactionDate   TIMESTAMP           NOT NULL,
    CONSTRAINT FKPurchaseHistoryPharmacyUID FOREIGN KEY (PharmacyUID) REFERENCES Pharmacy (UID),
    CONSTRAINT FKPurchaseHistoryMaskUID FOREIGN KEY (ProductID) REFERENCES Product (ProductID)
) PRIMARY KEY(UID, TransactionDate),
  INTERLEAVE IN PARENT User ON DELETE CASCADE;
