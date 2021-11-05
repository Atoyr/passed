USE [passed]
GO

SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[Profiles](
	[ID] [uniqueidentifier] NOT NULL,
	[Email] [nvarchar](128) NOT NULL,
	[FirstName] [nvarchar](64) NOT NULL,
	[MiddleName] [nvarchar](64) NOT NULL,
	[LastName] [nvarchar](64) NOT NULL,
	[Nickname] [nvarchar](32) NOT NULL,
	[ValidFlg] [bit] NOT NULL,
	[InsertDatetime] [datetime] NOT NULL,
	[ModifiedDatetime] [datetime] NOT NULL,
	[InsertAccountID] [uniqueidentifier] NOT NULL,
	[InsertSystemID] [uniqueidentifier] NOT NULL,
	[ModifiedAccountID] [uniqueidentifier] NOT NULL,
	[ModifiedSystemID] [uniqueidentifier] NOT NULL
 CONSTRAINT [PK_Profile] PRIMARY KEY CLUSTERED
(
	[ID], [Email] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON, FILLFACTOR = 90, OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF) ON [PRIMARY]
) ON [PRIMARY]
GO

ALTER TABLE [dbo].[Profiles] ADD  CONSTRAINT [DF__Profile__ID]  DEFAULT NewID() FOR [ID]
GO

ALTER TABLE [dbo].[Profiles] ADD  CONSTRAINT [DF__Profile__ValidFlg]  DEFAULT 1 FOR [ValidFlg]
GO

ALTER TABLE [dbo].[Profiles] ADD  CONSTRAINT [DF__Profile__InsertDatetime]  DEFAULT (getdate()) FOR [InsertDatetime]
GO

ALTER TABLE [dbo].[Profiles] ADD  CONSTRAINT [DF__Profile__ModifiedDatetime]  DEFAULT (getdate()) FOR [ModifiedDatetime]
GO