package config

// GetDefaultAuthConfig 返回默认的AuthConfig
func GetDefaultAuthConfig() AuthConfig {
	return AuthConfig{
		ExcludedPaths: []string{
			"/api/user/login",
			"/api/ping",
			"/api/metrics",
			"/api/server_info",
		},
	}
}
