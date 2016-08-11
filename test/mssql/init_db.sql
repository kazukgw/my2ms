USE [test]
GO

/****** Object:  Table [dbo].[users]    Script Date: 2016/08/11 13:17:41 ******/
SET ANSI_NULLS ON
GO

SET QUOTED_IDENTIFIER ON
GO

CREATE TABLE [dbo].[users](
	[user_id] [nvarchar](255) NOT NULL,
	[code] [int] NULL DEFAULT (NULL),
	[name] [nvarchar](255) NULL DEFAULT (NULL),
	[mail] [nvarchar](255) NULL DEFAULT (NULL),
 CONSTRAINT [user_id] PRIMARY KEY CLUSTERED
(
	[user_id] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]

GO

