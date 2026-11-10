package healthcheck

import (
	"context"

	"github.com/ThienKim52/golang-dev/internal/app/repository/link"
)

//go:generate mockery --name=HealthCheckService --filename=healthcheckservice.go --outpkg=mocks_genpass
type HealthCheckService interface {
	GetHealthCheck(ctx context.Context) (HealthCheckResult, error)
}

// healthCheckService implements HealthCheckService
type healthCheckService struct {
	repo        repository.LinkRepository
	serviceName string
	instanceID  string
}

type HealthCheckResult struct {
	Message     string `json:"message"`
	ServiceName string `json:"service_name"`
	InstanceID  string `json:"instance_id"`
}

// NewHealthCheck creates a new health check service
func NewHealthCheck(repo repository.LinkRepository, serviceName, instanceID string) HealthCheckService {
	return &healthCheckService{
		repo:        repo,
		serviceName: serviceName,
		instanceID:  instanceID,
	}
}

// GetHealthCheck returns the health check response
func (s *healthCheckService) GetHealthCheck(ctx context.Context) (HealthCheckResult, error) {
	// Check Redis connection if repository is available
	result := HealthCheckResult{
		Message:     "OK",
		ServiceName: s.serviceName,
		InstanceID:  s.instanceID,
	}
	//if ping repository get error, return error
	if s.repo != nil {
		if err := s.repo.Ping(ctx); err != nil {
			return HealthCheckResult{}, err
		}
	}

	return result, nil
}
