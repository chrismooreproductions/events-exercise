package events

import (
	"events-exercise/internal/accounts"
	"fmt"
	"sync"
)

// StreamProcessor is an interface that allows a client to process events from a stream.

// This is the original interface I was asked to implement, but I have updated it to be more
// representitive of its function.
type StreamProcessor interface {
	// Process processes an event from the stream.
	Process(any) error
	// Result method required to return the processed events
	Result() map[string]*User
}

type BadgeStatus string
type BadgeColour string

const (
	great    BadgeStatus = "great"
	amazing  BadgeStatus = "amazing"
	ultimate BadgeStatus = "ultimate"
	champion BadgeStatus = "champion"

	blue  BadgeColour = "blue"
	red   BadgeColour = "red"
	green BadgeColour = "green"
)

type User struct {
	Id         string              `json:"id"`
	FullName   string              `json:"full_name"`
	Email      string              `json:"email"`
	BadgeCount map[BadgeColour]int `json:"badge_count"`
}

// streamProcessor is an implementation of the StreamProcessor interface
type streamProcessor struct {
	sync.Mutex
	// accounts is public for inspection
	accounts map[string]*User
}

func toBC(bc accounts.BadgeColour) BadgeColour {
	switch bc {
	case accounts.BadgeColour_BLUE:
		return blue
	case accounts.BadgeColour_RED:
		return red
	case accounts.BadgeColour_GREEN:
		return green
	default:
		return ""
	}
}

func newUserWithBadgeCount(bc accounts.BadgeColour, count int) *User {
	return &User{
		BadgeCount: map[BadgeColour]int{
			toBC(bc): count,
		},
	}
}

// the String interface on a type is useful for printing a struct using `%s` in fmt.
func (u *User) String() string {
	return fmt.Sprintf("\nUser: %s (%s)\n  Email: %s\n  Badges: %v\n  Badge Status: %v\n", u.FullName, u.Id, u.Email, u.BadgeCount, u.BadgeStatus())
}

// badge status determined by the badge status requirements.
func (u *User) BadgeStatus() BadgeStatus {
	greenCnt := u.BadgeCount[green]
	blueCnt := u.BadgeCount[blue]
	redCnt := u.BadgeCount[red]

	if blueCnt >= 10 && redCnt >= 5 && greenCnt >= 1 {
		return champion
	}

	if blueCnt >= 6 && redCnt >= 3 {
		return ultimate
	}

	if blueCnt >= 3 {
		return amazing
	}

	return great
}

// This is the idiomatic way to new up a struct.
func NewStreamProcessor() *streamProcessor {
	return &streamProcessor{
		accounts: map[string]*User{},
	}
}

// Process processes an event
func (s *streamProcessor) Process(event any) error {
	s.Lock()
	defer s.Unlock()
	// Use the event to update your projection.
	return s.processBadges(event)
}

func (s *streamProcessor) Result() map[string]*User {
	s.Lock()
	defer s.Unlock()
	// return the accounts
	return s.accounts
}

// processBadges uses a type switch to process events based on their type.
func (s *streamProcessor) processBadges(event any) error {
	switch evt := event.(type) {
	case accounts.UserAccountCreatedEvent:
		return s.processUserAccountCreatedEvent(&evt)
	case accounts.UserAccountUpdatedEvent:
		return s.processUserAccountUpdatedEvent(&evt)
	case accounts.UserGainedBadgeEvent:
		return s.processUserGainedBadgeEvent(&evt)
	case accounts.UserLostBadgeEvent:
		return s.processUserLostBadgeEvent(&evt)
	default:
		return nil
	}
}

// Process the user account created event.
func (s *streamProcessor) processUserAccountCreatedEvent(event *accounts.UserAccountCreatedEvent) error {
	u, ok := s.accounts[event.UserId]
	if !ok {
		s.accounts[event.UserId] = &User{
			Id:         event.UserId,
			BadgeCount: map[BadgeColour]int{},
		}
		return nil
	}
	u.Id = event.UserId
	if u.BadgeCount == nil {
		u.BadgeCount = map[BadgeColour]int{}
	}
	return nil
}

func (u *User) applyOptionalFields(e *accounts.UserAccountUpdatedEvent) {

}

// Process the user account updated event.
func (s *streamProcessor) processUserAccountUpdatedEvent(event *accounts.UserAccountUpdatedEvent) error {
	u, ok := s.accounts[event.UserId]
	if !ok {
		s.accounts[event.UserId] = &User{
			Id:         event.UserId,
			FullName:   event.FullName,
			Email:      event.Email,
			BadgeCount: map[BadgeColour]int{},
		}
		return nil
	}

	// only update if the event has a value
	if event.FullName != "" {
		u.FullName = event.FullName
	}
	if event.Email != "" {
		u.Email = event.Email
	}
	return nil
}

// Process the user gained badge event.
func (s *streamProcessor) processUserGainedBadgeEvent(event *accounts.UserGainedBadgeEvent) error {
	usr, usrExists := s.accounts[event.UserId]
	if !usrExists {
		// this user doesn't exist yet but let's track their badges anyway
		s.accounts[event.UserId] = newUserWithBadgeCount(event.BadgeColour, 1)
		return nil
	}
	_, badgeColourExists := usr.BadgeCount[toBC(event.BadgeColour)]
	if !badgeColourExists {
		usr.BadgeCount[toBC(event.BadgeColour)] = 1
		return nil
	}
	usr.BadgeCount[toBC(event.BadgeColour)] = usr.BadgeCount[toBC(event.BadgeColour)] + 1
	return nil
}

// Process the user lost badge event
func (s *streamProcessor) processUserLostBadgeEvent(event *accounts.UserLostBadgeEvent) error {
	usr, ok := s.accounts[event.UserId]
	if !ok {
		// this user doesn't exist yet but let's track their badges anyway
		s.accounts[event.UserId] = newUserWithBadgeCount(event.BadgeColour, -1)
		return nil
	}
	_, badgeColourExists := usr.BadgeCount[toBC(event.BadgeColour)]
	if !badgeColourExists {
		usr.BadgeCount[toBC(event.BadgeColour)] = -1
		return nil
	}
	usr.BadgeCount[toBC(event.BadgeColour)] = usr.BadgeCount[toBC(event.BadgeColour)] - 1
	return nil
}
