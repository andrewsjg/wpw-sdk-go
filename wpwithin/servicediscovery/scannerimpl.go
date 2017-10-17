package servicediscovery

import (
	"fmt"
	"strings"
	"time"

	"encoding/json"

	"github.com/WPTechInnovation/wpw-sdk-go/wpwithin/types"
	log "github.com/sirupsen/logrus"
)

type scannerImpl struct {
	run       bool /* Used to stop scanning before timeout */
	stepSleep int  /* Time to sleep before connection reads */
	comm      Communicator
}

func (scanner *scannerImpl) ScanForServices(timeout int) (map[string]types.BroadcastMessage, error) {

	/*
		This function works by setting up a connection broadcast listener, returning a result object which includes a channel
		informing when scanning is finished.
		Inside the result is an error object and also a list of scanned services
		Error is != nil if there was a problem
	*/

	log.Debugf("Begin ScanForServices(timeout = %d)", timeout)

	result := make(map[string]types.BroadcastMessage, 0)
	// Enable the scanner to run
	scanner.run = true
	// Calculate when the operation will expire based on the timeout duration
	timeoutTime := time.Now().Add(time.Duration(timeout) * time.Millisecond)
	timedOut := false

	var srvConn Connection

	// Reading incoming messages - 2kb buffer
	buf := make([]byte, 2048)

	for scanner.run && !timedOut {

		_srvConn, err := scanner.comm.Listen()
		srvConn = _srvConn

		if err != nil {

			scanner.run = false
		}

		// Defer closing connection in go routine instead of main routine as it will be closed before the go routine starts.
		defer srvConn.Close()

		// Wait for incoming message
		//srvConn.SetProperty("ReadDeadLine", time.Now().Add(time.Duration(scanner.stepSleep)*time.Millisecond))
		srvConn.SetProperty("ReadDeadLine", time.Now().Add(time.Duration(timeout)*time.Millisecond))

		nRecv, addrRecv, err := srvConn.Read(buf)

		if err != nil {

			// Want to ignore timeout errors for scanning
			// Not an error in the sense that a timeout is expected if there is nothing broadcasted.
			if !strings.Contains(err.Error(), "i/o timeout") {

				log.Error(err)
			}
		}

		if nRecv > 0 { /* Did we actually receive any data? */

			log.Debugf("Did receive message from %s: %s", addrRecv, string(buf[0:nRecv]))

			var msg types.BroadcastMessage

			// Try to deserialize the message into a broadcast message
			// NB: Anybody can send a message here so not all messages are expected to be valid
			err = json.Unmarshal(buf[0:nRecv], &msg)

			if err != nil {

				// This is not neccessarily an error - could be a message from another source (ignore)
				log.WithFields(log.Fields{"Error: ": fmt.Sprintf("Err: %q", err.Error())}).Error("Did not decode message")
			} else {

				log.Infof("Did decode broadcast message: %#v", msg)

				result[msg.ServerID] = msg
			}
		}
		// Have we timed out? i.e. Is the current time greater or equal time out time
		timedOut = timeoutTime.Unix() <= time.Now().Unix()
	}

	log.WithFields(log.Fields{"Timed out": timedOut, "Run": scanner.run, "Found": len(result)}).Debug("Finish ScanForServices()")

	return result, nil
}

func (scanner *scannerImpl) ScanForService(timeout int, serviceName string) (*types.BroadcastMessage, error) {
	result, err := scanner.ScanForServices(timeout)
	if err != nil {
		return nil, err
	} else {
		for _, r := range result {
			if r.DeviceName == serviceName {
				return &r, nil
			}
		}
	}
	return nil, nil
}

func (scanner *scannerImpl) StopScanner() {

	scanner.run = false
}
