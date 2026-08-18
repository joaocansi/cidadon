package councillor

import "gorm.io/gorm"

type Model struct {
	gorm.Model
	Party string
}

func (*Model) TableName() string {
	return "councillors"
}

func (m *Model) ToDomain() *Councillor {
	return &Councillor{
		UpdatedAt: m.UpdatedAt,
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		Party:     m.Party,
	}
}
