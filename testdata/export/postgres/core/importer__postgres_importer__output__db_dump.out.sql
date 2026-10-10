CREATE SCHEMA "humanresources";

CREATE SCHEMA "person";

CREATE SCHEMA "production";

CREATE SCHEMA "purchasing";

CREATE SCHEMA "sales";

CREATE TABLE "humanresources"."department" (
  "departmentid" integer PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "groupname" "public.Name" NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "humanresources"."employee" (
  "businessentityid" integer PRIMARY KEY NOT NULL,
  "nationalidnumber" "character varying(15)" NOT NULL,
  "loginid" "character varying(256)" NOT NULL,
  "jobtitle" "character varying(50)" NOT NULL,
  "birthdate" date NOT NULL,
  "maritalstatus" character(1) NOT NULL,
  "gender" character(1) NOT NULL,
  "hiredate" date NOT NULL,
  "salariedflag" "public.Flag" NOT NULL DEFAULT true,
  "vacationhours" smallint NOT NULL DEFAULT 0,
  "sickleavehours" smallint NOT NULL DEFAULT 0,
  "currentflag" "public.Flag" NOT NULL DEFAULT true,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  "organizationnode" character varying DEFAULT ('/'::charactervarying),
  CONSTRAINT "CK_Employee_BirthDate" CHECK (((birthdate >= '1930-01-01'::date) AND (birthdate <= (now() - '18 years'::interval)))),
  CONSTRAINT "CK_Employee_Gender" CHECK ((upper((gender)::text) = ANY (ARRAY['M'::text, 'F'::text]))),
  CONSTRAINT "CK_Employee_HireDate" CHECK (((hiredate >= '1996-07-01'::date) AND (hiredate <= (now() + '1 day'::interval)))),
  CONSTRAINT "CK_Employee_MaritalStatus" CHECK ((upper((maritalstatus)::text) = ANY (ARRAY['M'::text, 'S'::text]))),
  CONSTRAINT "CK_Employee_SickLeaveHours" CHECK (((sickleavehours >= 0) AND (sickleavehours <= 120))),
  CONSTRAINT "CK_Employee_VacationHours" CHECK (((vacationhours >= '-40'::integer) AND (vacationhours <= 240)))
);

CREATE TABLE "humanresources"."employeedepartmenthistory" (
  "businessentityid" integer NOT NULL,
  "departmentid" smallint NOT NULL,
  "shiftid" smallint NOT NULL,
  "startdate" date NOT NULL,
  "enddate" date,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_EmployeeDepartmentHistory_EndDate" CHECK (((enddate >= startdate) OR (enddate IS NULL))),
  PRIMARY KEY ("businessentityid", "startdate", "departmentid", "shiftid")
);

CREATE TABLE "humanresources"."employeepayhistory" (
  "businessentityid" integer NOT NULL,
  "ratechangedate" timestamp NOT NULL,
  "rate" numeric NOT NULL,
  "payfrequency" smallint NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_EmployeePayHistory_PayFrequency" CHECK ((payfrequency = ANY (ARRAY[1, 2]))),
  CONSTRAINT "CK_EmployeePayHistory_Rate" CHECK (((rate >= 6.50) AND (rate <= 200.00))),
  PRIMARY KEY ("businessentityid", "ratechangedate")
);

CREATE TABLE "humanresources"."jobcandidate" (
  "jobcandidateid" integer PRIMARY KEY NOT NULL,
  "businessentityid" integer,
  "resume" xml,
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "humanresources"."shift" (
  "shiftid" integer PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "starttime" time NOT NULL,
  "endtime" time NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "person"."address" (
  "addressid" integer PRIMARY KEY NOT NULL,
  "addressline1" "character varying(60)" NOT NULL,
  "addressline2" "character varying(60)",
  "city" "character varying(30)" NOT NULL,
  "stateprovinceid" integer NOT NULL,
  "postalcode" "character varying(15)" NOT NULL,
  "spatiallocation" "character varying(44)",
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "person"."businessentityaddress" (
  "businessentityid" integer NOT NULL,
  "addressid" integer NOT NULL,
  "addresstypeid" integer NOT NULL,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ("businessentityid", "addressid", "addresstypeid")
);

CREATE TABLE "person"."countryregion" (
  "countryregioncode" "character varying(3)" PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "person"."emailaddress" (
  "businessentityid" integer NOT NULL,
  "emailaddressid" integer NOT NULL,
  "emailaddress" "character varying(50)",
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ("businessentityid", "emailaddressid")
);

CREATE TABLE "person"."person" (
  "businessentityid" integer PRIMARY KEY NOT NULL,
  "persontype" character(2) NOT NULL,
  "namestyle" "public.NameStyle" NOT NULL DEFAULT false,
  "title" "character varying(8)",
  "firstname" "public.Name" NOT NULL,
  "middlename" "public.Name",
  "lastname" "public.Name" NOT NULL,
  "suffix" "character varying(10)",
  "emailpromotion" integer NOT NULL DEFAULT 0,
  "additionalcontactinfo" xml,
  "demographics" xml,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_Person_EmailPromotion" CHECK (((emailpromotion >= 0) AND (emailpromotion <= 2))),
  CONSTRAINT "CK_Person_PersonType" CHECK (((persontype IS NULL) OR (upper((persontype)::text) = ANY (ARRAY['SC'::text, 'VC'::text, 'IN'::text, 'EM'::text, 'SP'::text, 'GC'::text]))))
);

CREATE TABLE "person"."personphone" (
  "businessentityid" integer NOT NULL,
  "phonenumber" "public.Phone" NOT NULL,
  "phonenumbertypeid" integer NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ("businessentityid", "phonenumber", "phonenumbertypeid")
);

CREATE TABLE "person"."phonenumbertype" (
  "phonenumbertypeid" integer PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "person"."stateprovince" (
  "stateprovinceid" integer PRIMARY KEY NOT NULL,
  "stateprovincecode" character(3) NOT NULL,
  "countryregioncode" "character varying(3)" NOT NULL,
  "isonlystateprovinceflag" "public.Flag" NOT NULL DEFAULT true,
  "name" "public.Name" NOT NULL,
  "territoryid" integer NOT NULL,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "person"."addresstype" (
  "addresstypeid" integer PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "person"."businessentity" (
  "businessentityid" integer PRIMARY KEY NOT NULL,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "person"."businessentitycontact" (
  "businessentityid" integer NOT NULL,
  "personid" integer NOT NULL,
  "contacttypeid" integer NOT NULL,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ("businessentityid", "personid", "contacttypeid")
);

CREATE TABLE "person"."contacttype" (
  "contacttypeid" integer PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "person"."password" (
  "businessentityid" integer PRIMARY KEY NOT NULL,
  "passwordhash" "character varying(128)" NOT NULL,
  "passwordsalt" "character varying(10)" NOT NULL,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "production"."billofmaterials" (
  "billofmaterialsid" integer PRIMARY KEY NOT NULL,
  "productassemblyid" integer,
  "componentid" integer NOT NULL,
  "startdate" timestamp NOT NULL DEFAULT (now()),
  "enddate" timestamp,
  "unitmeasurecode" character(3) NOT NULL,
  "bomlevel" smallint NOT NULL,
  "perassemblyqty" numeric(8,2) NOT NULL DEFAULT 1,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_BillOfMaterials_BOMLevel" CHECK ((((productassemblyid IS NULL) AND (bomlevel = 0) AND (perassemblyqty = 1.00)) OR ((productassemblyid IS NOT NULL) AND (bomlevel >= 1)))),
  CONSTRAINT "CK_BillOfMaterials_EndDate" CHECK (((enddate > startdate) OR (enddate IS NULL))),
  CONSTRAINT "CK_BillOfMaterials_PerAssemblyQty" CHECK ((perassemblyqty >= 1.00)),
  CONSTRAINT "CK_BillOfMaterials_ProductAssemblyID" CHECK ((productassemblyid <> componentid))
);

CREATE TABLE "production"."culture" (
  "cultureid" character(6) PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "production"."document" (
  "title" "character varying(50)" NOT NULL,
  "owner" integer NOT NULL,
  "folderflag" "public.Flag" NOT NULL DEFAULT false,
  "filename" "character varying(400)" NOT NULL,
  "fileextension" "character varying(8)",
  "revision" character(5) NOT NULL,
  "changenumber" integer NOT NULL DEFAULT 0,
  "status" smallint NOT NULL,
  "documentsummary" text,
  "document" bytea,
  "rowguid" uuid UNIQUE NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  "documentnode" character varying PRIMARY KEY NOT NULL DEFAULT ('/'::charactervarying),
  CONSTRAINT "CK_Document_Status" CHECK (((status >= 1) AND (status <= 3)))
);

CREATE TABLE "production"."illustration" (
  "illustrationid" integer PRIMARY KEY NOT NULL,
  "diagram" xml,
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "production"."location" (
  "locationid" integer PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "costrate" numeric NOT NULL DEFAULT 0,
  "availability" numeric(8,2) NOT NULL DEFAULT 0,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_Location_Availability" CHECK ((availability >= 0.00)),
  CONSTRAINT "CK_Location_CostRate" CHECK ((costrate >= 0.00))
);

CREATE TABLE "production"."product" (
  "productid" integer PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "productnumber" "character varying(25)" NOT NULL,
  "makeflag" "public.Flag" NOT NULL DEFAULT true,
  "finishedgoodsflag" "public.Flag" NOT NULL DEFAULT true,
  "color" "character varying(15)",
  "safetystocklevel" smallint NOT NULL,
  "reorderpoint" smallint NOT NULL,
  "standardcost" numeric NOT NULL,
  "listprice" numeric NOT NULL,
  "size" "character varying(5)",
  "sizeunitmeasurecode" character(3),
  "weightunitmeasurecode" character(3),
  "weight" numeric(8,2),
  "daystomanufacture" integer NOT NULL,
  "productline" character(2),
  "class" character(2),
  "style" character(2),
  "productsubcategoryid" integer,
  "productmodelid" integer,
  "sellstartdate" timestamp NOT NULL,
  "sellenddate" timestamp,
  "discontinueddate" timestamp,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_Product_Class" CHECK (((upper((class)::text) = ANY (ARRAY['L'::text, 'M'::text, 'H'::text])) OR (class IS NULL))),
  CONSTRAINT "CK_Product_DaysToManufacture" CHECK ((daystomanufacture >= 0)),
  CONSTRAINT "CK_Product_ListPrice" CHECK ((listprice >= 0.00)),
  CONSTRAINT "CK_Product_ProductLine" CHECK (((upper((productline)::text) = ANY (ARRAY['S'::text, 'T'::text, 'M'::text, 'R'::text])) OR (productline IS NULL))),
  CONSTRAINT "CK_Product_ReorderPoint" CHECK ((reorderpoint > 0)),
  CONSTRAINT "CK_Product_SafetyStockLevel" CHECK ((safetystocklevel > 0)),
  CONSTRAINT "CK_Product_SellEndDate" CHECK (((sellenddate >= sellstartdate) OR (sellenddate IS NULL))),
  CONSTRAINT "CK_Product_StandardCost" CHECK ((standardcost >= 0.00)),
  CONSTRAINT "CK_Product_Style" CHECK (((upper((style)::text) = ANY (ARRAY['W'::text, 'M'::text, 'U'::text])) OR (style IS NULL))),
  CONSTRAINT "CK_Product_Weight" CHECK ((weight > 0.00))
);

CREATE TABLE "production"."productcategory" (
  "productcategoryid" integer PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "production"."productcosthistory" (
  "productid" integer NOT NULL,
  "startdate" timestamp NOT NULL,
  "enddate" timestamp,
  "standardcost" numeric NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_ProductCostHistory_EndDate" CHECK (((enddate >= startdate) OR (enddate IS NULL))),
  CONSTRAINT "CK_ProductCostHistory_StandardCost" CHECK ((standardcost >= 0.00)),
  PRIMARY KEY ("productid", "startdate")
);

CREATE TABLE "production"."productdescription" (
  "productdescriptionid" integer PRIMARY KEY NOT NULL,
  "description" "character varying(400)" NOT NULL,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "production"."productdocument" (
  "productid" integer NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  "documentnode" character varying NOT NULL DEFAULT ('/'::charactervarying),
  PRIMARY KEY ("productid", "documentnode")
);

CREATE TABLE "production"."productinventory" (
  "productid" integer NOT NULL,
  "locationid" smallint NOT NULL,
  "shelf" "character varying(10)" NOT NULL,
  "bin" smallint NOT NULL,
  "quantity" smallint NOT NULL DEFAULT 0,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_ProductInventory_Bin" CHECK (((bin >= 0) AND (bin <= 100))),
  PRIMARY KEY ("productid", "locationid")
);

CREATE TABLE "production"."productlistpricehistory" (
  "productid" integer NOT NULL,
  "startdate" timestamp NOT NULL,
  "enddate" timestamp,
  "listprice" numeric NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_ProductListPriceHistory_EndDate" CHECK (((enddate >= startdate) OR (enddate IS NULL))),
  CONSTRAINT "CK_ProductListPriceHistory_ListPrice" CHECK ((listprice > 0.00)),
  PRIMARY KEY ("productid", "startdate")
);

CREATE TABLE "production"."productmodel" (
  "productmodelid" integer PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "catalogdescription" xml,
  "instructions" xml,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "production"."productmodelillustration" (
  "productmodelid" integer NOT NULL,
  "illustrationid" integer NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ("productmodelid", "illustrationid")
);

CREATE TABLE "production"."productmodelproductdescriptionculture" (
  "productmodelid" integer NOT NULL,
  "productdescriptionid" integer NOT NULL,
  "cultureid" character(6) NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ("productmodelid", "productdescriptionid", "cultureid")
);

CREATE TABLE "production"."productphoto" (
  "productphotoid" integer PRIMARY KEY NOT NULL,
  "thumbnailphoto" bytea,
  "thumbnailphotofilename" "character varying(50)",
  "largephoto" bytea,
  "largephotofilename" "character varying(50)",
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "production"."productproductphoto" (
  "productid" integer NOT NULL,
  "productphotoid" integer NOT NULL,
  "primary" "public.Flag" NOT NULL DEFAULT false,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ("productid", "productphotoid")
);

CREATE TABLE "production"."productreview" (
  "productreviewid" integer PRIMARY KEY NOT NULL,
  "productid" integer NOT NULL,
  "reviewername" "public.Name" NOT NULL,
  "reviewdate" timestamp NOT NULL DEFAULT (now()),
  "emailaddress" "character varying(50)" NOT NULL,
  "rating" integer NOT NULL,
  "comments" "character varying(3850)",
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_ProductReview_Rating" CHECK (((rating >= 1) AND (rating <= 5)))
);

CREATE TABLE "production"."productsubcategory" (
  "productsubcategoryid" integer PRIMARY KEY NOT NULL,
  "productcategoryid" integer NOT NULL,
  "name" "public.Name" NOT NULL,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "production"."scrapreason" (
  "scrapreasonid" integer PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "production"."transactionhistory" (
  "transactionid" integer PRIMARY KEY NOT NULL,
  "productid" integer NOT NULL,
  "referenceorderid" integer NOT NULL,
  "referenceorderlineid" integer NOT NULL DEFAULT 0,
  "transactiondate" timestamp NOT NULL DEFAULT (now()),
  "transactiontype" character(1) NOT NULL,
  "quantity" integer NOT NULL,
  "actualcost" numeric NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_TransactionHistory_TransactionType" CHECK ((upper((transactiontype)::text) = ANY (ARRAY['W'::text, 'S'::text, 'P'::text])))
);

CREATE TABLE "production"."transactionhistoryarchive" (
  "transactionid" integer PRIMARY KEY NOT NULL,
  "productid" integer NOT NULL,
  "referenceorderid" integer NOT NULL,
  "referenceorderlineid" integer NOT NULL DEFAULT 0,
  "transactiondate" timestamp NOT NULL DEFAULT (now()),
  "transactiontype" character(1) NOT NULL,
  "quantity" integer NOT NULL,
  "actualcost" numeric NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_TransactionHistoryArchive_TransactionType" CHECK ((upper((transactiontype)::text) = ANY (ARRAY['W'::text, 'S'::text, 'P'::text])))
);

CREATE TABLE "production"."unitmeasure" (
  "unitmeasurecode" character(3) PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "production"."workorder" (
  "workorderid" integer PRIMARY KEY NOT NULL,
  "productid" integer NOT NULL,
  "orderqty" integer NOT NULL,
  "scrappedqty" smallint NOT NULL,
  "startdate" timestamp NOT NULL,
  "enddate" timestamp,
  "duedate" timestamp NOT NULL,
  "scrapreasonid" smallint,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_WorkOrder_EndDate" CHECK (((enddate >= startdate) OR (enddate IS NULL))),
  CONSTRAINT "CK_WorkOrder_OrderQty" CHECK ((orderqty > 0)),
  CONSTRAINT "CK_WorkOrder_ScrappedQty" CHECK ((scrappedqty >= 0))
);

CREATE TABLE "production"."workorderrouting" (
  "workorderid" integer NOT NULL,
  "productid" integer NOT NULL,
  "operationsequence" smallint NOT NULL,
  "locationid" smallint NOT NULL,
  "scheduledstartdate" timestamp NOT NULL,
  "scheduledenddate" timestamp NOT NULL,
  "actualstartdate" timestamp,
  "actualenddate" timestamp,
  "actualresourcehrs" numeric(9,4),
  "plannedcost" numeric NOT NULL,
  "actualcost" numeric,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_WorkOrderRouting_ActualCost" CHECK ((actualcost > 0.00)),
  CONSTRAINT "CK_WorkOrderRouting_ActualEndDate" CHECK (((actualenddate >= actualstartdate) OR (actualenddate IS NULL) OR (actualstartdate IS NULL))),
  CONSTRAINT "CK_WorkOrderRouting_ActualResourceHrs" CHECK ((actualresourcehrs >= 0.0000)),
  CONSTRAINT "CK_WorkOrderRouting_PlannedCost" CHECK ((plannedcost > 0.00)),
  CONSTRAINT "CK_WorkOrderRouting_ScheduledEndDate" CHECK ((scheduledenddate >= scheduledstartdate)),
  PRIMARY KEY ("workorderid", "productid", "operationsequence")
);

CREATE TABLE "purchasing"."purchaseorderdetail" (
  "purchaseorderid" integer NOT NULL,
  "purchaseorderdetailid" integer NOT NULL,
  "duedate" timestamp NOT NULL,
  "orderqty" smallint NOT NULL,
  "productid" integer NOT NULL,
  "unitprice" numeric NOT NULL,
  "receivedqty" numeric(8,2) NOT NULL,
  "rejectedqty" numeric(8,2) NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_PurchaseOrderDetail_OrderQty" CHECK ((orderqty > 0)),
  CONSTRAINT "CK_PurchaseOrderDetail_ReceivedQty" CHECK ((receivedqty >= 0.00)),
  CONSTRAINT "CK_PurchaseOrderDetail_RejectedQty" CHECK ((rejectedqty >= 0.00)),
  CONSTRAINT "CK_PurchaseOrderDetail_UnitPrice" CHECK ((unitprice >= 0.00)),
  PRIMARY KEY ("purchaseorderid", "purchaseorderdetailid")
);

CREATE TABLE "purchasing"."purchaseorderheader" (
  "purchaseorderid" integer PRIMARY KEY NOT NULL,
  "revisionnumber" smallint NOT NULL DEFAULT 0,
  "status" smallint NOT NULL DEFAULT 1,
  "employeeid" integer NOT NULL,
  "vendorid" integer NOT NULL,
  "shipmethodid" integer NOT NULL,
  "orderdate" timestamp NOT NULL DEFAULT (now()),
  "shipdate" timestamp,
  "subtotal" numeric NOT NULL DEFAULT 0,
  "taxamt" numeric NOT NULL DEFAULT 0,
  "freight" numeric NOT NULL DEFAULT 0,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_PurchaseOrderHeader_Freight" CHECK ((freight >= 0.00)),
  CONSTRAINT "CK_PurchaseOrderHeader_ShipDate" CHECK (((shipdate >= orderdate) OR (shipdate IS NULL))),
  CONSTRAINT "CK_PurchaseOrderHeader_Status" CHECK (((status >= 1) AND (status <= 4))),
  CONSTRAINT "CK_PurchaseOrderHeader_SubTotal" CHECK ((subtotal >= 0.00)),
  CONSTRAINT "CK_PurchaseOrderHeader_TaxAmt" CHECK ((taxamt >= 0.00))
);

CREATE TABLE "purchasing"."productvendor" (
  "productid" integer NOT NULL,
  "businessentityid" integer NOT NULL,
  "averageleadtime" integer NOT NULL,
  "standardprice" numeric NOT NULL,
  "lastreceiptcost" numeric,
  "lastreceiptdate" timestamp,
  "minorderqty" integer NOT NULL,
  "maxorderqty" integer NOT NULL,
  "onorderqty" integer,
  "unitmeasurecode" character(3) NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_ProductVendor_AverageLeadTime" CHECK ((averageleadtime >= 1)),
  CONSTRAINT "CK_ProductVendor_LastReceiptCost" CHECK ((lastreceiptcost > 0.00)),
  CONSTRAINT "CK_ProductVendor_MaxOrderQty" CHECK ((maxorderqty >= 1)),
  CONSTRAINT "CK_ProductVendor_MinOrderQty" CHECK ((minorderqty >= 1)),
  CONSTRAINT "CK_ProductVendor_OnOrderQty" CHECK ((onorderqty >= 0)),
  CONSTRAINT "CK_ProductVendor_StandardPrice" CHECK ((standardprice > 0.00)),
  PRIMARY KEY ("productid", "businessentityid")
);

CREATE TABLE "purchasing"."shipmethod" (
  "shipmethodid" integer PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "shipbase" numeric NOT NULL DEFAULT 0,
  "shiprate" numeric NOT NULL DEFAULT 0,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_ShipMethod_ShipBase" CHECK ((shipbase > 0.00)),
  CONSTRAINT "CK_ShipMethod_ShipRate" CHECK ((shiprate > 0.00))
);

CREATE TABLE "purchasing"."vendor" (
  "businessentityid" integer PRIMARY KEY NOT NULL,
  "accountnumber" "public.AccountNumber" NOT NULL,
  "name" "public.Name" NOT NULL,
  "creditrating" smallint NOT NULL,
  "preferredvendorstatus" "public.Flag" NOT NULL DEFAULT true,
  "activeflag" "public.Flag" NOT NULL DEFAULT true,
  "purchasingwebserviceurl" "character varying(1024)",
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_Vendor_CreditRating" CHECK (((creditrating >= 1) AND (creditrating <= 5)))
);

CREATE TABLE "sales"."customer" (
  "customerid" integer PRIMARY KEY NOT NULL,
  "personid" integer,
  "storeid" integer,
  "territoryid" integer,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sales"."creditcard" (
  "creditcardid" integer PRIMARY KEY NOT NULL,
  "cardtype" "character varying(50)" NOT NULL,
  "cardnumber" "character varying(25)" NOT NULL,
  "expmonth" smallint NOT NULL,
  "expyear" smallint NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sales"."currencyrate" (
  "currencyrateid" integer PRIMARY KEY NOT NULL,
  "currencyratedate" timestamp NOT NULL,
  "fromcurrencycode" character(3) NOT NULL,
  "tocurrencycode" character(3) NOT NULL,
  "averagerate" numeric NOT NULL,
  "endofdayrate" numeric NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sales"."countryregioncurrency" (
  "countryregioncode" "character varying(3)" NOT NULL,
  "currencycode" character(3) NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ("countryregioncode", "currencycode")
);

CREATE TABLE "sales"."currency" (
  "currencycode" character(3) PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sales"."personcreditcard" (
  "businessentityid" integer NOT NULL,
  "creditcardid" integer NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ("businessentityid", "creditcardid")
);

CREATE TABLE "sales"."store" (
  "businessentityid" integer PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "salespersonid" integer,
  "demographics" xml,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sales"."shoppingcartitem" (
  "shoppingcartitemid" integer PRIMARY KEY NOT NULL,
  "shoppingcartid" "character varying(50)" NOT NULL,
  "quantity" integer NOT NULL DEFAULT 1,
  "productid" integer NOT NULL,
  "datecreated" timestamp NOT NULL DEFAULT (now()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_ShoppingCartItem_Quantity" CHECK ((quantity >= 1))
);

CREATE TABLE "sales"."specialoffer" (
  "specialofferid" integer PRIMARY KEY NOT NULL,
  "description" "character varying(255)" NOT NULL,
  "discountpct" numeric NOT NULL DEFAULT 0,
  "type" "character varying(50)" NOT NULL,
  "category" "character varying(50)" NOT NULL,
  "startdate" timestamp NOT NULL,
  "enddate" timestamp NOT NULL,
  "minqty" integer NOT NULL DEFAULT 0,
  "maxqty" integer,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_SpecialOffer_DiscountPct" CHECK ((discountpct >= 0.00)),
  CONSTRAINT "CK_SpecialOffer_EndDate" CHECK ((enddate >= startdate)),
  CONSTRAINT "CK_SpecialOffer_MaxQty" CHECK ((maxqty >= 0)),
  CONSTRAINT "CK_SpecialOffer_MinQty" CHECK ((minqty >= 0))
);

CREATE TABLE "sales"."salesorderdetail" (
  "salesorderid" integer NOT NULL,
  "salesorderdetailid" integer NOT NULL,
  "carriertrackingnumber" "character varying(25)",
  "orderqty" smallint NOT NULL,
  "productid" integer NOT NULL,
  "specialofferid" integer NOT NULL,
  "unitprice" numeric NOT NULL,
  "unitpricediscount" numeric NOT NULL DEFAULT 0,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_SalesOrderDetail_OrderQty" CHECK ((orderqty > 0)),
  CONSTRAINT "CK_SalesOrderDetail_UnitPrice" CHECK ((unitprice >= 0.00)),
  CONSTRAINT "CK_SalesOrderDetail_UnitPriceDiscount" CHECK ((unitpricediscount >= 0.00)),
  PRIMARY KEY ("salesorderid", "salesorderdetailid")
);

CREATE TABLE "sales"."salesorderheader" (
  "salesorderid" integer PRIMARY KEY NOT NULL,
  "revisionnumber" smallint NOT NULL DEFAULT 0,
  "orderdate" timestamp NOT NULL DEFAULT (now()),
  "duedate" timestamp NOT NULL,
  "shipdate" timestamp,
  "status" smallint NOT NULL DEFAULT 1,
  "onlineorderflag" "public.Flag" NOT NULL DEFAULT true,
  "purchaseordernumber" "public.OrderNumber",
  "accountnumber" "public.AccountNumber",
  "customerid" integer NOT NULL,
  "salespersonid" integer,
  "territoryid" integer,
  "billtoaddressid" integer NOT NULL,
  "shiptoaddressid" integer NOT NULL,
  "shipmethodid" integer NOT NULL,
  "creditcardid" integer,
  "creditcardapprovalcode" "character varying(15)",
  "currencyrateid" integer,
  "subtotal" numeric NOT NULL DEFAULT 0,
  "taxamt" numeric NOT NULL DEFAULT 0,
  "freight" numeric NOT NULL DEFAULT 0,
  "totaldue" numeric,
  "comment" "character varying(128)",
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_SalesOrderHeader_DueDate" CHECK ((duedate >= orderdate)),
  CONSTRAINT "CK_SalesOrderHeader_Freight" CHECK ((freight >= 0.00)),
  CONSTRAINT "CK_SalesOrderHeader_ShipDate" CHECK (((shipdate >= orderdate) OR (shipdate IS NULL))),
  CONSTRAINT "CK_SalesOrderHeader_Status" CHECK (((status >= 0) AND (status <= 8))),
  CONSTRAINT "CK_SalesOrderHeader_SubTotal" CHECK ((subtotal >= 0.00)),
  CONSTRAINT "CK_SalesOrderHeader_TaxAmt" CHECK ((taxamt >= 0.00))
);

CREATE TABLE "sales"."salesorderheadersalesreason" (
  "salesorderid" integer NOT NULL,
  "salesreasonid" integer NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ("salesorderid", "salesreasonid")
);

CREATE TABLE "sales"."specialofferproduct" (
  "specialofferid" integer NOT NULL,
  "productid" integer NOT NULL,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ("specialofferid", "productid")
);

CREATE TABLE "sales"."salesperson" (
  "businessentityid" integer PRIMARY KEY NOT NULL,
  "territoryid" integer,
  "salesquota" numeric,
  "bonus" numeric NOT NULL DEFAULT 0,
  "commissionpct" numeric NOT NULL DEFAULT 0,
  "salesytd" numeric NOT NULL DEFAULT 0,
  "saleslastyear" numeric NOT NULL DEFAULT 0,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_SalesPerson_Bonus" CHECK ((bonus >= 0.00)),
  CONSTRAINT "CK_SalesPerson_CommissionPct" CHECK ((commissionpct >= 0.00)),
  CONSTRAINT "CK_SalesPerson_SalesLastYear" CHECK ((saleslastyear >= 0.00)),
  CONSTRAINT "CK_SalesPerson_SalesQuota" CHECK ((salesquota > 0.00)),
  CONSTRAINT "CK_SalesPerson_SalesYTD" CHECK ((salesytd >= 0.00))
);

CREATE TABLE "sales"."salespersonquotahistory" (
  "businessentityid" integer NOT NULL,
  "quotadate" timestamp NOT NULL,
  "salesquota" numeric NOT NULL,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_SalesPersonQuotaHistory_SalesQuota" CHECK ((salesquota > 0.00)),
  PRIMARY KEY ("businessentityid", "quotadate")
);

CREATE TABLE "sales"."salesreason" (
  "salesreasonid" integer PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "reasontype" "public.Name" NOT NULL,
  "modifieddate" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "sales"."salesterritory" (
  "territoryid" integer PRIMARY KEY NOT NULL,
  "name" "public.Name" NOT NULL,
  "countryregioncode" "character varying(3)" NOT NULL,
  "group" "character varying(50)" NOT NULL,
  "salesytd" numeric NOT NULL DEFAULT 0,
  "saleslastyear" numeric NOT NULL DEFAULT 0,
  "costytd" numeric NOT NULL DEFAULT 0,
  "costlastyear" numeric NOT NULL DEFAULT 0,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_SalesTerritory_CostLastYear" CHECK ((costlastyear >= 0.00)),
  CONSTRAINT "CK_SalesTerritory_CostYTD" CHECK ((costytd >= 0.00)),
  CONSTRAINT "CK_SalesTerritory_SalesLastYear" CHECK ((saleslastyear >= 0.00)),
  CONSTRAINT "CK_SalesTerritory_SalesYTD" CHECK ((salesytd >= 0.00))
);

CREATE TABLE "sales"."salesterritoryhistory" (
  "businessentityid" integer NOT NULL,
  "territoryid" integer NOT NULL,
  "startdate" timestamp NOT NULL,
  "enddate" timestamp,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_SalesTerritoryHistory_EndDate" CHECK (((enddate >= startdate) OR (enddate IS NULL))),
  PRIMARY KEY ("businessentityid", "startdate", "territoryid")
);

CREATE TABLE "sales"."salestaxrate" (
  "salestaxrateid" integer PRIMARY KEY NOT NULL,
  "stateprovinceid" integer NOT NULL,
  "taxtype" smallint NOT NULL,
  "taxrate" numeric NOT NULL DEFAULT 0,
  "name" "public.Name" NOT NULL,
  "rowguid" uuid NOT NULL DEFAULT (public.uuid_generate_v1()),
  "modifieddate" timestamp NOT NULL DEFAULT (now()),
  CONSTRAINT "CK_SalesTaxRate_TaxType" CHECK (((taxtype >= 1) AND (taxtype <= 3)))
);

COMMENT ON TABLE "humanresources"."department" IS 'Lookup table containing the departments within the Adventure Works Cycles company.';

COMMENT ON COLUMN "humanresources"."department"."departmentid" IS 'Primary key for Department records.';

COMMENT ON COLUMN "humanresources"."department"."name" IS 'Name of the department.';

COMMENT ON COLUMN "humanresources"."department"."groupname" IS 'Name of the group to which the department belongs.';

COMMENT ON TABLE "humanresources"."employee" IS 'Employee information such as salary, department, and title.';

COMMENT ON COLUMN "humanresources"."employee"."businessentityid" IS 'Primary key for Employee records.  Foreign key to BusinessEntity.BusinessEntityID.';

COMMENT ON COLUMN "humanresources"."employee"."nationalidnumber" IS 'Unique national identification number such as a social security number.';

COMMENT ON COLUMN "humanresources"."employee"."loginid" IS 'Network login.';

COMMENT ON COLUMN "humanresources"."employee"."jobtitle" IS 'Work title such as Buyer or Sales Representative.';

COMMENT ON COLUMN "humanresources"."employee"."birthdate" IS 'Date of birth.';

COMMENT ON COLUMN "humanresources"."employee"."maritalstatus" IS 'M = Married, S = Single';

COMMENT ON COLUMN "humanresources"."employee"."gender" IS 'M = Male, F = Female';

COMMENT ON COLUMN "humanresources"."employee"."hiredate" IS 'Employee hired on this date.';

COMMENT ON COLUMN "humanresources"."employee"."salariedflag" IS 'Job classification. 0 = Hourly, not exempt from collective bargaining. 1 = Salaried, exempt from collective bargaining.';

COMMENT ON COLUMN "humanresources"."employee"."vacationhours" IS 'Number of available vacation hours.';

COMMENT ON COLUMN "humanresources"."employee"."sickleavehours" IS 'Number of available sick leave hours.';

COMMENT ON COLUMN "humanresources"."employee"."currentflag" IS '0 = Inactive, 1 = Active';

COMMENT ON COLUMN "humanresources"."employee"."organizationnode" IS 'Where the employee is located in corporate hierarchy.';

COMMENT ON TABLE "humanresources"."employeedepartmenthistory" IS 'Employee department transfers.';

COMMENT ON COLUMN "humanresources"."employeedepartmenthistory"."businessentityid" IS 'Employee identification number. Foreign key to Employee.BusinessEntityID.';

COMMENT ON COLUMN "humanresources"."employeedepartmenthistory"."departmentid" IS 'Department in which the employee worked including currently. Foreign key to Department.DepartmentID.';

COMMENT ON COLUMN "humanresources"."employeedepartmenthistory"."shiftid" IS 'Identifies which 8-hour shift the employee works. Foreign key to Shift.Shift.ID.';

COMMENT ON COLUMN "humanresources"."employeedepartmenthistory"."startdate" IS 'Date the employee started work in the department.';

COMMENT ON COLUMN "humanresources"."employeedepartmenthistory"."enddate" IS 'Date the employee left the department. NULL = Current department.';

COMMENT ON TABLE "humanresources"."employeepayhistory" IS 'Employee pay history.';

COMMENT ON COLUMN "humanresources"."employeepayhistory"."businessentityid" IS 'Employee identification number. Foreign key to Employee.BusinessEntityID.';

COMMENT ON COLUMN "humanresources"."employeepayhistory"."ratechangedate" IS 'Date the change in pay is effective';

COMMENT ON COLUMN "humanresources"."employeepayhistory"."rate" IS 'Salary hourly rate.';

COMMENT ON COLUMN "humanresources"."employeepayhistory"."payfrequency" IS '1 = Salary received monthly, 2 = Salary received biweekly';

COMMENT ON TABLE "humanresources"."jobcandidate" IS 'RÃ©sumÃ©s submitted to Human Resources by job applicants.';

COMMENT ON COLUMN "humanresources"."jobcandidate"."jobcandidateid" IS 'Primary key for JobCandidate records.';

COMMENT ON COLUMN "humanresources"."jobcandidate"."businessentityid" IS 'Employee identification number if applicant was hired. Foreign key to Employee.BusinessEntityID.';

COMMENT ON COLUMN "humanresources"."jobcandidate"."resume" IS 'RÃ©sumÃ© in XML format.';

COMMENT ON TABLE "humanresources"."shift" IS 'Work shift lookup table.';

COMMENT ON COLUMN "humanresources"."shift"."shiftid" IS 'Primary key for Shift records.';

COMMENT ON COLUMN "humanresources"."shift"."name" IS 'Shift description.';

COMMENT ON COLUMN "humanresources"."shift"."starttime" IS 'Shift start time.';

COMMENT ON COLUMN "humanresources"."shift"."endtime" IS 'Shift end time.';

COMMENT ON TABLE "person"."address" IS 'Street address information for customers, employees, and vendors.';

COMMENT ON COLUMN "person"."address"."addressid" IS 'Primary key for Address records.';

COMMENT ON COLUMN "person"."address"."addressline1" IS 'First street address line.';

COMMENT ON COLUMN "person"."address"."addressline2" IS 'Second street address line.';

COMMENT ON COLUMN "person"."address"."city" IS 'Name of the city.';

COMMENT ON COLUMN "person"."address"."stateprovinceid" IS 'Unique identification number for the state or province. Foreign key to StateProvince table.';

COMMENT ON COLUMN "person"."address"."postalcode" IS 'Postal code for the street address.';

COMMENT ON COLUMN "person"."address"."spatiallocation" IS 'Latitude and longitude of this address.';

COMMENT ON TABLE "person"."businessentityaddress" IS 'Cross-reference table mapping customers, vendors, and employees to their addresses.';

COMMENT ON COLUMN "person"."businessentityaddress"."businessentityid" IS 'Primary key. Foreign key to BusinessEntity.BusinessEntityID.';

COMMENT ON COLUMN "person"."businessentityaddress"."addressid" IS 'Primary key. Foreign key to Address.AddressID.';

COMMENT ON COLUMN "person"."businessentityaddress"."addresstypeid" IS 'Primary key. Foreign key to AddressType.AddressTypeID.';

COMMENT ON TABLE "person"."countryregion" IS 'Lookup table containing the ISO standard codes for countries and regions.';

COMMENT ON COLUMN "person"."countryregion"."countryregioncode" IS 'ISO standard code for countries and regions.';

COMMENT ON COLUMN "person"."countryregion"."name" IS 'Country or region name.';

COMMENT ON TABLE "person"."emailaddress" IS 'Where to send a person email.';

COMMENT ON COLUMN "person"."emailaddress"."businessentityid" IS 'Primary key. Person associated with this email address.  Foreign key to Person.BusinessEntityID';

COMMENT ON COLUMN "person"."emailaddress"."emailaddressid" IS 'Primary key. ID of this email address.';

COMMENT ON COLUMN "person"."emailaddress"."emailaddress" IS 'E-mail address for the person.';

COMMENT ON TABLE "person"."person" IS 'Human beings involved with AdventureWorks: employees, customer contacts, and vendor contacts.';

COMMENT ON COLUMN "person"."person"."businessentityid" IS 'Primary key for Person records.';

COMMENT ON COLUMN "person"."person"."persontype" IS 'Primary type of person: SC = Store Contact, IN = Individual (retail) customer, SP = Sales person, EM = Employee (non-sales), VC = Vendor contact, GC = General contact';

COMMENT ON COLUMN "person"."person"."namestyle" IS '0 = The data in FirstName and LastName are stored in western style (first name, last name) order.  1 = Eastern style (last name, first name) order.';

COMMENT ON COLUMN "person"."person"."title" IS 'A courtesy title. For example, Mr. or Ms.';

COMMENT ON COLUMN "person"."person"."firstname" IS 'First name of the person.';

COMMENT ON COLUMN "person"."person"."middlename" IS 'Middle name or middle initial of the person.';

COMMENT ON COLUMN "person"."person"."lastname" IS 'Last name of the person.';

COMMENT ON COLUMN "person"."person"."suffix" IS 'Surname suffix. For example, Sr. or Jr.';

COMMENT ON COLUMN "person"."person"."emailpromotion" IS '0 = Contact does not wish to receive e-mail promotions, 1 = Contact does wish to receive e-mail promotions from AdventureWorks, 2 = Contact does wish to receive e-mail promotions from AdventureWorks and selected partners.';

COMMENT ON COLUMN "person"."person"."additionalcontactinfo" IS 'Additional contact information about the person stored in xml format.';

COMMENT ON COLUMN "person"."person"."demographics" IS 'Personal information such as hobbies, and income collected from online shoppers. Used for sales analysis.';

COMMENT ON TABLE "person"."personphone" IS 'Telephone number and type of a person.';

COMMENT ON COLUMN "person"."personphone"."businessentityid" IS 'Business entity identification number. Foreign key to Person.BusinessEntityID.';

COMMENT ON COLUMN "person"."personphone"."phonenumber" IS 'Telephone number identification number.';

COMMENT ON COLUMN "person"."personphone"."phonenumbertypeid" IS 'Kind of phone number. Foreign key to PhoneNumberType.PhoneNumberTypeID.';

COMMENT ON TABLE "person"."phonenumbertype" IS 'Type of phone number of a person.';

COMMENT ON COLUMN "person"."phonenumbertype"."phonenumbertypeid" IS 'Primary key for telephone number type records.';

COMMENT ON COLUMN "person"."phonenumbertype"."name" IS 'Name of the telephone number type';

COMMENT ON TABLE "person"."stateprovince" IS 'State and province lookup table.';

COMMENT ON COLUMN "person"."stateprovince"."stateprovinceid" IS 'Primary key for StateProvince records.';

COMMENT ON COLUMN "person"."stateprovince"."stateprovincecode" IS 'ISO standard state or province code.';

COMMENT ON COLUMN "person"."stateprovince"."countryregioncode" IS 'ISO standard country or region code. Foreign key to CountryRegion.CountryRegionCode.';

COMMENT ON COLUMN "person"."stateprovince"."isonlystateprovinceflag" IS '0 = StateProvinceCode exists. 1 = StateProvinceCode unavailable, using CountryRegionCode.';

COMMENT ON COLUMN "person"."stateprovince"."name" IS 'State or province description.';

COMMENT ON COLUMN "person"."stateprovince"."territoryid" IS 'ID of the territory in which the state or province is located. Foreign key to SalesTerritory.SalesTerritoryID.';

COMMENT ON TABLE "person"."addresstype" IS 'Types of addresses stored in the Address table.';

COMMENT ON COLUMN "person"."addresstype"."addresstypeid" IS 'Primary key for AddressType records.';

COMMENT ON COLUMN "person"."addresstype"."name" IS 'Address type description. For example, Billing, Home, or Shipping.';

COMMENT ON TABLE "person"."businessentity" IS 'Source of the ID that connects vendors, customers, and employees with address and contact information.';

COMMENT ON COLUMN "person"."businessentity"."businessentityid" IS 'Primary key for all customers, vendors, and employees.';

COMMENT ON TABLE "person"."businessentitycontact" IS 'Cross-reference table mapping stores, vendors, and employees to people';

COMMENT ON COLUMN "person"."businessentitycontact"."businessentityid" IS 'Primary key. Foreign key to BusinessEntity.BusinessEntityID.';

COMMENT ON COLUMN "person"."businessentitycontact"."personid" IS 'Primary key. Foreign key to Person.BusinessEntityID.';

COMMENT ON COLUMN "person"."businessentitycontact"."contacttypeid" IS 'Primary key.  Foreign key to ContactType.ContactTypeID.';

COMMENT ON TABLE "person"."contacttype" IS 'Lookup table containing the types of business entity contacts.';

COMMENT ON COLUMN "person"."contacttype"."contacttypeid" IS 'Primary key for ContactType records.';

COMMENT ON COLUMN "person"."contacttype"."name" IS 'Contact type description.';

COMMENT ON TABLE "person"."password" IS 'One way hashed authentication information';

COMMENT ON COLUMN "person"."password"."passwordhash" IS 'Password for the e-mail account.';

COMMENT ON COLUMN "person"."password"."passwordsalt" IS 'Random value concatenated with the password string before the password is hashed.';

COMMENT ON TABLE "production"."billofmaterials" IS 'Items required to make bicycles and bicycle subassemblies. It identifies the heirarchical relationship between a parent product and its components.';

COMMENT ON COLUMN "production"."billofmaterials"."billofmaterialsid" IS 'Primary key for BillOfMaterials records.';

COMMENT ON COLUMN "production"."billofmaterials"."productassemblyid" IS 'Parent product identification number. Foreign key to Product.ProductID.';

COMMENT ON COLUMN "production"."billofmaterials"."componentid" IS 'Component identification number. Foreign key to Product.ProductID.';

COMMENT ON COLUMN "production"."billofmaterials"."startdate" IS 'Date the component started being used in the assembly item.';

COMMENT ON COLUMN "production"."billofmaterials"."enddate" IS 'Date the component stopped being used in the assembly item.';

COMMENT ON COLUMN "production"."billofmaterials"."unitmeasurecode" IS 'Standard code identifying the unit of measure for the quantity.';

COMMENT ON COLUMN "production"."billofmaterials"."bomlevel" IS 'Indicates the depth the component is from its parent (AssemblyID).';

COMMENT ON COLUMN "production"."billofmaterials"."perassemblyqty" IS 'Quantity of the component needed to create the assembly.';

COMMENT ON TABLE "production"."culture" IS 'Lookup table containing the languages in which some AdventureWorks data is stored.';

COMMENT ON COLUMN "production"."culture"."cultureid" IS 'Primary key for Culture records.';

COMMENT ON COLUMN "production"."culture"."name" IS 'Culture description.';

COMMENT ON TABLE "production"."document" IS 'Product maintenance documents.';

COMMENT ON COLUMN "production"."document"."title" IS 'Title of the document.';

COMMENT ON COLUMN "production"."document"."owner" IS 'Employee who controls the document.  Foreign key to Employee.BusinessEntityID';

COMMENT ON COLUMN "production"."document"."folderflag" IS '0 = This is a folder, 1 = This is a document.';

COMMENT ON COLUMN "production"."document"."filename" IS 'File name of the document';

COMMENT ON COLUMN "production"."document"."fileextension" IS 'File extension indicating the document type. For example, .doc or .txt.';

COMMENT ON COLUMN "production"."document"."revision" IS 'Revision number of the document.';

COMMENT ON COLUMN "production"."document"."changenumber" IS 'Engineering change approval number.';

COMMENT ON COLUMN "production"."document"."status" IS '1 = Pending approval, 2 = Approved, 3 = Obsolete';

COMMENT ON COLUMN "production"."document"."documentsummary" IS 'Document abstract.';

COMMENT ON COLUMN "production"."document"."document" IS 'Complete document.';

COMMENT ON COLUMN "production"."document"."rowguid" IS 'ROWGUIDCOL number uniquely identifying the record. Required for FileStream.';

COMMENT ON COLUMN "production"."document"."documentnode" IS 'Primary key for Document records.';

COMMENT ON TABLE "production"."illustration" IS 'Bicycle assembly diagrams.';

COMMENT ON COLUMN "production"."illustration"."illustrationid" IS 'Primary key for Illustration records.';

COMMENT ON COLUMN "production"."illustration"."diagram" IS 'Illustrations used in manufacturing instructions. Stored as XML.';

COMMENT ON TABLE "production"."location" IS 'Product inventory and manufacturing locations.';

COMMENT ON COLUMN "production"."location"."locationid" IS 'Primary key for Location records.';

COMMENT ON COLUMN "production"."location"."name" IS 'Location description.';

COMMENT ON COLUMN "production"."location"."costrate" IS 'Standard hourly cost of the manufacturing location.';

COMMENT ON COLUMN "production"."location"."availability" IS 'Work capacity (in hours) of the manufacturing location.';

COMMENT ON TABLE "production"."product" IS 'Products sold or used in the manfacturing of sold products.';

COMMENT ON COLUMN "production"."product"."productid" IS 'Primary key for Product records.';

COMMENT ON COLUMN "production"."product"."name" IS 'Name of the product.';

COMMENT ON COLUMN "production"."product"."productnumber" IS 'Unique product identification number.';

COMMENT ON COLUMN "production"."product"."makeflag" IS '0 = Product is purchased, 1 = Product is manufactured in-house.';

COMMENT ON COLUMN "production"."product"."finishedgoodsflag" IS '0 = Product is not a salable item. 1 = Product is salable.';

COMMENT ON COLUMN "production"."product"."color" IS 'Product color.';

COMMENT ON COLUMN "production"."product"."safetystocklevel" IS 'Minimum inventory quantity.';

COMMENT ON COLUMN "production"."product"."reorderpoint" IS 'Inventory level that triggers a purchase order or work order.';

COMMENT ON COLUMN "production"."product"."standardcost" IS 'Standard cost of the product.';

COMMENT ON COLUMN "production"."product"."listprice" IS 'Selling price.';

COMMENT ON COLUMN "production"."product"."size" IS 'Product size.';

COMMENT ON COLUMN "production"."product"."sizeunitmeasurecode" IS 'Unit of measure for Size column.';

COMMENT ON COLUMN "production"."product"."weightunitmeasurecode" IS 'Unit of measure for Weight column.';

COMMENT ON COLUMN "production"."product"."weight" IS 'Product weight.';

COMMENT ON COLUMN "production"."product"."daystomanufacture" IS 'Number of days required to manufacture the product.';

COMMENT ON COLUMN "production"."product"."productline" IS 'R = Road, M = Mountain, T = Touring, S = Standard';

COMMENT ON COLUMN "production"."product"."class" IS 'H = High, M = Medium, L = Low';

COMMENT ON COLUMN "production"."product"."style" IS 'W = Womens, M = Mens, U = Universal';

COMMENT ON COLUMN "production"."product"."productsubcategoryid" IS 'Product is a member of this product subcategory. Foreign key to ProductSubCategory.ProductSubCategoryID.';

COMMENT ON COLUMN "production"."product"."productmodelid" IS 'Product is a member of this product model. Foreign key to ProductModel.ProductModelID.';

COMMENT ON COLUMN "production"."product"."sellstartdate" IS 'Date the product was available for sale.';

COMMENT ON COLUMN "production"."product"."sellenddate" IS 'Date the product was no longer available for sale.';

COMMENT ON COLUMN "production"."product"."discontinueddate" IS 'Date the product was discontinued.';

COMMENT ON TABLE "production"."productcategory" IS 'High-level product categorization.';

COMMENT ON COLUMN "production"."productcategory"."productcategoryid" IS 'Primary key for ProductCategory records.';

COMMENT ON COLUMN "production"."productcategory"."name" IS 'Category description.';

COMMENT ON TABLE "production"."productcosthistory" IS 'Changes in the cost of a product over time.';

COMMENT ON COLUMN "production"."productcosthistory"."productid" IS 'Product identification number. Foreign key to Product.ProductID';

COMMENT ON COLUMN "production"."productcosthistory"."startdate" IS 'Product cost start date.';

COMMENT ON COLUMN "production"."productcosthistory"."enddate" IS 'Product cost end date.';

COMMENT ON COLUMN "production"."productcosthistory"."standardcost" IS 'Standard cost of the product.';

COMMENT ON TABLE "production"."productdescription" IS 'Product descriptions in several languages.';

COMMENT ON COLUMN "production"."productdescription"."productdescriptionid" IS 'Primary key for ProductDescription records.';

COMMENT ON COLUMN "production"."productdescription"."description" IS 'Description of the product.';

COMMENT ON TABLE "production"."productdocument" IS 'Cross-reference table mapping products to related product documents.';

COMMENT ON COLUMN "production"."productdocument"."productid" IS 'Product identification number. Foreign key to Product.ProductID.';

COMMENT ON COLUMN "production"."productdocument"."documentnode" IS 'Document identification number. Foreign key to Document.DocumentNode.';

COMMENT ON TABLE "production"."productinventory" IS 'Product inventory information.';

COMMENT ON COLUMN "production"."productinventory"."productid" IS 'Product identification number. Foreign key to Product.ProductID.';

COMMENT ON COLUMN "production"."productinventory"."locationid" IS 'Inventory location identification number. Foreign key to Location.LocationID.';

COMMENT ON COLUMN "production"."productinventory"."shelf" IS 'Storage compartment within an inventory location.';

COMMENT ON COLUMN "production"."productinventory"."bin" IS 'Storage container on a shelf in an inventory location.';

COMMENT ON COLUMN "production"."productinventory"."quantity" IS 'Quantity of products in the inventory location.';

COMMENT ON TABLE "production"."productlistpricehistory" IS 'Changes in the list price of a product over time.';

COMMENT ON COLUMN "production"."productlistpricehistory"."productid" IS 'Product identification number. Foreign key to Product.ProductID';

COMMENT ON COLUMN "production"."productlistpricehistory"."startdate" IS 'List price start date.';

COMMENT ON COLUMN "production"."productlistpricehistory"."enddate" IS 'List price end date';

COMMENT ON COLUMN "production"."productlistpricehistory"."listprice" IS 'Product list price.';

COMMENT ON TABLE "production"."productmodel" IS 'Product model classification.';

COMMENT ON COLUMN "production"."productmodel"."productmodelid" IS 'Primary key for ProductModel records.';

COMMENT ON COLUMN "production"."productmodel"."name" IS 'Product model description.';

COMMENT ON COLUMN "production"."productmodel"."catalogdescription" IS 'Detailed product catalog information in xml format.';

COMMENT ON COLUMN "production"."productmodel"."instructions" IS 'Manufacturing instructions in xml format.';

COMMENT ON TABLE "production"."productmodelillustration" IS 'Cross-reference table mapping product models and illustrations.';

COMMENT ON COLUMN "production"."productmodelillustration"."productmodelid" IS 'Primary key. Foreign key to ProductModel.ProductModelID.';

COMMENT ON COLUMN "production"."productmodelillustration"."illustrationid" IS 'Primary key. Foreign key to Illustration.IllustrationID.';

COMMENT ON TABLE "production"."productmodelproductdescriptionculture" IS 'Cross-reference table mapping product descriptions and the language the description is written in.';

COMMENT ON COLUMN "production"."productmodelproductdescriptionculture"."productmodelid" IS 'Primary key. Foreign key to ProductModel.ProductModelID.';

COMMENT ON COLUMN "production"."productmodelproductdescriptionculture"."productdescriptionid" IS 'Primary key. Foreign key to ProductDescription.ProductDescriptionID.';

COMMENT ON COLUMN "production"."productmodelproductdescriptionculture"."cultureid" IS 'Culture identification number. Foreign key to Culture.CultureID.';

COMMENT ON TABLE "production"."productphoto" IS 'Product images.';

COMMENT ON COLUMN "production"."productphoto"."productphotoid" IS 'Primary key for ProductPhoto records.';

COMMENT ON COLUMN "production"."productphoto"."thumbnailphoto" IS 'Small image of the product.';

COMMENT ON COLUMN "production"."productphoto"."thumbnailphotofilename" IS 'Small image file name.';

COMMENT ON COLUMN "production"."productphoto"."largephoto" IS 'Large image of the product.';

COMMENT ON COLUMN "production"."productphoto"."largephotofilename" IS 'Large image file name.';

COMMENT ON TABLE "production"."productproductphoto" IS 'Cross-reference table mapping products and product photos.';

COMMENT ON COLUMN "production"."productproductphoto"."productid" IS 'Product identification number. Foreign key to Product.ProductID.';

COMMENT ON COLUMN "production"."productproductphoto"."productphotoid" IS 'Product photo identification number. Foreign key to ProductPhoto.ProductPhotoID.';

COMMENT ON COLUMN "production"."productproductphoto"."primary" IS '0 = Photo is not the principal image. 1 = Photo is the principal image.';

COMMENT ON TABLE "production"."productreview" IS 'Customer reviews of products they have purchased.';

COMMENT ON COLUMN "production"."productreview"."productreviewid" IS 'Primary key for ProductReview records.';

COMMENT ON COLUMN "production"."productreview"."productid" IS 'Product identification number. Foreign key to Product.ProductID.';

COMMENT ON COLUMN "production"."productreview"."reviewername" IS 'Name of the reviewer.';

COMMENT ON COLUMN "production"."productreview"."reviewdate" IS 'Date review was submitted.';

COMMENT ON COLUMN "production"."productreview"."emailaddress" IS 'Reviewer''s e-mail address.';

COMMENT ON COLUMN "production"."productreview"."rating" IS 'Product rating given by the reviewer. Scale is 1 to 5 with 5 as the highest rating.';

COMMENT ON COLUMN "production"."productreview"."comments" IS 'Reviewer''s comments
WARNING: can have several lines!';

COMMENT ON TABLE "production"."productsubcategory" IS 'Product subcategories. See ProductCategory table.';

COMMENT ON COLUMN "production"."productsubcategory"."productsubcategoryid" IS 'Primary key for ProductSubcategory records.';

COMMENT ON COLUMN "production"."productsubcategory"."productcategoryid" IS 'Product category identification number. Foreign key to ProductCategory.ProductCategoryID.';

COMMENT ON COLUMN "production"."productsubcategory"."name" IS 'Subcategory description.';

COMMENT ON TABLE "production"."scrapreason" IS 'Manufacturing failure reasons lookup table.';

COMMENT ON COLUMN "production"."scrapreason"."scrapreasonid" IS 'Primary key for ScrapReason records.';

COMMENT ON COLUMN "production"."scrapreason"."name" IS 'Failure description.';

COMMENT ON TABLE "production"."transactionhistory" IS 'Record of each purchase order, sales order, or work order transaction year to date.';

COMMENT ON COLUMN "production"."transactionhistory"."transactionid" IS 'Primary key for TransactionHistory records.';

COMMENT ON COLUMN "production"."transactionhistory"."productid" IS 'Product identification number. Foreign key to Product.ProductID.';

COMMENT ON COLUMN "production"."transactionhistory"."referenceorderid" IS 'Purchase order, sales order, or work order identification number.';

COMMENT ON COLUMN "production"."transactionhistory"."referenceorderlineid" IS 'Line number associated with the purchase order, sales order, or work order.';

COMMENT ON COLUMN "production"."transactionhistory"."transactiondate" IS 'Date and time of the transaction.';

COMMENT ON COLUMN "production"."transactionhistory"."transactiontype" IS 'W = WorkOrder, S = SalesOrder, P = PurchaseOrder';

COMMENT ON COLUMN "production"."transactionhistory"."quantity" IS 'Product quantity.';

COMMENT ON COLUMN "production"."transactionhistory"."actualcost" IS 'Product cost.';

COMMENT ON TABLE "production"."transactionhistoryarchive" IS 'Transactions for previous years.';

COMMENT ON COLUMN "production"."transactionhistoryarchive"."transactionid" IS 'Primary key for TransactionHistoryArchive records.';

COMMENT ON COLUMN "production"."transactionhistoryarchive"."productid" IS 'Product identification number. Foreign key to Product.ProductID.';

COMMENT ON COLUMN "production"."transactionhistoryarchive"."referenceorderid" IS 'Purchase order, sales order, or work order identification number.';

COMMENT ON COLUMN "production"."transactionhistoryarchive"."referenceorderlineid" IS 'Line number associated with the purchase order, sales order, or work order.';

COMMENT ON COLUMN "production"."transactionhistoryarchive"."transactiondate" IS 'Date and time of the transaction.';

COMMENT ON COLUMN "production"."transactionhistoryarchive"."transactiontype" IS 'W = Work Order, S = Sales Order, P = Purchase Order';

COMMENT ON COLUMN "production"."transactionhistoryarchive"."quantity" IS 'Product quantity.';

COMMENT ON COLUMN "production"."transactionhistoryarchive"."actualcost" IS 'Product cost.';

COMMENT ON TABLE "production"."unitmeasure" IS 'Unit of measure lookup table.';

COMMENT ON COLUMN "production"."unitmeasure"."unitmeasurecode" IS 'Primary key.';

COMMENT ON COLUMN "production"."unitmeasure"."name" IS 'Unit of measure description.';

COMMENT ON TABLE "production"."workorder" IS 'Manufacturing work orders.';

COMMENT ON COLUMN "production"."workorder"."workorderid" IS 'Primary key for WorkOrder records.';

COMMENT ON COLUMN "production"."workorder"."productid" IS 'Product identification number. Foreign key to Product.ProductID.';

COMMENT ON COLUMN "production"."workorder"."orderqty" IS 'Product quantity to build.';

COMMENT ON COLUMN "production"."workorder"."scrappedqty" IS 'Quantity that failed inspection.';

COMMENT ON COLUMN "production"."workorder"."startdate" IS 'Work order start date.';

COMMENT ON COLUMN "production"."workorder"."enddate" IS 'Work order end date.';

COMMENT ON COLUMN "production"."workorder"."duedate" IS 'Work order due date.';

COMMENT ON COLUMN "production"."workorder"."scrapreasonid" IS 'Reason for inspection failure.';

COMMENT ON TABLE "production"."workorderrouting" IS 'Work order details.';

COMMENT ON COLUMN "production"."workorderrouting"."workorderid" IS 'Primary key. Foreign key to WorkOrder.WorkOrderID.';

COMMENT ON COLUMN "production"."workorderrouting"."productid" IS 'Primary key. Foreign key to Product.ProductID.';

COMMENT ON COLUMN "production"."workorderrouting"."operationsequence" IS 'Primary key. Indicates the manufacturing process sequence.';

COMMENT ON COLUMN "production"."workorderrouting"."locationid" IS 'Manufacturing location where the part is processed. Foreign key to Location.LocationID.';

COMMENT ON COLUMN "production"."workorderrouting"."scheduledstartdate" IS 'Planned manufacturing start date.';

COMMENT ON COLUMN "production"."workorderrouting"."scheduledenddate" IS 'Planned manufacturing end date.';

COMMENT ON COLUMN "production"."workorderrouting"."actualstartdate" IS 'Actual start date.';

COMMENT ON COLUMN "production"."workorderrouting"."actualenddate" IS 'Actual end date.';

COMMENT ON COLUMN "production"."workorderrouting"."actualresourcehrs" IS 'Number of manufacturing hours used.';

COMMENT ON COLUMN "production"."workorderrouting"."plannedcost" IS 'Estimated manufacturing cost.';

COMMENT ON COLUMN "production"."workorderrouting"."actualcost" IS 'Actual manufacturing cost.';

COMMENT ON TABLE "purchasing"."purchaseorderdetail" IS 'Individual products associated with a specific purchase order. See PurchaseOrderHeader.';

COMMENT ON COLUMN "purchasing"."purchaseorderdetail"."purchaseorderid" IS 'Primary key. Foreign key to PurchaseOrderHeader.PurchaseOrderID.';

COMMENT ON COLUMN "purchasing"."purchaseorderdetail"."purchaseorderdetailid" IS 'Primary key. One line number per purchased product.';

COMMENT ON COLUMN "purchasing"."purchaseorderdetail"."duedate" IS 'Date the product is expected to be received.';

COMMENT ON COLUMN "purchasing"."purchaseorderdetail"."orderqty" IS 'Quantity ordered.';

COMMENT ON COLUMN "purchasing"."purchaseorderdetail"."productid" IS 'Product identification number. Foreign key to Product.ProductID.';

COMMENT ON COLUMN "purchasing"."purchaseorderdetail"."unitprice" IS 'Vendor''s selling price of a single product.';

COMMENT ON COLUMN "purchasing"."purchaseorderdetail"."receivedqty" IS 'Quantity actually received from the vendor.';

COMMENT ON COLUMN "purchasing"."purchaseorderdetail"."rejectedqty" IS 'Quantity rejected during inspection.';

COMMENT ON TABLE "purchasing"."purchaseorderheader" IS 'General purchase order information. See PurchaseOrderDetail
for more details.';

COMMENT ON COLUMN "purchasing"."purchaseorderheader"."purchaseorderid" IS 'Primary key.';

COMMENT ON COLUMN "purchasing"."purchaseorderheader"."revisionnumber" IS 'Incremental number to track changes to the purchase order over time.';

COMMENT ON COLUMN "purchasing"."purchaseorderheader"."status" IS 'Order current status. 1 = Pending; 2 = Approved; 3 = Rejected; 4 = Complete';

COMMENT ON COLUMN "purchasing"."purchaseorderheader"."employeeid" IS 'Employee who created the purchase order. Foreign key to Employee.BusinessEntityID.';

COMMENT ON COLUMN "purchasing"."purchaseorderheader"."vendorid" IS 'Vendor with whom the purchase order is placed. Foreign key to Vendor.BusinessEntityID.';

COMMENT ON COLUMN "purchasing"."purchaseorderheader"."shipmethodid" IS 'Shipping method. Foreign key to ShipMethod.ShipMethodID.';

COMMENT ON COLUMN "purchasing"."purchaseorderheader"."orderdate" IS 'Purchase order creation date.';

COMMENT ON COLUMN "purchasing"."purchaseorderheader"."shipdate" IS 'Estimated shipment date from the vendor.';

COMMENT ON COLUMN "purchasing"."purchaseorderheader"."subtotal" IS 'Purchase order subtotal. Computed as SUM(PurchaseOrderDetail.LineTotal)for the appropriate PurchaseOrderID.';

COMMENT ON COLUMN "purchasing"."purchaseorderheader"."taxamt" IS 'Tax amount.';

COMMENT ON COLUMN "purchasing"."purchaseorderheader"."freight" IS 'Shipping cost.';

COMMENT ON TABLE "purchasing"."productvendor" IS 'Cross-reference table mapping vendors with the products they supply.';

COMMENT ON COLUMN "purchasing"."productvendor"."productid" IS 'Primary key. Foreign key to Product.ProductID.';

COMMENT ON COLUMN "purchasing"."productvendor"."businessentityid" IS 'Primary key. Foreign key to Vendor.BusinessEntityID.';

COMMENT ON COLUMN "purchasing"."productvendor"."averageleadtime" IS 'The average span of time (in days) between placing an order with the vendor and receiving the purchased product.';

COMMENT ON COLUMN "purchasing"."productvendor"."standardprice" IS 'The vendor''s usual selling price.';

COMMENT ON COLUMN "purchasing"."productvendor"."lastreceiptcost" IS 'The selling price when last purchased.';

COMMENT ON COLUMN "purchasing"."productvendor"."lastreceiptdate" IS 'Date the product was last received by the vendor.';

COMMENT ON COLUMN "purchasing"."productvendor"."minorderqty" IS 'The maximum quantity that should be ordered.';

COMMENT ON COLUMN "purchasing"."productvendor"."maxorderqty" IS 'The minimum quantity that should be ordered.';

COMMENT ON COLUMN "purchasing"."productvendor"."onorderqty" IS 'The quantity currently on order.';

COMMENT ON COLUMN "purchasing"."productvendor"."unitmeasurecode" IS 'The product''s unit of measure.';

COMMENT ON TABLE "purchasing"."shipmethod" IS 'Shipping company lookup table.';

COMMENT ON COLUMN "purchasing"."shipmethod"."shipmethodid" IS 'Primary key for ShipMethod records.';

COMMENT ON COLUMN "purchasing"."shipmethod"."name" IS 'Shipping company name.';

COMMENT ON COLUMN "purchasing"."shipmethod"."shipbase" IS 'Minimum shipping charge.';

COMMENT ON COLUMN "purchasing"."shipmethod"."shiprate" IS 'Shipping charge per pound.';

COMMENT ON TABLE "purchasing"."vendor" IS 'Companies from whom Adventure Works Cycles purchases parts or other goods.';

COMMENT ON COLUMN "purchasing"."vendor"."businessentityid" IS 'Primary key for Vendor records.  Foreign key to BusinessEntity.BusinessEntityID';

COMMENT ON COLUMN "purchasing"."vendor"."accountnumber" IS 'Vendor account (identification) number.';

COMMENT ON COLUMN "purchasing"."vendor"."name" IS 'Company name.';

COMMENT ON COLUMN "purchasing"."vendor"."creditrating" IS '1 = Superior, 2 = Excellent, 3 = Above average, 4 = Average, 5 = Below average';

COMMENT ON COLUMN "purchasing"."vendor"."preferredvendorstatus" IS '0 = Do not use if another vendor is available. 1 = Preferred over other vendors supplying the same product.';

COMMENT ON COLUMN "purchasing"."vendor"."activeflag" IS '0 = Vendor no longer used. 1 = Vendor is actively used.';

COMMENT ON COLUMN "purchasing"."vendor"."purchasingwebserviceurl" IS 'Vendor URL.';

COMMENT ON TABLE "sales"."customer" IS 'Current customer information. Also see the Person and Store tables.';

COMMENT ON COLUMN "sales"."customer"."customerid" IS 'Primary key.';

COMMENT ON COLUMN "sales"."customer"."personid" IS 'Foreign key to Person.BusinessEntityID';

COMMENT ON COLUMN "sales"."customer"."storeid" IS 'Foreign key to Store.BusinessEntityID';

COMMENT ON COLUMN "sales"."customer"."territoryid" IS 'ID of the territory in which the customer is located. Foreign key to SalesTerritory.SalesTerritoryID.';

COMMENT ON TABLE "sales"."creditcard" IS 'Customer credit card information.';

COMMENT ON COLUMN "sales"."creditcard"."creditcardid" IS 'Primary key for CreditCard records.';

COMMENT ON COLUMN "sales"."creditcard"."cardtype" IS 'Credit card name.';

COMMENT ON COLUMN "sales"."creditcard"."cardnumber" IS 'Credit card number.';

COMMENT ON COLUMN "sales"."creditcard"."expmonth" IS 'Credit card expiration month.';

COMMENT ON COLUMN "sales"."creditcard"."expyear" IS 'Credit card expiration year.';

COMMENT ON TABLE "sales"."currencyrate" IS 'Currency exchange rates.';

COMMENT ON COLUMN "sales"."currencyrate"."currencyrateid" IS 'Primary key for CurrencyRate records.';

COMMENT ON COLUMN "sales"."currencyrate"."currencyratedate" IS 'Date and time the exchange rate was obtained.';

COMMENT ON COLUMN "sales"."currencyrate"."fromcurrencycode" IS 'Exchange rate was converted from this currency code.';

COMMENT ON COLUMN "sales"."currencyrate"."tocurrencycode" IS 'Exchange rate was converted to this currency code.';

COMMENT ON COLUMN "sales"."currencyrate"."averagerate" IS 'Average exchange rate for the day.';

COMMENT ON COLUMN "sales"."currencyrate"."endofdayrate" IS 'Final exchange rate for the day.';

COMMENT ON TABLE "sales"."countryregioncurrency" IS 'Cross-reference table mapping ISO currency codes to a country or region.';

COMMENT ON COLUMN "sales"."countryregioncurrency"."countryregioncode" IS 'ISO code for countries and regions. Foreign key to CountryRegion.CountryRegionCode.';

COMMENT ON COLUMN "sales"."countryregioncurrency"."currencycode" IS 'ISO standard currency code. Foreign key to Currency.CurrencyCode.';

COMMENT ON TABLE "sales"."currency" IS 'Lookup table containing standard ISO currencies.';

COMMENT ON COLUMN "sales"."currency"."currencycode" IS 'The ISO code for the Currency.';

COMMENT ON COLUMN "sales"."currency"."name" IS 'Currency name.';

COMMENT ON TABLE "sales"."personcreditcard" IS 'Cross-reference table mapping people to their credit card information in the CreditCard table.';

COMMENT ON COLUMN "sales"."personcreditcard"."businessentityid" IS 'Business entity identification number. Foreign key to Person.BusinessEntityID.';

COMMENT ON COLUMN "sales"."personcreditcard"."creditcardid" IS 'Credit card identification number. Foreign key to CreditCard.CreditCardID.';

COMMENT ON TABLE "sales"."store" IS 'Customers (resellers) of Adventure Works products.';

COMMENT ON COLUMN "sales"."store"."businessentityid" IS 'Primary key. Foreign key to Customer.BusinessEntityID.';

COMMENT ON COLUMN "sales"."store"."name" IS 'Name of the store.';

COMMENT ON COLUMN "sales"."store"."salespersonid" IS 'ID of the sales person assigned to the customer. Foreign key to SalesPerson.BusinessEntityID.';

COMMENT ON COLUMN "sales"."store"."demographics" IS 'Demographic informationg about the store such as the number of employees, annual sales and store type.';

COMMENT ON TABLE "sales"."shoppingcartitem" IS 'Contains online customer orders until the order is submitted or cancelled.';

COMMENT ON COLUMN "sales"."shoppingcartitem"."shoppingcartitemid" IS 'Primary key for ShoppingCartItem records.';

COMMENT ON COLUMN "sales"."shoppingcartitem"."shoppingcartid" IS 'Shopping cart identification number.';

COMMENT ON COLUMN "sales"."shoppingcartitem"."quantity" IS 'Product quantity ordered.';

COMMENT ON COLUMN "sales"."shoppingcartitem"."productid" IS 'Product ordered. Foreign key to Product.ProductID.';

COMMENT ON COLUMN "sales"."shoppingcartitem"."datecreated" IS 'Date the time the record was created.';

COMMENT ON TABLE "sales"."specialoffer" IS 'Sale discounts lookup table.';

COMMENT ON COLUMN "sales"."specialoffer"."specialofferid" IS 'Primary key for SpecialOffer records.';

COMMENT ON COLUMN "sales"."specialoffer"."description" IS 'Discount description.';

COMMENT ON COLUMN "sales"."specialoffer"."discountpct" IS 'Discount precentage.';

COMMENT ON COLUMN "sales"."specialoffer"."type" IS 'Discount type category.';

COMMENT ON COLUMN "sales"."specialoffer"."category" IS 'Group the discount applies to such as Reseller or Customer.';

COMMENT ON COLUMN "sales"."specialoffer"."startdate" IS 'Discount start date.';

COMMENT ON COLUMN "sales"."specialoffer"."enddate" IS 'Discount end date.';

COMMENT ON COLUMN "sales"."specialoffer"."minqty" IS 'Minimum discount percent allowed.';

COMMENT ON COLUMN "sales"."specialoffer"."maxqty" IS 'Maximum discount percent allowed.';

COMMENT ON TABLE "sales"."salesorderdetail" IS 'Individual products associated with a specific sales order. See SalesOrderHeader.';

COMMENT ON COLUMN "sales"."salesorderdetail"."salesorderid" IS 'Primary key. Foreign key to SalesOrderHeader.SalesOrderID.';

COMMENT ON COLUMN "sales"."salesorderdetail"."salesorderdetailid" IS 'Primary key. One incremental unique number per product sold.';

COMMENT ON COLUMN "sales"."salesorderdetail"."carriertrackingnumber" IS 'Shipment tracking number supplied by the shipper.';

COMMENT ON COLUMN "sales"."salesorderdetail"."orderqty" IS 'Quantity ordered per product.';

COMMENT ON COLUMN "sales"."salesorderdetail"."productid" IS 'Product sold to customer. Foreign key to Product.ProductID.';

COMMENT ON COLUMN "sales"."salesorderdetail"."specialofferid" IS 'Promotional code. Foreign key to SpecialOffer.SpecialOfferID.';

COMMENT ON COLUMN "sales"."salesorderdetail"."unitprice" IS 'Selling price of a single product.';

COMMENT ON COLUMN "sales"."salesorderdetail"."unitpricediscount" IS 'Discount amount.';

COMMENT ON TABLE "sales"."salesorderheader" IS 'General sales order information.';

COMMENT ON COLUMN "sales"."salesorderheader"."salesorderid" IS 'Primary key.';

COMMENT ON COLUMN "sales"."salesorderheader"."revisionnumber" IS 'Incremental number to track changes to the sales order over time.';

COMMENT ON COLUMN "sales"."salesorderheader"."orderdate" IS 'Dates the sales order was created.';

COMMENT ON COLUMN "sales"."salesorderheader"."duedate" IS 'Date the order is due to the customer.';

COMMENT ON COLUMN "sales"."salesorderheader"."shipdate" IS 'Date the order was shipped to the customer.';

COMMENT ON COLUMN "sales"."salesorderheader"."status" IS 'Order current status. 1 = In process; 2 = Approved; 3 = Backordered; 4 = Rejected; 5 = Shipped; 6 = Cancelled';

COMMENT ON COLUMN "sales"."salesorderheader"."onlineorderflag" IS '0 = Order placed by sales person. 1 = Order placed online by customer.';

COMMENT ON COLUMN "sales"."salesorderheader"."purchaseordernumber" IS 'Customer purchase order number reference.';

COMMENT ON COLUMN "sales"."salesorderheader"."accountnumber" IS 'Financial accounting number reference.';

COMMENT ON COLUMN "sales"."salesorderheader"."customerid" IS 'Customer identification number. Foreign key to Customer.BusinessEntityID.';

COMMENT ON COLUMN "sales"."salesorderheader"."salespersonid" IS 'Sales person who created the sales order. Foreign key to SalesPerson.BusinessEntityID.';

COMMENT ON COLUMN "sales"."salesorderheader"."territoryid" IS 'Territory in which the sale was made. Foreign key to SalesTerritory.SalesTerritoryID.';

COMMENT ON COLUMN "sales"."salesorderheader"."billtoaddressid" IS 'Customer billing address. Foreign key to Address.AddressID.';

COMMENT ON COLUMN "sales"."salesorderheader"."shiptoaddressid" IS 'Customer shipping address. Foreign key to Address.AddressID.';

COMMENT ON COLUMN "sales"."salesorderheader"."shipmethodid" IS 'Shipping method. Foreign key to ShipMethod.ShipMethodID.';

COMMENT ON COLUMN "sales"."salesorderheader"."creditcardid" IS 'Credit card identification number. Foreign key to CreditCard.CreditCardID.';

COMMENT ON COLUMN "sales"."salesorderheader"."creditcardapprovalcode" IS 'Approval code provided by the credit card company.';

COMMENT ON COLUMN "sales"."salesorderheader"."currencyrateid" IS 'Currency exchange rate used. Foreign key to CurrencyRate.CurrencyRateID.';

COMMENT ON COLUMN "sales"."salesorderheader"."subtotal" IS 'Sales subtotal. Computed as SUM(SalesOrderDetail.LineTotal)for the appropriate SalesOrderID.';

COMMENT ON COLUMN "sales"."salesorderheader"."taxamt" IS 'Tax amount.';

COMMENT ON COLUMN "sales"."salesorderheader"."freight" IS 'Shipping cost.';

COMMENT ON COLUMN "sales"."salesorderheader"."totaldue" IS 'Total due from customer. Computed as Subtotal + TaxAmt + Freight.';

COMMENT ON COLUMN "sales"."salesorderheader"."comment" IS 'Sales representative comments.';

COMMENT ON TABLE "sales"."salesorderheadersalesreason" IS 'Cross-reference table mapping sales orders to sales reason codes.';

COMMENT ON COLUMN "sales"."salesorderheadersalesreason"."salesorderid" IS 'Primary key. Foreign key to SalesOrderHeader.SalesOrderID.';

COMMENT ON COLUMN "sales"."salesorderheadersalesreason"."salesreasonid" IS 'Primary key. Foreign key to SalesReason.SalesReasonID.';

COMMENT ON TABLE "sales"."specialofferproduct" IS 'Cross-reference table mapping products to special offer discounts.';

COMMENT ON COLUMN "sales"."specialofferproduct"."specialofferid" IS 'Primary key for SpecialOfferProduct records.';

COMMENT ON COLUMN "sales"."specialofferproduct"."productid" IS 'Product identification number. Foreign key to Product.ProductID.';

COMMENT ON TABLE "sales"."salesperson" IS 'Sales representative current information.';

COMMENT ON COLUMN "sales"."salesperson"."businessentityid" IS 'Primary key for SalesPerson records. Foreign key to Employee.BusinessEntityID';

COMMENT ON COLUMN "sales"."salesperson"."territoryid" IS 'Territory currently assigned to. Foreign key to SalesTerritory.SalesTerritoryID.';

COMMENT ON COLUMN "sales"."salesperson"."salesquota" IS 'Projected yearly sales.';

COMMENT ON COLUMN "sales"."salesperson"."bonus" IS 'Bonus due if quota is met.';

COMMENT ON COLUMN "sales"."salesperson"."commissionpct" IS 'Commision percent received per sale.';

COMMENT ON COLUMN "sales"."salesperson"."salesytd" IS 'Sales total year to date.';

COMMENT ON COLUMN "sales"."salesperson"."saleslastyear" IS 'Sales total of previous year.';

COMMENT ON TABLE "sales"."salespersonquotahistory" IS 'Sales performance tracking.';

COMMENT ON COLUMN "sales"."salespersonquotahistory"."businessentityid" IS 'Sales person identification number. Foreign key to SalesPerson.BusinessEntityID.';

COMMENT ON COLUMN "sales"."salespersonquotahistory"."quotadate" IS 'Sales quota date.';

COMMENT ON COLUMN "sales"."salespersonquotahistory"."salesquota" IS 'Sales quota amount.';

COMMENT ON TABLE "sales"."salesreason" IS 'Lookup table of customer purchase reasons.';

COMMENT ON COLUMN "sales"."salesreason"."salesreasonid" IS 'Primary key for SalesReason records.';

COMMENT ON COLUMN "sales"."salesreason"."name" IS 'Sales reason description.';

COMMENT ON COLUMN "sales"."salesreason"."reasontype" IS 'Category the sales reason belongs to.';

COMMENT ON TABLE "sales"."salesterritory" IS 'Sales territory lookup table.';

COMMENT ON COLUMN "sales"."salesterritory"."territoryid" IS 'Primary key for SalesTerritory records.';

COMMENT ON COLUMN "sales"."salesterritory"."name" IS 'Sales territory description';

COMMENT ON COLUMN "sales"."salesterritory"."countryregioncode" IS 'ISO standard country or region code. Foreign key to CountryRegion.CountryRegionCode.';

COMMENT ON COLUMN "sales"."salesterritory"."group" IS 'Geographic area to which the sales territory belong.';

COMMENT ON COLUMN "sales"."salesterritory"."salesytd" IS 'Sales in the territory year to date.';

COMMENT ON COLUMN "sales"."salesterritory"."saleslastyear" IS 'Sales in the territory the previous year.';

COMMENT ON COLUMN "sales"."salesterritory"."costytd" IS 'Business costs in the territory year to date.';

COMMENT ON COLUMN "sales"."salesterritory"."costlastyear" IS 'Business costs in the territory the previous year.';

COMMENT ON TABLE "sales"."salesterritoryhistory" IS 'Sales representative transfers to other sales territories.';

COMMENT ON COLUMN "sales"."salesterritoryhistory"."businessentityid" IS 'Primary key. The sales rep.  Foreign key to SalesPerson.BusinessEntityID.';

COMMENT ON COLUMN "sales"."salesterritoryhistory"."territoryid" IS 'Primary key. Territory identification number. Foreign key to SalesTerritory.SalesTerritoryID.';

COMMENT ON COLUMN "sales"."salesterritoryhistory"."startdate" IS 'Primary key. Date the sales representive started work in the territory.';

COMMENT ON COLUMN "sales"."salesterritoryhistory"."enddate" IS 'Date the sales representative left work in the territory.';

COMMENT ON TABLE "sales"."salestaxrate" IS 'Tax rate lookup table.';

COMMENT ON COLUMN "sales"."salestaxrate"."salestaxrateid" IS 'Primary key for SalesTaxRate records.';

COMMENT ON COLUMN "sales"."salestaxrate"."stateprovinceid" IS 'State, province, or country/region the sales tax applies to.';

COMMENT ON COLUMN "sales"."salestaxrate"."taxtype" IS '1 = Tax applied to retail transactions, 2 = Tax applied to wholesale transactions, 3 = Tax applied to all sales (retail and wholesale) transactions.';

COMMENT ON COLUMN "sales"."salestaxrate"."taxrate" IS 'Tax rate amount.';

COMMENT ON COLUMN "sales"."salestaxrate"."name" IS 'Tax rate description.';

ALTER TABLE "humanresources"."employeedepartmenthistory" ADD CONSTRAINT "FK_EmployeeDepartmentHistory_Department_DepartmentID" FOREIGN KEY ("departmentid") REFERENCES "humanresources"."department" ("departmentid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "humanresources"."employeedepartmenthistory" ADD CONSTRAINT "FK_EmployeeDepartmentHistory_Employee_BusinessEntityID" FOREIGN KEY ("businessentityid") REFERENCES "humanresources"."employee" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "humanresources"."employeedepartmenthistory" ADD CONSTRAINT "FK_EmployeeDepartmentHistory_Shift_ShiftID" FOREIGN KEY ("shiftid") REFERENCES "humanresources"."shift" ("shiftid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "humanresources"."employeepayhistory" ADD CONSTRAINT "FK_EmployeePayHistory_Employee_BusinessEntityID" FOREIGN KEY ("businessentityid") REFERENCES "humanresources"."employee" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "humanresources"."employee" ADD CONSTRAINT "FK_Employee_Person_BusinessEntityID" FOREIGN KEY ("businessentityid") REFERENCES "person"."person" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "humanresources"."jobcandidate" ADD CONSTRAINT "FK_JobCandidate_Employee_BusinessEntityID" FOREIGN KEY ("businessentityid") REFERENCES "humanresources"."employee" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "person"."address" ADD CONSTRAINT "FK_Address_StateProvince_StateProvinceID" FOREIGN KEY ("stateprovinceid") REFERENCES "person"."stateprovince" ("stateprovinceid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "person"."businessentityaddress" ADD CONSTRAINT "FK_BusinessEntityAddress_AddressType_AddressTypeID" FOREIGN KEY ("addresstypeid") REFERENCES "person"."addresstype" ("addresstypeid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "person"."businessentityaddress" ADD CONSTRAINT "FK_BusinessEntityAddress_Address_AddressID" FOREIGN KEY ("addressid") REFERENCES "person"."address" ("addressid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "person"."businessentityaddress" ADD CONSTRAINT "FK_BusinessEntityAddress_BusinessEntity_BusinessEntityID" FOREIGN KEY ("businessentityid") REFERENCES "person"."businessentity" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "person"."businessentitycontact" ADD CONSTRAINT "FK_BusinessEntityContact_BusinessEntity_BusinessEntityID" FOREIGN KEY ("businessentityid") REFERENCES "person"."businessentity" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "person"."businessentitycontact" ADD CONSTRAINT "FK_BusinessEntityContact_ContactType_ContactTypeID" FOREIGN KEY ("contacttypeid") REFERENCES "person"."contacttype" ("contacttypeid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "person"."businessentitycontact" ADD CONSTRAINT "FK_BusinessEntityContact_Person_PersonID" FOREIGN KEY ("personid") REFERENCES "person"."person" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "person"."emailaddress" ADD CONSTRAINT "FK_EmailAddress_Person_BusinessEntityID" FOREIGN KEY ("businessentityid") REFERENCES "person"."person" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "person"."password" ADD CONSTRAINT "FK_Password_Person_BusinessEntityID" FOREIGN KEY ("businessentityid") REFERENCES "person"."person" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "person"."personphone" ADD CONSTRAINT "FK_PersonPhone_Person_BusinessEntityID" FOREIGN KEY ("businessentityid") REFERENCES "person"."person" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "person"."personphone" ADD CONSTRAINT "FK_PersonPhone_PhoneNumberType_PhoneNumberTypeID" FOREIGN KEY ("phonenumbertypeid") REFERENCES "person"."phonenumbertype" ("phonenumbertypeid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "person"."person" ADD CONSTRAINT "FK_Person_BusinessEntity_BusinessEntityID" FOREIGN KEY ("businessentityid") REFERENCES "person"."businessentity" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "person"."stateprovince" ADD CONSTRAINT "FK_StateProvince_CountryRegion_CountryRegionCode" FOREIGN KEY ("countryregioncode") REFERENCES "person"."countryregion" ("countryregioncode") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "person"."stateprovince" ADD CONSTRAINT "FK_StateProvince_SalesTerritory_TerritoryID" FOREIGN KEY ("territoryid") REFERENCES "sales"."salesterritory" ("territoryid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."billofmaterials" ADD CONSTRAINT "FK_BillOfMaterials_Product_ComponentID" FOREIGN KEY ("componentid") REFERENCES "production"."product" ("productid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."billofmaterials" ADD CONSTRAINT "FK_BillOfMaterials_Product_ProductAssemblyID" FOREIGN KEY ("productassemblyid") REFERENCES "production"."product" ("productid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."billofmaterials" ADD CONSTRAINT "FK_BillOfMaterials_UnitMeasure_UnitMeasureCode" FOREIGN KEY ("unitmeasurecode") REFERENCES "production"."unitmeasure" ("unitmeasurecode") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."document" ADD CONSTRAINT "FK_Document_Employee_Owner" FOREIGN KEY ("owner") REFERENCES "humanresources"."employee" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."productcosthistory" ADD CONSTRAINT "FK_ProductCostHistory_Product_ProductID" FOREIGN KEY ("productid") REFERENCES "production"."product" ("productid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."productdocument" ADD CONSTRAINT "FK_ProductDocument_Document_DocumentNode" FOREIGN KEY ("documentnode") REFERENCES "production"."document" ("documentnode") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."productdocument" ADD CONSTRAINT "FK_ProductDocument_Product_ProductID" FOREIGN KEY ("productid") REFERENCES "production"."product" ("productid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."productinventory" ADD CONSTRAINT "FK_ProductInventory_Location_LocationID" FOREIGN KEY ("locationid") REFERENCES "production"."location" ("locationid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."productinventory" ADD CONSTRAINT "FK_ProductInventory_Product_ProductID" FOREIGN KEY ("productid") REFERENCES "production"."product" ("productid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."productlistpricehistory" ADD CONSTRAINT "FK_ProductListPriceHistory_Product_ProductID" FOREIGN KEY ("productid") REFERENCES "production"."product" ("productid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."productmodelillustration" ADD CONSTRAINT "FK_ProductModelIllustration_Illustration_IllustrationID" FOREIGN KEY ("illustrationid") REFERENCES "production"."illustration" ("illustrationid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."productmodelillustration" ADD CONSTRAINT "FK_ProductModelIllustration_ProductModel_ProductModelID" FOREIGN KEY ("productmodelid") REFERENCES "production"."productmodel" ("productmodelid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."productmodelproductdescriptionculture" ADD CONSTRAINT "FK_ProductModelProductDescriptionCulture_Culture_CultureID" FOREIGN KEY ("cultureid") REFERENCES "production"."culture" ("cultureid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."productmodelproductdescriptionculture" ADD CONSTRAINT "FK_ProductModelProductDescriptionCulture_ProductDescription_Pro" FOREIGN KEY ("productdescriptionid") REFERENCES "production"."productdescription" ("productdescriptionid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."productmodelproductdescriptionculture" ADD CONSTRAINT "FK_ProductModelProductDescriptionCulture_ProductModel_ProductMo" FOREIGN KEY ("productmodelid") REFERENCES "production"."productmodel" ("productmodelid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."productproductphoto" ADD CONSTRAINT "FK_ProductProductPhoto_ProductPhoto_ProductPhotoID" FOREIGN KEY ("productphotoid") REFERENCES "production"."productphoto" ("productphotoid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."productproductphoto" ADD CONSTRAINT "FK_ProductProductPhoto_Product_ProductID" FOREIGN KEY ("productid") REFERENCES "production"."product" ("productid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."productreview" ADD CONSTRAINT "FK_ProductReview_Product_ProductID" FOREIGN KEY ("productid") REFERENCES "production"."product" ("productid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."productsubcategory" ADD CONSTRAINT "FK_ProductSubcategory_ProductCategory_ProductCategoryID" FOREIGN KEY ("productcategoryid") REFERENCES "production"."productcategory" ("productcategoryid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."product" ADD CONSTRAINT "FK_Product_ProductModel_ProductModelID" FOREIGN KEY ("productmodelid") REFERENCES "production"."productmodel" ("productmodelid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."product" ADD CONSTRAINT "FK_Product_ProductSubcategory_ProductSubcategoryID" FOREIGN KEY ("productsubcategoryid") REFERENCES "production"."productsubcategory" ("productsubcategoryid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."product" ADD CONSTRAINT "FK_Product_UnitMeasure_SizeUnitMeasureCode" FOREIGN KEY ("sizeunitmeasurecode") REFERENCES "production"."unitmeasure" ("unitmeasurecode") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."product" ADD CONSTRAINT "FK_Product_UnitMeasure_WeightUnitMeasureCode" FOREIGN KEY ("weightunitmeasurecode") REFERENCES "production"."unitmeasure" ("unitmeasurecode") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."transactionhistory" ADD CONSTRAINT "FK_TransactionHistory_Product_ProductID" FOREIGN KEY ("productid") REFERENCES "production"."product" ("productid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."workorderrouting" ADD CONSTRAINT "FK_WorkOrderRouting_Location_LocationID" FOREIGN KEY ("locationid") REFERENCES "production"."location" ("locationid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."workorderrouting" ADD CONSTRAINT "FK_WorkOrderRouting_WorkOrder_WorkOrderID" FOREIGN KEY ("workorderid") REFERENCES "production"."workorder" ("workorderid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."workorder" ADD CONSTRAINT "FK_WorkOrder_Product_ProductID" FOREIGN KEY ("productid") REFERENCES "production"."product" ("productid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "production"."workorder" ADD CONSTRAINT "FK_WorkOrder_ScrapReason_ScrapReasonID" FOREIGN KEY ("scrapreasonid") REFERENCES "production"."scrapreason" ("scrapreasonid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "purchasing"."productvendor" ADD CONSTRAINT "FK_ProductVendor_Product_ProductID" FOREIGN KEY ("productid") REFERENCES "production"."product" ("productid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "purchasing"."productvendor" ADD CONSTRAINT "FK_ProductVendor_UnitMeasure_UnitMeasureCode" FOREIGN KEY ("unitmeasurecode") REFERENCES "production"."unitmeasure" ("unitmeasurecode") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "purchasing"."productvendor" ADD CONSTRAINT "FK_ProductVendor_Vendor_BusinessEntityID" FOREIGN KEY ("businessentityid") REFERENCES "purchasing"."vendor" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "purchasing"."purchaseorderdetail" ADD CONSTRAINT "FK_PurchaseOrderDetail_Product_ProductID" FOREIGN KEY ("productid") REFERENCES "production"."product" ("productid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "purchasing"."purchaseorderdetail" ADD CONSTRAINT "FK_PurchaseOrderDetail_PurchaseOrderHeader_PurchaseOrderID" FOREIGN KEY ("purchaseorderid") REFERENCES "purchasing"."purchaseorderheader" ("purchaseorderid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "purchasing"."purchaseorderheader" ADD CONSTRAINT "FK_PurchaseOrderHeader_Employee_EmployeeID" FOREIGN KEY ("employeeid") REFERENCES "humanresources"."employee" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "purchasing"."purchaseorderheader" ADD CONSTRAINT "FK_PurchaseOrderHeader_ShipMethod_ShipMethodID" FOREIGN KEY ("shipmethodid") REFERENCES "purchasing"."shipmethod" ("shipmethodid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "purchasing"."purchaseorderheader" ADD CONSTRAINT "FK_PurchaseOrderHeader_Vendor_VendorID" FOREIGN KEY ("vendorid") REFERENCES "purchasing"."vendor" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "purchasing"."vendor" ADD CONSTRAINT "FK_Vendor_BusinessEntity_BusinessEntityID" FOREIGN KEY ("businessentityid") REFERENCES "person"."businessentity" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."countryregioncurrency" ADD CONSTRAINT "FK_CountryRegionCurrency_CountryRegion_CountryRegionCode" FOREIGN KEY ("countryregioncode") REFERENCES "person"."countryregion" ("countryregioncode") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."countryregioncurrency" ADD CONSTRAINT "FK_CountryRegionCurrency_Currency_CurrencyCode" FOREIGN KEY ("currencycode") REFERENCES "sales"."currency" ("currencycode") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."currencyrate" ADD CONSTRAINT "FK_CurrencyRate_Currency_FromCurrencyCode" FOREIGN KEY ("fromcurrencycode") REFERENCES "sales"."currency" ("currencycode") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."currencyrate" ADD CONSTRAINT "FK_CurrencyRate_Currency_ToCurrencyCode" FOREIGN KEY ("tocurrencycode") REFERENCES "sales"."currency" ("currencycode") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."customer" ADD CONSTRAINT "FK_Customer_Person_PersonID" FOREIGN KEY ("personid") REFERENCES "person"."person" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."customer" ADD CONSTRAINT "FK_Customer_SalesTerritory_TerritoryID" FOREIGN KEY ("territoryid") REFERENCES "sales"."salesterritory" ("territoryid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."customer" ADD CONSTRAINT "FK_Customer_Store_StoreID" FOREIGN KEY ("storeid") REFERENCES "sales"."store" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."personcreditcard" ADD CONSTRAINT "FK_PersonCreditCard_CreditCard_CreditCardID" FOREIGN KEY ("creditcardid") REFERENCES "sales"."creditcard" ("creditcardid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."personcreditcard" ADD CONSTRAINT "FK_PersonCreditCard_Person_BusinessEntityID" FOREIGN KEY ("businessentityid") REFERENCES "person"."person" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salesorderdetail" ADD CONSTRAINT "FK_SalesOrderDetail_SalesOrderHeader_SalesOrderID" FOREIGN KEY ("salesorderid") REFERENCES "sales"."salesorderheader" ("salesorderid") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salesorderdetail" ADD CONSTRAINT "FK_SalesOrderDetail_SpecialOfferProduct_SpecialOfferIDProductID" FOREIGN KEY ("specialofferid", "productid") REFERENCES "sales"."specialofferproduct" ("specialofferid", "productid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salesorderheadersalesreason" ADD CONSTRAINT "FK_SalesOrderHeaderSalesReason_SalesOrderHeader_SalesOrderID" FOREIGN KEY ("salesorderid") REFERENCES "sales"."salesorderheader" ("salesorderid") ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salesorderheadersalesreason" ADD CONSTRAINT "FK_SalesOrderHeaderSalesReason_SalesReason_SalesReasonID" FOREIGN KEY ("salesreasonid") REFERENCES "sales"."salesreason" ("salesreasonid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salesorderheader" ADD CONSTRAINT "FK_SalesOrderHeader_Address_BillToAddressID" FOREIGN KEY ("billtoaddressid") REFERENCES "person"."address" ("addressid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salesorderheader" ADD CONSTRAINT "FK_SalesOrderHeader_Address_ShipToAddressID" FOREIGN KEY ("shiptoaddressid") REFERENCES "person"."address" ("addressid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salesorderheader" ADD CONSTRAINT "FK_SalesOrderHeader_CreditCard_CreditCardID" FOREIGN KEY ("creditcardid") REFERENCES "sales"."creditcard" ("creditcardid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salesorderheader" ADD CONSTRAINT "FK_SalesOrderHeader_CurrencyRate_CurrencyRateID" FOREIGN KEY ("currencyrateid") REFERENCES "sales"."currencyrate" ("currencyrateid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salesorderheader" ADD CONSTRAINT "FK_SalesOrderHeader_Customer_CustomerID" FOREIGN KEY ("customerid") REFERENCES "sales"."customer" ("customerid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salesorderheader" ADD CONSTRAINT "FK_SalesOrderHeader_SalesPerson_SalesPersonID" FOREIGN KEY ("salespersonid") REFERENCES "sales"."salesperson" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salesorderheader" ADD CONSTRAINT "FK_SalesOrderHeader_SalesTerritory_TerritoryID" FOREIGN KEY ("territoryid") REFERENCES "sales"."salesterritory" ("territoryid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salesorderheader" ADD CONSTRAINT "FK_SalesOrderHeader_ShipMethod_ShipMethodID" FOREIGN KEY ("shipmethodid") REFERENCES "purchasing"."shipmethod" ("shipmethodid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salespersonquotahistory" ADD CONSTRAINT "FK_SalesPersonQuotaHistory_SalesPerson_BusinessEntityID" FOREIGN KEY ("businessentityid") REFERENCES "sales"."salesperson" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salesperson" ADD CONSTRAINT "FK_SalesPerson_Employee_BusinessEntityID" FOREIGN KEY ("businessentityid") REFERENCES "humanresources"."employee" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salesperson" ADD CONSTRAINT "FK_SalesPerson_SalesTerritory_TerritoryID" FOREIGN KEY ("territoryid") REFERENCES "sales"."salesterritory" ("territoryid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salestaxrate" ADD CONSTRAINT "FK_SalesTaxRate_StateProvince_StateProvinceID" FOREIGN KEY ("stateprovinceid") REFERENCES "person"."stateprovince" ("stateprovinceid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salesterritoryhistory" ADD CONSTRAINT "FK_SalesTerritoryHistory_SalesPerson_BusinessEntityID" FOREIGN KEY ("businessentityid") REFERENCES "sales"."salesperson" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salesterritoryhistory" ADD CONSTRAINT "FK_SalesTerritoryHistory_SalesTerritory_TerritoryID" FOREIGN KEY ("territoryid") REFERENCES "sales"."salesterritory" ("territoryid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."salesterritory" ADD CONSTRAINT "FK_SalesTerritory_CountryRegion_CountryRegionCode" FOREIGN KEY ("countryregioncode") REFERENCES "person"."countryregion" ("countryregioncode") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."shoppingcartitem" ADD CONSTRAINT "FK_ShoppingCartItem_Product_ProductID" FOREIGN KEY ("productid") REFERENCES "production"."product" ("productid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."specialofferproduct" ADD CONSTRAINT "FK_SpecialOfferProduct_Product_ProductID" FOREIGN KEY ("productid") REFERENCES "production"."product" ("productid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."specialofferproduct" ADD CONSTRAINT "FK_SpecialOfferProduct_SpecialOffer_SpecialOfferID" FOREIGN KEY ("specialofferid") REFERENCES "sales"."specialoffer" ("specialofferid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."store" ADD CONSTRAINT "FK_Store_BusinessEntity_BusinessEntityID" FOREIGN KEY ("businessentityid") REFERENCES "person"."businessentity" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "sales"."store" ADD CONSTRAINT "FK_Store_SalesPerson_SalesPersonID" FOREIGN KEY ("salespersonid") REFERENCES "sales"."salesperson" ("businessentityid") DEFERRABLE INITIALLY IMMEDIATE;
