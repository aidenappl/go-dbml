USE [Shop]
GO
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Customers](
	[CustomerID] [int] IDENTITY(1,1) NOT NULL,
	[Name] [nvarchar](100) NOT NULL,
	[Email] [nvarchar](max) NULL,
	[Balance] [decimal](18, 2) NOT NULL CONSTRAINT [DF_Customers_Balance] DEFAULT ((0)),
	[Created] [datetime2](7) NOT NULL DEFAULT (getdate()),
	[Flag] [bit] NULL,
 CONSTRAINT [PK_Customers] PRIMARY KEY CLUSTERED
(
	[CustomerID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF) ON [PRIMARY]
) ON [PRIMARY] TEXTIMAGE_ON [PRIMARY]
GO
CREATE TABLE [dbo].[Orders](
	[OrderID] [int] IDENTITY(1,1) NOT NULL,
	[CustomerID] [int] NOT NULL,
	[Total] [money] NULL,
	[Status] [varchar](20) NOT NULL,
 CONSTRAINT [PK_Orders] PRIMARY KEY NONCLUSTERED ([OrderID] ASC)
)
GO
ALTER TABLE [dbo].[Orders]  WITH CHECK ADD  CONSTRAINT [FK_Orders_Customers] FOREIGN KEY([CustomerID])
REFERENCES [dbo].[Customers] ([CustomerID])
ON DELETE CASCADE
GO
ALTER TABLE [dbo].[Orders] CHECK CONSTRAINT [FK_Orders_Customers]
GO
ALTER TABLE [dbo].[Orders] ADD CONSTRAINT [DF_Orders_Status] DEFAULT ('new') FOR [Status]
GO
CREATE NONCLUSTERED INDEX [IX_Orders_Customer] ON [dbo].[Orders] ([CustomerID] ASC)
GO
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'Customer master' , @level0type=N'SCHEMA',@level0name=N'dbo', @level1type=N'TABLE',@level1name=N'Customers'
GO
