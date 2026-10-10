CREATE TABLE "Customer" (
    "CustomerId" SERIAL NOT NULL,
    "FirstName" VARCHAR(40) NOT NULL,
    "LastName" VARCHAR(20) NOT NULL,
    "Company" VARCHAR(80),
    "SupportRepId" INT,
    CONSTRAINT "PK_Customer" PRIMARY KEY  ("CustomerId")
);
CREATE TABLE "Employee" (
    "EmployeeId" INT NOT NULL,
    "LastName" VARCHAR(20) NOT NULL,
    "ReportsTo" INT,
    "HireDate" TIMESTAMP,
    CONSTRAINT "PK_Employee" PRIMARY KEY  ("EmployeeId")
);
ALTER TABLE "Customer" ADD CONSTRAINT "FK_CustomerSupportRepId"
    FOREIGN KEY ("SupportRepId") REFERENCES "Employee" ("EmployeeId") ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE "Employee" ADD CONSTRAINT "FK_EmployeeReportsTo"
    FOREIGN KEY ("ReportsTo") REFERENCES "Employee" ("EmployeeId") ON DELETE NO ACTION ON UPDATE NO ACTION;
CREATE INDEX "IFK_CustomerSupportRepId" ON "Customer" ("SupportRepId");
