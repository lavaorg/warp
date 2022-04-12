package wkit

import (
	"testing"
	"time"
)

var events *EventItem = NewEventItem("MyEvents")

// testing  -- test object without remote access
//
// - create an event object
// - create several readers
//   - for each reader
//     - walk to event object
//     - open event object
//     - block on read
//     - repeat 5 times
//     - clunk object
// - start a producer
// - generate 5 events
// - verify all clients clunked
//

const (
	readerCount = 2
	eventCount  = 2
)

func TestEventItem(t *testing.T) {
	//warp9.LogDebug(true)
	createReaders(t, events)
	time.Sleep(2 * time.Second)
	produceEvents(t, events)
	time.Sleep(2 * time.Second)
	verify(t, events)
}

func verify(t *testing.T, events *EventItem) {
	cnt := len(events.subscribers)
	if cnt != 0 {
		t.Errorf("readers didn't clunk:%v", cnt)
	}
}

func createReaders(t *testing.T, events *EventItem) {
	for x := 0; x < readerCount; x++ {
		go reader(t, events)
	}
}

func produceEvents(t *testing.T, events *EventItem) {

	event := []byte("event")
	// publish eventCount events
	for x := 0; x < eventCount; x++ {
		t.Logf("publish: %v\n", string(event))
		events.Publish(event)
	}
}

func reader(t *testing.T, events *EventItem) {

	x, err := events.Walked() //create subscription
	if err != nil {
		t.Fatalf("Walked failed:%v", err)
	}
	evt := x.(*EventItem) //let it panic if bad
	evt.Open(0)

	for x := 0; x < eventCount; x++ {
		var buf []byte = make([]byte, 100)
		n, err := evt.Read(buf, 0, uint32(len(buf)))
		if err != nil {
			t.Errorf("read failed:%v", err)
		}
		s := string(buf[:n])
		if s != "event" {
			t.Errorf("read returned bad event:[%v][%v][%v]\n", n, s, buf[:n])
		}
	}
	evt.Clunk()
}
