package healthcheck

import "gorm.io/gorm"

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

type ReadyResponse struct {
	Ok       bool   `json:"ok"`
	Database string `json:"database"`
	Message string `json:"message"`
}

func (s *Service) Ready() ReadyResponse {
	pql, err := s.db.DB()

	if err != nil {
		return ReadyResponse{
			Ok:       false,
			Database: "down",
			Message: err.Error(),
		}
	}

	err = pql.Ping()

	if err != nil {
		return  ReadyResponse{
			Ok: false,
			Database: "dowm",
			Message: err.Error(),
		}
	}

	return ReadyResponse{
		Ok: true,
		Database: "up",
		Message: "ready",
	}
}
