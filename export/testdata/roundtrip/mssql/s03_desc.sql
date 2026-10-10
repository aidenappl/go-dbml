CREATE TABLE [HR].[Employees] (
    [Id] INT NOT NULL IDENTITY(100, 1),
    [ManagerId] INT NULL,
    [FullName] NVARCHAR (200) NOT NULL,
    [HireDate] DATE NOT NULL,
    CONSTRAINT [PK_Employees] PRIMARY KEY ([Id]),
    CONSTRAINT [FK_Employees_Manager] FOREIGN KEY ([ManagerId]) REFERENCES [HR].[Employees] ([Id])
);
GO
EXEC sp_addextendedproperty @name = N'MS_Description', @value = N'Staff list', @level0type = N'SCHEMA', @level0name = N'HR', @level1type = N'TABLE', @level1name = N'Employees';
GO
EXEC sp_addextendedproperty @name = N'MS_Description', @value = N'Legal name', @level0type = N'SCHEMA', @level0name = N'HR', @level1type = N'TABLE', @level1name = N'Employees', @level2type = N'COLUMN', @level2name = N'FullName';
GO
CREATE NONCLUSTERED INDEX [IX_Employees_Manager] ON [HR].[Employees] ([ManagerId] ASC, [HireDate] DESC);
GO
