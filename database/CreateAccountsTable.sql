USE [passed]
GO

SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[Accounts](
	[ID] [uniqueidentifier] NOT NULL,
	[ProfileID] [uniqueidentifier] NOT NULL,
	[Signature] [varbinary](4192) NOT NULL,
	[Private] [varbinary](4192) NOT NULL,
	[Public] [varbinary](4192) NOT NULL,
	[ValidFlg] [bit] NOT NULL,
	[InsertDatetime] [datetime] NOT NULL,
	[ModifiedDatetime] [datetime] NOT NULL,
	[InsertAccountID] [uniqueidentifier] NOT NULL,
	[InsertSystemID] [uniqueidentifier] NOT NULL,
	[ModifiedAccountID] [uniqueidentifier] NOT NULL,
	[ModifiedSystemID] [uniqueidentifier] NOT NULL
 CONSTRAINT [PK_Account] PRIMARY KEY CLUSTERED
(
	[ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, FILLFACTOR = 90, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO

ALTER TABLE [dbo].[Accounts] ADD CONSTRAINT [DF__Account__ID]  DEFAULT (NewID()) FOR [ID]
GO

ALTER TABLE [dbo].[Accounts] ADD CONSTRAINT [DF__Account__InsertDatetime]  DEFAULT (getdate()) FOR [InsertDatetime]
GO

ALTER TABLE [dbo].[Accounts] ADD CONSTRAINT [DF__Account__ModifiedDatetime]  DEFAULT (getdate()) FOR [ModifiedDatetime]
GO