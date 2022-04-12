package wkit

import (
	"sync"

	"github.com/lavaorg/warp/warp9"
)

// Event is a buffer of bytes the content and encoding is decided by the
// EventItem owner. The suggestion is if a complex structure use cbor
// encoding, but this is not enforced.
type Event []byte

// Publisher provides a method for the EventItem owner to use to publish
// an event to the available subscribers.
//
type Publisher interface {
	Publish(event Event) error
}

// EventItem is a singleton object that provides remote clients the ability
// to join an event queue and read events.  Each client opens the EventItem
// to register for events. Each client receives it's own queue of events publihed
// and will receive all events. The event queue will be destroyed when the session
// ends (e.g. client invokes Clunk() or is disconnected. A client will perform the
// following:
//  Open - starts a "subscription" to the EventItem
//  Read - will block waiting for the next available event
//  Clunk - will destroy the subscription (and all outstanding events dropped)
//  disconnect will destroy the subscription
//
type EventItem struct {
	*BaseItem
	done            chan struct{}
	subscribers     map[*EventItem]bool //all subscribers to this instance of EventItem
	subscribersLock *sync.RWMutex
	// above items are shared across all Endpoint instances

	subscription *subscriptionType //unique instance for each client fid
}

// every client fid gets a event queue to receive events.
type subscriptionType struct {
	sync.Mutex
	readCh  chan bool
	closeCh chan bool
	events  []Event
}

// NewEventItem constructs a EventItem object. This object can
// be placed into a Warp namespace tree.
//
func NewEventItem(name string) *EventItem {
	return &EventItem{
		NewBaseItem(name, false),
		make(chan struct{}),
		make(map[*EventItem]bool),
		&sync.RWMutex{},
		nil, //no subscriptions
	}
}

// newSubscription constructs a subscriptionType object.
func newSubscription() *subscriptionType {
	return &subscriptionType{
		readCh:  make(chan bool, 1), //cap:1 to avoid block
		closeCh: make(chan bool, 2), //cap:2 to handle walk and clunk
	}
}

// Read will return the next available event; blocking until one is available.
// off: ignored
// if the object is clunk'ed while in a read a (0,Eof)
func (e *EventItem) Read(obuf []byte, off uint64, count uint32) (uint32, error) {
	found := false
	var event Event
	for !found {
		e.subscription.Lock()
		if len(e.subscription.events) != 0 {
			event = e.subscription.events[0]                  //grab the head
			e.subscription.events = e.subscription.events[1:] //slice the tail
			found = true
		}
		e.subscription.Unlock()
		if !found {
			select {
			case <-e.subscription.readCh:
			case <-e.subscription.closeCh:
				warp9.Error("Got clunk while reading. Closing read %p", e)
				return 0, warp9.WarpErrorEOF
			}
		}
	}
	elen := uint32(len(event))
	if elen > count {
		return 0, warp9.ErrorMsg(201, "event:event exceeds count request")
	}
	c := copy(obuf, event)
	if c != int(elen) {
		return 0, warp9.ErrorMsg(201, "event:event exceeds buffer size")
	}
	return elen, nil
}

// Walked will lone the endpoint unique for this subscriber and create a new
// event queue to be associated with the subscriber.
//
func (e *EventItem) Walked() (Item, error) {
	subscriber := *e //this clones the primordial EventItem for this subscriber
	e.subscribersLock.Lock()
	subscriber.subscription = newSubscription()
	warp9.Debug("subscription created: id:%p map len:%d", &subscriber, len(e.subscribers))
	e.subscribers[&subscriber] = true
	e.subscribersLock.Unlock()
	return &subscriber, nil
}

// Clunk removes the subscriber's subscription.
//
func (e *EventItem) Clunk() error {
	e.subscribersLock.Lock()
	_, found := e.subscribers[e]
	if found {
		delete(e.subscribers, e)
	}
	e.subscribersLock.Unlock()
	e.subscription.closeCh <- true //release any readers
	warp9.Debug("Clunk %p closed", e)
	return nil
}

// Publish will place event on the queue of each current
// subscription, implicitly unblocking any subscriber reads.
//
func (e *EventItem) Publish(event Event) error {

	//strConsumers := ""
	e.subscribersLock.RLock()
	for subscriber, _ := range e.subscribers {
		subscriber.subscription.Lock()
		subscriber.subscription.events = append(subscriber.subscription.events, event)
		subscriber.subscription.Unlock()
		select {
		case subscriber.subscription.readCh <- true:
			//strConsumers += fmt.Sprintf("%p ", subscriber)
		default:
		}
	}
	e.subscribersLock.RUnlock()

	return nil
}
