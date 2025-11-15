package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/salex06/pr-service/internal/service"
)

// StatsHandler представляет контроллер, который
// отвечает за получение и отправку ответов на запросы
// статистики по приложению
type StatsHandler struct {
	statsService *service.StatsService
}

// NewStatHandler конструирует и возвращает объект StatsHandler
func NewStatHandler(svc *service.StatsService) *StatsHandler {
	return &StatsHandler{
		statsService: svc,
	}
}

// HandleGetStatsRequest получает запрос на сбор статистики работы приложения
// и формирует ответ
func (sh *StatsHandler) HandleGetStatsRequest(c *gin.Context) {
	stat, err := sh.statsService.GetStat()
	if err != nil {
		c.JSON(err.Status, err.Error)
		return
	}

	c.JSON(http.StatusOK, stat)
}
