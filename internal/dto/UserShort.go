// Package dto - пакет определяющий формы представления сущностей и
// структуры ответов/запросов
package dto

// UserShort является формой представления сущности User
// с уникальным идентификатором и флагом активности
type UserShort struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}
