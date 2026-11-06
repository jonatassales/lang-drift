package main

import (
	"strings"

	"github.com/jonatassales/tao/lib/cli/src/application/session"
	clisession "github.com/jonatassales/tao/lib/cli/src/presentation/cli/session"
)

var (
	inspectLineIDs = map[string]string{
		"services/signup-api": "inspectGo",
		"services/notify-api": "inspectNode",
		"platform/nats":       "inspectNats",
		"platform/redis":      "inspectRedis",
		"apps/signup":         "inspectNext",
	}
	componentLineIDs = map[string]string{
		"services/signup-api": "go",
		"services/notify-api": "node",
		"platform/nats":       "nats",
		"platform/redis":      "redis",
		"apps/signup":         "next",
	}
	waitingLineIDs = map[string]string{
		"spec":      "waitingSpec",
		"tickets":   "waitingTickets",
		"implement": "waitingImplement",
	}
	confirmedLineIDs = map[string]string{
		"spec":      "confirmedSpec",
		"tickets":   "confirmedTickets",
		"implement": "confirmedImplement",
	}
	stageLineIDs = map[string]string{
		"grilling":    "grilling",
		"to-spec":     "toSpec",
		"to-tickets":  "toTickets",
		"implement":   "implement",
		"code-review": "codeReview",
	}
)

// printedLineOf turns one observed session event into a printedLine keyed by
// the simulation's line id.
func printedLineOf(event session.Event) printedLine {
	lineID := event.ID
	name, _ := event.Data["name"].(string)
	componentPath, _ := event.Data["path"].(string)
	stage, _ := event.Data["stage"].(string)
	skill, _ := event.Data["skill"].(string)

	switch event.ID {
	case session.ReadDocument:
		switch {
		case strings.HasPrefix(name, "0002-PRD"):
			lineID = "readPrd"
		case strings.HasPrefix(name, "0005-DD"):
			lineID = "readDesignDoc"
		case strings.HasPrefix(name, "0011-SPEC"):
			lineID = "readSpec"
		}
	case session.Inspect:
		if mapped, known := inspectLineIDs[componentPath]; known {
			lineID = mapped
		}
	case session.Component:
		if mapped, known := componentLineIDs[componentPath]; known {
			lineID = mapped
		}
	case session.Waiting:
		if mapped, known := waitingLineIDs[stage]; known {
			lineID = mapped
		}
	case session.Confirmed:
		if mapped, known := confirmedLineIDs[stage]; known {
			lineID = mapped
		}
	case session.Stage:
		lineID = skill
		if mapped, known := stageLineIDs[skill]; known {
			lineID = mapped
		}
	}

	return printedLine{
		id:   lineID,
		kind: clisession.KindOf(event),
		text: clisession.Text(event),
	}
}
