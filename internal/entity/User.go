// Package entity - пакет, определяющий сущности предметной области
package entity

// User представляет сущность пользователя -
// участника команды с уникальным идентификатором,
// именем и флагом активности
type User struct {
	UserID   string
	Username string
	TeamName string
	IsActive bool
}
