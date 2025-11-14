package pull_request

var (
	writableColumns = []string{
		"pr_id",
		"pr_name",
		"author_id",
		"status",
	}

	readableColumns = []string{
		"pr_id",
		"pr_name",
		"author_id",
		"status",
		"created_at",
		"merged_at",
	}
)
