CREATE SCHEMA `humanresources`;

CREATE SCHEMA `person`;

CREATE SCHEMA `production`;

CREATE SCHEMA `purchasing`;

CREATE SCHEMA `sales`;

CREATE TABLE `humanresources`.`department` (
  `departmentid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for Department records.',
  `name` public.Name NOT NULL COMMENT 'Name of the department.',
  `groupname` public.Name NOT NULL COMMENT 'Name of the group to which the department belongs.',
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `humanresources`.`employee` (
  `businessentityid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for Employee records.  Foreign key to BusinessEntity.BusinessEntityID.',
  `nationalidnumber` character varying(15) NOT NULL COMMENT 'Unique national identification number such as a social security number.',
  `loginid` character varying(256) NOT NULL COMMENT 'Network login.',
  `jobtitle` character varying(50) NOT NULL COMMENT 'Work title such as Buyer or Sales Representative.',
  `birthdate` date NOT NULL COMMENT 'Date of birth.',
  `maritalstatus` character(1) NOT NULL COMMENT 'M = Married, S = Single',
  `gender` character(1) NOT NULL COMMENT 'M = Male, F = Female',
  `hiredate` date NOT NULL COMMENT 'Employee hired on this date.',
  `salariedflag` public.Flag NOT NULL DEFAULT true COMMENT 'Job classification. 0 = Hourly, not exempt from collective bargaining. 1 = Salaried, exempt from collective bargaining.',
  `vacationhours` smallint NOT NULL DEFAULT 0 COMMENT 'Number of available vacation hours.',
  `sickleavehours` smallint NOT NULL DEFAULT 0 COMMENT 'Number of available sick leave hours.',
  `currentflag` public.Flag NOT NULL DEFAULT true COMMENT '0 = Inactive, 1 = Active',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  `organizationnode` character varying DEFAULT ('/'::charactervarying) COMMENT 'Where the employee is located in corporate hierarchy.',
  CONSTRAINT `CK_Employee_BirthDate` CHECK (((birthdate >= '1930-01-01'::date) AND (birthdate <= (now() - '18 years'::interval)))),
  CONSTRAINT `CK_Employee_Gender` CHECK ((upper((gender)::text) = ANY (ARRAY['M'::text, 'F'::text]))),
  CONSTRAINT `CK_Employee_HireDate` CHECK (((hiredate >= '1996-07-01'::date) AND (hiredate <= (now() + '1 day'::interval)))),
  CONSTRAINT `CK_Employee_MaritalStatus` CHECK ((upper((maritalstatus)::text) = ANY (ARRAY['M'::text, 'S'::text]))),
  CONSTRAINT `CK_Employee_SickLeaveHours` CHECK (((sickleavehours >= 0) AND (sickleavehours <= 120))),
  CONSTRAINT `CK_Employee_VacationHours` CHECK (((vacationhours >= '-40'::integer) AND (vacationhours <= 240)))
);

CREATE TABLE `humanresources`.`employeedepartmenthistory` (
  `businessentityid` integer NOT NULL COMMENT 'Employee identification number. Foreign key to Employee.BusinessEntityID.',
  `departmentid` smallint NOT NULL COMMENT 'Department in which the employee worked including currently. Foreign key to Department.DepartmentID.',
  `shiftid` smallint NOT NULL COMMENT 'Identifies which 8-hour shift the employee works. Foreign key to Shift.Shift.ID.',
  `startdate` date NOT NULL COMMENT 'Date the employee started work in the department.',
  `enddate` date COMMENT 'Date the employee left the department. NULL = Current department.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_EmployeeDepartmentHistory_EndDate` CHECK (((enddate >= startdate) OR (enddate IS NULL))),
  PRIMARY KEY (`businessentityid`, `startdate`, `departmentid`, `shiftid`)
);

CREATE TABLE `humanresources`.`employeepayhistory` (
  `businessentityid` integer NOT NULL COMMENT 'Employee identification number. Foreign key to Employee.BusinessEntityID.',
  `ratechangedate` timestamp NOT NULL COMMENT 'Date the change in pay is effective',
  `rate` numeric NOT NULL COMMENT 'Salary hourly rate.',
  `payfrequency` smallint NOT NULL COMMENT '1 = Salary received monthly, 2 = Salary received biweekly',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_EmployeePayHistory_PayFrequency` CHECK ((payfrequency = ANY (ARRAY[1, 2]))),
  CONSTRAINT `CK_EmployeePayHistory_Rate` CHECK (((rate >= 6.50) AND (rate <= 200.00))),
  PRIMARY KEY (`businessentityid`, `ratechangedate`)
);

CREATE TABLE `humanresources`.`jobcandidate` (
  `jobcandidateid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for JobCandidate records.',
  `businessentityid` integer COMMENT 'Employee identification number if applicant was hired. Foreign key to Employee.BusinessEntityID.',
  `resume` xml COMMENT 'RÃ©sumÃ© in XML format.',
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `humanresources`.`shift` (
  `shiftid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for Shift records.',
  `name` public.Name NOT NULL COMMENT 'Shift description.',
  `starttime` time NOT NULL COMMENT 'Shift start time.',
  `endtime` time NOT NULL COMMENT 'Shift end time.',
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `person`.`address` (
  `addressid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for Address records.',
  `addressline1` character varying(60) NOT NULL COMMENT 'First street address line.',
  `addressline2` character varying(60) COMMENT 'Second street address line.',
  `city` character varying(30) NOT NULL COMMENT 'Name of the city.',
  `stateprovinceid` integer NOT NULL COMMENT 'Unique identification number for the state or province. Foreign key to StateProvince table.',
  `postalcode` character varying(15) NOT NULL COMMENT 'Postal code for the street address.',
  `spatiallocation` character varying(44) COMMENT 'Latitude and longitude of this address.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `person`.`businessentityaddress` (
  `businessentityid` integer NOT NULL COMMENT 'Primary key. Foreign key to BusinessEntity.BusinessEntityID.',
  `addressid` integer NOT NULL COMMENT 'Primary key. Foreign key to Address.AddressID.',
  `addresstypeid` integer NOT NULL COMMENT 'Primary key. Foreign key to AddressType.AddressTypeID.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY (`businessentityid`, `addressid`, `addresstypeid`)
);

CREATE TABLE `person`.`countryregion` (
  `countryregioncode` character varying(3) PRIMARY KEY NOT NULL COMMENT 'ISO standard code for countries and regions.',
  `name` public.Name NOT NULL COMMENT 'Country or region name.',
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `person`.`emailaddress` (
  `businessentityid` integer NOT NULL COMMENT 'Primary key. Person associated with this email address.  Foreign key to Person.BusinessEntityID',
  `emailaddressid` integer NOT NULL COMMENT 'Primary key. ID of this email address.',
  `emailaddress` character varying(50) COMMENT 'E-mail address for the person.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY (`businessentityid`, `emailaddressid`)
);

CREATE TABLE `person`.`person` (
  `businessentityid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for Person records.',
  `persontype` character(2) NOT NULL COMMENT 'Primary type of person: SC = Store Contact, IN = Individual (retail) customer, SP = Sales person, EM = Employee (non-sales), VC = Vendor contact, GC = General contact',
  `namestyle` public.NameStyle NOT NULL DEFAULT false COMMENT '0 = The data in FirstName and LastName are stored in western style (first name, last name) order.  1 = Eastern style (last name, first name) order.',
  `title` character varying(8) COMMENT 'A courtesy title. For example, Mr. or Ms.',
  `firstname` public.Name NOT NULL COMMENT 'First name of the person.',
  `middlename` public.Name COMMENT 'Middle name or middle initial of the person.',
  `lastname` public.Name NOT NULL COMMENT 'Last name of the person.',
  `suffix` character varying(10) COMMENT 'Surname suffix. For example, Sr. or Jr.',
  `emailpromotion` integer NOT NULL DEFAULT 0 COMMENT '0 = Contact does not wish to receive e-mail promotions, 1 = Contact does wish to receive e-mail promotions from AdventureWorks, 2 = Contact does wish to receive e-mail promotions from AdventureWorks and selected partners.',
  `additionalcontactinfo` xml COMMENT 'Additional contact information about the person stored in xml format.',
  `demographics` xml COMMENT 'Personal information such as hobbies, and income collected from online shoppers. Used for sales analysis.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_Person_EmailPromotion` CHECK (((emailpromotion >= 0) AND (emailpromotion <= 2))),
  CONSTRAINT `CK_Person_PersonType` CHECK (((persontype IS NULL) OR (upper((persontype)::text) = ANY (ARRAY['SC'::text, 'VC'::text, 'IN'::text, 'EM'::text, 'SP'::text, 'GC'::text]))))
);

CREATE TABLE `person`.`personphone` (
  `businessentityid` integer NOT NULL COMMENT 'Business entity identification number. Foreign key to Person.BusinessEntityID.',
  `phonenumber` public.Phone NOT NULL COMMENT 'Telephone number identification number.',
  `phonenumbertypeid` integer NOT NULL COMMENT 'Kind of phone number. Foreign key to PhoneNumberType.PhoneNumberTypeID.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY (`businessentityid`, `phonenumber`, `phonenumbertypeid`)
);

CREATE TABLE `person`.`phonenumbertype` (
  `phonenumbertypeid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for telephone number type records.',
  `name` public.Name NOT NULL COMMENT 'Name of the telephone number type',
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `person`.`stateprovince` (
  `stateprovinceid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for StateProvince records.',
  `stateprovincecode` character(3) NOT NULL COMMENT 'ISO standard state or province code.',
  `countryregioncode` character varying(3) NOT NULL COMMENT 'ISO standard country or region code. Foreign key to CountryRegion.CountryRegionCode.',
  `isonlystateprovinceflag` public.Flag NOT NULL DEFAULT true COMMENT '0 = StateProvinceCode exists. 1 = StateProvinceCode unavailable, using CountryRegionCode.',
  `name` public.Name NOT NULL COMMENT 'State or province description.',
  `territoryid` integer NOT NULL COMMENT 'ID of the territory in which the state or province is located. Foreign key to SalesTerritory.SalesTerritoryID.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `person`.`addresstype` (
  `addresstypeid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for AddressType records.',
  `name` public.Name NOT NULL COMMENT 'Address type description. For example, Billing, Home, or Shipping.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `person`.`businessentity` (
  `businessentityid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for all customers, vendors, and employees.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `person`.`businessentitycontact` (
  `businessentityid` integer NOT NULL COMMENT 'Primary key. Foreign key to BusinessEntity.BusinessEntityID.',
  `personid` integer NOT NULL COMMENT 'Primary key. Foreign key to Person.BusinessEntityID.',
  `contacttypeid` integer NOT NULL COMMENT 'Primary key.  Foreign key to ContactType.ContactTypeID.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY (`businessentityid`, `personid`, `contacttypeid`)
);

CREATE TABLE `person`.`contacttype` (
  `contacttypeid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for ContactType records.',
  `name` public.Name NOT NULL COMMENT 'Contact type description.',
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `person`.`password` (
  `businessentityid` integer PRIMARY KEY NOT NULL,
  `passwordhash` character varying(128) NOT NULL COMMENT 'Password for the e-mail account.',
  `passwordsalt` character varying(10) NOT NULL COMMENT 'Random value concatenated with the password string before the password is hashed.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `production`.`billofmaterials` (
  `billofmaterialsid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for BillOfMaterials records.',
  `productassemblyid` integer COMMENT 'Parent product identification number. Foreign key to Product.ProductID.',
  `componentid` integer NOT NULL COMMENT 'Component identification number. Foreign key to Product.ProductID.',
  `startdate` timestamp NOT NULL DEFAULT (now()) COMMENT 'Date the component started being used in the assembly item.',
  `enddate` timestamp COMMENT 'Date the component stopped being used in the assembly item.',
  `unitmeasurecode` character(3) NOT NULL COMMENT 'Standard code identifying the unit of measure for the quantity.',
  `bomlevel` smallint NOT NULL COMMENT 'Indicates the depth the component is from its parent (AssemblyID).',
  `perassemblyqty` numeric(8,2) NOT NULL DEFAULT 1 COMMENT 'Quantity of the component needed to create the assembly.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_BillOfMaterials_BOMLevel` CHECK ((((productassemblyid IS NULL) AND (bomlevel = 0) AND (perassemblyqty = 1.00)) OR ((productassemblyid IS NOT NULL) AND (bomlevel >= 1)))),
  CONSTRAINT `CK_BillOfMaterials_EndDate` CHECK (((enddate > startdate) OR (enddate IS NULL))),
  CONSTRAINT `CK_BillOfMaterials_PerAssemblyQty` CHECK ((perassemblyqty >= 1.00)),
  CONSTRAINT `CK_BillOfMaterials_ProductAssemblyID` CHECK ((productassemblyid <> componentid))
);

CREATE TABLE `production`.`culture` (
  `cultureid` character(6) PRIMARY KEY NOT NULL COMMENT 'Primary key for Culture records.',
  `name` public.Name NOT NULL COMMENT 'Culture description.',
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `production`.`document` (
  `title` character varying(50) NOT NULL COMMENT 'Title of the document.',
  `owner` integer NOT NULL COMMENT 'Employee who controls the document.  Foreign key to Employee.BusinessEntityID',
  `folderflag` public.Flag NOT NULL DEFAULT false COMMENT '0 = This is a folder, 1 = This is a document.',
  `filename` character varying(400) NOT NULL COMMENT 'File name of the document',
  `fileextension` character varying(8) COMMENT 'File extension indicating the document type. For example, .doc or .txt.',
  `revision` character(5) NOT NULL COMMENT 'Revision number of the document.',
  `changenumber` integer NOT NULL DEFAULT 0 COMMENT 'Engineering change approval number.',
  `status` smallint NOT NULL COMMENT '1 = Pending approval, 2 = Approved, 3 = Obsolete',
  `documentsummary` text COMMENT 'Document abstract.',
  `document` bytea COMMENT 'Complete document.',
  `rowguid` uuid UNIQUE NOT NULL DEFAULT (public.uuid_generate_v1()) COMMENT 'ROWGUIDCOL number uniquely identifying the record. Required for FileStream.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  `documentnode` character varying PRIMARY KEY NOT NULL DEFAULT ('/'::charactervarying) COMMENT 'Primary key for Document records.',
  CONSTRAINT `CK_Document_Status` CHECK (((status >= 1) AND (status <= 3)))
);

CREATE TABLE `production`.`illustration` (
  `illustrationid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for Illustration records.',
  `diagram` xml COMMENT 'Illustrations used in manufacturing instructions. Stored as XML.',
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `production`.`location` (
  `locationid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for Location records.',
  `name` public.Name NOT NULL COMMENT 'Location description.',
  `costrate` numeric NOT NULL DEFAULT 0 COMMENT 'Standard hourly cost of the manufacturing location.',
  `availability` numeric(8,2) NOT NULL DEFAULT 0 COMMENT 'Work capacity (in hours) of the manufacturing location.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_Location_Availability` CHECK ((availability >= 0.00)),
  CONSTRAINT `CK_Location_CostRate` CHECK ((costrate >= 0.00))
);

CREATE TABLE `production`.`product` (
  `productid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for Product records.',
  `name` public.Name NOT NULL COMMENT 'Name of the product.',
  `productnumber` character varying(25) NOT NULL COMMENT 'Unique product identification number.',
  `makeflag` public.Flag NOT NULL DEFAULT true COMMENT '0 = Product is purchased, 1 = Product is manufactured in-house.',
  `finishedgoodsflag` public.Flag NOT NULL DEFAULT true COMMENT '0 = Product is not a salable item. 1 = Product is salable.',
  `color` character varying(15) COMMENT 'Product color.',
  `safetystocklevel` smallint NOT NULL COMMENT 'Minimum inventory quantity.',
  `reorderpoint` smallint NOT NULL COMMENT 'Inventory level that triggers a purchase order or work order.',
  `standardcost` numeric NOT NULL COMMENT 'Standard cost of the product.',
  `listprice` numeric NOT NULL COMMENT 'Selling price.',
  `size` character varying(5) COMMENT 'Product size.',
  `sizeunitmeasurecode` character(3) COMMENT 'Unit of measure for Size column.',
  `weightunitmeasurecode` character(3) COMMENT 'Unit of measure for Weight column.',
  `weight` numeric(8,2) COMMENT 'Product weight.',
  `daystomanufacture` integer NOT NULL COMMENT 'Number of days required to manufacture the product.',
  `productline` character(2) COMMENT 'R = Road, M = Mountain, T = Touring, S = Standard',
  `class` character(2) COMMENT 'H = High, M = Medium, L = Low',
  `style` character(2) COMMENT 'W = Womens, M = Mens, U = Universal',
  `productsubcategoryid` integer COMMENT 'Product is a member of this product subcategory. Foreign key to ProductSubCategory.ProductSubCategoryID.',
  `productmodelid` integer COMMENT 'Product is a member of this product model. Foreign key to ProductModel.ProductModelID.',
  `sellstartdate` timestamp NOT NULL COMMENT 'Date the product was available for sale.',
  `sellenddate` timestamp COMMENT 'Date the product was no longer available for sale.',
  `discontinueddate` timestamp COMMENT 'Date the product was discontinued.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_Product_Class` CHECK (((upper((class)::text) = ANY (ARRAY['L'::text, 'M'::text, 'H'::text])) OR (class IS NULL))),
  CONSTRAINT `CK_Product_DaysToManufacture` CHECK ((daystomanufacture >= 0)),
  CONSTRAINT `CK_Product_ListPrice` CHECK ((listprice >= 0.00)),
  CONSTRAINT `CK_Product_ProductLine` CHECK (((upper((productline)::text) = ANY (ARRAY['S'::text, 'T'::text, 'M'::text, 'R'::text])) OR (productline IS NULL))),
  CONSTRAINT `CK_Product_ReorderPoint` CHECK ((reorderpoint > 0)),
  CONSTRAINT `CK_Product_SafetyStockLevel` CHECK ((safetystocklevel > 0)),
  CONSTRAINT `CK_Product_SellEndDate` CHECK (((sellenddate >= sellstartdate) OR (sellenddate IS NULL))),
  CONSTRAINT `CK_Product_StandardCost` CHECK ((standardcost >= 0.00)),
  CONSTRAINT `CK_Product_Style` CHECK (((upper((style)::text) = ANY (ARRAY['W'::text, 'M'::text, 'U'::text])) OR (style IS NULL))),
  CONSTRAINT `CK_Product_Weight` CHECK ((weight > 0.00))
);

CREATE TABLE `production`.`productcategory` (
  `productcategoryid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for ProductCategory records.',
  `name` public.Name NOT NULL COMMENT 'Category description.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `production`.`productcosthistory` (
  `productid` integer NOT NULL COMMENT 'Product identification number. Foreign key to Product.ProductID',
  `startdate` timestamp NOT NULL COMMENT 'Product cost start date.',
  `enddate` timestamp COMMENT 'Product cost end date.',
  `standardcost` numeric NOT NULL COMMENT 'Standard cost of the product.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_ProductCostHistory_EndDate` CHECK (((enddate >= startdate) OR (enddate IS NULL))),
  CONSTRAINT `CK_ProductCostHistory_StandardCost` CHECK ((standardcost >= 0.00)),
  PRIMARY KEY (`productid`, `startdate`)
);

CREATE TABLE `production`.`productdescription` (
  `productdescriptionid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for ProductDescription records.',
  `description` character varying(400) NOT NULL COMMENT 'Description of the product.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `production`.`productdocument` (
  `productid` integer NOT NULL COMMENT 'Product identification number. Foreign key to Product.ProductID.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  `documentnode` character varying NOT NULL DEFAULT ('/'::charactervarying) COMMENT 'Document identification number. Foreign key to Document.DocumentNode.',
  PRIMARY KEY (`productid`, `documentnode`)
);

CREATE TABLE `production`.`productinventory` (
  `productid` integer NOT NULL COMMENT 'Product identification number. Foreign key to Product.ProductID.',
  `locationid` smallint NOT NULL COMMENT 'Inventory location identification number. Foreign key to Location.LocationID.',
  `shelf` character varying(10) NOT NULL COMMENT 'Storage compartment within an inventory location.',
  `bin` smallint NOT NULL COMMENT 'Storage container on a shelf in an inventory location.',
  `quantity` smallint NOT NULL DEFAULT 0 COMMENT 'Quantity of products in the inventory location.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_ProductInventory_Bin` CHECK (((bin >= 0) AND (bin <= 100))),
  PRIMARY KEY (`productid`, `locationid`)
);

CREATE TABLE `production`.`productlistpricehistory` (
  `productid` integer NOT NULL COMMENT 'Product identification number. Foreign key to Product.ProductID',
  `startdate` timestamp NOT NULL COMMENT 'List price start date.',
  `enddate` timestamp COMMENT 'List price end date',
  `listprice` numeric NOT NULL COMMENT 'Product list price.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_ProductListPriceHistory_EndDate` CHECK (((enddate >= startdate) OR (enddate IS NULL))),
  CONSTRAINT `CK_ProductListPriceHistory_ListPrice` CHECK ((listprice > 0.00)),
  PRIMARY KEY (`productid`, `startdate`)
);

CREATE TABLE `production`.`productmodel` (
  `productmodelid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for ProductModel records.',
  `name` public.Name NOT NULL COMMENT 'Product model description.',
  `catalogdescription` xml COMMENT 'Detailed product catalog information in xml format.',
  `instructions` xml COMMENT 'Manufacturing instructions in xml format.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `production`.`productmodelillustration` (
  `productmodelid` integer NOT NULL COMMENT 'Primary key. Foreign key to ProductModel.ProductModelID.',
  `illustrationid` integer NOT NULL COMMENT 'Primary key. Foreign key to Illustration.IllustrationID.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY (`productmodelid`, `illustrationid`)
);

CREATE TABLE `production`.`productmodelproductdescriptionculture` (
  `productmodelid` integer NOT NULL COMMENT 'Primary key. Foreign key to ProductModel.ProductModelID.',
  `productdescriptionid` integer NOT NULL COMMENT 'Primary key. Foreign key to ProductDescription.ProductDescriptionID.',
  `cultureid` character(6) NOT NULL COMMENT 'Culture identification number. Foreign key to Culture.CultureID.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY (`productmodelid`, `productdescriptionid`, `cultureid`)
);

CREATE TABLE `production`.`productphoto` (
  `productphotoid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for ProductPhoto records.',
  `thumbnailphoto` bytea COMMENT 'Small image of the product.',
  `thumbnailphotofilename` character varying(50) COMMENT 'Small image file name.',
  `largephoto` bytea COMMENT 'Large image of the product.',
  `largephotofilename` character varying(50) COMMENT 'Large image file name.',
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `production`.`productproductphoto` (
  `productid` integer NOT NULL COMMENT 'Product identification number. Foreign key to Product.ProductID.',
  `productphotoid` integer NOT NULL COMMENT 'Product photo identification number. Foreign key to ProductPhoto.ProductPhotoID.',
  `primary` public.Flag NOT NULL DEFAULT false COMMENT '0 = Photo is not the principal image. 1 = Photo is the principal image.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY (`productid`, `productphotoid`)
);

CREATE TABLE `production`.`productreview` (
  `productreviewid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for ProductReview records.',
  `productid` integer NOT NULL COMMENT 'Product identification number. Foreign key to Product.ProductID.',
  `reviewername` public.Name NOT NULL COMMENT 'Name of the reviewer.',
  `reviewdate` timestamp NOT NULL DEFAULT (now()) COMMENT 'Date review was submitted.',
  `emailaddress` character varying(50) NOT NULL COMMENT 'Reviewer''s e-mail address.',
  `rating` integer NOT NULL COMMENT 'Product rating given by the reviewer. Scale is 1 to 5 with 5 as the highest rating.',
  `comments` character varying(3850) COMMENT 'Reviewer''s comments
WARNING: can have several lines!',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_ProductReview_Rating` CHECK (((rating >= 1) AND (rating <= 5)))
);

CREATE TABLE `production`.`productsubcategory` (
  `productsubcategoryid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for ProductSubcategory records.',
  `productcategoryid` integer NOT NULL COMMENT 'Product category identification number. Foreign key to ProductCategory.ProductCategoryID.',
  `name` public.Name NOT NULL COMMENT 'Subcategory description.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `production`.`scrapreason` (
  `scrapreasonid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for ScrapReason records.',
  `name` public.Name NOT NULL COMMENT 'Failure description.',
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `production`.`transactionhistory` (
  `transactionid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for TransactionHistory records.',
  `productid` integer NOT NULL COMMENT 'Product identification number. Foreign key to Product.ProductID.',
  `referenceorderid` integer NOT NULL COMMENT 'Purchase order, sales order, or work order identification number.',
  `referenceorderlineid` integer NOT NULL DEFAULT 0 COMMENT 'Line number associated with the purchase order, sales order, or work order.',
  `transactiondate` timestamp NOT NULL DEFAULT (now()) COMMENT 'Date and time of the transaction.',
  `transactiontype` character(1) NOT NULL COMMENT 'W = WorkOrder, S = SalesOrder, P = PurchaseOrder',
  `quantity` integer NOT NULL COMMENT 'Product quantity.',
  `actualcost` numeric NOT NULL COMMENT 'Product cost.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_TransactionHistory_TransactionType` CHECK ((upper((transactiontype)::text) = ANY (ARRAY['W'::text, 'S'::text, 'P'::text])))
);

CREATE TABLE `production`.`transactionhistoryarchive` (
  `transactionid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for TransactionHistoryArchive records.',
  `productid` integer NOT NULL COMMENT 'Product identification number. Foreign key to Product.ProductID.',
  `referenceorderid` integer NOT NULL COMMENT 'Purchase order, sales order, or work order identification number.',
  `referenceorderlineid` integer NOT NULL DEFAULT 0 COMMENT 'Line number associated with the purchase order, sales order, or work order.',
  `transactiondate` timestamp NOT NULL DEFAULT (now()) COMMENT 'Date and time of the transaction.',
  `transactiontype` character(1) NOT NULL COMMENT 'W = Work Order, S = Sales Order, P = Purchase Order',
  `quantity` integer NOT NULL COMMENT 'Product quantity.',
  `actualcost` numeric NOT NULL COMMENT 'Product cost.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_TransactionHistoryArchive_TransactionType` CHECK ((upper((transactiontype)::text) = ANY (ARRAY['W'::text, 'S'::text, 'P'::text])))
);

CREATE TABLE `production`.`unitmeasure` (
  `unitmeasurecode` character(3) PRIMARY KEY NOT NULL COMMENT 'Primary key.',
  `name` public.Name NOT NULL COMMENT 'Unit of measure description.',
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `production`.`workorder` (
  `workorderid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for WorkOrder records.',
  `productid` integer NOT NULL COMMENT 'Product identification number. Foreign key to Product.ProductID.',
  `orderqty` integer NOT NULL COMMENT 'Product quantity to build.',
  `scrappedqty` smallint NOT NULL COMMENT 'Quantity that failed inspection.',
  `startdate` timestamp NOT NULL COMMENT 'Work order start date.',
  `enddate` timestamp COMMENT 'Work order end date.',
  `duedate` timestamp NOT NULL COMMENT 'Work order due date.',
  `scrapreasonid` smallint COMMENT 'Reason for inspection failure.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_WorkOrder_EndDate` CHECK (((enddate >= startdate) OR (enddate IS NULL))),
  CONSTRAINT `CK_WorkOrder_OrderQty` CHECK ((orderqty > 0)),
  CONSTRAINT `CK_WorkOrder_ScrappedQty` CHECK ((scrappedqty >= 0))
);

CREATE TABLE `production`.`workorderrouting` (
  `workorderid` integer NOT NULL COMMENT 'Primary key. Foreign key to WorkOrder.WorkOrderID.',
  `productid` integer NOT NULL COMMENT 'Primary key. Foreign key to Product.ProductID.',
  `operationsequence` smallint NOT NULL COMMENT 'Primary key. Indicates the manufacturing process sequence.',
  `locationid` smallint NOT NULL COMMENT 'Manufacturing location where the part is processed. Foreign key to Location.LocationID.',
  `scheduledstartdate` timestamp NOT NULL COMMENT 'Planned manufacturing start date.',
  `scheduledenddate` timestamp NOT NULL COMMENT 'Planned manufacturing end date.',
  `actualstartdate` timestamp COMMENT 'Actual start date.',
  `actualenddate` timestamp COMMENT 'Actual end date.',
  `actualresourcehrs` numeric(9,4) COMMENT 'Number of manufacturing hours used.',
  `plannedcost` numeric NOT NULL COMMENT 'Estimated manufacturing cost.',
  `actualcost` numeric COMMENT 'Actual manufacturing cost.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_WorkOrderRouting_ActualCost` CHECK ((actualcost > 0.00)),
  CONSTRAINT `CK_WorkOrderRouting_ActualEndDate` CHECK (((actualenddate >= actualstartdate) OR (actualenddate IS NULL) OR (actualstartdate IS NULL))),
  CONSTRAINT `CK_WorkOrderRouting_ActualResourceHrs` CHECK ((actualresourcehrs >= 0.0000)),
  CONSTRAINT `CK_WorkOrderRouting_PlannedCost` CHECK ((plannedcost > 0.00)),
  CONSTRAINT `CK_WorkOrderRouting_ScheduledEndDate` CHECK ((scheduledenddate >= scheduledstartdate)),
  PRIMARY KEY (`workorderid`, `productid`, `operationsequence`)
);

CREATE TABLE `purchasing`.`purchaseorderdetail` (
  `purchaseorderid` integer NOT NULL COMMENT 'Primary key. Foreign key to PurchaseOrderHeader.PurchaseOrderID.',
  `purchaseorderdetailid` integer NOT NULL COMMENT 'Primary key. One line number per purchased product.',
  `duedate` timestamp NOT NULL COMMENT 'Date the product is expected to be received.',
  `orderqty` smallint NOT NULL COMMENT 'Quantity ordered.',
  `productid` integer NOT NULL COMMENT 'Product identification number. Foreign key to Product.ProductID.',
  `unitprice` numeric NOT NULL COMMENT 'Vendor''s selling price of a single product.',
  `receivedqty` numeric(8,2) NOT NULL COMMENT 'Quantity actually received from the vendor.',
  `rejectedqty` numeric(8,2) NOT NULL COMMENT 'Quantity rejected during inspection.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_PurchaseOrderDetail_OrderQty` CHECK ((orderqty > 0)),
  CONSTRAINT `CK_PurchaseOrderDetail_ReceivedQty` CHECK ((receivedqty >= 0.00)),
  CONSTRAINT `CK_PurchaseOrderDetail_RejectedQty` CHECK ((rejectedqty >= 0.00)),
  CONSTRAINT `CK_PurchaseOrderDetail_UnitPrice` CHECK ((unitprice >= 0.00)),
  PRIMARY KEY (`purchaseorderid`, `purchaseorderdetailid`)
);

CREATE TABLE `purchasing`.`purchaseorderheader` (
  `purchaseorderid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key.',
  `revisionnumber` smallint NOT NULL DEFAULT 0 COMMENT 'Incremental number to track changes to the purchase order over time.',
  `status` smallint NOT NULL DEFAULT 1 COMMENT 'Order current status. 1 = Pending; 2 = Approved; 3 = Rejected; 4 = Complete',
  `employeeid` integer NOT NULL COMMENT 'Employee who created the purchase order. Foreign key to Employee.BusinessEntityID.',
  `vendorid` integer NOT NULL COMMENT 'Vendor with whom the purchase order is placed. Foreign key to Vendor.BusinessEntityID.',
  `shipmethodid` integer NOT NULL COMMENT 'Shipping method. Foreign key to ShipMethod.ShipMethodID.',
  `orderdate` timestamp NOT NULL DEFAULT (now()) COMMENT 'Purchase order creation date.',
  `shipdate` timestamp COMMENT 'Estimated shipment date from the vendor.',
  `subtotal` numeric NOT NULL DEFAULT 0 COMMENT 'Purchase order subtotal. Computed as SUM(PurchaseOrderDetail.LineTotal)for the appropriate PurchaseOrderID.',
  `taxamt` numeric NOT NULL DEFAULT 0 COMMENT 'Tax amount.',
  `freight` numeric NOT NULL DEFAULT 0 COMMENT 'Shipping cost.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_PurchaseOrderHeader_Freight` CHECK ((freight >= 0.00)),
  CONSTRAINT `CK_PurchaseOrderHeader_ShipDate` CHECK (((shipdate >= orderdate) OR (shipdate IS NULL))),
  CONSTRAINT `CK_PurchaseOrderHeader_Status` CHECK (((status >= 1) AND (status <= 4))),
  CONSTRAINT `CK_PurchaseOrderHeader_SubTotal` CHECK ((subtotal >= 0.00)),
  CONSTRAINT `CK_PurchaseOrderHeader_TaxAmt` CHECK ((taxamt >= 0.00))
);

CREATE TABLE `purchasing`.`productvendor` (
  `productid` integer NOT NULL COMMENT 'Primary key. Foreign key to Product.ProductID.',
  `businessentityid` integer NOT NULL COMMENT 'Primary key. Foreign key to Vendor.BusinessEntityID.',
  `averageleadtime` integer NOT NULL COMMENT 'The average span of time (in days) between placing an order with the vendor and receiving the purchased product.',
  `standardprice` numeric NOT NULL COMMENT 'The vendor''s usual selling price.',
  `lastreceiptcost` numeric COMMENT 'The selling price when last purchased.',
  `lastreceiptdate` timestamp COMMENT 'Date the product was last received by the vendor.',
  `minorderqty` integer NOT NULL COMMENT 'The maximum quantity that should be ordered.',
  `maxorderqty` integer NOT NULL COMMENT 'The minimum quantity that should be ordered.',
  `onorderqty` integer COMMENT 'The quantity currently on order.',
  `unitmeasurecode` character(3) NOT NULL COMMENT 'The product''s unit of measure.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_ProductVendor_AverageLeadTime` CHECK ((averageleadtime >= 1)),
  CONSTRAINT `CK_ProductVendor_LastReceiptCost` CHECK ((lastreceiptcost > 0.00)),
  CONSTRAINT `CK_ProductVendor_MaxOrderQty` CHECK ((maxorderqty >= 1)),
  CONSTRAINT `CK_ProductVendor_MinOrderQty` CHECK ((minorderqty >= 1)),
  CONSTRAINT `CK_ProductVendor_OnOrderQty` CHECK ((onorderqty >= 0)),
  CONSTRAINT `CK_ProductVendor_StandardPrice` CHECK ((standardprice > 0.00)),
  PRIMARY KEY (`productid`, `businessentityid`)
);

CREATE TABLE `purchasing`.`shipmethod` (
  `shipmethodid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for ShipMethod records.',
  `name` public.Name NOT NULL COMMENT 'Shipping company name.',
  `shipbase` numeric NOT NULL DEFAULT 0 COMMENT 'Minimum shipping charge.',
  `shiprate` numeric NOT NULL DEFAULT 0 COMMENT 'Shipping charge per pound.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_ShipMethod_ShipBase` CHECK ((shipbase > 0.00)),
  CONSTRAINT `CK_ShipMethod_ShipRate` CHECK ((shiprate > 0.00))
);

CREATE TABLE `purchasing`.`vendor` (
  `businessentityid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for Vendor records.  Foreign key to BusinessEntity.BusinessEntityID',
  `accountnumber` public.AccountNumber NOT NULL COMMENT 'Vendor account (identification) number.',
  `name` public.Name NOT NULL COMMENT 'Company name.',
  `creditrating` smallint NOT NULL COMMENT '1 = Superior, 2 = Excellent, 3 = Above average, 4 = Average, 5 = Below average',
  `preferredvendorstatus` public.Flag NOT NULL DEFAULT true COMMENT '0 = Do not use if another vendor is available. 1 = Preferred over other vendors supplying the same product.',
  `activeflag` public.Flag NOT NULL DEFAULT true COMMENT '0 = Vendor no longer used. 1 = Vendor is actively used.',
  `purchasingwebserviceurl` character varying(1024) COMMENT 'Vendor URL.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_Vendor_CreditRating` CHECK (((creditrating >= 1) AND (creditrating <= 5)))
);

CREATE TABLE `sales`.`customer` (
  `customerid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key.',
  `personid` integer COMMENT 'Foreign key to Person.BusinessEntityID',
  `storeid` integer COMMENT 'Foreign key to Store.BusinessEntityID',
  `territoryid` integer COMMENT 'ID of the territory in which the customer is located. Foreign key to SalesTerritory.SalesTerritoryID.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `sales`.`creditcard` (
  `creditcardid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for CreditCard records.',
  `cardtype` character varying(50) NOT NULL COMMENT 'Credit card name.',
  `cardnumber` character varying(25) NOT NULL COMMENT 'Credit card number.',
  `expmonth` smallint NOT NULL COMMENT 'Credit card expiration month.',
  `expyear` smallint NOT NULL COMMENT 'Credit card expiration year.',
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `sales`.`currencyrate` (
  `currencyrateid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for CurrencyRate records.',
  `currencyratedate` timestamp NOT NULL COMMENT 'Date and time the exchange rate was obtained.',
  `fromcurrencycode` character(3) NOT NULL COMMENT 'Exchange rate was converted from this currency code.',
  `tocurrencycode` character(3) NOT NULL COMMENT 'Exchange rate was converted to this currency code.',
  `averagerate` numeric NOT NULL COMMENT 'Average exchange rate for the day.',
  `endofdayrate` numeric NOT NULL COMMENT 'Final exchange rate for the day.',
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `sales`.`countryregioncurrency` (
  `countryregioncode` character varying(3) NOT NULL COMMENT 'ISO code for countries and regions. Foreign key to CountryRegion.CountryRegionCode.',
  `currencycode` character(3) NOT NULL COMMENT 'ISO standard currency code. Foreign key to Currency.CurrencyCode.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY (`countryregioncode`, `currencycode`)
);

CREATE TABLE `sales`.`currency` (
  `currencycode` character(3) PRIMARY KEY NOT NULL COMMENT 'The ISO code for the Currency.',
  `name` public.Name NOT NULL COMMENT 'Currency name.',
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `sales`.`personcreditcard` (
  `businessentityid` integer NOT NULL COMMENT 'Business entity identification number. Foreign key to Person.BusinessEntityID.',
  `creditcardid` integer NOT NULL COMMENT 'Credit card identification number. Foreign key to CreditCard.CreditCardID.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY (`businessentityid`, `creditcardid`)
);

CREATE TABLE `sales`.`store` (
  `businessentityid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key. Foreign key to Customer.BusinessEntityID.',
  `name` public.Name NOT NULL COMMENT 'Name of the store.',
  `salespersonid` integer COMMENT 'ID of the sales person assigned to the customer. Foreign key to SalesPerson.BusinessEntityID.',
  `demographics` xml COMMENT 'Demographic informationg about the store such as the number of employees, annual sales and store type.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `sales`.`shoppingcartitem` (
  `shoppingcartitemid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for ShoppingCartItem records.',
  `shoppingcartid` character varying(50) NOT NULL COMMENT 'Shopping cart identification number.',
  `quantity` integer NOT NULL DEFAULT 1 COMMENT 'Product quantity ordered.',
  `productid` integer NOT NULL COMMENT 'Product ordered. Foreign key to Product.ProductID.',
  `datecreated` timestamp NOT NULL DEFAULT (now()) COMMENT 'Date the time the record was created.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_ShoppingCartItem_Quantity` CHECK ((quantity >= 1))
);

CREATE TABLE `sales`.`specialoffer` (
  `specialofferid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for SpecialOffer records.',
  `description` character varying(255) NOT NULL COMMENT 'Discount description.',
  `discountpct` numeric NOT NULL DEFAULT 0 COMMENT 'Discount precentage.',
  `type` character varying(50) NOT NULL COMMENT 'Discount type category.',
  `category` character varying(50) NOT NULL COMMENT 'Group the discount applies to such as Reseller or Customer.',
  `startdate` timestamp NOT NULL COMMENT 'Discount start date.',
  `enddate` timestamp NOT NULL COMMENT 'Discount end date.',
  `minqty` integer NOT NULL DEFAULT 0 COMMENT 'Minimum discount percent allowed.',
  `maxqty` integer COMMENT 'Maximum discount percent allowed.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_SpecialOffer_DiscountPct` CHECK ((discountpct >= 0.00)),
  CONSTRAINT `CK_SpecialOffer_EndDate` CHECK ((enddate >= startdate)),
  CONSTRAINT `CK_SpecialOffer_MaxQty` CHECK ((maxqty >= 0)),
  CONSTRAINT `CK_SpecialOffer_MinQty` CHECK ((minqty >= 0))
);

CREATE TABLE `sales`.`salesorderdetail` (
  `salesorderid` integer NOT NULL COMMENT 'Primary key. Foreign key to SalesOrderHeader.SalesOrderID.',
  `salesorderdetailid` integer NOT NULL COMMENT 'Primary key. One incremental unique number per product sold.',
  `carriertrackingnumber` character varying(25) COMMENT 'Shipment tracking number supplied by the shipper.',
  `orderqty` smallint NOT NULL COMMENT 'Quantity ordered per product.',
  `productid` integer NOT NULL COMMENT 'Product sold to customer. Foreign key to Product.ProductID.',
  `specialofferid` integer NOT NULL COMMENT 'Promotional code. Foreign key to SpecialOffer.SpecialOfferID.',
  `unitprice` numeric NOT NULL COMMENT 'Selling price of a single product.',
  `unitpricediscount` numeric NOT NULL DEFAULT 0 COMMENT 'Discount amount.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_SalesOrderDetail_OrderQty` CHECK ((orderqty > 0)),
  CONSTRAINT `CK_SalesOrderDetail_UnitPrice` CHECK ((unitprice >= 0.00)),
  CONSTRAINT `CK_SalesOrderDetail_UnitPriceDiscount` CHECK ((unitpricediscount >= 0.00)),
  PRIMARY KEY (`salesorderid`, `salesorderdetailid`)
);

CREATE TABLE `sales`.`salesorderheader` (
  `salesorderid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key.',
  `revisionnumber` smallint NOT NULL DEFAULT 0 COMMENT 'Incremental number to track changes to the sales order over time.',
  `orderdate` timestamp NOT NULL DEFAULT (now()) COMMENT 'Dates the sales order was created.',
  `duedate` timestamp NOT NULL COMMENT 'Date the order is due to the customer.',
  `shipdate` timestamp COMMENT 'Date the order was shipped to the customer.',
  `status` smallint NOT NULL DEFAULT 1 COMMENT 'Order current status. 1 = In process; 2 = Approved; 3 = Backordered; 4 = Rejected; 5 = Shipped; 6 = Cancelled',
  `onlineorderflag` public.Flag NOT NULL DEFAULT true COMMENT '0 = Order placed by sales person. 1 = Order placed online by customer.',
  `purchaseordernumber` public.OrderNumber COMMENT 'Customer purchase order number reference.',
  `accountnumber` public.AccountNumber COMMENT 'Financial accounting number reference.',
  `customerid` integer NOT NULL COMMENT 'Customer identification number. Foreign key to Customer.BusinessEntityID.',
  `salespersonid` integer COMMENT 'Sales person who created the sales order. Foreign key to SalesPerson.BusinessEntityID.',
  `territoryid` integer COMMENT 'Territory in which the sale was made. Foreign key to SalesTerritory.SalesTerritoryID.',
  `billtoaddressid` integer NOT NULL COMMENT 'Customer billing address. Foreign key to Address.AddressID.',
  `shiptoaddressid` integer NOT NULL COMMENT 'Customer shipping address. Foreign key to Address.AddressID.',
  `shipmethodid` integer NOT NULL COMMENT 'Shipping method. Foreign key to ShipMethod.ShipMethodID.',
  `creditcardid` integer COMMENT 'Credit card identification number. Foreign key to CreditCard.CreditCardID.',
  `creditcardapprovalcode` character varying(15) COMMENT 'Approval code provided by the credit card company.',
  `currencyrateid` integer COMMENT 'Currency exchange rate used. Foreign key to CurrencyRate.CurrencyRateID.',
  `subtotal` numeric NOT NULL DEFAULT 0 COMMENT 'Sales subtotal. Computed as SUM(SalesOrderDetail.LineTotal)for the appropriate SalesOrderID.',
  `taxamt` numeric NOT NULL DEFAULT 0 COMMENT 'Tax amount.',
  `freight` numeric NOT NULL DEFAULT 0 COMMENT 'Shipping cost.',
  `totaldue` numeric COMMENT 'Total due from customer. Computed as Subtotal + TaxAmt + Freight.',
  `comment` character varying(128) COMMENT 'Sales representative comments.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_SalesOrderHeader_DueDate` CHECK ((duedate >= orderdate)),
  CONSTRAINT `CK_SalesOrderHeader_Freight` CHECK ((freight >= 0.00)),
  CONSTRAINT `CK_SalesOrderHeader_ShipDate` CHECK (((shipdate >= orderdate) OR (shipdate IS NULL))),
  CONSTRAINT `CK_SalesOrderHeader_Status` CHECK (((status >= 0) AND (status <= 8))),
  CONSTRAINT `CK_SalesOrderHeader_SubTotal` CHECK ((subtotal >= 0.00)),
  CONSTRAINT `CK_SalesOrderHeader_TaxAmt` CHECK ((taxamt >= 0.00))
);

CREATE TABLE `sales`.`salesorderheadersalesreason` (
  `salesorderid` integer NOT NULL COMMENT 'Primary key. Foreign key to SalesOrderHeader.SalesOrderID.',
  `salesreasonid` integer NOT NULL COMMENT 'Primary key. Foreign key to SalesReason.SalesReasonID.',
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY (`salesorderid`, `salesreasonid`)
);

CREATE TABLE `sales`.`specialofferproduct` (
  `specialofferid` integer NOT NULL COMMENT 'Primary key for SpecialOfferProduct records.',
  `productid` integer NOT NULL COMMENT 'Product identification number. Foreign key to Product.ProductID.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY (`specialofferid`, `productid`)
);

CREATE TABLE `sales`.`salesperson` (
  `businessentityid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for SalesPerson records. Foreign key to Employee.BusinessEntityID',
  `territoryid` integer COMMENT 'Territory currently assigned to. Foreign key to SalesTerritory.SalesTerritoryID.',
  `salesquota` numeric COMMENT 'Projected yearly sales.',
  `bonus` numeric NOT NULL DEFAULT 0 COMMENT 'Bonus due if quota is met.',
  `commissionpct` numeric NOT NULL DEFAULT 0 COMMENT 'Commision percent received per sale.',
  `salesytd` numeric NOT NULL DEFAULT 0 COMMENT 'Sales total year to date.',
  `saleslastyear` numeric NOT NULL DEFAULT 0 COMMENT 'Sales total of previous year.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_SalesPerson_Bonus` CHECK ((bonus >= 0.00)),
  CONSTRAINT `CK_SalesPerson_CommissionPct` CHECK ((commissionpct >= 0.00)),
  CONSTRAINT `CK_SalesPerson_SalesLastYear` CHECK ((saleslastyear >= 0.00)),
  CONSTRAINT `CK_SalesPerson_SalesQuota` CHECK ((salesquota > 0.00)),
  CONSTRAINT `CK_SalesPerson_SalesYTD` CHECK ((salesytd >= 0.00))
);

CREATE TABLE `sales`.`salespersonquotahistory` (
  `businessentityid` integer NOT NULL COMMENT 'Sales person identification number. Foreign key to SalesPerson.BusinessEntityID.',
  `quotadate` timestamp NOT NULL COMMENT 'Sales quota date.',
  `salesquota` numeric NOT NULL COMMENT 'Sales quota amount.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_SalesPersonQuotaHistory_SalesQuota` CHECK ((salesquota > 0.00)),
  PRIMARY KEY (`businessentityid`, `quotadate`)
);

CREATE TABLE `sales`.`salesreason` (
  `salesreasonid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for SalesReason records.',
  `name` public.Name NOT NULL COMMENT 'Sales reason description.',
  `reasontype` public.Name NOT NULL COMMENT 'Category the sales reason belongs to.',
  `modifieddate` timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE `sales`.`salesterritory` (
  `territoryid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for SalesTerritory records.',
  `name` public.Name NOT NULL COMMENT 'Sales territory description',
  `countryregioncode` character varying(3) NOT NULL COMMENT 'ISO standard country or region code. Foreign key to CountryRegion.CountryRegionCode.',
  `group` character varying(50) NOT NULL COMMENT 'Geographic area to which the sales territory belong.',
  `salesytd` numeric NOT NULL DEFAULT 0 COMMENT 'Sales in the territory year to date.',
  `saleslastyear` numeric NOT NULL DEFAULT 0 COMMENT 'Sales in the territory the previous year.',
  `costytd` numeric NOT NULL DEFAULT 0 COMMENT 'Business costs in the territory year to date.',
  `costlastyear` numeric NOT NULL DEFAULT 0 COMMENT 'Business costs in the territory the previous year.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_SalesTerritory_CostLastYear` CHECK ((costlastyear >= 0.00)),
  CONSTRAINT `CK_SalesTerritory_CostYTD` CHECK ((costytd >= 0.00)),
  CONSTRAINT `CK_SalesTerritory_SalesLastYear` CHECK ((saleslastyear >= 0.00)),
  CONSTRAINT `CK_SalesTerritory_SalesYTD` CHECK ((salesytd >= 0.00))
);

CREATE TABLE `sales`.`salesterritoryhistory` (
  `businessentityid` integer NOT NULL COMMENT 'Primary key. The sales rep.  Foreign key to SalesPerson.BusinessEntityID.',
  `territoryid` integer NOT NULL COMMENT 'Primary key. Territory identification number. Foreign key to SalesTerritory.SalesTerritoryID.',
  `startdate` timestamp NOT NULL COMMENT 'Primary key. Date the sales representive started work in the territory.',
  `enddate` timestamp COMMENT 'Date the sales representative left work in the territory.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_SalesTerritoryHistory_EndDate` CHECK (((enddate >= startdate) OR (enddate IS NULL))),
  PRIMARY KEY (`businessentityid`, `startdate`, `territoryid`)
);

CREATE TABLE `sales`.`salestaxrate` (
  `salestaxrateid` integer PRIMARY KEY NOT NULL COMMENT 'Primary key for SalesTaxRate records.',
  `stateprovinceid` integer NOT NULL COMMENT 'State, province, or country/region the sales tax applies to.',
  `taxtype` smallint NOT NULL COMMENT '1 = Tax applied to retail transactions, 2 = Tax applied to wholesale transactions, 3 = Tax applied to all sales (retail and wholesale) transactions.',
  `taxrate` numeric NOT NULL DEFAULT 0 COMMENT 'Tax rate amount.',
  `name` public.Name NOT NULL COMMENT 'Tax rate description.',
  `rowguid` uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  `modifieddate` timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT `CK_SalesTaxRate_TaxType` CHECK (((taxtype >= 1) AND (taxtype <= 3)))
);

ALTER TABLE `humanresources`.`department` COMMENT = 'Lookup table containing the departments within the Adventure Works Cycles company.';

ALTER TABLE `humanresources`.`employee` COMMENT = 'Employee information such as salary, department, and title.';

ALTER TABLE `humanresources`.`employeedepartmenthistory` COMMENT = 'Employee department transfers.';

ALTER TABLE `humanresources`.`employeepayhistory` COMMENT = 'Employee pay history.';

ALTER TABLE `humanresources`.`jobcandidate` COMMENT = 'RÃ©sumÃ©s submitted to Human Resources by job applicants.';

ALTER TABLE `humanresources`.`shift` COMMENT = 'Work shift lookup table.';

ALTER TABLE `person`.`address` COMMENT = 'Street address information for customers, employees, and vendors.';

ALTER TABLE `person`.`businessentityaddress` COMMENT = 'Cross-reference table mapping customers, vendors, and employees to their addresses.';

ALTER TABLE `person`.`countryregion` COMMENT = 'Lookup table containing the ISO standard codes for countries and regions.';

ALTER TABLE `person`.`emailaddress` COMMENT = 'Where to send a person email.';

ALTER TABLE `person`.`person` COMMENT = 'Human beings involved with AdventureWorks: employees, customer contacts, and vendor contacts.';

ALTER TABLE `person`.`personphone` COMMENT = 'Telephone number and type of a person.';

ALTER TABLE `person`.`phonenumbertype` COMMENT = 'Type of phone number of a person.';

ALTER TABLE `person`.`stateprovince` COMMENT = 'State and province lookup table.';

ALTER TABLE `person`.`addresstype` COMMENT = 'Types of addresses stored in the Address table.';

ALTER TABLE `person`.`businessentity` COMMENT = 'Source of the ID that connects vendors, customers, and employees with address and contact information.';

ALTER TABLE `person`.`businessentitycontact` COMMENT = 'Cross-reference table mapping stores, vendors, and employees to people';

ALTER TABLE `person`.`contacttype` COMMENT = 'Lookup table containing the types of business entity contacts.';

ALTER TABLE `person`.`password` COMMENT = 'One way hashed authentication information';

ALTER TABLE `production`.`billofmaterials` COMMENT = 'Items required to make bicycles and bicycle subassemblies. It identifies the heirarchical relationship between a parent product and its components.';

ALTER TABLE `production`.`culture` COMMENT = 'Lookup table containing the languages in which some AdventureWorks data is stored.';

ALTER TABLE `production`.`document` COMMENT = 'Product maintenance documents.';

ALTER TABLE `production`.`illustration` COMMENT = 'Bicycle assembly diagrams.';

ALTER TABLE `production`.`location` COMMENT = 'Product inventory and manufacturing locations.';

ALTER TABLE `production`.`product` COMMENT = 'Products sold or used in the manfacturing of sold products.';

ALTER TABLE `production`.`productcategory` COMMENT = 'High-level product categorization.';

ALTER TABLE `production`.`productcosthistory` COMMENT = 'Changes in the cost of a product over time.';

ALTER TABLE `production`.`productdescription` COMMENT = 'Product descriptions in several languages.';

ALTER TABLE `production`.`productdocument` COMMENT = 'Cross-reference table mapping products to related product documents.';

ALTER TABLE `production`.`productinventory` COMMENT = 'Product inventory information.';

ALTER TABLE `production`.`productlistpricehistory` COMMENT = 'Changes in the list price of a product over time.';

ALTER TABLE `production`.`productmodel` COMMENT = 'Product model classification.';

ALTER TABLE `production`.`productmodelillustration` COMMENT = 'Cross-reference table mapping product models and illustrations.';

ALTER TABLE `production`.`productmodelproductdescriptionculture` COMMENT = 'Cross-reference table mapping product descriptions and the language the description is written in.';

ALTER TABLE `production`.`productphoto` COMMENT = 'Product images.';

ALTER TABLE `production`.`productproductphoto` COMMENT = 'Cross-reference table mapping products and product photos.';

ALTER TABLE `production`.`productreview` COMMENT = 'Customer reviews of products they have purchased.';

ALTER TABLE `production`.`productsubcategory` COMMENT = 'Product subcategories. See ProductCategory table.';

ALTER TABLE `production`.`scrapreason` COMMENT = 'Manufacturing failure reasons lookup table.';

ALTER TABLE `production`.`transactionhistory` COMMENT = 'Record of each purchase order, sales order, or work order transaction year to date.';

ALTER TABLE `production`.`transactionhistoryarchive` COMMENT = 'Transactions for previous years.';

ALTER TABLE `production`.`unitmeasure` COMMENT = 'Unit of measure lookup table.';

ALTER TABLE `production`.`workorder` COMMENT = 'Manufacturing work orders.';

ALTER TABLE `production`.`workorderrouting` COMMENT = 'Work order details.';

ALTER TABLE `purchasing`.`purchaseorderdetail` COMMENT = 'Individual products associated with a specific purchase order. See PurchaseOrderHeader.';

ALTER TABLE `purchasing`.`purchaseorderheader` COMMENT = 'General purchase order information. See PurchaseOrderDetail
for more details.';

ALTER TABLE `purchasing`.`productvendor` COMMENT = 'Cross-reference table mapping vendors with the products they supply.';

ALTER TABLE `purchasing`.`shipmethod` COMMENT = 'Shipping company lookup table.';

ALTER TABLE `purchasing`.`vendor` COMMENT = 'Companies from whom Adventure Works Cycles purchases parts or other goods.';

ALTER TABLE `sales`.`customer` COMMENT = 'Current customer information. Also see the Person and Store tables.';

ALTER TABLE `sales`.`creditcard` COMMENT = 'Customer credit card information.';

ALTER TABLE `sales`.`currencyrate` COMMENT = 'Currency exchange rates.';

ALTER TABLE `sales`.`countryregioncurrency` COMMENT = 'Cross-reference table mapping ISO currency codes to a country or region.';

ALTER TABLE `sales`.`currency` COMMENT = 'Lookup table containing standard ISO currencies.';

ALTER TABLE `sales`.`personcreditcard` COMMENT = 'Cross-reference table mapping people to their credit card information in the CreditCard table.';

ALTER TABLE `sales`.`store` COMMENT = 'Customers (resellers) of Adventure Works products.';

ALTER TABLE `sales`.`shoppingcartitem` COMMENT = 'Contains online customer orders until the order is submitted or cancelled.';

ALTER TABLE `sales`.`specialoffer` COMMENT = 'Sale discounts lookup table.';

ALTER TABLE `sales`.`salesorderdetail` COMMENT = 'Individual products associated with a specific sales order. See SalesOrderHeader.';

ALTER TABLE `sales`.`salesorderheader` COMMENT = 'General sales order information.';

ALTER TABLE `sales`.`salesorderheadersalesreason` COMMENT = 'Cross-reference table mapping sales orders to sales reason codes.';

ALTER TABLE `sales`.`specialofferproduct` COMMENT = 'Cross-reference table mapping products to special offer discounts.';

ALTER TABLE `sales`.`salesperson` COMMENT = 'Sales representative current information.';

ALTER TABLE `sales`.`salespersonquotahistory` COMMENT = 'Sales performance tracking.';

ALTER TABLE `sales`.`salesreason` COMMENT = 'Lookup table of customer purchase reasons.';

ALTER TABLE `sales`.`salesterritory` COMMENT = 'Sales territory lookup table.';

ALTER TABLE `sales`.`salesterritoryhistory` COMMENT = 'Sales representative transfers to other sales territories.';

ALTER TABLE `sales`.`salestaxrate` COMMENT = 'Tax rate lookup table.';

ALTER TABLE `humanresources`.`employeedepartmenthistory` ADD CONSTRAINT `FK_EmployeeDepartmentHistory_Department_DepartmentID` FOREIGN KEY (`departmentid`) REFERENCES `humanresources`.`department` (`departmentid`);

ALTER TABLE `humanresources`.`employeedepartmenthistory` ADD CONSTRAINT `FK_EmployeeDepartmentHistory_Employee_BusinessEntityID` FOREIGN KEY (`businessentityid`) REFERENCES `humanresources`.`employee` (`businessentityid`);

ALTER TABLE `humanresources`.`employeedepartmenthistory` ADD CONSTRAINT `FK_EmployeeDepartmentHistory_Shift_ShiftID` FOREIGN KEY (`shiftid`) REFERENCES `humanresources`.`shift` (`shiftid`);

ALTER TABLE `humanresources`.`employeepayhistory` ADD CONSTRAINT `FK_EmployeePayHistory_Employee_BusinessEntityID` FOREIGN KEY (`businessentityid`) REFERENCES `humanresources`.`employee` (`businessentityid`);

ALTER TABLE `humanresources`.`employee` ADD CONSTRAINT `FK_Employee_Person_BusinessEntityID` FOREIGN KEY (`businessentityid`) REFERENCES `person`.`person` (`businessentityid`);

ALTER TABLE `humanresources`.`jobcandidate` ADD CONSTRAINT `FK_JobCandidate_Employee_BusinessEntityID` FOREIGN KEY (`businessentityid`) REFERENCES `humanresources`.`employee` (`businessentityid`);

ALTER TABLE `person`.`address` ADD CONSTRAINT `FK_Address_StateProvince_StateProvinceID` FOREIGN KEY (`stateprovinceid`) REFERENCES `person`.`stateprovince` (`stateprovinceid`);

ALTER TABLE `person`.`businessentityaddress` ADD CONSTRAINT `FK_BusinessEntityAddress_AddressType_AddressTypeID` FOREIGN KEY (`addresstypeid`) REFERENCES `person`.`addresstype` (`addresstypeid`);

ALTER TABLE `person`.`businessentityaddress` ADD CONSTRAINT `FK_BusinessEntityAddress_Address_AddressID` FOREIGN KEY (`addressid`) REFERENCES `person`.`address` (`addressid`);

ALTER TABLE `person`.`businessentityaddress` ADD CONSTRAINT `FK_BusinessEntityAddress_BusinessEntity_BusinessEntityID` FOREIGN KEY (`businessentityid`) REFERENCES `person`.`businessentity` (`businessentityid`);

ALTER TABLE `person`.`businessentitycontact` ADD CONSTRAINT `FK_BusinessEntityContact_BusinessEntity_BusinessEntityID` FOREIGN KEY (`businessentityid`) REFERENCES `person`.`businessentity` (`businessentityid`);

ALTER TABLE `person`.`businessentitycontact` ADD CONSTRAINT `FK_BusinessEntityContact_ContactType_ContactTypeID` FOREIGN KEY (`contacttypeid`) REFERENCES `person`.`contacttype` (`contacttypeid`);

ALTER TABLE `person`.`businessentitycontact` ADD CONSTRAINT `FK_BusinessEntityContact_Person_PersonID` FOREIGN KEY (`personid`) REFERENCES `person`.`person` (`businessentityid`);

ALTER TABLE `person`.`emailaddress` ADD CONSTRAINT `FK_EmailAddress_Person_BusinessEntityID` FOREIGN KEY (`businessentityid`) REFERENCES `person`.`person` (`businessentityid`);

ALTER TABLE `person`.`password` ADD CONSTRAINT `FK_Password_Person_BusinessEntityID` FOREIGN KEY (`businessentityid`) REFERENCES `person`.`person` (`businessentityid`);

ALTER TABLE `person`.`personphone` ADD CONSTRAINT `FK_PersonPhone_Person_BusinessEntityID` FOREIGN KEY (`businessentityid`) REFERENCES `person`.`person` (`businessentityid`);

ALTER TABLE `person`.`personphone` ADD CONSTRAINT `FK_PersonPhone_PhoneNumberType_PhoneNumberTypeID` FOREIGN KEY (`phonenumbertypeid`) REFERENCES `person`.`phonenumbertype` (`phonenumbertypeid`);

ALTER TABLE `person`.`person` ADD CONSTRAINT `FK_Person_BusinessEntity_BusinessEntityID` FOREIGN KEY (`businessentityid`) REFERENCES `person`.`businessentity` (`businessentityid`);

ALTER TABLE `person`.`stateprovince` ADD CONSTRAINT `FK_StateProvince_CountryRegion_CountryRegionCode` FOREIGN KEY (`countryregioncode`) REFERENCES `person`.`countryregion` (`countryregioncode`);

ALTER TABLE `person`.`stateprovince` ADD CONSTRAINT `FK_StateProvince_SalesTerritory_TerritoryID` FOREIGN KEY (`territoryid`) REFERENCES `sales`.`salesterritory` (`territoryid`);

ALTER TABLE `production`.`billofmaterials` ADD CONSTRAINT `FK_BillOfMaterials_Product_ComponentID` FOREIGN KEY (`componentid`) REFERENCES `production`.`product` (`productid`);

ALTER TABLE `production`.`billofmaterials` ADD CONSTRAINT `FK_BillOfMaterials_Product_ProductAssemblyID` FOREIGN KEY (`productassemblyid`) REFERENCES `production`.`product` (`productid`);

ALTER TABLE `production`.`billofmaterials` ADD CONSTRAINT `FK_BillOfMaterials_UnitMeasure_UnitMeasureCode` FOREIGN KEY (`unitmeasurecode`) REFERENCES `production`.`unitmeasure` (`unitmeasurecode`);

ALTER TABLE `production`.`document` ADD CONSTRAINT `FK_Document_Employee_Owner` FOREIGN KEY (`owner`) REFERENCES `humanresources`.`employee` (`businessentityid`);

ALTER TABLE `production`.`productcosthistory` ADD CONSTRAINT `FK_ProductCostHistory_Product_ProductID` FOREIGN KEY (`productid`) REFERENCES `production`.`product` (`productid`);

ALTER TABLE `production`.`productdocument` ADD CONSTRAINT `FK_ProductDocument_Document_DocumentNode` FOREIGN KEY (`documentnode`) REFERENCES `production`.`document` (`documentnode`);

ALTER TABLE `production`.`productdocument` ADD CONSTRAINT `FK_ProductDocument_Product_ProductID` FOREIGN KEY (`productid`) REFERENCES `production`.`product` (`productid`);

ALTER TABLE `production`.`productinventory` ADD CONSTRAINT `FK_ProductInventory_Location_LocationID` FOREIGN KEY (`locationid`) REFERENCES `production`.`location` (`locationid`);

ALTER TABLE `production`.`productinventory` ADD CONSTRAINT `FK_ProductInventory_Product_ProductID` FOREIGN KEY (`productid`) REFERENCES `production`.`product` (`productid`);

ALTER TABLE `production`.`productlistpricehistory` ADD CONSTRAINT `FK_ProductListPriceHistory_Product_ProductID` FOREIGN KEY (`productid`) REFERENCES `production`.`product` (`productid`);

ALTER TABLE `production`.`productmodelillustration` ADD CONSTRAINT `FK_ProductModelIllustration_Illustration_IllustrationID` FOREIGN KEY (`illustrationid`) REFERENCES `production`.`illustration` (`illustrationid`);

ALTER TABLE `production`.`productmodelillustration` ADD CONSTRAINT `FK_ProductModelIllustration_ProductModel_ProductModelID` FOREIGN KEY (`productmodelid`) REFERENCES `production`.`productmodel` (`productmodelid`);

ALTER TABLE `production`.`productmodelproductdescriptionculture` ADD CONSTRAINT `FK_ProductModelProductDescriptionCulture_Culture_CultureID` FOREIGN KEY (`cultureid`) REFERENCES `production`.`culture` (`cultureid`);

ALTER TABLE `production`.`productmodelproductdescriptionculture` ADD CONSTRAINT `FK_ProductModelProductDescriptionCulture_ProductDescription_Pro` FOREIGN KEY (`productdescriptionid`) REFERENCES `production`.`productdescription` (`productdescriptionid`);

ALTER TABLE `production`.`productmodelproductdescriptionculture` ADD CONSTRAINT `FK_ProductModelProductDescriptionCulture_ProductModel_ProductMo` FOREIGN KEY (`productmodelid`) REFERENCES `production`.`productmodel` (`productmodelid`);

ALTER TABLE `production`.`productproductphoto` ADD CONSTRAINT `FK_ProductProductPhoto_ProductPhoto_ProductPhotoID` FOREIGN KEY (`productphotoid`) REFERENCES `production`.`productphoto` (`productphotoid`);

ALTER TABLE `production`.`productproductphoto` ADD CONSTRAINT `FK_ProductProductPhoto_Product_ProductID` FOREIGN KEY (`productid`) REFERENCES `production`.`product` (`productid`);

ALTER TABLE `production`.`productreview` ADD CONSTRAINT `FK_ProductReview_Product_ProductID` FOREIGN KEY (`productid`) REFERENCES `production`.`product` (`productid`);

ALTER TABLE `production`.`productsubcategory` ADD CONSTRAINT `FK_ProductSubcategory_ProductCategory_ProductCategoryID` FOREIGN KEY (`productcategoryid`) REFERENCES `production`.`productcategory` (`productcategoryid`);

ALTER TABLE `production`.`product` ADD CONSTRAINT `FK_Product_ProductModel_ProductModelID` FOREIGN KEY (`productmodelid`) REFERENCES `production`.`productmodel` (`productmodelid`);

ALTER TABLE `production`.`product` ADD CONSTRAINT `FK_Product_ProductSubcategory_ProductSubcategoryID` FOREIGN KEY (`productsubcategoryid`) REFERENCES `production`.`productsubcategory` (`productsubcategoryid`);

ALTER TABLE `production`.`product` ADD CONSTRAINT `FK_Product_UnitMeasure_SizeUnitMeasureCode` FOREIGN KEY (`sizeunitmeasurecode`) REFERENCES `production`.`unitmeasure` (`unitmeasurecode`);

ALTER TABLE `production`.`product` ADD CONSTRAINT `FK_Product_UnitMeasure_WeightUnitMeasureCode` FOREIGN KEY (`weightunitmeasurecode`) REFERENCES `production`.`unitmeasure` (`unitmeasurecode`);

ALTER TABLE `production`.`transactionhistory` ADD CONSTRAINT `FK_TransactionHistory_Product_ProductID` FOREIGN KEY (`productid`) REFERENCES `production`.`product` (`productid`);

ALTER TABLE `production`.`workorderrouting` ADD CONSTRAINT `FK_WorkOrderRouting_Location_LocationID` FOREIGN KEY (`locationid`) REFERENCES `production`.`location` (`locationid`);

ALTER TABLE `production`.`workorderrouting` ADD CONSTRAINT `FK_WorkOrderRouting_WorkOrder_WorkOrderID` FOREIGN KEY (`workorderid`) REFERENCES `production`.`workorder` (`workorderid`);

ALTER TABLE `production`.`workorder` ADD CONSTRAINT `FK_WorkOrder_Product_ProductID` FOREIGN KEY (`productid`) REFERENCES `production`.`product` (`productid`);

ALTER TABLE `production`.`workorder` ADD CONSTRAINT `FK_WorkOrder_ScrapReason_ScrapReasonID` FOREIGN KEY (`scrapreasonid`) REFERENCES `production`.`scrapreason` (`scrapreasonid`);

ALTER TABLE `purchasing`.`productvendor` ADD CONSTRAINT `FK_ProductVendor_Product_ProductID` FOREIGN KEY (`productid`) REFERENCES `production`.`product` (`productid`);

ALTER TABLE `purchasing`.`productvendor` ADD CONSTRAINT `FK_ProductVendor_UnitMeasure_UnitMeasureCode` FOREIGN KEY (`unitmeasurecode`) REFERENCES `production`.`unitmeasure` (`unitmeasurecode`);

ALTER TABLE `purchasing`.`productvendor` ADD CONSTRAINT `FK_ProductVendor_Vendor_BusinessEntityID` FOREIGN KEY (`businessentityid`) REFERENCES `purchasing`.`vendor` (`businessentityid`);

ALTER TABLE `purchasing`.`purchaseorderdetail` ADD CONSTRAINT `FK_PurchaseOrderDetail_Product_ProductID` FOREIGN KEY (`productid`) REFERENCES `production`.`product` (`productid`);

ALTER TABLE `purchasing`.`purchaseorderdetail` ADD CONSTRAINT `FK_PurchaseOrderDetail_PurchaseOrderHeader_PurchaseOrderID` FOREIGN KEY (`purchaseorderid`) REFERENCES `purchasing`.`purchaseorderheader` (`purchaseorderid`);

ALTER TABLE `purchasing`.`purchaseorderheader` ADD CONSTRAINT `FK_PurchaseOrderHeader_Employee_EmployeeID` FOREIGN KEY (`employeeid`) REFERENCES `humanresources`.`employee` (`businessentityid`);

ALTER TABLE `purchasing`.`purchaseorderheader` ADD CONSTRAINT `FK_PurchaseOrderHeader_ShipMethod_ShipMethodID` FOREIGN KEY (`shipmethodid`) REFERENCES `purchasing`.`shipmethod` (`shipmethodid`);

ALTER TABLE `purchasing`.`purchaseorderheader` ADD CONSTRAINT `FK_PurchaseOrderHeader_Vendor_VendorID` FOREIGN KEY (`vendorid`) REFERENCES `purchasing`.`vendor` (`businessentityid`);

ALTER TABLE `purchasing`.`vendor` ADD CONSTRAINT `FK_Vendor_BusinessEntity_BusinessEntityID` FOREIGN KEY (`businessentityid`) REFERENCES `person`.`businessentity` (`businessentityid`);

ALTER TABLE `sales`.`countryregioncurrency` ADD CONSTRAINT `FK_CountryRegionCurrency_CountryRegion_CountryRegionCode` FOREIGN KEY (`countryregioncode`) REFERENCES `person`.`countryregion` (`countryregioncode`);

ALTER TABLE `sales`.`countryregioncurrency` ADD CONSTRAINT `FK_CountryRegionCurrency_Currency_CurrencyCode` FOREIGN KEY (`currencycode`) REFERENCES `sales`.`currency` (`currencycode`);

ALTER TABLE `sales`.`currencyrate` ADD CONSTRAINT `FK_CurrencyRate_Currency_FromCurrencyCode` FOREIGN KEY (`fromcurrencycode`) REFERENCES `sales`.`currency` (`currencycode`);

ALTER TABLE `sales`.`currencyrate` ADD CONSTRAINT `FK_CurrencyRate_Currency_ToCurrencyCode` FOREIGN KEY (`tocurrencycode`) REFERENCES `sales`.`currency` (`currencycode`);

ALTER TABLE `sales`.`customer` ADD CONSTRAINT `FK_Customer_Person_PersonID` FOREIGN KEY (`personid`) REFERENCES `person`.`person` (`businessentityid`);

ALTER TABLE `sales`.`customer` ADD CONSTRAINT `FK_Customer_SalesTerritory_TerritoryID` FOREIGN KEY (`territoryid`) REFERENCES `sales`.`salesterritory` (`territoryid`);

ALTER TABLE `sales`.`customer` ADD CONSTRAINT `FK_Customer_Store_StoreID` FOREIGN KEY (`storeid`) REFERENCES `sales`.`store` (`businessentityid`);

ALTER TABLE `sales`.`personcreditcard` ADD CONSTRAINT `FK_PersonCreditCard_CreditCard_CreditCardID` FOREIGN KEY (`creditcardid`) REFERENCES `sales`.`creditcard` (`creditcardid`);

ALTER TABLE `sales`.`personcreditcard` ADD CONSTRAINT `FK_PersonCreditCard_Person_BusinessEntityID` FOREIGN KEY (`businessentityid`) REFERENCES `person`.`person` (`businessentityid`);

ALTER TABLE `sales`.`salesorderdetail` ADD CONSTRAINT `FK_SalesOrderDetail_SalesOrderHeader_SalesOrderID` FOREIGN KEY (`salesorderid`) REFERENCES `sales`.`salesorderheader` (`salesorderid`) ON DELETE CASCADE;

ALTER TABLE `sales`.`salesorderdetail` ADD CONSTRAINT `FK_SalesOrderDetail_SpecialOfferProduct_SpecialOfferIDProductID` FOREIGN KEY (`specialofferid`, `productid`) REFERENCES `sales`.`specialofferproduct` (`specialofferid`, `productid`);

ALTER TABLE `sales`.`salesorderheadersalesreason` ADD CONSTRAINT `FK_SalesOrderHeaderSalesReason_SalesOrderHeader_SalesOrderID` FOREIGN KEY (`salesorderid`) REFERENCES `sales`.`salesorderheader` (`salesorderid`) ON DELETE CASCADE;

ALTER TABLE `sales`.`salesorderheadersalesreason` ADD CONSTRAINT `FK_SalesOrderHeaderSalesReason_SalesReason_SalesReasonID` FOREIGN KEY (`salesreasonid`) REFERENCES `sales`.`salesreason` (`salesreasonid`);

ALTER TABLE `sales`.`salesorderheader` ADD CONSTRAINT `FK_SalesOrderHeader_Address_BillToAddressID` FOREIGN KEY (`billtoaddressid`) REFERENCES `person`.`address` (`addressid`);

ALTER TABLE `sales`.`salesorderheader` ADD CONSTRAINT `FK_SalesOrderHeader_Address_ShipToAddressID` FOREIGN KEY (`shiptoaddressid`) REFERENCES `person`.`address` (`addressid`);

ALTER TABLE `sales`.`salesorderheader` ADD CONSTRAINT `FK_SalesOrderHeader_CreditCard_CreditCardID` FOREIGN KEY (`creditcardid`) REFERENCES `sales`.`creditcard` (`creditcardid`);

ALTER TABLE `sales`.`salesorderheader` ADD CONSTRAINT `FK_SalesOrderHeader_CurrencyRate_CurrencyRateID` FOREIGN KEY (`currencyrateid`) REFERENCES `sales`.`currencyrate` (`currencyrateid`);

ALTER TABLE `sales`.`salesorderheader` ADD CONSTRAINT `FK_SalesOrderHeader_Customer_CustomerID` FOREIGN KEY (`customerid`) REFERENCES `sales`.`customer` (`customerid`);

ALTER TABLE `sales`.`salesorderheader` ADD CONSTRAINT `FK_SalesOrderHeader_SalesPerson_SalesPersonID` FOREIGN KEY (`salespersonid`) REFERENCES `sales`.`salesperson` (`businessentityid`);

ALTER TABLE `sales`.`salesorderheader` ADD CONSTRAINT `FK_SalesOrderHeader_SalesTerritory_TerritoryID` FOREIGN KEY (`territoryid`) REFERENCES `sales`.`salesterritory` (`territoryid`);

ALTER TABLE `sales`.`salesorderheader` ADD CONSTRAINT `FK_SalesOrderHeader_ShipMethod_ShipMethodID` FOREIGN KEY (`shipmethodid`) REFERENCES `purchasing`.`shipmethod` (`shipmethodid`);

ALTER TABLE `sales`.`salespersonquotahistory` ADD CONSTRAINT `FK_SalesPersonQuotaHistory_SalesPerson_BusinessEntityID` FOREIGN KEY (`businessentityid`) REFERENCES `sales`.`salesperson` (`businessentityid`);

ALTER TABLE `sales`.`salesperson` ADD CONSTRAINT `FK_SalesPerson_Employee_BusinessEntityID` FOREIGN KEY (`businessentityid`) REFERENCES `humanresources`.`employee` (`businessentityid`);

ALTER TABLE `sales`.`salesperson` ADD CONSTRAINT `FK_SalesPerson_SalesTerritory_TerritoryID` FOREIGN KEY (`territoryid`) REFERENCES `sales`.`salesterritory` (`territoryid`);

ALTER TABLE `sales`.`salestaxrate` ADD CONSTRAINT `FK_SalesTaxRate_StateProvince_StateProvinceID` FOREIGN KEY (`stateprovinceid`) REFERENCES `person`.`stateprovince` (`stateprovinceid`);

ALTER TABLE `sales`.`salesterritoryhistory` ADD CONSTRAINT `FK_SalesTerritoryHistory_SalesPerson_BusinessEntityID` FOREIGN KEY (`businessentityid`) REFERENCES `sales`.`salesperson` (`businessentityid`);

ALTER TABLE `sales`.`salesterritoryhistory` ADD CONSTRAINT `FK_SalesTerritoryHistory_SalesTerritory_TerritoryID` FOREIGN KEY (`territoryid`) REFERENCES `sales`.`salesterritory` (`territoryid`);

ALTER TABLE `sales`.`salesterritory` ADD CONSTRAINT `FK_SalesTerritory_CountryRegion_CountryRegionCode` FOREIGN KEY (`countryregioncode`) REFERENCES `person`.`countryregion` (`countryregioncode`);

ALTER TABLE `sales`.`shoppingcartitem` ADD CONSTRAINT `FK_ShoppingCartItem_Product_ProductID` FOREIGN KEY (`productid`) REFERENCES `production`.`product` (`productid`);

ALTER TABLE `sales`.`specialofferproduct` ADD CONSTRAINT `FK_SpecialOfferProduct_Product_ProductID` FOREIGN KEY (`productid`) REFERENCES `production`.`product` (`productid`);

ALTER TABLE `sales`.`specialofferproduct` ADD CONSTRAINT `FK_SpecialOfferProduct_SpecialOffer_SpecialOfferID` FOREIGN KEY (`specialofferid`) REFERENCES `sales`.`specialoffer` (`specialofferid`);

ALTER TABLE `sales`.`store` ADD CONSTRAINT `FK_Store_BusinessEntity_BusinessEntityID` FOREIGN KEY (`businessentityid`) REFERENCES `person`.`businessentity` (`businessentityid`);

ALTER TABLE `sales`.`store` ADD CONSTRAINT `FK_Store_SalesPerson_SalesPersonID` FOREIGN KEY (`salespersonid`) REFERENCES `sales`.`salesperson` (`businessentityid`);
