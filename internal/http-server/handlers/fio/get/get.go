package get

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/render"

	"github.com/go-chi/chi/middleware"
	http_server "github.com/kldd0/fio-service/internal/http-server"
	"github.com/kldd0/fio-service/internal/logs"
	"github.com/kldd0/fio-service/internal/model/domain_models"
	"github.com/kldd0/fio-service/internal/storage"
	"go.uber.org/zap"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	http_server.Response
	People []domain_models.FioStruct
}

type PeopleGetter interface {
	Get(ctx context.Context, filter string, target interface{}, limit, offset int) ([]domain_models.FioStruct, error)
}

func New(log *zap.Logger, peopleGetter PeopleGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.fio.get.New"

		logs.Logger.With(
			zap.String("op", op),
			zap.String("request_id", middleware.GetReqID(r.Context())),
		)

		filter := r.URL.Query().Get("filter")
		eq := r.URL.Query().Get("eq")
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = 10
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			offset = 0
		}

		result, err := peopleGetter.Get(r.Context(), filter, eq, limit, offset)
		if errors.Is(err, storage.ErrEntryDoesntExist) {
			log.Info(
				"entry doesn't exist",
				zap.String("filter", filter),
				zap.String("target", eq),
				zap.Int("limit", limit),
				zap.Int("offset", offset),
			)

			render.JSON(w, r, http_server.Error("entry doesn't exist"))

			return
		}

		fmt.Println(result)

		responseOK(w, r, result)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, data []domain_models.FioStruct) {
	render.JSON(w, r, Response{
		Response: http_server.OK(),
		People:   data,
	})
}
