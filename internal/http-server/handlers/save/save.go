package save

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"

	"github.com/S-a-b-r/url-shortener/internal/lib/api/response"
	"github.com/S-a-b-r/url-shortener/internal/lib/random"
	"github.com/S-a-b-r/url-shortener/internal/storage"
)

const aliasLength = 6

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", slog.String("error", err.Error()))

			render.JSON(w, r, response.Error("failed to decode package"))

			return
		}

		log.Info("saving url", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			log.Error("failed to validate request", slog.String("error", err.Error()))

			render.JSON(w, r, response.ValidationError(err.(validator.ValidationErrors)))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, response.Error("url already exists"))

			return
		}

		log.Info("url added", slog.String("url", req.URL), slog.Int64("id", id))

		render.JSON(w, r, response.Ok(alias))
	}
}
