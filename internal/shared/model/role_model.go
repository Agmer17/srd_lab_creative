package model

import (
	"time"

	"github.com/Agmer17/srd_lab_creative/internal/db/sqlcgen"
	"github.com/google/uuid"
)

type ProjectRole struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
}

func MapToProjectRoleModel(p sqlcgen.Role) ProjectRole {
	return ProjectRole{
		Id:        p.ID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
	}
}

func GenListToRoleModel(pl []sqlcgen.Role) []ProjectRole {

	var res []ProjectRole = make([]ProjectRole, len(pl))
	for i, v := range pl {
		res[i] = MapToProjectRoleModel(v)
	}

	return res
}
