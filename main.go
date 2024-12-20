package main

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		url := "https://gos.1cfresh.com/a/shrm/3041589/odata/standart.odata/Catalog_Номенклатура?$format=json"

		// Создаем запрос
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			slog.Error("Failed to create request", "error", err)
			render.JSON(w, r, map[string]string{"error": "Failed to create request"})
			return
		}

		req.Header.Set("Authorization", "Bearer ")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			slog.Error("Request failed", "error", err)
			render.JSON(w, r, map[string]string{"error": "Request failed"})
			return
		}
		defer resp.Body.Close()

		// Проверяем статус ответа
		if resp.StatusCode != http.StatusOK {
			slog.Error("Non-200 status code", "status", resp.StatusCode)
			render.JSON(w, r, map[string]string{
				"error":       "Unexpected status code",
				"status_code": http.StatusText(resp.StatusCode),
			})
			return
		}

		// Читаем тело ответа
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Failed to read response body", "error", err)
			render.JSON(w, r, map[string]string{"error": "Failed to read response body"})
			return
		}

		// Отправляем данные в браузер
		render.Data(w, r, body)
	})

	// Запускаем сервер
	slog.Info("Server is running on http://localhost:5000")
	http.ListenAndServe(":5000", r)
}
