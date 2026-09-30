package handler

import "log/slog"

// Application holds the shared dependencies of the handlers.
// The database will be added here later.
type Application struct {
	errorLog *slog.Logger
	infoLog  *slog.Logger
}

func NewApplication(errorLog, infoLog *slog.Logger) *Application {
	return &Application{errorLog: errorLog, infoLog: infoLog}
}
