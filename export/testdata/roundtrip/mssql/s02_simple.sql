CREATE TABLE dbo.Products (
    ProductID int IDENTITY(1,1) PRIMARY KEY,
    Name nvarchar(50) NOT NULL,
    Price money NOT NULL DEFAULT 0,
    Discontinued bit NOT NULL DEFAULT 0,
    RowGuid uniqueidentifier NOT NULL DEFAULT NEWID(),
    Notes nvarchar(max) NULL
);
CREATE UNIQUE NONCLUSTERED INDEX IX_Products_Name ON dbo.Products (Name);
CREATE TABLE sales.Reviews (
    ReviewID int IDENTITY PRIMARY KEY,
    ProductID int NOT NULL FOREIGN KEY REFERENCES dbo.Products(ProductID),
    Rating tinyint NOT NULL CHECK (Rating BETWEEN 1 AND 5),
    Body nvarchar(2000)
);
