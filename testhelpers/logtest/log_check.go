package logtest

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type Matcher struct {
	Filters []EntryPredicate
}

func (m Matcher) MatchEntries(r *Recording) []logrus.Entry {
	var results []logrus.Entry
	for _, entry := range r.Entries {
		if m.isMatch(entry) {
			results = append(results, entry)
		}
	}

	return results
}

func (m Matcher) isMatch(entry logrus.Entry) bool {
	for _, filter := range m.Filters {
		if !filter.Matches(entry) {
			return false
		}
	}
	return true
}

type CheckError struct {
	Entry     logrus.Entry
	Validator EntryPredicate
}

func (c CheckError) Error() string {
	return fmt.Sprintf("Failed validator: %s - %+v", c.Validator.Description, c.Entry)
}

type Check struct {
	Filters    []EntryPredicate
	Validators []EntryPredicate
}

func (c Check) Check(r *Recording) []error {
	matcher := Matcher{
		Filters: c.Filters,
	}

	var results []error

	for _, entry := range matcher.MatchEntries(r) {
		for _, predicate := range c.Validators {
			if !predicate.Matches(entry) {
				results = append(results, CheckError{
					Entry:     entry,
					Validator: predicate,
				})
			}
		}
	}

	return results
}