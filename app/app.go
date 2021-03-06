package app

import (
	"context"
	"time"

	"github.com/indrenicloud/tricloud-agent/app/logg"

	"github.com/indrenicloud/tricloud-agent/app/conn"
)

var WAITTIME time.Duration = 10 * time.Second

// Run runs
func Run() {

	ErrorChannel := make(chan struct{})
	var Connection *conn.Connection
	registered := false
	for {
		if !registered {
			registered = conn.RegisterAgent()
		}

		Connection = conn.NewConnection(context.Background(), ErrorChannel)
		if Connection == nil {
			goto sleep
		}
		Connection.Run()

		clearChannel(ErrorChannel)

		select {
		case <-ErrorChannel:
			logg.Debug("Error channel called")
			Connection.Close()
			goto sleep
		}
	sleep:
		logg.Debug("Connection Error sleeping before reconnecting")
		time.Sleep(WAITTIME)
		clearChannel(ErrorChannel)
	}

}

func clearChannel(c chan struct{}) {
	select {
	case <-c:
	default:
	}
}
