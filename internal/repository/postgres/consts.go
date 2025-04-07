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
)
