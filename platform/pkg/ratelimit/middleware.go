// Package ratelimit содержит HTTP middleware для распределённого rate limiting.
//
// Библиотека redis_rate реализует GCRA (Generic Cell Rate Algorithm) —
// вариант token bucket, атомарно исполняемый в Redis через Lua-скрипт.
// Благодаря Redis лимит общий для всех инстансов сервиса, а не per-instance.
package ratelimit

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// meterName — имя Meter'а в OTel (отображается как otel_scope_name в Prometheus)
const meterName = "github.com/krapagen/my_microservices_rocket/platform/pkg/ratelimit"

// Middleware возвращает HTTP middleware, ограничивающий частоту запросов.
//
// rate — лимит запросов в секунду, burst — допустимый всплеск.
// Ключ лимита — путь запроса (r.URL.Path), поэтому лимиты по разным
// endpoint'ам независимы.
//
// Fail-open: при ошибке Redis (limiter.Allow вернул error) запрос
// пропускается дальше — временное отсутствие лимитирования лучше
// полной остановки сервиса из-за проблем с Redis.
func Middleware(limiter *redis_rate.Limiter, rate, burst int) func(http.Handler) http.Handler {
	limit := redis_rate.Limit{
		Rate:   rate,
		Burst:  burst,
		Period: time.Second,
	}

	// Счётчик решений rate limiter'а. Middleware стоит ПЕРЕД otelhttp, поэтому
	// отклонённые 429-запросы иначе не оставили бы следов в метриках.
	// result: allowed | rejected | fail_open
	requestsTotal, err := otel.Meter(meterName).Int64Counter(
		"rate_limit_requests_total",
		metric.WithDescription("запросы через rate limiter по решению и path"),
	)
	if err != nil {
		slog.Error("rate limiter: не удалось создать метрику", "error", err)
	}
	// instance = hostname — различает реплики сервиса. Без него все реплики пишут
	// в одну серию Prometheus, их кумулятивные значения интерливятся, и rate/irate
	// показывают фантомную нагрузку на пустом сервере
	host, hostErr := os.Hostname()
	if hostErr != nil {
		host = "unknown"
	}
	addResult := func(r *http.Request, result string) {
		if requestsTotal != nil {
			requestsTotal.Add(r.Context(), 1, metric.WithAttributes(
				attribute.String("path", r.URL.Path),
				attribute.String("result", result),
				attribute.String("instance", host),
			))
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res, err := limiter.Allow(r.Context(), r.URL.Path, limit)
			if err != nil {
				slog.ErrorContext(
					r.Context(),
					"rate limiter недоступен, запрос пропущен (fail-open)",
					"error", err,
					"path", r.URL.Path,
				)

				addResult(r, "fail_open")
				next.ServeHTTP(w, r)
				return
			}

			if res.Allowed == 0 {
				slog.WarnContext(
					r.Context(),
					"запрос отклонён rate limiter",
					"path", r.URL.Path,
					"retry_after", res.RetryAfter,
				)

				addResult(r, "rejected")

				// +1 округляет вверх: Retry-After=0 выглядел бы как "можно сразу"
				w.Header().Set("Retry-After", strconv.FormatInt(int64(res.RetryAfter.Seconds())+1, 10))
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			addResult(r, "allowed")
			next.ServeHTTP(w, r)
		})
	}
}
