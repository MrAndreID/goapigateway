package types

var Patch = []map[string]any{
	{
		"endpoint": "/api/v1/user/{id}",
		"backend": map[string]string{
			"host": url["goapi"],
			"path": "/api/v1/user/{id}",
		},
	},
}
