SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[Accounts](
	[ID] [uniqueidentifier] NOT NULL,
	[Email] [nvarchar](128) NOT NULL,
	[Signature] [varbinary](4192) NOT NULL,
	[Private] [varbinary](4192) NOT NULL,
	[Public] [varbinary](4192) NOT NULL,
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

	CONSTRAINT [PK__Accounts] PRIMARY KEY NONCLUSTERED
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
	) ON [PRIMARY],

	INDEX [IX__Accounts__Cluster] CLUSTERED
	(
		[ID],
		[Email],
		[ValidFlg]
	)
) ON [PRIMARY]
WITH
(
		SYSTEM_VERSIONING = ON
)
GO

ALTER TABLE [dbo].[Accounts] ADD  CONSTRAINT [DF__Accounts__ID]  DEFAULT (NewID()) FOR [ID]
GO

ALTER TABLE [dbo].[Accounts] ADD  CONSTRAINT [DF__Accounts__ValidFlg]  DEFAULT 1 FOR [ValidFlg]
GO

ALTER TABLE [dbo].[Accounts] ADD CONSTRAINT [DF__Accounts__InsertAt]  DEFAULT (getdate()) FOR [InsertAt]
GO

ALTER TABLE [dbo].[Accounts] ADD CONSTRAINT [DF__Accounts__UpdateAt]  DEFAULT (getdate()) FOR [UpdateAt]
GO