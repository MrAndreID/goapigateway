package types

var Put = []map[string]any{
	{
		"endpoint": "/api/v1/user/{id}",
		"backend": map[string]string{
			"method": "PATCH",
			"host":   url["goapi"],
			"path":   "/api/v1/user/{id}",
		},
	},
}
