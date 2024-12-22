package prom

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Создаем метрики для мониторинга api
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint"}, // Метки для метода и эндпоинта
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of request duration",
			Buckets: prometheus.DefBuckets, // Стандартные интервалы
		},
		[]string{"method", "endpoint"},
	)

	dbUp = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_up",
			Help: "Database availability: 1 if up, 0 if down",
		},
	)

	dbQueryDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Histogram of database query durations",
			Buckets: prometheus.DefBuckets,
		},
	)
)

// Регистрация метрик
func InitPrometheus() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

// Middleware для сбора метрик
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(r.Method, r.URL.Path))
		defer timer.ObserveDuration()

		// Увеличиваем счетчик запросов
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path).Inc()

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}

// Обработчик для метрик
func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}

// проверяет состояние соединения с базой данных
func CheckDBState(db *pgxpool.Pool) {
	start := time.Now()
	err := db.Ping(context.Background()) // Проверка соединения
	duration := time.Since(start).Seconds()

	// Записываем продолжительность запроса
	dbQueryDuration.Observe(duration)

	if err != nil {
		dbUp.Set(0) // БД недоступна
	} else {
		dbUp.Set(1) // БД доступна
	}
}

// InitDBMetrics регистрирует метрики БД
func InitDBMetrics() {
	prometheus.MustRegister(dbUp)
	prometheus.MustRegister(dbQueryDuration)
}
