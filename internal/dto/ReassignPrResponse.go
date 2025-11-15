package dto

// ReassignPrResponse представляет структуру ответа на
// запрос переназначения сотрудника с информацией о PR
// и идентификатором вновь назначенного пользователя
type ReassignPrResponse struct {
	Pr         PullRequest `json:"pr"`
	ReplacedBy string      `json:"replaced_by"`
}
