package bootstrap

import "github.com/Agmer17/srd_lab_creative/internal/auth"

type ServiceConfigs struct {
	AuthService *auth.AuthService
}

func NewServiceConfigs(googleClientId string, googleSecret string, rpf *RepositoryConfigs) *ServiceConfigs {
	authService := auth.NewAuthService(googleClientId, googleSecret, rpf.AuthRepository)
	return &ServiceConfigs{
		AuthService: authService,
	}
}
