SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[Profiles](
	[ID] [uniqueidentifier] NOT NULL,
	[AccountID] [uniqueidentifier] NOT NULL,
	[FirstName] [nvarchar](64) NOT NULL,
	[MiddleName] [nvarchar](64) NOT NULL,
	[LastName] [nvarchar](64) NOT NULL,
	[Nickname] [nvarchar](32) NOT NULL,
	[ValidFlg] [bit] NOT NULL,
	[InsertAt] [datetime] NOT NULL,
	[UpdateAt] [datetime] NOT NULL,
	[InsertAccountID] [uniqueidentifier] NOT NULL,
	[InsertSystemID] [uniqueidentifier] NOT NULL,
	[UpdateAccountID] [uniqueidentifier] NOT NULL,
	[UpdateSystemID] [uniqueidentifier] NOT NULL,
	[ValidFrom] [datetime2] GENERATED ALWAYS AS ROW START,
	[ValidTo] [datetime2] GENERATED ALWAYS AS ROW END,
	PERIOD FOR SYSTEM_TIME (ValidFrom, ValidTo),

	CONSTRAINT [PK__Profiles] PRIMARY KEY NONCLUSTERED
	(
		[ID] ASC
	)
	WITH
	(
		PAD_INDEX = OFF,
		 STATISTICS_NORECOMPUTE = OFF,
		IGNORE_DUP_KEY = OFF,
		ALLOW_ROW_LOCKS = ON,
		ALLOW_PAGE_LOCKS = ON,
		FILLFACTOR = 90,
		OPTIMIZE_FOR_SEQUENTIAL_KEY = OFF
	) ON [PRIMARY]
) ON [PRIMARY]
WITH
(
		SYSTEM_VERSIONING = ON
)
GO

ALTER TABLE [dbo].[Profiles] ADD  CONSTRAINT [DF__Profiles__ID]  DEFAULT (NewID()) FOR [ID]
GO

ALTER TABLE [dbo].[Profiles] ADD  CONSTRAINT [DF__Profiles__ValidFlg]  DEFAULT 1 FOR [ValidFlg]
GO

ALTER TABLE [dbo].[Profiles] ADD  CONSTRAINT [DF__Profiles__InsertAt]  DEFAULT (getdate()) FOR [InsertAt]
GO

ALTER TABLE [dbo].[Profiles] ADD  CONSTRAINT [DF__Profiles__UpdateAt]  DEFAULT (getdate()) FOR [UpdateAt]
GO