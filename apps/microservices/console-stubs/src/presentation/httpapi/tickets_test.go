package httpapi

import (
	"testing"
)

func ticketPage(t *testing.T, handler Handler, route string) ([]map[string]any, *string) {
	t.Helper()
	code, page := call(t, handler, "GET", route, nil)
	if code != 200 {
		t.Fatalf("%s answered %d: %v", route, code, page)
	}
	items := []map[string]any{}
	for _, raw := range page["items"].([]any) {
		items = append(items, raw.(map[string]any))
	}
	if page["nextCursor"] == nil {
		return items, nil
	}
	cursor := page["nextCursor"].(string)
	return items, &cursor
}

func allTickets(t *testing.T, handler Handler, product, reached string) []map[string]any {
	t.Helper()
	route := "/console/products/" + product + "/tickets"
	if reached != "" {
		route += "?reached=" + reached
	}
	all := []map[string]any{}
	seen := map[string]bool{}
	next := &route
	for pages := 0; next != nil; pages++ {
		if pages > 10 {
			t.Fatal("ticket pages never end")
		}
		items, cursor := ticketPage(t, handler, *next)
		if len(items) > 20 {
			t.Fatalf("page of %d tickets", len(items))
		}
		for _, item := range items {
			id := item["id"].(string)
			if seen[id] {
				t.Fatalf("ticket %s repeated across pages", id)
			}
			seen[id] = true
			all = append(all, item)
		}
		next = nil
		if cursor != nil {
			separator := "?"
			if reached != "" {
				separator = "&"
			}
			value := route + separator + "cursor=" + *cursor
			next = &value
		}
	}
	return all
}

var phaseOrder = map[string]int{"ticket": 0, "analysis": 1, "plan": 2, "workflow": 3}

func TestTicketsPageByPhaseWithoutDuplicates(t *testing.T) {
	handler, store := fixture(t)
	product := store.Snapshot().Products[0].ID

	first, _ := ticketPage(t, handler, "/console/products/"+product+"/tickets")
	if len(first) != 20 {
		t.Fatalf("first page holds %d tickets", len(first))
	}
	every := allTickets(t, handler, product, "")
	if len(every) < 40 {
		t.Fatalf("expected a long list, got %d tickets", len(every))
	}
	ids := map[string]bool{}
	for _, ticket := range every {
		ids[ticket["id"].(string)] = true
	}

	for phase, order := range phaseOrder {
		reached := allTickets(t, handler, product, phase)
		if len(reached) == 0 {
			t.Fatalf("no ticket has reached %s", phase)
		}
		for _, ticket := range reached {
			if phaseOrder[ticket["phase"].(string)] < order {
				t.Fatalf("ticket at %s listed under %s", ticket["phase"], phase)
			}
			if !ids[ticket["id"].(string)] {
				t.Fatalf("ticket under %s missing from the full list", phase)
			}
			if ticket[phase] == nil && phase != "ticket" {
				t.Fatalf("ticket under %s has no %s section", phase, phase)
			}
		}
	}

	outcomes := map[string]bool{}
	workflows := allTickets(t, handler, product, "workflow")
	if len(workflows) <= 20 {
		t.Fatalf("the workflows fit one page: %d", len(workflows))
	}
	for _, ticket := range workflows {
		outcomes[ticket["workflow"].(map[string]any)["outcome"].(string)] = true
	}
	for _, outcome := range []string{"running", "paused", "readyForReview", "readyToMerge"} {
		if !outcomes[outcome] {
			t.Fatalf("no workflow is %s", outcome)
		}
	}

	for _, route := range []string{
		"/console/products/" + product + "/tickets?reached=tests",
		"/console/products/" + product + "/tickets?cursor=invalid",
	} {
		if code, _ := call(t, handler, "GET", route, nil); code != 400 {
			t.Fatalf("%s accepted", route)
		}
	}
}

func TestTicketDetailAndIsolation(t *testing.T) {
	handler, store := fixture(t)
	product := store.Snapshot().Products[0].ID
	first, _ := ticketPage(t, handler, "/console/products/"+product+"/tickets")

	code, detail := call(t, handler, "GET", "/console/tickets/"+first[0]["id"].(string), nil)
	if code != 200 || detail["ticket"].(map[string]any)["id"] != first[0]["id"] {
		t.Fatalf("ticket detail: %d %v", code, detail)
	}
	if code, _ := call(t, handler, "GET", "/console/tickets/missing", nil); code != 404 {
		t.Fatal("unknown ticket found")
	}

	other := createOtherProduct(t, handler)
	items, cursor := ticketPage(t, handler, "/console/products/"+other+"/tickets")
	if len(items) != 0 || cursor != nil {
		t.Fatal("a created project has tickets")
	}
}

func TestEarlyAccessTicketCarriesTheWebsiteFigures(t *testing.T) {
	handler, store := fixture(t)
	product := store.Snapshot().Products[0].ID
	var early map[string]any
	for _, ticket := range allTickets(t, handler, product, "workflow") {
		if ticket["slug"] == "early-access-runtime" {
			early = ticket
		}
	}
	if early == nil {
		t.Fatal("the early-access ticket has not reached the workflow")
	}
	if early["artifacts"].(map[string]any)["prd"] != "0002-PRD-early-access.prd.md" {
		t.Fatalf("artifacts %v", early["artifacts"])
	}
	plan := early["plan"].(map[string]any)
	suggested := plan["suggested"].(map[string]any)
	bounded := plan["bounded"].(map[string]any)
	if suggested["cpuMillis"] != float64(1450) || suggested["memoryMiB"] != float64(2304) {
		t.Fatalf("suggested %v", suggested)
	}
	if bounded["cpuMillis"] != float64(2000) || bounded["memoryMiB"] != float64(3072) {
		t.Fatalf("bounded %v", bounded)
	}
	episode := early["analysis"].(map[string]any)["episodes"].([]any)[0].(map[string]any)
	if episode["id"] != "ep-0412" || episode["similarity"] != 0.81 {
		t.Fatalf("first episode %v", episode)
	}
	workflow := early["workflow"].(map[string]any)
	if workflow["estimateCents"] != float64(950) || workflow["budgetCents"] != float64(2500) {
		t.Fatalf("workflow %v", workflow)
	}
}
