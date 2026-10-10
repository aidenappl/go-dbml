CREATE TABLE [dbo].[Acct] (
  [Id] INT IDENTITY(1,1) NOT NULL,
  [Na]]me] NVARCHAR(50) NOT NULL CONSTRAINT [DF_name] DEFAULT (N'O''Brien'),
  [Kind] VARCHAR(10) NOT NULL CHECK ([Kind] IN ('a''b', 'c')),
  [Total] AS ([Id] * 2) PERSISTED,
  [Parent] INT NULL,
  CONSTRAINT [PK_Acct] PRIMARY KEY CLUSTERED ([Id] ASC),
  CONSTRAINT [FK_Parent] FOREIGN KEY ([Parent]) REFERENCES [dbo].[Acct] ([Id])
)
GO
EXEC sys.sp_addextendedproperty @name=N'MS_Description', @value=N'it''s a table' , @level0type=N'SCHEMA',@level0name=N'dbo', @level1type=N'TABLE',@level1name=N'Acct'
GO
