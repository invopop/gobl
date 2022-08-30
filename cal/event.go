package cal

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
)

// AlarmAction defines what type of action should be taken for an
// alarm.
type AlarmAction org.Key

// List of supported alarm actions.
const (
	AlarmActionAudio   AlarmAction = "audio"
	AlarmActionDisplay AlarmAction = "display"
	AlarmActionEmail   AlarmAction = "email"
)

// DefAlarmAction holds the definition of an alarm action
type DefAlarmAction struct {
	// Key to match against
	Key AlarmAction `json:"key" jsonschema:"title=Key"`
	// Description of the Note Key
	Description string `json:"description" jsonschema:"title=Description"`
}

// AlarmActionDefinitions holds a list of alarm actions with their descriptions
// for the use case.
var AlarmActionDefinitions = []DefAlarmAction{
	{AlarmActionAudio, "Play an audio alarm."},
	{AlarmActionDisplay, "Present the alarm with the provided description."},
	{AlarmActionEmail, "Send an email notification at the trigger time including the details of the event."},
}

// Event models a calendar event, based on the iCalendar
// format.
type Event struct {
	// Unique Identifier
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`

	// Sequence defines the version of this event.
	Sequence int64 `json:"sequence,omitempty" jsonschema:"title=Sequence"`

	// Who organized this event.
	Organizer *org.Person `json:"organizer,omitempty" jsonschema:"title=Organizer"`

	// The people who will be attending the event.
	Attendees []*org.Person `json:"attendess,omitempty" jsonschema:"title=Attendees"`

	// Notifications that should be triggered to alert about the event.
	Alarms []*Alarm `json:"alarms,omitempty" jsonschema:"title=Alarms"`

	// What is the event about?
	Description string `json:"description,omitempty" jsonschema:"title=Description"`

	// Any additional semi-structured data for the event.
	Meta org.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Alarm defines what should happen at a given time.
type Alarm struct {
	// Unique Identifier
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`

	// What actions should be taken when the alarm triggers?
	Actions []AlarmAction `json:"actions,omitempty" jsonschema:"title=Actions"`

	// Time in seconds before or after event that the alarm should be
	// triggered to use instead of the `At` property.
	Trigger int64 `json:"trigger,omitempty" jsonschema:"title=Trigger"`

	// Description of the alarm
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (k AlarmAction) JSONSchema() *jsonschema.Schema {
	s := &jsonschema.Schema{
		Title:       "Alarm Action",
		Type:        "string", // they're all strings
		OneOf:       make([]*jsonschema.Schema, len(AlarmActionDefinitions)),
		Description: "The action to perform when the alarm is triggered",
	}
	for i, v := range AlarmActionDefinitions {
		s.OneOf[i] = &jsonschema.Schema{
			Const:       org.Key(v.Key).String(),
			Description: v.Description,
		}
	}
	return s
}
