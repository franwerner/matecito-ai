package transport

import "log/slog"

// WithCorrelation attaches domain-correlation fields (project/change/agent —
// the real unit of work, not an HTTP request-id) to logger, per
// observability/logging. Phase 1 has no real values yet: project/change
// arrive in Phase 3 (identity + registration), agent in Phase 4 (MCP write
// event-log). Empty values are omitted rather than logged blank, and the
// helper never fails or invents a value (spec Stub A).
func WithCorrelation(logger *slog.Logger, project, change, agent string) *slog.Logger {
	args := make([]any, 0, 6)
	if project != "" {
		args = append(args, "project", project)
	}
	if change != "" {
		args = append(args, "change", change)
	}
	if agent != "" {
		args = append(args, "agent", agent)
	}
	if len(args) == 0 {
		return logger
	}
	return logger.With(args...)
}
