package pull_request

var (
	writableColumns = []string{
		"pr_id",
		"pr_name",
		"author_id",
		"status_id",
	}

	readableColumns = []string{
		"pr_id",
		"pr_name",
		"author_id",
		"status_id",
		"created_at",
		"merged_at",
	}
)
