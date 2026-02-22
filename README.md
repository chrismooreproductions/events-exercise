# Events Exercise

##

I've done the suggested update to the task and added the http server.

To run, you should just be able to `cd cmd && go run .` The defaults in the config will
start the app on `http://localhost:9001`.

Of course in production we'd deploy this as a docker image.

Then you can make an http request to the server to receive the results as json:

```bash
curl http://localhost:9001/api/events
```

Pipe it to `jq` if you're feeling fancy:

```bash
curl http://localhost:9001/api/events | jq
```

Expected output:

```json
{
    "users": [
        {
            "id": "0af3a961-5146-46b5-93f8-95c0ab687007",
            "full_name": "Anthony Swiss-Jones",
            "email": "anthony.swiss@test.com",
            "badge_count": {
                "blue": 3,
                "green": 2,
                "red": 1
            },
            "badge_status": "amazing"
        },
        {
            "id": "6c5031e7-ff1c-4986-ac27-05a2737cd2f4",
            "full_name": "Neves Firmino",
            "email": "neves88@test.com",
            "badge_count": {
                "blue": 8,
                "green": 1,
                "red": 5
            },
            "badge_status": "ultimate"
        },
        {
            "id": "d60f3e10-b707-4c76-b165-da38b95aa4b9",
            "full_name": "Rosetta Brandon",
            "email": "rbrandon@test.com",
            "badge_count": {
                "blue": 6,
                "red": 4
            },
            "badge_status": "ultimate"
        }
    ]
}
```

## Introduction

The events exercise represents a simplified version of the event-sourcing concepts that we use internally at Utility Warehouse. At UW we consume steams of events (from Apache Kafka), and write changes from said events into a database (a projection) to build a view of the "current state" of our data, instead of it's time-series representation.

For example, an event stream might contain a number of `EmailChangedEvents` to describe how an account's associated email has changed over time. However, our microservice may only be interested in the user's current email address. If we update our database every time that we encounter an event of this type, once we have processed the full event stream the database should contain the current email address for a given account and will continue to update for the future.

## Guidance / Rules
1. Don't stress! This is a pairing exercise and you'll have friendly UW engineers to help nudge you in the right direction and offer advice if you need it.
2. We're not looking for a perfect solution, instead, we're looking to understand more about your thought process, how you approach problem solving, and how you choose to structure your solution.
3. This challenge is designed to be flexible and somewhat open-ended. You may not finish it all in the time that we have, and you shouldn't feel bad about that. Do what you can, and don't be upset if we don't get through the entire challenge.
4. Please do your best to vocalise your thought process. We're very interested in what you're thinking when you are designing and implementing your solution. We will also want to you to screen-share your editor so that we can see your solution.
5. Feel free to google, use reference, and anything else you need to assist with your solution. We appreciate that engineers aren't walking encyclopedia's of code and sometimes we need to jog our memories. Treat the exercise as you would a typical piece of work.

## Challenge

At UW we write our data projections into relational databases, search indexes, or a other forms of longer-term data stores. To simplify this challenge, we will instead write code to create an "in-memory" (for example slice, map or other type of collection) projection of account data from a stream of events that describe how a user account changes over time.

With your projection in place, it should contain the current (most up-to-date) ID, full name, email address and number of badges of each type for each of the accounts in our data stream.

### Events

- `UserAccountCreatedEvent` - _Is received when an account has been created on our system. It contains only a user ID, a UUID string which acts as a unique identifier for the new account._
- `UserAccountUpdatedEvent` - _Is received when the information (name, email) is first set or changed for a given user account identified by the associated user ID. Note that only populated fields should be updated in the projection._
- `UserGainedBadgeEvent` - _A user has gained a badge of a certain colour. Users can have any number of badges of each colour. Accounts are created with zero badges of each colour._
- `UserLostBadgeEvent` - _A user has lost a badge of a certain colour. The number of badges of each type may not fall below zero._

Events are [modelled in protobuf](https://developers.google.com/protocol-buffers), but have been compiled into Go structures found within the `internal/accounts/accounts.pb.go` file for your convenience.

You'll find the event stream code within `internal/events/stream.go`. By providing an implementation of the `StreamReader` interface to the `Read()` function of the `Stream` object, you'll be able to read events one at a time. For example...

```go
type Reader struct {}

// Implement the StreamReader interface
func (r *Reader) Read(event any) error {
	// Use the event to update your projection.
	return nil
}

func main() {
	// Create a new stream instance.
	stream := events.Stream{}
	
	// Start reading fom the stream.
	err := stream.Read(&Reader{})
	if err != nil {
		log.Fatal(err)
	}
}
```


### Badge Status

Accounts are given a particular status based on the number of badges of each type that they hold. Here's a list of the possible badge status values:

- **Great** - _Default status._
- **Amazing** - _3+ Blue Badges._
- **Ultimate** - _6+ Blue Badges and 3+ Red Badges._
- **Champion** - _10+ Blue Badges, 5+ Red Badges and 1+ Green Badges._

Badge status should be added to the account's projection.


## Additional Challenges

If you have extra time, or you're looking to continue the challenge in your own time, here are some additional challenges.

1. Consider how you would test the logic that creates your projection. How would you change your implementation to make this possible? Attempt to write some unit/acceptance tests to confirm that your projection logic functions as expected.
2. Build a simple REST API to expose the data in your projection. For example, a simple endpoint that retrieves a JSON payload for user including name, email and badge status for a given account ID.
