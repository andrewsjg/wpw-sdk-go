package servicediscovery

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/types"
)

type broadcasterImpl struct {
	run       bool
	stepSleep int
	host      string
	port      int
	comm      Communicator
}

func (bcast *broadcasterImpl) StartBroadcast(msg types.BroadcastMessage, timeoutMillis int) error {

	log.Debug("Start svc broadcast")

	var conn Connection

	// Enable broadcaster to run
	bcast.run = true

	// Determine when the broadcast operation will expire
	timeoutTime := time.Now().Add(time.Duration(timeoutMillis) * time.Millisecond)
	timedOut := false

	for bcast.run && !timedOut {

		log.Debug("Broadcasting..")

		_conn, err := bcast.comm.Connect(bcast.host, int(bcast.port))
		conn = _conn

		if err != nil {

			log.Error(err)

		} else {

			log.Debug("Did open UDP broadcast socket")

			jsonBytes, err := json.Marshal(msg)

			if err != nil {

				log.Error(err)
				continue
			}

			_, err = conn.Write(jsonBytes)

			if err != nil {

				log.Error(err)
				continue
			}

			log.Debug("Did successfully write broadcast message")
		}

		// Sleep before broadcasting another message - don't want to flood network
		time.Sleep(time.Duration(bcast.stepSleep) * time.Millisecond)

		if timeoutMillis > 0 {
			// Determine if operation has timed out yet
			timedOut = timeoutTime.Unix() <= time.Now().Unix()
		}
	}

	// Clean up connection if still active
	if conn != nil {

		err := conn.Close()

		if err != nil {

			log.Error(err)
		} else {

			log.Debug("Did close UDP broadcast socket")
		}
	}

	log.WithFields(log.Fields{"Timed out": timedOut, "Run": bcast.run}).Debug("Finished broadcasting")

	return nil
}

func (bcast *broadcasterImpl) StopBroadcast() error {

	log.Debug("Stop svc broadcast")

	bcast.run = false

	return nil
}
