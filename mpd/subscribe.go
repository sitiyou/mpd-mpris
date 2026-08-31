package mpd

import (
	"context"
	"time"

	"github.com/fhs/gompd/v2/mpd"
	"github.com/pkg/errors"
)

// ErrPollTimeout is returned by Watcher.Poll when no event arrived within the
// given timeout.
var ErrPollTimeout = errors.New("mpd: poll timeout")

// Watcher is our implementation of the watcher.
// It automatically subscribes to MPRIS-related events and
// `Poll` can be used to wait for any event.
type Watcher struct {
	*mpd.Watcher
}

var (
	// See https://mpd.readthedocs.io/en/latest/protocol.html#command-idle
	eventsToSubscribe = []string{
		"playlist", // the queue (i.e. the current playlist) has been modified
		"player",   // the player has been started, stopped or seeked or tags of the currently playing song have changed (e.g. received from stream)
		"mixer",    // the volume has been changed
		"options",  // options like repeat, random, crossfade, replay gain
	}
)

// NewWatcher creates a new watcher with the given parameters.
func NewWatcher(net, addr, passwd string) (*Watcher, error) {
	w, err := mpd.NewWatcher(net, addr, passwd, eventsToSubscribe...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Watcher{w}, nil
}

// Poll waits for the next event, or errors out. If no event arrives within
// timeout, it returns ErrPollTimeout; timeout <= 0 waits indefinitely.
func (w *Watcher) Poll(ctx context.Context, timeout time.Duration) error {
	if timeout <= 0 {
		select {
		case <-w.Event:
			return nil
		case err := <-w.Error:
			return errors.Wrap(err, "polling for events")
		case <-ctx.Done():
			return context.Canceled
		}
	}
	select {
	case <-w.Event:
		return nil
	case err := <-w.Error:
		return errors.Wrap(err, "polling for events")
	case <-ctx.Done():
		return context.Canceled
	case <-time.After(timeout):
		return ErrPollTimeout
	}
}
