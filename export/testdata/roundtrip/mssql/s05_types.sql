CREATE TABLE [dbo].[Docs] (
    [DocId] [bigint] IDENTITY(1,1) NOT NULL,
    [Title] [nvarchar](400) NOT NULL,
    [Content] [varbinary](max) NULL,
    [Xml] [xml] NULL,
    [Version] [rowversion] NOT NULL,
    [Score] [float] NULL,
    [Ratio] [real] NULL,
    [Small] [smallmoney] NULL,
    [Created] [smalldatetime] NOT NULL CONSTRAINT [DF_Docs_Created] DEFAULT (getdate()),
    [Flag] [char](1) NOT NULL DEFAULT ('N'),
    CONSTRAINT [PK_Docs] PRIMARY KEY NONCLUSTERED ([DocId])
)
GO
CREATE CLUSTERED INDEX [CIX_Docs_Created] ON [dbo].[Docs] ([Created])
GO
