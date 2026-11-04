// Package collection holds the rules of the deterministic shell that gathers
// a plan's sources: how a ticket is read by the tracker's convention, which
// links are followed and how, the limits, what is refused, and the context
// hash. Prose never executes; every rule here is lexical.
package collection

import (
	"errors"
	"path"
	"regexp"
	"strings"
)

// Ticket is what the tracker's convention says about a ticket file: its
// number and slug from the file name, its title from the heading and its
// status from the status line. Anything absent stays empty.
type Ticket struct {
	Path   string
	Number string
	Slug   string
	Title  string
	Status string
}

var (
	ticketName    = regexp.MustCompile(`^(\d{4})-(.+)\.ticket\.md$`)
	headingLine   = regexp.MustCompile(`(?m)^# +(?:\d{4}: +)?(.+?)\s*$`)
	statusLine    = regexp.MustCompile(`(?m)^\*\*Status:\*\* +(\S+)\s*$`)
	ticketPrefix  = regexp.MustCompile(`(?i)^(\d{4}-)?(ticket-)?`)
	notLabelRunes = regexp.MustCompile(`[^a-z0-9]+`)
)

// ParseTicket reads a ticket by the tracker's convention.
func ParseTicket(relativePath string, content []byte) Ticket {
	ticket := Ticket{Path: relativePath}

	name := path.Base(relativePath)
	if match := ticketName.FindStringSubmatch(name); match != nil {
		ticket.Number, ticket.Slug = match[1], match[2]
	} else {
		ticket.Slug = strings.TrimSuffix(strings.TrimSuffix(name, ".md"), ".ticket")
	}

	if match := headingLine.FindSubmatch(content); match != nil {
		ticket.Title = string(match[1])
	}

	if match := statusLine.FindSubmatch(content); match != nil {
		ticket.Status = string(match[1])
	}

	return ticket
}

// Name is the ticket as the lines name it: its file name.
func (ticket Ticket) Name() string { return path.Base(ticket.Path) }

// Reference is the ticket as the ledger names it: number and slug.
func (ticket Ticket) Reference() string {
	if ticket.Number == "" {
		return ticket.Slug
	}

	return ticket.Number + "-" + ticket.Slug
}

// ErrNoRuntimeID says a slug leaves nothing a runtime can be named after.
var ErrNoRuntimeID = errors.New("the ticket's slug leaves nothing to name a runtime")

// RuntimeID derives a runtime id from a ticket slug: the number and a
// TICKET marker stripped, lowercased, every other character a dash, at
// most 63 characters.
func RuntimeID(slug string) (string, error) {
	text := strings.ToLower(ticketPrefix.ReplaceAllString(strings.ToLower(slug), ""))
	text = strings.Trim(notLabelRunes.ReplaceAllString(text, "-"), "-")

	if len(text) > 63 {
		text = strings.Trim(text[:63], "-")
	}

	if text == "" {
		return "", ErrNoRuntimeID
	}

	return text, nil
}
