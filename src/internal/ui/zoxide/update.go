package zoxide

import "log/slog"

// Apply updates the zoxide modal with query results
func (msg UpdateMsg) Apply(m *Model) Cmd {
	// Ignore stale results - only apply if query matches current input
	currentQuery := m.textInput.Value()
	if msg.query != currentQuery {
		slog.Debug("Ignoring stale zoxide query result",
			"msgQuery", msg.query,
			"currentQuery", currentQuery,
			"id", msg.reqID)
		return nil
	}

	slog.Debug("Applying zoxide query results",
		"query", msg.query,
		"resultCount", len(msg.results),
		"id", msg.reqID)

	m.results = msg.results
	m.cursor = 0
	m.renderIndex = 0

	return nil
}
