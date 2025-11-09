package main

import (
	"log/slog"
	"net/http"
	"os"

	m "vrc-world-matching"
)

func main() {
	// Cloud Logging用のログ設定
	ops := slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				a.Key = "severity"
				level := a.Value.Any().(slog.Level)
				if level == slog.LevelWarn {
					a.Value = slog.StringValue("WARNING")
				}
			}

			return a
		},
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &ops))
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /worlds", m.AuthMiddleware(m.GetWorldList))
	mux.HandleFunc("GET /worlds/{world_id}", m.AuthMiddleware(m.GetWorld))
	mux.HandleFunc("POST /worlds/{world_id}", m.AuthMiddleware(m.PostWorld))
	mux.HandleFunc("DELETE /worlds/{world_id}", m.AuthMiddleware(m.DeleteWorld))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
