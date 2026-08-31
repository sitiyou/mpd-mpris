package mpris

import (
	"fmt"
	"time"
)

// Option represents a togglable option.
type Option func(*Instance)

// NoInstance registers the instance's name without the instance# part.
func NoInstance() Option {
	return func(ins *Instance) {
		ins.name = "org.mpris.MediaPlayer2.mpd"
	}
}

// InstanceName gives a custom name after "mpd" for the MPRIS instance.
func InstanceName(name string) Option {
	return func(ins *Instance) {
		ins.name = fmt.Sprintf("org.mpris.MediaPlayer2.mpd.%s", name)
	}
}

func IsLocal(isLocal bool) Option {
	return func(ins *Instance) {
		if isLocal {
			ins.displayName = "MPD"
		}
	}
}

// PollInterval sets the interval for re-syncing the position with MPD.
func PollInterval(d time.Duration) Option {
	return func(ins *Instance) {
		ins.pollInterval = d
	}
}
