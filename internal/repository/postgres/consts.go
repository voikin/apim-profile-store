package postgres

const (
	ApplicationTableName = "applications"

	ApplicationColumnID        = "id"
	ApplicationColumnName      = "name"
	ApplicationColumnCreatedAt = "created_at"

	ApplicationProfileTableName = "application_profiles"

	ApplicationProfileColumnID            = "id"
	ApplicationProfileColumnApplicationID = "application_id"
	ApplicationProfileColumnVersion       = "version"
	ApplicationProfileColumnGraphID       = "graph_id"
	ApplicationProfileColumnCreatedAt     = "created_at"

	ApplicationProfileVersionsTableName = "application_profile_versions"

	ApplicationProfileVersionsColumnApplicationID = "application_id"
	ApplicationProfileVersionsColumnLastVersion   = "last_version"
	ApplicationProfileVersionsColumnUpdatedAt     = "updated_at"
)

var (
	ApplicationAllColumns = []string{ //nolint:gochecknoglobals // global by design
		ApplicationColumnID,
		ApplicationColumnName,
		ApplicationColumnCreatedAt,
	}

	ApplicationProfileAllColumns = []string{ //nolint:gochecknoglobals // global by design
		ApplicationProfileColumnID,
		ApplicationProfileColumnApplicationID,
		ApplicationProfileColumnVersion,
		ApplicationProfileColumnGraphID,
		ApplicationProfileColumnCreatedAt,
	}

	ApplicationProfileVersionsAllColumns = []string{ //nolint:gochecknoglobals // global by design
		ApplicationProfileVersionsColumnApplicationID,
		ApplicationProfileVersionsColumnLastVersion,
		ApplicationProfileVersionsColumnUpdatedAt,
	}
)

const (
	pgUniqueViolationCode = "23505"
)
