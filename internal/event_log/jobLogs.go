package event_log

import (
	"time"
	"context"
	"internal/event_log/types"
	jobTypes "internal/jobs/types"
)

func LogCreatedJob(ctx context.Context, job *jobTypes.LogableJob, editor *types.Editor) *types.NormalizedLoggedEvent {
	changesMap, err := structToMap(job, "m"); if err != nil {
		return nil
	}

	changes := make(map[field]types.Change)
	for key, value := range changesMap {
		changes[key] = types.Change{Old:nil, New:value}
	}

	event := &logEvent {
		EventType: created,
		Timestamp: time.Now().Unix(),
		Editor: editor.Name,
		EditorID: editor.Oid,
		Changes: changes,
	}
	return event.log(ctx, editor.Collection).normalize()
}
