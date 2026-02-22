package events_test

import (
	"events-exercise/internal/events"
	"events-exercise/internal/logger"
	"testing"
)

func TestEventProcessorResultMatchesExpected(t *testing.T) {
	l := logger.Logger
	p := events.NewStreamProcessor()
	events.NewStreamReader(p, l)

	expected := []struct {
		id       string
		fullName string
		email    string
		badges   map[string]int
	}{
		{
			id:       "0af3a961-5146-46b5-93f8-95c0ab687007",
			fullName: "Anthony Swiss-Jones",
			email:    "anthony.swiss@test.com",
			badges:   map[string]int{"blue": 3, "green": 2, "red": 1},
		},
		{
			id:       "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
			fullName: "Neves Firmino",
			email:    "neves.firmino@test.com",
			badges:   map[string]int{"blue": 8, "green": 1, "red": 5},
		},
		{
			id:       "d60f3e10-b707-4c76-b165-da38b95aa4b9",
			fullName: "Rosetta Brandon",
			email:    "rbrandon@test.com",
			badges:   map[string]int{"blue": 6, "red": 4},
		},
	}

	result := p.Result()

	for _, tc := range expected {
		t.Run(tc.id, func(t *testing.T) {
			user, ok := result[tc.id]
			if !ok {
				t.Fatalf("expected user %s to exist", tc.id)
			}
			if user.FullName != tc.fullName {
				t.Errorf("expected full_name %q, got %q", tc.fullName, user.FullName)
			}
			if user.Email != tc.email {
				t.Errorf("expected email %q, got %q", tc.email, user.Email)
			}
			for colour, count := range tc.badges {
				if user.BadgeCount[events.BadgeColour(colour)] != count {
					t.Errorf("expected %d %s badges, got %d", count, colour, user.BadgeCount[events.BadgeColour(colour)])
				}
			}
			// Check for unexpected badge colours
			for colour := range user.BadgeCount {
				if _, ok := tc.badges[string(colour)]; !ok {
					t.Errorf("unexpected badge colour %q for user %s", colour, tc.id)
				}
			}
		})
	}
}
