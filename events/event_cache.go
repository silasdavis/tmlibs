package events

// An EventCache buffers events for a Fireable
// All events are cached. Filtering happens on Flush
type EventCache struct {
	evsw   Fireable
	events []eventInfo
}

// Create a new EventCache with an EventSwitch as backend
func NewEventCache(evsw Fireable) *EventCache {
	return &EventCache{
		evsw: evsw,
	}
}

// a cached event
type eventInfo struct {
	event string
	data  EventData
}

// Cache an event to be fired upon finality.
func (evc *EventCache) FireEvent(event string, data EventData) {
	// append to list (go will grow our backing array exponentially)
	evc.events = append(evc.events, eventInfo{event, data})
}

// Fire events by running evsw.FireEvent on all cached events. Blocks.
// Clears cached events
func (evc *EventCache) Flush() {
	for _, ei := range evc.events {
		evc.evsw.FireEvent(ei.event, ei.data)
	}
	// Clear the buffer by re-slicing its length to zero
	if cap(evc.events) > len(evc.events)<<1 {
		// Trim the backing array capacity when it is more than double the length of the slice to avoid tying up memory
		// after a spike in the number of events to buffer
		evc.events = evc.events[:0:len(evc.events)]
	} else {
		// Re-slice the length to 0 to clear buffer but hang on to spare capacity in backing array that has been added
		// in previous cache round
		evc.events = evc.events[:0]
	}
}
