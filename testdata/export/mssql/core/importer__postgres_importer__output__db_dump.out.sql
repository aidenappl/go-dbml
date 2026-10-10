CREATE SCHEMA [humanresources]
GO

CREATE SCHEMA [person]
GO

CREATE SCHEMA [production]
GO

CREATE SCHEMA [purchasing]
GO

CREATE SCHEMA [sales]
GO

CREATE TABLE [humanresources].[department] (
  [departmentid] integer PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [groupname] public.Name NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [humanresources].[employee] (
  [businessentityid] integer PRIMARY KEY NOT NULL,
  [nationalidnumber] character varying(15) NOT NULL,
  [loginid] character varying(256) NOT NULL,
  [jobtitle] character varying(50) NOT NULL,
  [birthdate] date NOT NULL,
  [maritalstatus] character(1) NOT NULL,
  [gender] character(1) NOT NULL,
  [hiredate] date NOT NULL,
  [salariedflag] public.Flag NOT NULL DEFAULT (true),
  [vacationhours] smallint NOT NULL DEFAULT (0),
  [sickleavehours] smallint NOT NULL DEFAULT (0),
  [currentflag] public.Flag NOT NULL DEFAULT (true),
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  [organizationnode] character varying DEFAULT ('/'::charactervarying),
  CONSTRAINT [CK_Employee_BirthDate] CHECK (((birthdate >= '1930-01-01'::date) AND (birthdate <= (now() - '18 years'::interval)))),
  CONSTRAINT [CK_Employee_Gender] CHECK ((upper((gender)::text) = ANY (ARRAY['M'::text, 'F'::text]))),
  CONSTRAINT [CK_Employee_HireDate] CHECK (((hiredate >= '1996-07-01'::date) AND (hiredate <= (now() + '1 day'::interval)))),
  CONSTRAINT [CK_Employee_MaritalStatus] CHECK ((upper((maritalstatus)::text) = ANY (ARRAY['M'::text, 'S'::text]))),
  CONSTRAINT [CK_Employee_SickLeaveHours] CHECK (((sickleavehours >= 0) AND (sickleavehours <= 120))),
  CONSTRAINT [CK_Employee_VacationHours] CHECK (((vacationhours >= '-40'::integer) AND (vacationhours <= 240)))
)
GO

CREATE TABLE [humanresources].[employeedepartmenthistory] (
  [businessentityid] integer NOT NULL,
  [departmentid] smallint NOT NULL,
  [shiftid] smallint NOT NULL,
  [startdate] date NOT NULL,
  [enddate] date,
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_EmployeeDepartmentHistory_EndDate] CHECK (((enddate >= startdate) OR (enddate IS NULL))),
  PRIMARY KEY ([businessentityid], [startdate], [departmentid], [shiftid])
)
GO

CREATE TABLE [humanresources].[employeepayhistory] (
  [businessentityid] integer NOT NULL,
  [ratechangedate] timestamp NOT NULL,
  [rate] numeric NOT NULL,
  [payfrequency] smallint NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_EmployeePayHistory_PayFrequency] CHECK ((payfrequency = ANY (ARRAY[1, 2]))),
  CONSTRAINT [CK_EmployeePayHistory_Rate] CHECK (((rate >= 6.50) AND (rate <= 200.00))),
  PRIMARY KEY ([businessentityid], [ratechangedate])
)
GO

CREATE TABLE [humanresources].[jobcandidate] (
  [jobcandidateid] integer PRIMARY KEY NOT NULL,
  [businessentityid] integer,
  [resume] xml,
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [humanresources].[shift] (
  [shiftid] integer PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [starttime] time NOT NULL,
  [endtime] time NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [person].[address] (
  [addressid] integer PRIMARY KEY NOT NULL,
  [addressline1] character varying(60) NOT NULL,
  [addressline2] character varying(60),
  [city] character varying(30) NOT NULL,
  [stateprovinceid] integer NOT NULL,
  [postalcode] character varying(15) NOT NULL,
  [spatiallocation] character varying(44),
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [person].[businessentityaddress] (
  [businessentityid] integer NOT NULL,
  [addressid] integer NOT NULL,
  [addresstypeid] integer NOT NULL,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ([businessentityid], [addressid], [addresstypeid])
)
GO

CREATE TABLE [person].[countryregion] (
  [countryregioncode] character varying(3) PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [person].[emailaddress] (
  [businessentityid] integer NOT NULL,
  [emailaddressid] integer NOT NULL,
  [emailaddress] character varying(50),
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ([businessentityid], [emailaddressid])
)
GO

CREATE TABLE [person].[person] (
  [businessentityid] integer PRIMARY KEY NOT NULL,
  [persontype] character(2) NOT NULL,
  [namestyle] public.NameStyle NOT NULL DEFAULT (false),
  [title] character varying(8),
  [firstname] public.Name NOT NULL,
  [middlename] public.Name,
  [lastname] public.Name NOT NULL,
  [suffix] character varying(10),
  [emailpromotion] integer NOT NULL DEFAULT (0),
  [additionalcontactinfo] xml,
  [demographics] xml,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_Person_EmailPromotion] CHECK (((emailpromotion >= 0) AND (emailpromotion <= 2))),
  CONSTRAINT [CK_Person_PersonType] CHECK (((persontype IS NULL) OR (upper((persontype)::text) = ANY (ARRAY['SC'::text, 'VC'::text, 'IN'::text, 'EM'::text, 'SP'::text, 'GC'::text]))))
)
GO

CREATE TABLE [person].[personphone] (
  [businessentityid] integer NOT NULL,
  [phonenumber] public.Phone NOT NULL,
  [phonenumbertypeid] integer NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ([businessentityid], [phonenumber], [phonenumbertypeid])
)
GO

CREATE TABLE [person].[phonenumbertype] (
  [phonenumbertypeid] integer PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [person].[stateprovince] (
  [stateprovinceid] integer PRIMARY KEY NOT NULL,
  [stateprovincecode] character(3) NOT NULL,
  [countryregioncode] character varying(3) NOT NULL,
  [isonlystateprovinceflag] public.Flag NOT NULL DEFAULT (true),
  [name] public.Name NOT NULL,
  [territoryid] integer NOT NULL,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [person].[addresstype] (
  [addresstypeid] integer PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [person].[businessentity] (
  [businessentityid] integer PRIMARY KEY NOT NULL,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [person].[businessentitycontact] (
  [businessentityid] integer NOT NULL,
  [personid] integer NOT NULL,
  [contacttypeid] integer NOT NULL,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ([businessentityid], [personid], [contacttypeid])
)
GO

CREATE TABLE [person].[contacttype] (
  [contacttypeid] integer PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [person].[password] (
  [businessentityid] integer PRIMARY KEY NOT NULL,
  [passwordhash] character varying(128) NOT NULL,
  [passwordsalt] character varying(10) NOT NULL,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [production].[billofmaterials] (
  [billofmaterialsid] integer PRIMARY KEY NOT NULL,
  [productassemblyid] integer,
  [componentid] integer NOT NULL,
  [startdate] timestamp NOT NULL DEFAULT (now()),
  [enddate] timestamp,
  [unitmeasurecode] character(3) NOT NULL,
  [bomlevel] smallint NOT NULL,
  [perassemblyqty] numeric(8,2) NOT NULL DEFAULT (1),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_BillOfMaterials_BOMLevel] CHECK ((((productassemblyid IS NULL) AND (bomlevel = 0) AND (perassemblyqty = 1.00)) OR ((productassemblyid IS NOT NULL) AND (bomlevel >= 1)))),
  CONSTRAINT [CK_BillOfMaterials_EndDate] CHECK (((enddate > startdate) OR (enddate IS NULL))),
  CONSTRAINT [CK_BillOfMaterials_PerAssemblyQty] CHECK ((perassemblyqty >= 1.00)),
  CONSTRAINT [CK_BillOfMaterials_ProductAssemblyID] CHECK ((productassemblyid <> componentid))
)
GO

CREATE TABLE [production].[culture] (
  [cultureid] character(6) PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [production].[document] (
  [title] character varying(50) NOT NULL,
  [owner] integer NOT NULL,
  [folderflag] public.Flag NOT NULL DEFAULT (false),
  [filename] character varying(400) NOT NULL,
  [fileextension] character varying(8),
  [revision] character(5) NOT NULL,
  [changenumber] integer NOT NULL DEFAULT (0),
  [status] smallint NOT NULL,
  [documentsummary] text,
  [document] bytea,
  [rowguid] uuid UNIQUE NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  [documentnode] character varying PRIMARY KEY NOT NULL DEFAULT ('/'::charactervarying),
  CONSTRAINT [CK_Document_Status] CHECK (((status >= 1) AND (status <= 3)))
)
GO

CREATE TABLE [production].[illustration] (
  [illustrationid] integer PRIMARY KEY NOT NULL,
  [diagram] xml,
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [production].[location] (
  [locationid] integer PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [costrate] numeric NOT NULL DEFAULT (0),
  [availability] numeric(8,2) NOT NULL DEFAULT (0),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_Location_Availability] CHECK ((availability >= 0.00)),
  CONSTRAINT [CK_Location_CostRate] CHECK ((costrate >= 0.00))
)
GO

CREATE TABLE [production].[product] (
  [productid] integer PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [productnumber] character varying(25) NOT NULL,
  [makeflag] public.Flag NOT NULL DEFAULT (true),
  [finishedgoodsflag] public.Flag NOT NULL DEFAULT (true),
  [color] character varying(15),
  [safetystocklevel] smallint NOT NULL,
  [reorderpoint] smallint NOT NULL,
  [standardcost] numeric NOT NULL,
  [listprice] numeric NOT NULL,
  [size] character varying(5),
  [sizeunitmeasurecode] character(3),
  [weightunitmeasurecode] character(3),
  [weight] numeric(8,2),
  [daystomanufacture] integer NOT NULL,
  [productline] character(2),
  [class] character(2),
  [style] character(2),
  [productsubcategoryid] integer,
  [productmodelid] integer,
  [sellstartdate] timestamp NOT NULL,
  [sellenddate] timestamp,
  [discontinueddate] timestamp,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_Product_Class] CHECK (((upper((class)::text) = ANY (ARRAY['L'::text, 'M'::text, 'H'::text])) OR (class IS NULL))),
  CONSTRAINT [CK_Product_DaysToManufacture] CHECK ((daystomanufacture >= 0)),
  CONSTRAINT [CK_Product_ListPrice] CHECK ((listprice >= 0.00)),
  CONSTRAINT [CK_Product_ProductLine] CHECK (((upper((productline)::text) = ANY (ARRAY['S'::text, 'T'::text, 'M'::text, 'R'::text])) OR (productline IS NULL))),
  CONSTRAINT [CK_Product_ReorderPoint] CHECK ((reorderpoint > 0)),
  CONSTRAINT [CK_Product_SafetyStockLevel] CHECK ((safetystocklevel > 0)),
  CONSTRAINT [CK_Product_SellEndDate] CHECK (((sellenddate >= sellstartdate) OR (sellenddate IS NULL))),
  CONSTRAINT [CK_Product_StandardCost] CHECK ((standardcost >= 0.00)),
  CONSTRAINT [CK_Product_Style] CHECK (((upper((style)::text) = ANY (ARRAY['W'::text, 'M'::text, 'U'::text])) OR (style IS NULL))),
  CONSTRAINT [CK_Product_Weight] CHECK ((weight > 0.00))
)
GO

CREATE TABLE [production].[productcategory] (
  [productcategoryid] integer PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [production].[productcosthistory] (
  [productid] integer NOT NULL,
  [startdate] timestamp NOT NULL,
  [enddate] timestamp,
  [standardcost] numeric NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_ProductCostHistory_EndDate] CHECK (((enddate >= startdate) OR (enddate IS NULL))),
  CONSTRAINT [CK_ProductCostHistory_StandardCost] CHECK ((standardcost >= 0.00)),
  PRIMARY KEY ([productid], [startdate])
)
GO

CREATE TABLE [production].[productdescription] (
  [productdescriptionid] integer PRIMARY KEY NOT NULL,
  [description] character varying(400) NOT NULL,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [production].[productdocument] (
  [productid] integer NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  [documentnode] character varying NOT NULL DEFAULT ('/'::charactervarying),
  PRIMARY KEY ([productid], [documentnode])
)
GO

CREATE TABLE [production].[productinventory] (
  [productid] integer NOT NULL,
  [locationid] smallint NOT NULL,
  [shelf] character varying(10) NOT NULL,
  [bin] smallint NOT NULL,
  [quantity] smallint NOT NULL DEFAULT (0),
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_ProductInventory_Bin] CHECK (((bin >= 0) AND (bin <= 100))),
  PRIMARY KEY ([productid], [locationid])
)
GO

CREATE TABLE [production].[productlistpricehistory] (
  [productid] integer NOT NULL,
  [startdate] timestamp NOT NULL,
  [enddate] timestamp,
  [listprice] numeric NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_ProductListPriceHistory_EndDate] CHECK (((enddate >= startdate) OR (enddate IS NULL))),
  CONSTRAINT [CK_ProductListPriceHistory_ListPrice] CHECK ((listprice > 0.00)),
  PRIMARY KEY ([productid], [startdate])
)
GO

CREATE TABLE [production].[productmodel] (
  [productmodelid] integer PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [catalogdescription] xml,
  [instructions] xml,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [production].[productmodelillustration] (
  [productmodelid] integer NOT NULL,
  [illustrationid] integer NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ([productmodelid], [illustrationid])
)
GO

CREATE TABLE [production].[productmodelproductdescriptionculture] (
  [productmodelid] integer NOT NULL,
  [productdescriptionid] integer NOT NULL,
  [cultureid] character(6) NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ([productmodelid], [productdescriptionid], [cultureid])
)
GO

CREATE TABLE [production].[productphoto] (
  [productphotoid] integer PRIMARY KEY NOT NULL,
  [thumbnailphoto] bytea,
  [thumbnailphotofilename] character varying(50),
  [largephoto] bytea,
  [largephotofilename] character varying(50),
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [production].[productproductphoto] (
  [productid] integer NOT NULL,
  [productphotoid] integer NOT NULL,
  [primary] public.Flag NOT NULL DEFAULT (false),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ([productid], [productphotoid])
)
GO

CREATE TABLE [production].[productreview] (
  [productreviewid] integer PRIMARY KEY NOT NULL,
  [productid] integer NOT NULL,
  [reviewername] public.Name NOT NULL,
  [reviewdate] timestamp NOT NULL DEFAULT (now()),
  [emailaddress] character varying(50) NOT NULL,
  [rating] integer NOT NULL,
  [comments] character varying(3850),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_ProductReview_Rating] CHECK (((rating >= 1) AND (rating <= 5)))
)
GO

CREATE TABLE [production].[productsubcategory] (
  [productsubcategoryid] integer PRIMARY KEY NOT NULL,
  [productcategoryid] integer NOT NULL,
  [name] public.Name NOT NULL,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [production].[scrapreason] (
  [scrapreasonid] integer PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [production].[transactionhistory] (
  [transactionid] integer PRIMARY KEY NOT NULL,
  [productid] integer NOT NULL,
  [referenceorderid] integer NOT NULL,
  [referenceorderlineid] integer NOT NULL DEFAULT (0),
  [transactiondate] timestamp NOT NULL DEFAULT (now()),
  [transactiontype] character(1) NOT NULL,
  [quantity] integer NOT NULL,
  [actualcost] numeric NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_TransactionHistory_TransactionType] CHECK ((upper((transactiontype)::text) = ANY (ARRAY['W'::text, 'S'::text, 'P'::text])))
)
GO

CREATE TABLE [production].[transactionhistoryarchive] (
  [transactionid] integer PRIMARY KEY NOT NULL,
  [productid] integer NOT NULL,
  [referenceorderid] integer NOT NULL,
  [referenceorderlineid] integer NOT NULL DEFAULT (0),
  [transactiondate] timestamp NOT NULL DEFAULT (now()),
  [transactiontype] character(1) NOT NULL,
  [quantity] integer NOT NULL,
  [actualcost] numeric NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_TransactionHistoryArchive_TransactionType] CHECK ((upper((transactiontype)::text) = ANY (ARRAY['W'::text, 'S'::text, 'P'::text])))
)
GO

CREATE TABLE [production].[unitmeasure] (
  [unitmeasurecode] character(3) PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [production].[workorder] (
  [workorderid] integer PRIMARY KEY NOT NULL,
  [productid] integer NOT NULL,
  [orderqty] integer NOT NULL,
  [scrappedqty] smallint NOT NULL,
  [startdate] timestamp NOT NULL,
  [enddate] timestamp,
  [duedate] timestamp NOT NULL,
  [scrapreasonid] smallint,
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_WorkOrder_EndDate] CHECK (((enddate >= startdate) OR (enddate IS NULL))),
  CONSTRAINT [CK_WorkOrder_OrderQty] CHECK ((orderqty > 0)),
  CONSTRAINT [CK_WorkOrder_ScrappedQty] CHECK ((scrappedqty >= 0))
)
GO

CREATE TABLE [production].[workorderrouting] (
  [workorderid] integer NOT NULL,
  [productid] integer NOT NULL,
  [operationsequence] smallint NOT NULL,
  [locationid] smallint NOT NULL,
  [scheduledstartdate] timestamp NOT NULL,
  [scheduledenddate] timestamp NOT NULL,
  [actualstartdate] timestamp,
  [actualenddate] timestamp,
  [actualresourcehrs] numeric(9,4),
  [plannedcost] numeric NOT NULL,
  [actualcost] numeric,
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_WorkOrderRouting_ActualCost] CHECK ((actualcost > 0.00)),
  CONSTRAINT [CK_WorkOrderRouting_ActualEndDate] CHECK (((actualenddate >= actualstartdate) OR (actualenddate IS NULL) OR (actualstartdate IS NULL))),
  CONSTRAINT [CK_WorkOrderRouting_ActualResourceHrs] CHECK ((actualresourcehrs >= 0.0000)),
  CONSTRAINT [CK_WorkOrderRouting_PlannedCost] CHECK ((plannedcost > 0.00)),
  CONSTRAINT [CK_WorkOrderRouting_ScheduledEndDate] CHECK ((scheduledenddate >= scheduledstartdate)),
  PRIMARY KEY ([workorderid], [productid], [operationsequence])
)
GO

CREATE TABLE [purchasing].[purchaseorderdetail] (
  [purchaseorderid] integer NOT NULL,
  [purchaseorderdetailid] integer NOT NULL,
  [duedate] timestamp NOT NULL,
  [orderqty] smallint NOT NULL,
  [productid] integer NOT NULL,
  [unitprice] numeric NOT NULL,
  [receivedqty] numeric(8,2) NOT NULL,
  [rejectedqty] numeric(8,2) NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_PurchaseOrderDetail_OrderQty] CHECK ((orderqty > 0)),
  CONSTRAINT [CK_PurchaseOrderDetail_ReceivedQty] CHECK ((receivedqty >= 0.00)),
  CONSTRAINT [CK_PurchaseOrderDetail_RejectedQty] CHECK ((rejectedqty >= 0.00)),
  CONSTRAINT [CK_PurchaseOrderDetail_UnitPrice] CHECK ((unitprice >= 0.00)),
  PRIMARY KEY ([purchaseorderid], [purchaseorderdetailid])
)
GO

CREATE TABLE [purchasing].[purchaseorderheader] (
  [purchaseorderid] integer PRIMARY KEY NOT NULL,
  [revisionnumber] smallint NOT NULL DEFAULT (0),
  [status] smallint NOT NULL DEFAULT (1),
  [employeeid] integer NOT NULL,
  [vendorid] integer NOT NULL,
  [shipmethodid] integer NOT NULL,
  [orderdate] timestamp NOT NULL DEFAULT (now()),
  [shipdate] timestamp,
  [subtotal] numeric NOT NULL DEFAULT (0),
  [taxamt] numeric NOT NULL DEFAULT (0),
  [freight] numeric NOT NULL DEFAULT (0),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_PurchaseOrderHeader_Freight] CHECK ((freight >= 0.00)),
  CONSTRAINT [CK_PurchaseOrderHeader_ShipDate] CHECK (((shipdate >= orderdate) OR (shipdate IS NULL))),
  CONSTRAINT [CK_PurchaseOrderHeader_Status] CHECK (((status >= 1) AND (status <= 4))),
  CONSTRAINT [CK_PurchaseOrderHeader_SubTotal] CHECK ((subtotal >= 0.00)),
  CONSTRAINT [CK_PurchaseOrderHeader_TaxAmt] CHECK ((taxamt >= 0.00))
)
GO

CREATE TABLE [purchasing].[productvendor] (
  [productid] integer NOT NULL,
  [businessentityid] integer NOT NULL,
  [averageleadtime] integer NOT NULL,
  [standardprice] numeric NOT NULL,
  [lastreceiptcost] numeric,
  [lastreceiptdate] timestamp,
  [minorderqty] integer NOT NULL,
  [maxorderqty] integer NOT NULL,
  [onorderqty] integer,
  [unitmeasurecode] character(3) NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_ProductVendor_AverageLeadTime] CHECK ((averageleadtime >= 1)),
  CONSTRAINT [CK_ProductVendor_LastReceiptCost] CHECK ((lastreceiptcost > 0.00)),
  CONSTRAINT [CK_ProductVendor_MaxOrderQty] CHECK ((maxorderqty >= 1)),
  CONSTRAINT [CK_ProductVendor_MinOrderQty] CHECK ((minorderqty >= 1)),
  CONSTRAINT [CK_ProductVendor_OnOrderQty] CHECK ((onorderqty >= 0)),
  CONSTRAINT [CK_ProductVendor_StandardPrice] CHECK ((standardprice > 0.00)),
  PRIMARY KEY ([productid], [businessentityid])
)
GO

CREATE TABLE [purchasing].[shipmethod] (
  [shipmethodid] integer PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [shipbase] numeric NOT NULL DEFAULT (0),
  [shiprate] numeric NOT NULL DEFAULT (0),
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_ShipMethod_ShipBase] CHECK ((shipbase > 0.00)),
  CONSTRAINT [CK_ShipMethod_ShipRate] CHECK ((shiprate > 0.00))
)
GO

CREATE TABLE [purchasing].[vendor] (
  [businessentityid] integer PRIMARY KEY NOT NULL,
  [accountnumber] public.AccountNumber NOT NULL,
  [name] public.Name NOT NULL,
  [creditrating] smallint NOT NULL,
  [preferredvendorstatus] public.Flag NOT NULL DEFAULT (true),
  [activeflag] public.Flag NOT NULL DEFAULT (true),
  [purchasingwebserviceurl] character varying(1024),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_Vendor_CreditRating] CHECK (((creditrating >= 1) AND (creditrating <= 5)))
)
GO

CREATE TABLE [sales].[customer] (
  [customerid] integer PRIMARY KEY NOT NULL,
  [personid] integer,
  [storeid] integer,
  [territoryid] integer,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [sales].[creditcard] (
  [creditcardid] integer PRIMARY KEY NOT NULL,
  [cardtype] character varying(50) NOT NULL,
  [cardnumber] character varying(25) NOT NULL,
  [expmonth] smallint NOT NULL,
  [expyear] smallint NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [sales].[currencyrate] (
  [currencyrateid] integer PRIMARY KEY NOT NULL,
  [currencyratedate] timestamp NOT NULL,
  [fromcurrencycode] character(3) NOT NULL,
  [tocurrencycode] character(3) NOT NULL,
  [averagerate] numeric NOT NULL,
  [endofdayrate] numeric NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [sales].[countryregioncurrency] (
  [countryregioncode] character varying(3) NOT NULL,
  [currencycode] character(3) NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ([countryregioncode], [currencycode])
)
GO

CREATE TABLE [sales].[currency] (
  [currencycode] character(3) PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [sales].[personcreditcard] (
  [businessentityid] integer NOT NULL,
  [creditcardid] integer NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ([businessentityid], [creditcardid])
)
GO

CREATE TABLE [sales].[store] (
  [businessentityid] integer PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [salespersonid] integer,
  [demographics] xml,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [sales].[shoppingcartitem] (
  [shoppingcartitemid] integer PRIMARY KEY NOT NULL,
  [shoppingcartid] character varying(50) NOT NULL,
  [quantity] integer NOT NULL DEFAULT (1),
  [productid] integer NOT NULL,
  [datecreated] timestamp NOT NULL DEFAULT (now()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_ShoppingCartItem_Quantity] CHECK ((quantity >= 1))
)
GO

CREATE TABLE [sales].[specialoffer] (
  [specialofferid] integer PRIMARY KEY NOT NULL,
  [description] character varying(255) NOT NULL,
  [discountpct] numeric NOT NULL DEFAULT (0),
  [type] character varying(50) NOT NULL,
  [category] character varying(50) NOT NULL,
  [startdate] timestamp NOT NULL,
  [enddate] timestamp NOT NULL,
  [minqty] integer NOT NULL DEFAULT (0),
  [maxqty] integer,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_SpecialOffer_DiscountPct] CHECK ((discountpct >= 0.00)),
  CONSTRAINT [CK_SpecialOffer_EndDate] CHECK ((enddate >= startdate)),
  CONSTRAINT [CK_SpecialOffer_MaxQty] CHECK ((maxqty >= 0)),
  CONSTRAINT [CK_SpecialOffer_MinQty] CHECK ((minqty >= 0))
)
GO

CREATE TABLE [sales].[salesorderdetail] (
  [salesorderid] integer NOT NULL,
  [salesorderdetailid] integer NOT NULL,
  [carriertrackingnumber] character varying(25),
  [orderqty] smallint NOT NULL,
  [productid] integer NOT NULL,
  [specialofferid] integer NOT NULL,
  [unitprice] numeric NOT NULL,
  [unitpricediscount] numeric NOT NULL DEFAULT (0),
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_SalesOrderDetail_OrderQty] CHECK ((orderqty > 0)),
  CONSTRAINT [CK_SalesOrderDetail_UnitPrice] CHECK ((unitprice >= 0.00)),
  CONSTRAINT [CK_SalesOrderDetail_UnitPriceDiscount] CHECK ((unitpricediscount >= 0.00)),
  PRIMARY KEY ([salesorderid], [salesorderdetailid])
)
GO

CREATE TABLE [sales].[salesorderheader] (
  [salesorderid] integer PRIMARY KEY NOT NULL,
  [revisionnumber] smallint NOT NULL DEFAULT (0),
  [orderdate] timestamp NOT NULL DEFAULT (now()),
  [duedate] timestamp NOT NULL,
  [shipdate] timestamp,
  [status] smallint NOT NULL DEFAULT (1),
  [onlineorderflag] public.Flag NOT NULL DEFAULT (true),
  [purchaseordernumber] public.OrderNumber,
  [accountnumber] public.AccountNumber,
  [customerid] integer NOT NULL,
  [salespersonid] integer,
  [territoryid] integer,
  [billtoaddressid] integer NOT NULL,
  [shiptoaddressid] integer NOT NULL,
  [shipmethodid] integer NOT NULL,
  [creditcardid] integer,
  [creditcardapprovalcode] character varying(15),
  [currencyrateid] integer,
  [subtotal] numeric NOT NULL DEFAULT (0),
  [taxamt] numeric NOT NULL DEFAULT (0),
  [freight] numeric NOT NULL DEFAULT (0),
  [totaldue] numeric,
  [comment] character varying(128),
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_SalesOrderHeader_DueDate] CHECK ((duedate >= orderdate)),
  CONSTRAINT [CK_SalesOrderHeader_Freight] CHECK ((freight >= 0.00)),
  CONSTRAINT [CK_SalesOrderHeader_ShipDate] CHECK (((shipdate >= orderdate) OR (shipdate IS NULL))),
  CONSTRAINT [CK_SalesOrderHeader_Status] CHECK (((status >= 0) AND (status <= 8))),
  CONSTRAINT [CK_SalesOrderHeader_SubTotal] CHECK ((subtotal >= 0.00)),
  CONSTRAINT [CK_SalesOrderHeader_TaxAmt] CHECK ((taxamt >= 0.00))
)
GO

CREATE TABLE [sales].[salesorderheadersalesreason] (
  [salesorderid] integer NOT NULL,
  [salesreasonid] integer NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ([salesorderid], [salesreasonid])
)
GO

CREATE TABLE [sales].[specialofferproduct] (
  [specialofferid] integer NOT NULL,
  [productid] integer NOT NULL,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ([specialofferid], [productid])
)
GO

CREATE TABLE [sales].[salesperson] (
  [businessentityid] integer PRIMARY KEY NOT NULL,
  [territoryid] integer,
  [salesquota] numeric,
  [bonus] numeric NOT NULL DEFAULT (0),
  [commissionpct] numeric NOT NULL DEFAULT (0),
  [salesytd] numeric NOT NULL DEFAULT (0),
  [saleslastyear] numeric NOT NULL DEFAULT (0),
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_SalesPerson_Bonus] CHECK ((bonus >= 0.00)),
  CONSTRAINT [CK_SalesPerson_CommissionPct] CHECK ((commissionpct >= 0.00)),
  CONSTRAINT [CK_SalesPerson_SalesLastYear] CHECK ((saleslastyear >= 0.00)),
  CONSTRAINT [CK_SalesPerson_SalesQuota] CHECK ((salesquota > 0.00)),
  CONSTRAINT [CK_SalesPerson_SalesYTD] CHECK ((salesytd >= 0.00))
)
GO

CREATE TABLE [sales].[salespersonquotahistory] (
  [businessentityid] integer NOT NULL,
  [quotadate] timestamp NOT NULL,
  [salesquota] numeric NOT NULL,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_SalesPersonQuotaHistory_SalesQuota] CHECK ((salesquota > 0.00)),
  PRIMARY KEY ([businessentityid], [quotadate])
)
GO

CREATE TABLE [sales].[salesreason] (
  [salesreasonid] integer PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [reasontype] public.Name NOT NULL,
  [modifieddate] timestamp NOT NULL DEFAULT (now())
)
GO

CREATE TABLE [sales].[salesterritory] (
  [territoryid] integer PRIMARY KEY NOT NULL,
  [name] public.Name NOT NULL,
  [countryregioncode] character varying(3) NOT NULL,
  [group] character varying(50) NOT NULL,
  [salesytd] numeric NOT NULL DEFAULT (0),
  [saleslastyear] numeric NOT NULL DEFAULT (0),
  [costytd] numeric NOT NULL DEFAULT (0),
  [costlastyear] numeric NOT NULL DEFAULT (0),
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_SalesTerritory_CostLastYear] CHECK ((costlastyear >= 0.00)),
  CONSTRAINT [CK_SalesTerritory_CostYTD] CHECK ((costytd >= 0.00)),
  CONSTRAINT [CK_SalesTerritory_SalesLastYear] CHECK ((saleslastyear >= 0.00)),
  CONSTRAINT [CK_SalesTerritory_SalesYTD] CHECK ((salesytd >= 0.00))
)
GO

CREATE TABLE [sales].[salesterritoryhistory] (
  [businessentityid] integer NOT NULL,
  [territoryid] integer NOT NULL,
  [startdate] timestamp NOT NULL,
  [enddate] timestamp,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_SalesTerritoryHistory_EndDate] CHECK (((enddate >= startdate) OR (enddate IS NULL))),
  PRIMARY KEY ([businessentityid], [startdate], [territoryid])
)
GO

CREATE TABLE [sales].[salestaxrate] (
  [salestaxrateid] integer PRIMARY KEY NOT NULL,
  [stateprovinceid] integer NOT NULL,
  [taxtype] smallint NOT NULL,
  [taxrate] numeric NOT NULL DEFAULT (0),
  [name] public.Name NOT NULL,
  [rowguid] uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  [modifieddate] timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT [CK_SalesTaxRate_TaxType] CHECK (((taxtype >= 1) AND (taxtype <= 3)))
)
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Lookup table containing the departments within the Adventure Works Cycles company.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'department';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for Department records.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'department',
@level2type = N'Column', @level2name = 'departmentid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Name of the department.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'department',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Name of the group to which the department belongs.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'department',
@level2type = N'Column', @level2name = 'groupname';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Employee information such as salary, department, and title.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employee';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for Employee records.  Foreign key to BusinessEntity.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employee',
@level2type = N'Column', @level2name = 'businessentityid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Unique national identification number such as a social security number.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employee',
@level2type = N'Column', @level2name = 'nationalidnumber';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Network login.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employee',
@level2type = N'Column', @level2name = 'loginid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Work title such as Buyer or Sales Representative.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employee',
@level2type = N'Column', @level2name = 'jobtitle';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date of birth.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employee',
@level2type = N'Column', @level2name = 'birthdate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'M = Married, S = Single',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employee',
@level2type = N'Column', @level2name = 'maritalstatus';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'M = Male, F = Female',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employee',
@level2type = N'Column', @level2name = 'gender';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Employee hired on this date.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employee',
@level2type = N'Column', @level2name = 'hiredate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Job classification. 0 = Hourly, not exempt from collective bargaining. 1 = Salaried, exempt from collective bargaining.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employee',
@level2type = N'Column', @level2name = 'salariedflag';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Number of available vacation hours.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employee',
@level2type = N'Column', @level2name = 'vacationhours';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Number of available sick leave hours.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employee',
@level2type = N'Column', @level2name = 'sickleavehours';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '0 = Inactive, 1 = Active',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employee',
@level2type = N'Column', @level2name = 'currentflag';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Where the employee is located in corporate hierarchy.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employee',
@level2type = N'Column', @level2name = 'organizationnode';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Employee department transfers.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employeedepartmenthistory';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Employee identification number. Foreign key to Employee.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employeedepartmenthistory',
@level2type = N'Column', @level2name = 'businessentityid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Department in which the employee worked including currently. Foreign key to Department.DepartmentID.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employeedepartmenthistory',
@level2type = N'Column', @level2name = 'departmentid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Identifies which 8-hour shift the employee works. Foreign key to Shift.Shift.ID.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employeedepartmenthistory',
@level2type = N'Column', @level2name = 'shiftid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date the employee started work in the department.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employeedepartmenthistory',
@level2type = N'Column', @level2name = 'startdate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date the employee left the department. NULL = Current department.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employeedepartmenthistory',
@level2type = N'Column', @level2name = 'enddate';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Employee pay history.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employeepayhistory';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Employee identification number. Foreign key to Employee.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employeepayhistory',
@level2type = N'Column', @level2name = 'businessentityid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date the change in pay is effective',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employeepayhistory',
@level2type = N'Column', @level2name = 'ratechangedate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Salary hourly rate.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employeepayhistory',
@level2type = N'Column', @level2name = 'rate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '1 = Salary received monthly, 2 = Salary received biweekly',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'employeepayhistory',
@level2type = N'Column', @level2name = 'payfrequency';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'RÃ©sumÃ©s submitted to Human Resources by job applicants.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'jobcandidate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for JobCandidate records.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'jobcandidate',
@level2type = N'Column', @level2name = 'jobcandidateid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Employee identification number if applicant was hired. Foreign key to Employee.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'jobcandidate',
@level2type = N'Column', @level2name = 'businessentityid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'RÃ©sumÃ© in XML format.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'jobcandidate',
@level2type = N'Column', @level2name = 'resume';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Work shift lookup table.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'shift';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for Shift records.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'shift',
@level2type = N'Column', @level2name = 'shiftid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Shift description.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'shift',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Shift start time.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'shift',
@level2type = N'Column', @level2name = 'starttime';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Shift end time.',
@level0type = N'Schema', @level0name = 'humanresources',
@level1type = N'Table',  @level1name = 'shift',
@level2type = N'Column', @level2name = 'endtime';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Street address information for customers, employees, and vendors.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'address';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for Address records.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'address',
@level2type = N'Column', @level2name = 'addressid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'First street address line.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'address',
@level2type = N'Column', @level2name = 'addressline1';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Second street address line.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'address',
@level2type = N'Column', @level2name = 'addressline2';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Name of the city.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'address',
@level2type = N'Column', @level2name = 'city';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Unique identification number for the state or province. Foreign key to StateProvince table.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'address',
@level2type = N'Column', @level2name = 'stateprovinceid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Postal code for the street address.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'address',
@level2type = N'Column', @level2name = 'postalcode';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Latitude and longitude of this address.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'address',
@level2type = N'Column', @level2name = 'spatiallocation';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Cross-reference table mapping customers, vendors, and employees to their addresses.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'businessentityaddress';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to BusinessEntity.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'businessentityaddress',
@level2type = N'Column', @level2name = 'businessentityid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to Address.AddressID.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'businessentityaddress',
@level2type = N'Column', @level2name = 'addressid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to AddressType.AddressTypeID.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'businessentityaddress',
@level2type = N'Column', @level2name = 'addresstypeid';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Lookup table containing the ISO standard codes for countries and regions.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'countryregion';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'ISO standard code for countries and regions.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'countryregion',
@level2type = N'Column', @level2name = 'countryregioncode';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Country or region name.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'countryregion',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Where to send a person email.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'emailaddress';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Person associated with this email address.  Foreign key to Person.BusinessEntityID',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'emailaddress',
@level2type = N'Column', @level2name = 'businessentityid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. ID of this email address.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'emailaddress',
@level2type = N'Column', @level2name = 'emailaddressid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'E-mail address for the person.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'emailaddress',
@level2type = N'Column', @level2name = 'emailaddress';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Human beings involved with AdventureWorks: employees, customer contacts, and vendor contacts.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'person';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for Person records.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'person',
@level2type = N'Column', @level2name = 'businessentityid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary type of person: SC = Store Contact, IN = Individual (retail) customer, SP = Sales person, EM = Employee (non-sales), VC = Vendor contact, GC = General contact',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'person',
@level2type = N'Column', @level2name = 'persontype';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '0 = The data in FirstName and LastName are stored in western style (first name, last name) order.  1 = Eastern style (last name, first name) order.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'person',
@level2type = N'Column', @level2name = 'namestyle';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'A courtesy title. For example, Mr. or Ms.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'person',
@level2type = N'Column', @level2name = 'title';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'First name of the person.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'person',
@level2type = N'Column', @level2name = 'firstname';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Middle name or middle initial of the person.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'person',
@level2type = N'Column', @level2name = 'middlename';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Last name of the person.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'person',
@level2type = N'Column', @level2name = 'lastname';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Surname suffix. For example, Sr. or Jr.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'person',
@level2type = N'Column', @level2name = 'suffix';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '0 = Contact does not wish to receive e-mail promotions, 1 = Contact does wish to receive e-mail promotions from AdventureWorks, 2 = Contact does wish to receive e-mail promotions from AdventureWorks and selected partners.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'person',
@level2type = N'Column', @level2name = 'emailpromotion';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Additional contact information about the person stored in xml format.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'person',
@level2type = N'Column', @level2name = 'additionalcontactinfo';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Personal information such as hobbies, and income collected from online shoppers. Used for sales analysis.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'person',
@level2type = N'Column', @level2name = 'demographics';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Telephone number and type of a person.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'personphone';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Business entity identification number. Foreign key to Person.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'personphone',
@level2type = N'Column', @level2name = 'businessentityid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Telephone number identification number.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'personphone',
@level2type = N'Column', @level2name = 'phonenumber';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Kind of phone number. Foreign key to PhoneNumberType.PhoneNumberTypeID.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'personphone',
@level2type = N'Column', @level2name = 'phonenumbertypeid';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Type of phone number of a person.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'phonenumbertype';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for telephone number type records.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'phonenumbertype',
@level2type = N'Column', @level2name = 'phonenumbertypeid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Name of the telephone number type',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'phonenumbertype',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'State and province lookup table.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'stateprovince';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for StateProvince records.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'stateprovince',
@level2type = N'Column', @level2name = 'stateprovinceid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'ISO standard state or province code.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'stateprovince',
@level2type = N'Column', @level2name = 'stateprovincecode';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'ISO standard country or region code. Foreign key to CountryRegion.CountryRegionCode.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'stateprovince',
@level2type = N'Column', @level2name = 'countryregioncode';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '0 = StateProvinceCode exists. 1 = StateProvinceCode unavailable, using CountryRegionCode.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'stateprovince',
@level2type = N'Column', @level2name = 'isonlystateprovinceflag';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'State or province description.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'stateprovince',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'ID of the territory in which the state or province is located. Foreign key to SalesTerritory.SalesTerritoryID.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'stateprovince',
@level2type = N'Column', @level2name = 'territoryid';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Types of addresses stored in the Address table.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'addresstype';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for AddressType records.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'addresstype',
@level2type = N'Column', @level2name = 'addresstypeid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Address type description. For example, Billing, Home, or Shipping.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'addresstype',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Source of the ID that connects vendors, customers, and employees with address and contact information.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'businessentity';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for all customers, vendors, and employees.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'businessentity',
@level2type = N'Column', @level2name = 'businessentityid';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Cross-reference table mapping stores, vendors, and employees to people',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'businessentitycontact';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to BusinessEntity.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'businessentitycontact',
@level2type = N'Column', @level2name = 'businessentityid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to Person.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'businessentitycontact',
@level2type = N'Column', @level2name = 'personid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key.  Foreign key to ContactType.ContactTypeID.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'businessentitycontact',
@level2type = N'Column', @level2name = 'contacttypeid';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Lookup table containing the types of business entity contacts.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'contacttype';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for ContactType records.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'contacttype',
@level2type = N'Column', @level2name = 'contacttypeid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Contact type description.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'contacttype',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'One way hashed authentication information',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'password';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Password for the e-mail account.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'password',
@level2type = N'Column', @level2name = 'passwordhash';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Random value concatenated with the password string before the password is hashed.',
@level0type = N'Schema', @level0name = 'person',
@level1type = N'Table',  @level1name = 'password',
@level2type = N'Column', @level2name = 'passwordsalt';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Items required to make bicycles and bicycle subassemblies. It identifies the heirarchical relationship between a parent product and its components.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'billofmaterials';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for BillOfMaterials records.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'billofmaterials',
@level2type = N'Column', @level2name = 'billofmaterialsid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Parent product identification number. Foreign key to Product.ProductID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'billofmaterials',
@level2type = N'Column', @level2name = 'productassemblyid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Component identification number. Foreign key to Product.ProductID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'billofmaterials',
@level2type = N'Column', @level2name = 'componentid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date the component started being used in the assembly item.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'billofmaterials',
@level2type = N'Column', @level2name = 'startdate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date the component stopped being used in the assembly item.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'billofmaterials',
@level2type = N'Column', @level2name = 'enddate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Standard code identifying the unit of measure for the quantity.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'billofmaterials',
@level2type = N'Column', @level2name = 'unitmeasurecode';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Indicates the depth the component is from its parent (AssemblyID).',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'billofmaterials',
@level2type = N'Column', @level2name = 'bomlevel';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Quantity of the component needed to create the assembly.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'billofmaterials',
@level2type = N'Column', @level2name = 'perassemblyqty';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Lookup table containing the languages in which some AdventureWorks data is stored.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'culture';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for Culture records.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'culture',
@level2type = N'Column', @level2name = 'cultureid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Culture description.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'culture',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Product maintenance documents.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'document';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Title of the document.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'document',
@level2type = N'Column', @level2name = 'title';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Employee who controls the document.  Foreign key to Employee.BusinessEntityID',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'document',
@level2type = N'Column', @level2name = 'owner';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '0 = This is a folder, 1 = This is a document.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'document',
@level2type = N'Column', @level2name = 'folderflag';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'File name of the document',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'document',
@level2type = N'Column', @level2name = 'filename';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'File extension indicating the document type. For example, .doc or .txt.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'document',
@level2type = N'Column', @level2name = 'fileextension';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Revision number of the document.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'document',
@level2type = N'Column', @level2name = 'revision';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Engineering change approval number.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'document',
@level2type = N'Column', @level2name = 'changenumber';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '1 = Pending approval, 2 = Approved, 3 = Obsolete',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'document',
@level2type = N'Column', @level2name = 'status';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Document abstract.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'document',
@level2type = N'Column', @level2name = 'documentsummary';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Complete document.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'document',
@level2type = N'Column', @level2name = 'document';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'ROWGUIDCOL number uniquely identifying the record. Required for FileStream.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'document',
@level2type = N'Column', @level2name = 'rowguid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for Document records.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'document',
@level2type = N'Column', @level2name = 'documentnode';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Bicycle assembly diagrams.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'illustration';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for Illustration records.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'illustration',
@level2type = N'Column', @level2name = 'illustrationid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Illustrations used in manufacturing instructions. Stored as XML.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'illustration',
@level2type = N'Column', @level2name = 'diagram';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Product inventory and manufacturing locations.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'location';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for Location records.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'location',
@level2type = N'Column', @level2name = 'locationid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Location description.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'location',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Standard hourly cost of the manufacturing location.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'location',
@level2type = N'Column', @level2name = 'costrate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Work capacity (in hours) of the manufacturing location.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'location',
@level2type = N'Column', @level2name = 'availability';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Products sold or used in the manfacturing of sold products.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for Product records.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'productid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Name of the product.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Unique product identification number.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'productnumber';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '0 = Product is purchased, 1 = Product is manufactured in-house.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'makeflag';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '0 = Product is not a salable item. 1 = Product is salable.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'finishedgoodsflag';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product color.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'color';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Minimum inventory quantity.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'safetystocklevel';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Inventory level that triggers a purchase order or work order.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'reorderpoint';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Standard cost of the product.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'standardcost';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Selling price.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'listprice';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product size.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'size';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Unit of measure for Size column.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'sizeunitmeasurecode';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Unit of measure for Weight column.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'weightunitmeasurecode';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product weight.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'weight';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Number of days required to manufacture the product.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'daystomanufacture';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'R = Road, M = Mountain, T = Touring, S = Standard',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'productline';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'H = High, M = Medium, L = Low',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'class';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'W = Womens, M = Mens, U = Universal',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'style';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product is a member of this product subcategory. Foreign key to ProductSubCategory.ProductSubCategoryID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'productsubcategoryid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product is a member of this product model. Foreign key to ProductModel.ProductModelID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'productmodelid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date the product was available for sale.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'sellstartdate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date the product was no longer available for sale.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'sellenddate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date the product was discontinued.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'product',
@level2type = N'Column', @level2name = 'discontinueddate';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'High-level product categorization.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productcategory';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for ProductCategory records.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productcategory',
@level2type = N'Column', @level2name = 'productcategoryid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Category description.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productcategory',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Changes in the cost of a product over time.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productcosthistory';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product identification number. Foreign key to Product.ProductID',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productcosthistory',
@level2type = N'Column', @level2name = 'productid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product cost start date.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productcosthistory',
@level2type = N'Column', @level2name = 'startdate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product cost end date.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productcosthistory',
@level2type = N'Column', @level2name = 'enddate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Standard cost of the product.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productcosthistory',
@level2type = N'Column', @level2name = 'standardcost';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Product descriptions in several languages.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productdescription';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for ProductDescription records.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productdescription',
@level2type = N'Column', @level2name = 'productdescriptionid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Description of the product.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productdescription',
@level2type = N'Column', @level2name = 'description';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Cross-reference table mapping products to related product documents.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productdocument';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product identification number. Foreign key to Product.ProductID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productdocument',
@level2type = N'Column', @level2name = 'productid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Document identification number. Foreign key to Document.DocumentNode.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productdocument',
@level2type = N'Column', @level2name = 'documentnode';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Product inventory information.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productinventory';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product identification number. Foreign key to Product.ProductID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productinventory',
@level2type = N'Column', @level2name = 'productid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Inventory location identification number. Foreign key to Location.LocationID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productinventory',
@level2type = N'Column', @level2name = 'locationid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Storage compartment within an inventory location.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productinventory',
@level2type = N'Column', @level2name = 'shelf';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Storage container on a shelf in an inventory location.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productinventory',
@level2type = N'Column', @level2name = 'bin';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Quantity of products in the inventory location.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productinventory',
@level2type = N'Column', @level2name = 'quantity';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Changes in the list price of a product over time.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productlistpricehistory';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product identification number. Foreign key to Product.ProductID',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productlistpricehistory',
@level2type = N'Column', @level2name = 'productid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'List price start date.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productlistpricehistory',
@level2type = N'Column', @level2name = 'startdate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'List price end date',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productlistpricehistory',
@level2type = N'Column', @level2name = 'enddate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product list price.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productlistpricehistory',
@level2type = N'Column', @level2name = 'listprice';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Product model classification.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productmodel';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for ProductModel records.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productmodel',
@level2type = N'Column', @level2name = 'productmodelid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product model description.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productmodel',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Detailed product catalog information in xml format.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productmodel',
@level2type = N'Column', @level2name = 'catalogdescription';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Manufacturing instructions in xml format.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productmodel',
@level2type = N'Column', @level2name = 'instructions';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Cross-reference table mapping product models and illustrations.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productmodelillustration';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to ProductModel.ProductModelID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productmodelillustration',
@level2type = N'Column', @level2name = 'productmodelid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to Illustration.IllustrationID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productmodelillustration',
@level2type = N'Column', @level2name = 'illustrationid';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Cross-reference table mapping product descriptions and the language the description is written in.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productmodelproductdescriptionculture';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to ProductModel.ProductModelID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productmodelproductdescriptionculture',
@level2type = N'Column', @level2name = 'productmodelid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to ProductDescription.ProductDescriptionID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productmodelproductdescriptionculture',
@level2type = N'Column', @level2name = 'productdescriptionid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Culture identification number. Foreign key to Culture.CultureID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productmodelproductdescriptionculture',
@level2type = N'Column', @level2name = 'cultureid';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Product images.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productphoto';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for ProductPhoto records.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productphoto',
@level2type = N'Column', @level2name = 'productphotoid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Small image of the product.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productphoto',
@level2type = N'Column', @level2name = 'thumbnailphoto';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Small image file name.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productphoto',
@level2type = N'Column', @level2name = 'thumbnailphotofilename';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Large image of the product.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productphoto',
@level2type = N'Column', @level2name = 'largephoto';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Large image file name.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productphoto',
@level2type = N'Column', @level2name = 'largephotofilename';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Cross-reference table mapping products and product photos.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productproductphoto';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product identification number. Foreign key to Product.ProductID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productproductphoto',
@level2type = N'Column', @level2name = 'productid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product photo identification number. Foreign key to ProductPhoto.ProductPhotoID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productproductphoto',
@level2type = N'Column', @level2name = 'productphotoid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '0 = Photo is not the principal image. 1 = Photo is the principal image.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productproductphoto',
@level2type = N'Column', @level2name = 'primary';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Customer reviews of products they have purchased.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productreview';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for ProductReview records.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productreview',
@level2type = N'Column', @level2name = 'productreviewid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product identification number. Foreign key to Product.ProductID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productreview',
@level2type = N'Column', @level2name = 'productid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Name of the reviewer.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productreview',
@level2type = N'Column', @level2name = 'reviewername';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date review was submitted.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productreview',
@level2type = N'Column', @level2name = 'reviewdate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Reviewer''s e-mail address.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productreview',
@level2type = N'Column', @level2name = 'emailaddress';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product rating given by the reviewer. Scale is 1 to 5 with 5 as the highest rating.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productreview',
@level2type = N'Column', @level2name = 'rating';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Reviewer''s comments
WARNING: can have several lines!',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productreview',
@level2type = N'Column', @level2name = 'comments';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Product subcategories. See ProductCategory table.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productsubcategory';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for ProductSubcategory records.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productsubcategory',
@level2type = N'Column', @level2name = 'productsubcategoryid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product category identification number. Foreign key to ProductCategory.ProductCategoryID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productsubcategory',
@level2type = N'Column', @level2name = 'productcategoryid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Subcategory description.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'productsubcategory',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Manufacturing failure reasons lookup table.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'scrapreason';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for ScrapReason records.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'scrapreason',
@level2type = N'Column', @level2name = 'scrapreasonid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Failure description.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'scrapreason',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Record of each purchase order, sales order, or work order transaction year to date.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistory';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for TransactionHistory records.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistory',
@level2type = N'Column', @level2name = 'transactionid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product identification number. Foreign key to Product.ProductID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistory',
@level2type = N'Column', @level2name = 'productid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Purchase order, sales order, or work order identification number.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistory',
@level2type = N'Column', @level2name = 'referenceorderid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Line number associated with the purchase order, sales order, or work order.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistory',
@level2type = N'Column', @level2name = 'referenceorderlineid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date and time of the transaction.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistory',
@level2type = N'Column', @level2name = 'transactiondate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'W = WorkOrder, S = SalesOrder, P = PurchaseOrder',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistory',
@level2type = N'Column', @level2name = 'transactiontype';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product quantity.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistory',
@level2type = N'Column', @level2name = 'quantity';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product cost.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistory',
@level2type = N'Column', @level2name = 'actualcost';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Transactions for previous years.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistoryarchive';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for TransactionHistoryArchive records.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistoryarchive',
@level2type = N'Column', @level2name = 'transactionid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product identification number. Foreign key to Product.ProductID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistoryarchive',
@level2type = N'Column', @level2name = 'productid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Purchase order, sales order, or work order identification number.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistoryarchive',
@level2type = N'Column', @level2name = 'referenceorderid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Line number associated with the purchase order, sales order, or work order.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistoryarchive',
@level2type = N'Column', @level2name = 'referenceorderlineid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date and time of the transaction.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistoryarchive',
@level2type = N'Column', @level2name = 'transactiondate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'W = Work Order, S = Sales Order, P = Purchase Order',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistoryarchive',
@level2type = N'Column', @level2name = 'transactiontype';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product quantity.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistoryarchive',
@level2type = N'Column', @level2name = 'quantity';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product cost.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'transactionhistoryarchive',
@level2type = N'Column', @level2name = 'actualcost';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Unit of measure lookup table.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'unitmeasure';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'unitmeasure',
@level2type = N'Column', @level2name = 'unitmeasurecode';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Unit of measure description.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'unitmeasure',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Manufacturing work orders.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorder';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for WorkOrder records.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorder',
@level2type = N'Column', @level2name = 'workorderid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product identification number. Foreign key to Product.ProductID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorder',
@level2type = N'Column', @level2name = 'productid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product quantity to build.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorder',
@level2type = N'Column', @level2name = 'orderqty';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Quantity that failed inspection.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorder',
@level2type = N'Column', @level2name = 'scrappedqty';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Work order start date.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorder',
@level2type = N'Column', @level2name = 'startdate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Work order end date.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorder',
@level2type = N'Column', @level2name = 'enddate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Work order due date.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorder',
@level2type = N'Column', @level2name = 'duedate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Reason for inspection failure.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorder',
@level2type = N'Column', @level2name = 'scrapreasonid';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Work order details.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorderrouting';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to WorkOrder.WorkOrderID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorderrouting',
@level2type = N'Column', @level2name = 'workorderid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to Product.ProductID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorderrouting',
@level2type = N'Column', @level2name = 'productid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Indicates the manufacturing process sequence.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorderrouting',
@level2type = N'Column', @level2name = 'operationsequence';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Manufacturing location where the part is processed. Foreign key to Location.LocationID.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorderrouting',
@level2type = N'Column', @level2name = 'locationid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Planned manufacturing start date.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorderrouting',
@level2type = N'Column', @level2name = 'scheduledstartdate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Planned manufacturing end date.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorderrouting',
@level2type = N'Column', @level2name = 'scheduledenddate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Actual start date.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorderrouting',
@level2type = N'Column', @level2name = 'actualstartdate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Actual end date.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorderrouting',
@level2type = N'Column', @level2name = 'actualenddate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Number of manufacturing hours used.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorderrouting',
@level2type = N'Column', @level2name = 'actualresourcehrs';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Estimated manufacturing cost.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorderrouting',
@level2type = N'Column', @level2name = 'plannedcost';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Actual manufacturing cost.',
@level0type = N'Schema', @level0name = 'production',
@level1type = N'Table',  @level1name = 'workorderrouting',
@level2type = N'Column', @level2name = 'actualcost';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Individual products associated with a specific purchase order. See PurchaseOrderHeader.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderdetail';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to PurchaseOrderHeader.PurchaseOrderID.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderdetail',
@level2type = N'Column', @level2name = 'purchaseorderid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. One line number per purchased product.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderdetail',
@level2type = N'Column', @level2name = 'purchaseorderdetailid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date the product is expected to be received.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderdetail',
@level2type = N'Column', @level2name = 'duedate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Quantity ordered.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderdetail',
@level2type = N'Column', @level2name = 'orderqty';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product identification number. Foreign key to Product.ProductID.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderdetail',
@level2type = N'Column', @level2name = 'productid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Vendor''s selling price of a single product.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderdetail',
@level2type = N'Column', @level2name = 'unitprice';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Quantity actually received from the vendor.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderdetail',
@level2type = N'Column', @level2name = 'receivedqty';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Quantity rejected during inspection.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderdetail',
@level2type = N'Column', @level2name = 'rejectedqty';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'General purchase order information. See PurchaseOrderDetail
for more details.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderheader';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderheader',
@level2type = N'Column', @level2name = 'purchaseorderid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Incremental number to track changes to the purchase order over time.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderheader',
@level2type = N'Column', @level2name = 'revisionnumber';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Order current status. 1 = Pending; 2 = Approved; 3 = Rejected; 4 = Complete',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderheader',
@level2type = N'Column', @level2name = 'status';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Employee who created the purchase order. Foreign key to Employee.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderheader',
@level2type = N'Column', @level2name = 'employeeid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Vendor with whom the purchase order is placed. Foreign key to Vendor.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderheader',
@level2type = N'Column', @level2name = 'vendorid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Shipping method. Foreign key to ShipMethod.ShipMethodID.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderheader',
@level2type = N'Column', @level2name = 'shipmethodid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Purchase order creation date.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderheader',
@level2type = N'Column', @level2name = 'orderdate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Estimated shipment date from the vendor.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderheader',
@level2type = N'Column', @level2name = 'shipdate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Purchase order subtotal. Computed as SUM(PurchaseOrderDetail.LineTotal)for the appropriate PurchaseOrderID.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderheader',
@level2type = N'Column', @level2name = 'subtotal';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Tax amount.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderheader',
@level2type = N'Column', @level2name = 'taxamt';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Shipping cost.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'purchaseorderheader',
@level2type = N'Column', @level2name = 'freight';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Cross-reference table mapping vendors with the products they supply.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'productvendor';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to Product.ProductID.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'productvendor',
@level2type = N'Column', @level2name = 'productid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to Vendor.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'productvendor',
@level2type = N'Column', @level2name = 'businessentityid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'The average span of time (in days) between placing an order with the vendor and receiving the purchased product.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'productvendor',
@level2type = N'Column', @level2name = 'averageleadtime';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'The vendor''s usual selling price.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'productvendor',
@level2type = N'Column', @level2name = 'standardprice';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'The selling price when last purchased.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'productvendor',
@level2type = N'Column', @level2name = 'lastreceiptcost';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date the product was last received by the vendor.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'productvendor',
@level2type = N'Column', @level2name = 'lastreceiptdate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'The maximum quantity that should be ordered.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'productvendor',
@level2type = N'Column', @level2name = 'minorderqty';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'The minimum quantity that should be ordered.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'productvendor',
@level2type = N'Column', @level2name = 'maxorderqty';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'The quantity currently on order.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'productvendor',
@level2type = N'Column', @level2name = 'onorderqty';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'The product''s unit of measure.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'productvendor',
@level2type = N'Column', @level2name = 'unitmeasurecode';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Shipping company lookup table.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'shipmethod';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for ShipMethod records.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'shipmethod',
@level2type = N'Column', @level2name = 'shipmethodid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Shipping company name.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'shipmethod',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Minimum shipping charge.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'shipmethod',
@level2type = N'Column', @level2name = 'shipbase';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Shipping charge per pound.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'shipmethod',
@level2type = N'Column', @level2name = 'shiprate';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Companies from whom Adventure Works Cycles purchases parts or other goods.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'vendor';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for Vendor records.  Foreign key to BusinessEntity.BusinessEntityID',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'vendor',
@level2type = N'Column', @level2name = 'businessentityid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Vendor account (identification) number.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'vendor',
@level2type = N'Column', @level2name = 'accountnumber';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Company name.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'vendor',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '1 = Superior, 2 = Excellent, 3 = Above average, 4 = Average, 5 = Below average',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'vendor',
@level2type = N'Column', @level2name = 'creditrating';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '0 = Do not use if another vendor is available. 1 = Preferred over other vendors supplying the same product.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'vendor',
@level2type = N'Column', @level2name = 'preferredvendorstatus';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '0 = Vendor no longer used. 1 = Vendor is actively used.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'vendor',
@level2type = N'Column', @level2name = 'activeflag';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Vendor URL.',
@level0type = N'Schema', @level0name = 'purchasing',
@level1type = N'Table',  @level1name = 'vendor',
@level2type = N'Column', @level2name = 'purchasingwebserviceurl';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Current customer information. Also see the Person and Store tables.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'customer';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'customer',
@level2type = N'Column', @level2name = 'customerid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Foreign key to Person.BusinessEntityID',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'customer',
@level2type = N'Column', @level2name = 'personid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Foreign key to Store.BusinessEntityID',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'customer',
@level2type = N'Column', @level2name = 'storeid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'ID of the territory in which the customer is located. Foreign key to SalesTerritory.SalesTerritoryID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'customer',
@level2type = N'Column', @level2name = 'territoryid';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Customer credit card information.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'creditcard';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for CreditCard records.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'creditcard',
@level2type = N'Column', @level2name = 'creditcardid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Credit card name.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'creditcard',
@level2type = N'Column', @level2name = 'cardtype';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Credit card number.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'creditcard',
@level2type = N'Column', @level2name = 'cardnumber';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Credit card expiration month.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'creditcard',
@level2type = N'Column', @level2name = 'expmonth';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Credit card expiration year.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'creditcard',
@level2type = N'Column', @level2name = 'expyear';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Currency exchange rates.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'currencyrate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for CurrencyRate records.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'currencyrate',
@level2type = N'Column', @level2name = 'currencyrateid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date and time the exchange rate was obtained.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'currencyrate',
@level2type = N'Column', @level2name = 'currencyratedate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Exchange rate was converted from this currency code.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'currencyrate',
@level2type = N'Column', @level2name = 'fromcurrencycode';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Exchange rate was converted to this currency code.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'currencyrate',
@level2type = N'Column', @level2name = 'tocurrencycode';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Average exchange rate for the day.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'currencyrate',
@level2type = N'Column', @level2name = 'averagerate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Final exchange rate for the day.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'currencyrate',
@level2type = N'Column', @level2name = 'endofdayrate';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Cross-reference table mapping ISO currency codes to a country or region.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'countryregioncurrency';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'ISO code for countries and regions. Foreign key to CountryRegion.CountryRegionCode.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'countryregioncurrency',
@level2type = N'Column', @level2name = 'countryregioncode';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'ISO standard currency code. Foreign key to Currency.CurrencyCode.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'countryregioncurrency',
@level2type = N'Column', @level2name = 'currencycode';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Lookup table containing standard ISO currencies.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'currency';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'The ISO code for the Currency.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'currency',
@level2type = N'Column', @level2name = 'currencycode';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Currency name.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'currency',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Cross-reference table mapping people to their credit card information in the CreditCard table.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'personcreditcard';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Business entity identification number. Foreign key to Person.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'personcreditcard',
@level2type = N'Column', @level2name = 'businessentityid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Credit card identification number. Foreign key to CreditCard.CreditCardID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'personcreditcard',
@level2type = N'Column', @level2name = 'creditcardid';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Customers (resellers) of Adventure Works products.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'store';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to Customer.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'store',
@level2type = N'Column', @level2name = 'businessentityid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Name of the store.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'store',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'ID of the sales person assigned to the customer. Foreign key to SalesPerson.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'store',
@level2type = N'Column', @level2name = 'salespersonid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Demographic informationg about the store such as the number of employees, annual sales and store type.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'store',
@level2type = N'Column', @level2name = 'demographics';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Contains online customer orders until the order is submitted or cancelled.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'shoppingcartitem';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for ShoppingCartItem records.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'shoppingcartitem',
@level2type = N'Column', @level2name = 'shoppingcartitemid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Shopping cart identification number.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'shoppingcartitem',
@level2type = N'Column', @level2name = 'shoppingcartid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product quantity ordered.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'shoppingcartitem',
@level2type = N'Column', @level2name = 'quantity';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product ordered. Foreign key to Product.ProductID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'shoppingcartitem',
@level2type = N'Column', @level2name = 'productid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date the time the record was created.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'shoppingcartitem',
@level2type = N'Column', @level2name = 'datecreated';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Sale discounts lookup table.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'specialoffer';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for SpecialOffer records.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'specialoffer',
@level2type = N'Column', @level2name = 'specialofferid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Discount description.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'specialoffer',
@level2type = N'Column', @level2name = 'description';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Discount precentage.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'specialoffer',
@level2type = N'Column', @level2name = 'discountpct';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Discount type category.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'specialoffer',
@level2type = N'Column', @level2name = 'type';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Group the discount applies to such as Reseller or Customer.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'specialoffer',
@level2type = N'Column', @level2name = 'category';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Discount start date.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'specialoffer',
@level2type = N'Column', @level2name = 'startdate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Discount end date.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'specialoffer',
@level2type = N'Column', @level2name = 'enddate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Minimum discount percent allowed.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'specialoffer',
@level2type = N'Column', @level2name = 'minqty';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Maximum discount percent allowed.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'specialoffer',
@level2type = N'Column', @level2name = 'maxqty';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Individual products associated with a specific sales order. See SalesOrderHeader.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderdetail';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to SalesOrderHeader.SalesOrderID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderdetail',
@level2type = N'Column', @level2name = 'salesorderid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. One incremental unique number per product sold.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderdetail',
@level2type = N'Column', @level2name = 'salesorderdetailid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Shipment tracking number supplied by the shipper.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderdetail',
@level2type = N'Column', @level2name = 'carriertrackingnumber';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Quantity ordered per product.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderdetail',
@level2type = N'Column', @level2name = 'orderqty';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product sold to customer. Foreign key to Product.ProductID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderdetail',
@level2type = N'Column', @level2name = 'productid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Promotional code. Foreign key to SpecialOffer.SpecialOfferID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderdetail',
@level2type = N'Column', @level2name = 'specialofferid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Selling price of a single product.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderdetail',
@level2type = N'Column', @level2name = 'unitprice';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Discount amount.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderdetail',
@level2type = N'Column', @level2name = 'unitpricediscount';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'General sales order information.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'salesorderid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Incremental number to track changes to the sales order over time.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'revisionnumber';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Dates the sales order was created.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'orderdate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date the order is due to the customer.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'duedate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date the order was shipped to the customer.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'shipdate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Order current status. 1 = In process; 2 = Approved; 3 = Backordered; 4 = Rejected; 5 = Shipped; 6 = Cancelled',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'status';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '0 = Order placed by sales person. 1 = Order placed online by customer.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'onlineorderflag';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Customer purchase order number reference.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'purchaseordernumber';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Financial accounting number reference.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'accountnumber';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Customer identification number. Foreign key to Customer.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'customerid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Sales person who created the sales order. Foreign key to SalesPerson.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'salespersonid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Territory in which the sale was made. Foreign key to SalesTerritory.SalesTerritoryID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'territoryid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Customer billing address. Foreign key to Address.AddressID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'billtoaddressid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Customer shipping address. Foreign key to Address.AddressID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'shiptoaddressid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Shipping method. Foreign key to ShipMethod.ShipMethodID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'shipmethodid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Credit card identification number. Foreign key to CreditCard.CreditCardID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'creditcardid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Approval code provided by the credit card company.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'creditcardapprovalcode';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Currency exchange rate used. Foreign key to CurrencyRate.CurrencyRateID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'currencyrateid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Sales subtotal. Computed as SUM(SalesOrderDetail.LineTotal)for the appropriate SalesOrderID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'subtotal';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Tax amount.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'taxamt';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Shipping cost.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'freight';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Total due from customer. Computed as Subtotal + TaxAmt + Freight.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'totaldue';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Sales representative comments.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheader',
@level2type = N'Column', @level2name = 'comment';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Cross-reference table mapping sales orders to sales reason codes.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheadersalesreason';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to SalesOrderHeader.SalesOrderID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheadersalesreason',
@level2type = N'Column', @level2name = 'salesorderid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Foreign key to SalesReason.SalesReasonID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesorderheadersalesreason',
@level2type = N'Column', @level2name = 'salesreasonid';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Cross-reference table mapping products to special offer discounts.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'specialofferproduct';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for SpecialOfferProduct records.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'specialofferproduct',
@level2type = N'Column', @level2name = 'specialofferid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Product identification number. Foreign key to Product.ProductID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'specialofferproduct',
@level2type = N'Column', @level2name = 'productid';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Sales representative current information.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesperson';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for SalesPerson records. Foreign key to Employee.BusinessEntityID',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesperson',
@level2type = N'Column', @level2name = 'businessentityid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Territory currently assigned to. Foreign key to SalesTerritory.SalesTerritoryID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesperson',
@level2type = N'Column', @level2name = 'territoryid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Projected yearly sales.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesperson',
@level2type = N'Column', @level2name = 'salesquota';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Bonus due if quota is met.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesperson',
@level2type = N'Column', @level2name = 'bonus';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Commision percent received per sale.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesperson',
@level2type = N'Column', @level2name = 'commissionpct';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Sales total year to date.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesperson',
@level2type = N'Column', @level2name = 'salesytd';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Sales total of previous year.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesperson',
@level2type = N'Column', @level2name = 'saleslastyear';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Sales performance tracking.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salespersonquotahistory';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Sales person identification number. Foreign key to SalesPerson.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salespersonquotahistory',
@level2type = N'Column', @level2name = 'businessentityid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Sales quota date.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salespersonquotahistory',
@level2type = N'Column', @level2name = 'quotadate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Sales quota amount.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salespersonquotahistory',
@level2type = N'Column', @level2name = 'salesquota';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Lookup table of customer purchase reasons.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesreason';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for SalesReason records.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesreason',
@level2type = N'Column', @level2name = 'salesreasonid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Sales reason description.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesreason',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Category the sales reason belongs to.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesreason',
@level2type = N'Column', @level2name = 'reasontype';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Sales territory lookup table.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesterritory';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for SalesTerritory records.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesterritory',
@level2type = N'Column', @level2name = 'territoryid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Sales territory description',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesterritory',
@level2type = N'Column', @level2name = 'name';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'ISO standard country or region code. Foreign key to CountryRegion.CountryRegionCode.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesterritory',
@level2type = N'Column', @level2name = 'countryregioncode';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Geographic area to which the sales territory belong.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesterritory',
@level2type = N'Column', @level2name = 'group';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Sales in the territory year to date.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesterritory',
@level2type = N'Column', @level2name = 'salesytd';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Sales in the territory the previous year.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesterritory',
@level2type = N'Column', @level2name = 'saleslastyear';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Business costs in the territory year to date.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesterritory',
@level2type = N'Column', @level2name = 'costytd';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Business costs in the territory the previous year.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesterritory',
@level2type = N'Column', @level2name = 'costlastyear';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Sales representative transfers to other sales territories.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesterritoryhistory';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. The sales rep.  Foreign key to SalesPerson.BusinessEntityID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesterritoryhistory',
@level2type = N'Column', @level2name = 'businessentityid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Territory identification number. Foreign key to SalesTerritory.SalesTerritoryID.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesterritoryhistory',
@level2type = N'Column', @level2name = 'territoryid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key. Date the sales representive started work in the territory.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesterritoryhistory',
@level2type = N'Column', @level2name = 'startdate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Date the sales representative left work in the territory.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salesterritoryhistory',
@level2type = N'Column', @level2name = 'enddate';
GO

EXEC sp_addextendedproperty
@name = N'Table_Description',
@value = 'Tax rate lookup table.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salestaxrate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Primary key for SalesTaxRate records.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salestaxrate',
@level2type = N'Column', @level2name = 'salestaxrateid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'State, province, or country/region the sales tax applies to.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salestaxrate',
@level2type = N'Column', @level2name = 'stateprovinceid';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = '1 = Tax applied to retail transactions, 2 = Tax applied to wholesale transactions, 3 = Tax applied to all sales (retail and wholesale) transactions.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salestaxrate',
@level2type = N'Column', @level2name = 'taxtype';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Tax rate amount.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salestaxrate',
@level2type = N'Column', @level2name = 'taxrate';
GO

EXEC sp_addextendedproperty
@name = N'Column_Description',
@value = 'Tax rate description.',
@level0type = N'Schema', @level0name = 'sales',
@level1type = N'Table',  @level1name = 'salestaxrate',
@level2type = N'Column', @level2name = 'name';
GO

ALTER TABLE [humanresources].[employeedepartmenthistory] ADD CONSTRAINT [FK_EmployeeDepartmentHistory_Department_DepartmentID] FOREIGN KEY ([departmentid]) REFERENCES [humanresources].[department] ([departmentid])
GO

ALTER TABLE [humanresources].[employeedepartmenthistory] ADD CONSTRAINT [FK_EmployeeDepartmentHistory_Employee_BusinessEntityID] FOREIGN KEY ([businessentityid]) REFERENCES [humanresources].[employee] ([businessentityid])
GO

ALTER TABLE [humanresources].[employeedepartmenthistory] ADD CONSTRAINT [FK_EmployeeDepartmentHistory_Shift_ShiftID] FOREIGN KEY ([shiftid]) REFERENCES [humanresources].[shift] ([shiftid])
GO

ALTER TABLE [humanresources].[employeepayhistory] ADD CONSTRAINT [FK_EmployeePayHistory_Employee_BusinessEntityID] FOREIGN KEY ([businessentityid]) REFERENCES [humanresources].[employee] ([businessentityid])
GO

ALTER TABLE [humanresources].[employee] ADD CONSTRAINT [FK_Employee_Person_BusinessEntityID] FOREIGN KEY ([businessentityid]) REFERENCES [person].[person] ([businessentityid])
GO

ALTER TABLE [humanresources].[jobcandidate] ADD CONSTRAINT [FK_JobCandidate_Employee_BusinessEntityID] FOREIGN KEY ([businessentityid]) REFERENCES [humanresources].[employee] ([businessentityid])
GO

ALTER TABLE [person].[address] ADD CONSTRAINT [FK_Address_StateProvince_StateProvinceID] FOREIGN KEY ([stateprovinceid]) REFERENCES [person].[stateprovince] ([stateprovinceid])
GO

ALTER TABLE [person].[businessentityaddress] ADD CONSTRAINT [FK_BusinessEntityAddress_AddressType_AddressTypeID] FOREIGN KEY ([addresstypeid]) REFERENCES [person].[addresstype] ([addresstypeid])
GO

ALTER TABLE [person].[businessentityaddress] ADD CONSTRAINT [FK_BusinessEntityAddress_Address_AddressID] FOREIGN KEY ([addressid]) REFERENCES [person].[address] ([addressid])
GO

ALTER TABLE [person].[businessentityaddress] ADD CONSTRAINT [FK_BusinessEntityAddress_BusinessEntity_BusinessEntityID] FOREIGN KEY ([businessentityid]) REFERENCES [person].[businessentity] ([businessentityid])
GO

ALTER TABLE [person].[businessentitycontact] ADD CONSTRAINT [FK_BusinessEntityContact_BusinessEntity_BusinessEntityID] FOREIGN KEY ([businessentityid]) REFERENCES [person].[businessentity] ([businessentityid])
GO

ALTER TABLE [person].[businessentitycontact] ADD CONSTRAINT [FK_BusinessEntityContact_ContactType_ContactTypeID] FOREIGN KEY ([contacttypeid]) REFERENCES [person].[contacttype] ([contacttypeid])
GO

ALTER TABLE [person].[businessentitycontact] ADD CONSTRAINT [FK_BusinessEntityContact_Person_PersonID] FOREIGN KEY ([personid]) REFERENCES [person].[person] ([businessentityid])
GO

ALTER TABLE [person].[emailaddress] ADD CONSTRAINT [FK_EmailAddress_Person_BusinessEntityID] FOREIGN KEY ([businessentityid]) REFERENCES [person].[person] ([businessentityid])
GO

ALTER TABLE [person].[password] ADD CONSTRAINT [FK_Password_Person_BusinessEntityID] FOREIGN KEY ([businessentityid]) REFERENCES [person].[person] ([businessentityid])
GO

ALTER TABLE [person].[personphone] ADD CONSTRAINT [FK_PersonPhone_Person_BusinessEntityID] FOREIGN KEY ([businessentityid]) REFERENCES [person].[person] ([businessentityid])
GO

ALTER TABLE [person].[personphone] ADD CONSTRAINT [FK_PersonPhone_PhoneNumberType_PhoneNumberTypeID] FOREIGN KEY ([phonenumbertypeid]) REFERENCES [person].[phonenumbertype] ([phonenumbertypeid])
GO

ALTER TABLE [person].[person] ADD CONSTRAINT [FK_Person_BusinessEntity_BusinessEntityID] FOREIGN KEY ([businessentityid]) REFERENCES [person].[businessentity] ([businessentityid])
GO

ALTER TABLE [person].[stateprovince] ADD CONSTRAINT [FK_StateProvince_CountryRegion_CountryRegionCode] FOREIGN KEY ([countryregioncode]) REFERENCES [person].[countryregion] ([countryregioncode])
GO

ALTER TABLE [person].[stateprovince] ADD CONSTRAINT [FK_StateProvince_SalesTerritory_TerritoryID] FOREIGN KEY ([territoryid]) REFERENCES [sales].[salesterritory] ([territoryid])
GO

ALTER TABLE [production].[billofmaterials] ADD CONSTRAINT [FK_BillOfMaterials_Product_ComponentID] FOREIGN KEY ([componentid]) REFERENCES [production].[product] ([productid])
GO

ALTER TABLE [production].[billofmaterials] ADD CONSTRAINT [FK_BillOfMaterials_Product_ProductAssemblyID] FOREIGN KEY ([productassemblyid]) REFERENCES [production].[product] ([productid])
GO

ALTER TABLE [production].[billofmaterials] ADD CONSTRAINT [FK_BillOfMaterials_UnitMeasure_UnitMeasureCode] FOREIGN KEY ([unitmeasurecode]) REFERENCES [production].[unitmeasure] ([unitmeasurecode])
GO

ALTER TABLE [production].[document] ADD CONSTRAINT [FK_Document_Employee_Owner] FOREIGN KEY ([owner]) REFERENCES [humanresources].[employee] ([businessentityid])
GO

ALTER TABLE [production].[productcosthistory] ADD CONSTRAINT [FK_ProductCostHistory_Product_ProductID] FOREIGN KEY ([productid]) REFERENCES [production].[product] ([productid])
GO

ALTER TABLE [production].[productdocument] ADD CONSTRAINT [FK_ProductDocument_Document_DocumentNode] FOREIGN KEY ([documentnode]) REFERENCES [production].[document] ([documentnode])
GO

ALTER TABLE [production].[productdocument] ADD CONSTRAINT [FK_ProductDocument_Product_ProductID] FOREIGN KEY ([productid]) REFERENCES [production].[product] ([productid])
GO

ALTER TABLE [production].[productinventory] ADD CONSTRAINT [FK_ProductInventory_Location_LocationID] FOREIGN KEY ([locationid]) REFERENCES [production].[location] ([locationid])
GO

ALTER TABLE [production].[productinventory] ADD CONSTRAINT [FK_ProductInventory_Product_ProductID] FOREIGN KEY ([productid]) REFERENCES [production].[product] ([productid])
GO

ALTER TABLE [production].[productlistpricehistory] ADD CONSTRAINT [FK_ProductListPriceHistory_Product_ProductID] FOREIGN KEY ([productid]) REFERENCES [production].[product] ([productid])
GO

ALTER TABLE [production].[productmodelillustration] ADD CONSTRAINT [FK_ProductModelIllustration_Illustration_IllustrationID] FOREIGN KEY ([illustrationid]) REFERENCES [production].[illustration] ([illustrationid])
GO

ALTER TABLE [production].[productmodelillustration] ADD CONSTRAINT [FK_ProductModelIllustration_ProductModel_ProductModelID] FOREIGN KEY ([productmodelid]) REFERENCES [production].[productmodel] ([productmodelid])
GO

ALTER TABLE [production].[productmodelproductdescriptionculture] ADD CONSTRAINT [FK_ProductModelProductDescriptionCulture_Culture_CultureID] FOREIGN KEY ([cultureid]) REFERENCES [production].[culture] ([cultureid])
GO

ALTER TABLE [production].[productmodelproductdescriptionculture] ADD CONSTRAINT [FK_ProductModelProductDescriptionCulture_ProductDescription_Pro] FOREIGN KEY ([productdescriptionid]) REFERENCES [production].[productdescription] ([productdescriptionid])
GO

ALTER TABLE [production].[productmodelproductdescriptionculture] ADD CONSTRAINT [FK_ProductModelProductDescriptionCulture_ProductModel_ProductMo] FOREIGN KEY ([productmodelid]) REFERENCES [production].[productmodel] ([productmodelid])
GO

ALTER TABLE [production].[productproductphoto] ADD CONSTRAINT [FK_ProductProductPhoto_ProductPhoto_ProductPhotoID] FOREIGN KEY ([productphotoid]) REFERENCES [production].[productphoto] ([productphotoid])
GO

ALTER TABLE [production].[productproductphoto] ADD CONSTRAINT [FK_ProductProductPhoto_Product_ProductID] FOREIGN KEY ([productid]) REFERENCES [production].[product] ([productid])
GO

ALTER TABLE [production].[productreview] ADD CONSTRAINT [FK_ProductReview_Product_ProductID] FOREIGN KEY ([productid]) REFERENCES [production].[product] ([productid])
GO

ALTER TABLE [production].[productsubcategory] ADD CONSTRAINT [FK_ProductSubcategory_ProductCategory_ProductCategoryID] FOREIGN KEY ([productcategoryid]) REFERENCES [production].[productcategory] ([productcategoryid])
GO

ALTER TABLE [production].[product] ADD CONSTRAINT [FK_Product_ProductModel_ProductModelID] FOREIGN KEY ([productmodelid]) REFERENCES [production].[productmodel] ([productmodelid])
GO

ALTER TABLE [production].[product] ADD CONSTRAINT [FK_Product_ProductSubcategory_ProductSubcategoryID] FOREIGN KEY ([productsubcategoryid]) REFERENCES [production].[productsubcategory] ([productsubcategoryid])
GO

ALTER TABLE [production].[product] ADD CONSTRAINT [FK_Product_UnitMeasure_SizeUnitMeasureCode] FOREIGN KEY ([sizeunitmeasurecode]) REFERENCES [production].[unitmeasure] ([unitmeasurecode])
GO

ALTER TABLE [production].[product] ADD CONSTRAINT [FK_Product_UnitMeasure_WeightUnitMeasureCode] FOREIGN KEY ([weightunitmeasurecode]) REFERENCES [production].[unitmeasure] ([unitmeasurecode])
GO

ALTER TABLE [production].[transactionhistory] ADD CONSTRAINT [FK_TransactionHistory_Product_ProductID] FOREIGN KEY ([productid]) REFERENCES [production].[product] ([productid])
GO

ALTER TABLE [production].[workorderrouting] ADD CONSTRAINT [FK_WorkOrderRouting_Location_LocationID] FOREIGN KEY ([locationid]) REFERENCES [production].[location] ([locationid])
GO

ALTER TABLE [production].[workorderrouting] ADD CONSTRAINT [FK_WorkOrderRouting_WorkOrder_WorkOrderID] FOREIGN KEY ([workorderid]) REFERENCES [production].[workorder] ([workorderid])
GO

ALTER TABLE [production].[workorder] ADD CONSTRAINT [FK_WorkOrder_Product_ProductID] FOREIGN KEY ([productid]) REFERENCES [production].[product] ([productid])
GO

ALTER TABLE [production].[workorder] ADD CONSTRAINT [FK_WorkOrder_ScrapReason_ScrapReasonID] FOREIGN KEY ([scrapreasonid]) REFERENCES [production].[scrapreason] ([scrapreasonid])
GO

ALTER TABLE [purchasing].[productvendor] ADD CONSTRAINT [FK_ProductVendor_Product_ProductID] FOREIGN KEY ([productid]) REFERENCES [production].[product] ([productid])
GO

ALTER TABLE [purchasing].[productvendor] ADD CONSTRAINT [FK_ProductVendor_UnitMeasure_UnitMeasureCode] FOREIGN KEY ([unitmeasurecode]) REFERENCES [production].[unitmeasure] ([unitmeasurecode])
GO

ALTER TABLE [purchasing].[productvendor] ADD CONSTRAINT [FK_ProductVendor_Vendor_BusinessEntityID] FOREIGN KEY ([businessentityid]) REFERENCES [purchasing].[vendor] ([businessentityid])
GO

ALTER TABLE [purchasing].[purchaseorderdetail] ADD CONSTRAINT [FK_PurchaseOrderDetail_Product_ProductID] FOREIGN KEY ([productid]) REFERENCES [production].[product] ([productid])
GO

ALTER TABLE [purchasing].[purchaseorderdetail] ADD CONSTRAINT [FK_PurchaseOrderDetail_PurchaseOrderHeader_PurchaseOrderID] FOREIGN KEY ([purchaseorderid]) REFERENCES [purchasing].[purchaseorderheader] ([purchaseorderid])
GO

ALTER TABLE [purchasing].[purchaseorderheader] ADD CONSTRAINT [FK_PurchaseOrderHeader_Employee_EmployeeID] FOREIGN KEY ([employeeid]) REFERENCES [humanresources].[employee] ([businessentityid])
GO

ALTER TABLE [purchasing].[purchaseorderheader] ADD CONSTRAINT [FK_PurchaseOrderHeader_ShipMethod_ShipMethodID] FOREIGN KEY ([shipmethodid]) REFERENCES [purchasing].[shipmethod] ([shipmethodid])
GO

ALTER TABLE [purchasing].[purchaseorderheader] ADD CONSTRAINT [FK_PurchaseOrderHeader_Vendor_VendorID] FOREIGN KEY ([vendorid]) REFERENCES [purchasing].[vendor] ([businessentityid])
GO

ALTER TABLE [purchasing].[vendor] ADD CONSTRAINT [FK_Vendor_BusinessEntity_BusinessEntityID] FOREIGN KEY ([businessentityid]) REFERENCES [person].[businessentity] ([businessentityid])
GO

ALTER TABLE [sales].[countryregioncurrency] ADD CONSTRAINT [FK_CountryRegionCurrency_CountryRegion_CountryRegionCode] FOREIGN KEY ([countryregioncode]) REFERENCES [person].[countryregion] ([countryregioncode])
GO

ALTER TABLE [sales].[countryregioncurrency] ADD CONSTRAINT [FK_CountryRegionCurrency_Currency_CurrencyCode] FOREIGN KEY ([currencycode]) REFERENCES [sales].[currency] ([currencycode])
GO

ALTER TABLE [sales].[currencyrate] ADD CONSTRAINT [FK_CurrencyRate_Currency_FromCurrencyCode] FOREIGN KEY ([fromcurrencycode]) REFERENCES [sales].[currency] ([currencycode])
GO

ALTER TABLE [sales].[currencyrate] ADD CONSTRAINT [FK_CurrencyRate_Currency_ToCurrencyCode] FOREIGN KEY ([tocurrencycode]) REFERENCES [sales].[currency] ([currencycode])
GO

ALTER TABLE [sales].[customer] ADD CONSTRAINT [FK_Customer_Person_PersonID] FOREIGN KEY ([personid]) REFERENCES [person].[person] ([businessentityid])
GO

ALTER TABLE [sales].[customer] ADD CONSTRAINT [FK_Customer_SalesTerritory_TerritoryID] FOREIGN KEY ([territoryid]) REFERENCES [sales].[salesterritory] ([territoryid])
GO

ALTER TABLE [sales].[customer] ADD CONSTRAINT [FK_Customer_Store_StoreID] FOREIGN KEY ([storeid]) REFERENCES [sales].[store] ([businessentityid])
GO

ALTER TABLE [sales].[personcreditcard] ADD CONSTRAINT [FK_PersonCreditCard_CreditCard_CreditCardID] FOREIGN KEY ([creditcardid]) REFERENCES [sales].[creditcard] ([creditcardid])
GO

ALTER TABLE [sales].[personcreditcard] ADD CONSTRAINT [FK_PersonCreditCard_Person_BusinessEntityID] FOREIGN KEY ([businessentityid]) REFERENCES [person].[person] ([businessentityid])
GO

ALTER TABLE [sales].[salesorderdetail] ADD CONSTRAINT [FK_SalesOrderDetail_SalesOrderHeader_SalesOrderID] FOREIGN KEY ([salesorderid]) REFERENCES [sales].[salesorderheader] ([salesorderid]) ON DELETE CASCADE
GO

ALTER TABLE [sales].[salesorderdetail] ADD CONSTRAINT [FK_SalesOrderDetail_SpecialOfferProduct_SpecialOfferIDProductID] FOREIGN KEY ([specialofferid], [productid]) REFERENCES [sales].[specialofferproduct] ([specialofferid], [productid])
GO

ALTER TABLE [sales].[salesorderheadersalesreason] ADD CONSTRAINT [FK_SalesOrderHeaderSalesReason_SalesOrderHeader_SalesOrderID] FOREIGN KEY ([salesorderid]) REFERENCES [sales].[salesorderheader] ([salesorderid]) ON DELETE CASCADE
GO

ALTER TABLE [sales].[salesorderheadersalesreason] ADD CONSTRAINT [FK_SalesOrderHeaderSalesReason_SalesReason_SalesReasonID] FOREIGN KEY ([salesreasonid]) REFERENCES [sales].[salesreason] ([salesreasonid])
GO

ALTER TABLE [sales].[salesorderheader] ADD CONSTRAINT [FK_SalesOrderHeader_Address_BillToAddressID] FOREIGN KEY ([billtoaddressid]) REFERENCES [person].[address] ([addressid])
GO

ALTER TABLE [sales].[salesorderheader] ADD CONSTRAINT [FK_SalesOrderHeader_Address_ShipToAddressID] FOREIGN KEY ([shiptoaddressid]) REFERENCES [person].[address] ([addressid])
GO

ALTER TABLE [sales].[salesorderheader] ADD CONSTRAINT [FK_SalesOrderHeader_CreditCard_CreditCardID] FOREIGN KEY ([creditcardid]) REFERENCES [sales].[creditcard] ([creditcardid])
GO

ALTER TABLE [sales].[salesorderheader] ADD CONSTRAINT [FK_SalesOrderHeader_CurrencyRate_CurrencyRateID] FOREIGN KEY ([currencyrateid]) REFERENCES [sales].[currencyrate] ([currencyrateid])
GO

ALTER TABLE [sales].[salesorderheader] ADD CONSTRAINT [FK_SalesOrderHeader_Customer_CustomerID] FOREIGN KEY ([customerid]) REFERENCES [sales].[customer] ([customerid])
GO

ALTER TABLE [sales].[salesorderheader] ADD CONSTRAINT [FK_SalesOrderHeader_SalesPerson_SalesPersonID] FOREIGN KEY ([salespersonid]) REFERENCES [sales].[salesperson] ([businessentityid])
GO

ALTER TABLE [sales].[salesorderheader] ADD CONSTRAINT [FK_SalesOrderHeader_SalesTerritory_TerritoryID] FOREIGN KEY ([territoryid]) REFERENCES [sales].[salesterritory] ([territoryid])
GO

ALTER TABLE [sales].[salesorderheader] ADD CONSTRAINT [FK_SalesOrderHeader_ShipMethod_ShipMethodID] FOREIGN KEY ([shipmethodid]) REFERENCES [purchasing].[shipmethod] ([shipmethodid])
GO

ALTER TABLE [sales].[salespersonquotahistory] ADD CONSTRAINT [FK_SalesPersonQuotaHistory_SalesPerson_BusinessEntityID] FOREIGN KEY ([businessentityid]) REFERENCES [sales].[salesperson] ([businessentityid])
GO

ALTER TABLE [sales].[salesperson] ADD CONSTRAINT [FK_SalesPerson_Employee_BusinessEntityID] FOREIGN KEY ([businessentityid]) REFERENCES [humanresources].[employee] ([businessentityid])
GO

ALTER TABLE [sales].[salesperson] ADD CONSTRAINT [FK_SalesPerson_SalesTerritory_TerritoryID] FOREIGN KEY ([territoryid]) REFERENCES [sales].[salesterritory] ([territoryid])
GO

ALTER TABLE [sales].[salestaxrate] ADD CONSTRAINT [FK_SalesTaxRate_StateProvince_StateProvinceID] FOREIGN KEY ([stateprovinceid]) REFERENCES [person].[stateprovince] ([stateprovinceid])
GO

ALTER TABLE [sales].[salesterritoryhistory] ADD CONSTRAINT [FK_SalesTerritoryHistory_SalesPerson_BusinessEntityID] FOREIGN KEY ([businessentityid]) REFERENCES [sales].[salesperson] ([businessentityid])
GO

ALTER TABLE [sales].[salesterritoryhistory] ADD CONSTRAINT [FK_SalesTerritoryHistory_SalesTerritory_TerritoryID] FOREIGN KEY ([territoryid]) REFERENCES [sales].[salesterritory] ([territoryid])
GO

ALTER TABLE [sales].[salesterritory] ADD CONSTRAINT [FK_SalesTerritory_CountryRegion_CountryRegionCode] FOREIGN KEY ([countryregioncode]) REFERENCES [person].[countryregion] ([countryregioncode])
GO

ALTER TABLE [sales].[shoppingcartitem] ADD CONSTRAINT [FK_ShoppingCartItem_Product_ProductID] FOREIGN KEY ([productid]) REFERENCES [production].[product] ([productid])
GO

ALTER TABLE [sales].[specialofferproduct] ADD CONSTRAINT [FK_SpecialOfferProduct_Product_ProductID] FOREIGN KEY ([productid]) REFERENCES [production].[product] ([productid])
GO

ALTER TABLE [sales].[specialofferproduct] ADD CONSTRAINT [FK_SpecialOfferProduct_SpecialOffer_SpecialOfferID] FOREIGN KEY ([specialofferid]) REFERENCES [sales].[specialoffer] ([specialofferid])
GO

ALTER TABLE [sales].[store] ADD CONSTRAINT [FK_Store_BusinessEntity_BusinessEntityID] FOREIGN KEY ([businessentityid]) REFERENCES [person].[businessentity] ([businessentityid])
GO

ALTER TABLE [sales].[store] ADD CONSTRAINT [FK_Store_SalesPerson_SalesPersonID] FOREIGN KEY ([salespersonid]) REFERENCES [sales].[salesperson] ([businessentityid])
GO
