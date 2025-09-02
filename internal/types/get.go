package types

var Get = []map[string]any{
	{
		"endpoint": "/api/v1/user",
		"backend": map[string]string{
			"host": url["goapi"],
			"path": "/api/v1/user",
		},
	},
}
