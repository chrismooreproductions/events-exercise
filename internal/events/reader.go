package events

import (
	"events-exercise/internal/accounts"
	"log/slog"
)

type StreamReader interface {
	Read() error
}

// This is the original Stream struct, updated to provide a more idiomatic way of providing
// the processor (using the StreamProcessor interface).

// I would suggest this is a preferrential pattern to passing the StreamProcessor directly
// to the Read() method.
type streamReader struct {
	logger *slog.Logger

	processor StreamProcessor
}

// NewStreamReader returns a new StreamReader.

// Again, this is the idiomatic way to new up a struct of a given type.
func NewStreamReader(p StreamProcessor, l *slog.Logger) {
	sr := &streamReader{processor: p, logger: l}

	numReaders := 3
	evtChan := make(chan any, numReaders)

	for i := 0; i < numReaders; i++ {
		go sr.Read(evtChan)
	}

	for _, evt := range eventSequence() {
		evtChan <- evt
	}
	close(evtChan)
}

// Read reads events from the stream and processes them via the StreamProcessor.

// In the original task it was really confusing having a Reader interface reading
// something within another Reader interface.

// i.e.
// func (s *StreamReader) Read() error {
//   ...
//     err := Read()
// }
// is something we should not do.

// We should call the processor a processor because that's what it is - this is the
// Reader, and it offloads events to the Processor for processing correctly.

// This looks a lot nicer to me. The candidate could be presented with the StreamProcessor
// interface and asked to implement it, and we can stub it for them in `processor.go`.
func (s *streamReader) Read(ch <-chan any) {
	for evt := range ch {
		if err := s.processor.Process(evt); err != nil {
			s.logger.Error("processing event", "event", evt)
		}
	}
}

// eventSequence returns a sequence of hard-coded events.
func eventSequence() []any {
	return []any{
		accounts.UserAccountCreatedEvent{
			UserId: "0af3a961-5146-46b5-93f8-95c0ab687007",
		},
		accounts.UserAccountCreatedEvent{
			UserId: "d60f3e10-b707-4c76-b165-da38b95aa4b9",
		},
		accounts.UserAccountCreatedEvent{
			UserId: "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "0af3a961-5146-46b5-93f8-95c0ab687007",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "0af3a961-5146-46b5-93f8-95c0ab687007",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "0af3a961-5146-46b5-93f8-95c0ab687007",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "d60f3e10-b707-4c76-b165-da38b95aa4b9",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserAccountUpdatedEvent{
			UserId:   "d60f3e10-b707-4c76-b165-da38b95aa4b9",
			FullName: "Rosetta Brandon",
			Email:    "rbrandon@test.com",
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "d60f3e10-b707-4c76-b165-da38b95aa4b9",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "d60f3e10-b707-4c76-b165-da38b95aa4b9",
			BadgeColour: accounts.BadgeColour_RED,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "0af3a961-5146-46b5-93f8-95c0ab687007",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserAccountUpdatedEvent{
			UserId:   "0af3a961-5146-46b5-93f8-95c0ab687007",
			FullName: "Anthony Swiss",
			Email:    "anthony.swiss@test.com",
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "0af3a961-5146-46b5-93f8-95c0ab687007",
			BadgeColour: accounts.BadgeColour_GREEN,
		},
		accounts.UserLostBadgeEvent{
			UserId:      "0af3a961-5146-46b5-93f8-95c0ab687007",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "d60f3e10-b707-4c76-b165-da38b95aa4b9",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "d60f3e10-b707-4c76-b165-da38b95aa4b9",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "0af3a961-5146-46b5-93f8-95c0ab687007",
			BadgeColour: accounts.BadgeColour_RED,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "0af3a961-5146-46b5-93f8-95c0ab687007",
			BadgeColour: accounts.BadgeColour_GREEN,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "d60f3e10-b707-4c76-b165-da38b95aa4b9",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "d60f3e10-b707-4c76-b165-da38b95aa4b9",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "d60f3e10-b707-4c76-b165-da38b95aa4b9",
			BadgeColour: accounts.BadgeColour_RED,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "d60f3e10-b707-4c76-b165-da38b95aa4b9",
			BadgeColour: accounts.BadgeColour_RED,
		},
		accounts.UserLostBadgeEvent{
			UserId:      "d60f3e10-b707-4c76-b165-da38b95aa4b9",
			BadgeColour: accounts.BadgeColour_RED,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "d60f3e10-b707-4c76-b165-da38b95aa4b9",
			BadgeColour: accounts.BadgeColour_RED,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "d60f3e10-b707-4c76-b165-da38b95aa4b9",
			BadgeColour: accounts.BadgeColour_RED,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserAccountUpdatedEvent{
			UserId:   "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
			FullName: "Neves Firmino",
			Email:    "neves88@test.com",
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
			BadgeColour: accounts.BadgeColour_BLUE,
		},
		accounts.UserAccountUpdatedEvent{
			UserId: "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
			Email:  "neves.firmino@test.com",
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
			BadgeColour: accounts.BadgeColour_RED,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
			BadgeColour: accounts.BadgeColour_RED,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
			BadgeColour: accounts.BadgeColour_RED,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
			BadgeColour: accounts.BadgeColour_RED,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
			BadgeColour: accounts.BadgeColour_RED,
		},
		accounts.UserGainedBadgeEvent{
			UserId:      "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
			BadgeColour: accounts.BadgeColour_GREEN,
		},
		accounts.UserAccountUpdatedEvent{
			UserId:   "0af3a961-5146-46b5-93f8-95c0ab687007",
			FullName: "Anthony Swiss-Jones",
		},
	}
}
