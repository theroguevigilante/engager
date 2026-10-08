// Package demo holds scripted behavioural telemetry used to demonstrate the
// adaptation loop for different kinds of disengagement. The scripts are only
// inputs; every decision is still made by the real deterministic engine.
package demo

// Telemetry is an aggregate behavioural batch for one attempt.
type Telemetry struct {
	Session, Active, Idle, TabSwitches, Interactions int
	Attempts                                         int
	Status                                           string // completed | abandoned
}

type Scenario struct {
	Key   string    `json:"key"`
	Label string    `json:"label"`
	Story string    `json:"story"`
	Pre   Telemetry `json:"-"`
	Post  Telemetry `json:"-"`
}

var engagedAfter = Telemetry{Session: 90, Active: 80, Idle: 10, TabSwitches: 0, Interactions: 40, Attempts: 1, Status: "completed"}

var Scenarios = []Scenario{
	{
		Key: "struggling", Label: "Struggling",
		Story: "Works hard, retries several times, but can't crack it.",
		Pre:   Telemetry{Session: 120, Active: 100, Idle: 20, TabSwitches: 1, Interactions: 110, Attempts: 3, Status: "completed"},
		Post:  engagedAfter,
	},
	{
		Key: "bored", Label: "Bored",
		Story: "Races through in seconds with barely any effort: not enough challenge.",
		Pre:   Telemetry{Session: 15, Active: 12, Idle: 3, TabSwitches: 0, Interactions: 4, Attempts: 1, Status: "completed"},
		Post:  engagedAfter,
	},
	{
		Key: "drifting", Label: "Drifting away",
		Story: "Stares at the screen, keeps switching tabs, barely interacts.",
		Pre:   Telemetry{Session: 100, Active: 30, Idle: 70, TabSwitches: 4, Interactions: 5, Attempts: 1, Status: "completed"},
		Post:  engagedAfter,
	},
}

func Find(key string) Scenario {
	for _, s := range Scenarios {
		if s.Key == key {
			return s
		}
	}
	return Scenarios[0]
}
