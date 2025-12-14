package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// RequestLogger — это middleware, который инжектит логгер в контекст
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 1. Получаем или генерируем Request ID
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}

		// 2. Создаем sub-logger с контекстом запроса
		// Мы берем глобальный логгер и добавляем к нему поля
		l := log.With().
			Str("req_id", reqID).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Logger()

		// 3. Создаем новый контекст с внедренным логгером
		ctx := l.WithContext(r.Context())

		// 4. Передаем управление дальше, подменяя контекст в запросе
		next.ServeHTTP(w, r.WithContext(ctx))

		// 5. Логируем завершение запроса (Access Log)
		// Обрати внимание: мы берем логгер УЖЕ из контекста
		// zerolog.Ctx(ctx) вернет тот самый логгер с req_id
		zerolog.Ctx(ctx).Info().
			Dur("duration", time.Since(start)).
			Msg("request completed")
	})
}
