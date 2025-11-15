package dto

// ErrorCode представляет тип,
// определяющий возможную ошибку при
// обработке запроса
type ErrorCode string

// Константы, определяющие тип ошибки
const (
	TeamExists  ErrorCode = "TEAM_EXISTS"
	PrExists    ErrorCode = "PR_EXISTS"
	PrMerged    ErrorCode = "PR_MERGED"
	NotAssigned ErrorCode = "NOT_ASSIGNED"
	NoCandidate ErrorCode = "NO_CANDIDATE"
	NotFound    ErrorCode = "NOT_FOUND"
)

// ErrorResponse определяет структуру ответа
// на запрос, обработка которого завершилась неудачно
type ErrorResponse struct {
	Status int               `json:"-"`
	Error  map[string]string `json:"error"` // code-message
}
