package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/rifflock/lfshook"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin"
	"github.com/wptechinnovation/wpw-sdk-go/wpwithin/rpc"
)

var applicationVersion string
var applicationBuildDate string
var applicationPlatform string

const applicationName string = "Worldpay Within rpc-agent"

/*
	A simple program to enable the WPWithin Core RPC interface.
	The intention is that this program is called by language wrappers in order to gain RPC access to the core.
*/

const exitOK = 0
const exitGeneralErr = 1

// Log level selectors passed in as argument
const levelPanic = "panic"
const levelFatal = "fatal"
const levelError = "error"
const levelWarn = "warn"
const levelInfo = "info"
const levelDebug = "debug"

// General constants
const logfilePerms = 0755
const rpcMinPort = 1
const defaultArgPort = 0 // Defaulting to this should cause error (desired) (force port specifer// )
const defaultArgTransportBuffer = 8192
const defaultArgFramed = false
const defaultArgBuffered = false
const defaultArgSecure = false
const defaultArgHost = "127.0.0.1"
const defaultArgProtocol = "binary"
const defaultArgCallbackPort = 0 // Default 0 means callback feature not to be used
const defaultLogFile = "wpw.log"

const argNameLogLevel = "loglevel"
const argNameLogfile = "logfile"
const argNameProtocol = "protocol"
const argNameFramed = "framed"
const argNameBuffered = "buffered"
const argNameHost = "host"
const argNamePort = "port"
const argNameSecure = "secure"
const argNameBuffer = "buffer"
const argNameCallbackPort = "callbackport"
const argNameVersion = "version"

var flagLogFile string
var flagLogLevel string
var flagProtocol string
var flagFramed bool
var flagBuffered bool
var flagHost string
var flagPort int
var flagSecure bool
var flagBuffer int
var flagCallbackPort int
var flagVersion bool

// Globally scoped vars
var sdk wpwithin.WPWithin
var rpcConfig rpc.Configuration

func main() {

	log.Debug("Before flag.parse()")
	flag.Parse()
	log.Debug("After flag.parse()")

	log.Debug("Will call checkAppInfo()")
	checkAppInfo()
	log.Debug("After call checkAppInfo()")

	log.Info(getVersionInfo())

	// Start off by setting logging to a high level
	// This way we can catch output during initial setup of args and logging via arguments.
	log.SetLevel(log.DebugLevel)
	log.SetOutput(os.Stdout)

	log.Debug("Begin main()")

	log.Debug("Before initLogFile()")
	initLogfile(flagLogLevel, flagLogFile)
	log.Debug("After initLogFile()")

	log.Debug("Before startRPC()")
	startRPC()
	log.Debug("After startRPC()")

	fmt.Println("Program end...")
	log.Debug("Application end - End main()")

	os.Exit(exitOK)
}

func init() {

	// Log config args
	flag.StringVar(&flagLogLevel, argNameLogLevel, levelWarn, "Log level")
	flag.StringVar(&flagLogFile, argNameLogfile, defaultLogFile, "Log file, if set, outputs to file, if not, not logfile.")

	// Program specific arguments
	flag.StringVar(&flagProtocol, argNameProtocol, defaultArgProtocol, "Transport protocol.")
	flag.BoolVar(&flagFramed, argNameFramed, defaultArgFramed, "Framed transmission - bool.")
	flag.BoolVar(&flagBuffered, argNameBuffered, defaultArgBuffered, "Buffered transmission - bool.")
	flag.StringVar(&flagHost, argNameHost, defaultArgHost, "Listening host.")
	flag.IntVar(&flagPort, argNamePort, defaultArgPort, "Port to listen on. Required.")
	flag.BoolVar(&flagSecure, argNameSecure, defaultArgSecure, "Secured transport - bool.")
	flag.IntVar(&flagBuffer, argNameBuffer, defaultArgTransportBuffer, "Buffer size.")
	flag.IntVar(&flagCallbackPort, argNameCallbackPort, defaultArgCallbackPort, "Callback Port")

	flag.BoolVar(&flagVersion, argNameVersion, false, "Print application version info")
}

func startRPC() {

	log.Debug("Before startRPC()")

	log.Debug("Before assign RPC config.")
	rpcConfig = rpc.Configuration{}
	rpcConfig.Protocol = flagProtocol
	rpcConfig.Framed = flagFramed
	rpcConfig.Buffered = flagBuffered
	rpcConfig.Host = flagHost
	rpcConfig.Port = flagPort
	rpcConfig.Secure = flagSecure
	rpcConfig.BufferSize = flagBuffer
	rpcConfig.CallbackPort = flagCallbackPort
	log.Debug("After assign RPC config.")

	// Validate required (with no defaults)
	if rpcConfig.Port < rpcMinPort {

		fmt.Println("Port value must be greater than zero")

		log.WithFields(log.Fields{"Port": rpcConfig}).Fatal("Invalid listening port provided")

		os.Exit(exitGeneralErr)
	}

	log.WithField("Configuration: ", fmt.Sprintf("%+v", rpcConfig)).Debug("Before rpc.NewService")
	rpc, err := rpc.NewService(rpcConfig, sdk)
	log.Debug("After rpc.NewService")

	if err != nil {

		log.WithFields(log.Fields{"Error": err.Error()}).Fatal("Error creating new RPC service")

		os.Exit(exitGeneralErr)
	}

	log.WithFields(log.Fields{"port": rpcConfig.Port}).Debug("Attempting to start RPC interface on port")
	if err := rpc.Start(); err != nil {

		log.WithFields(log.Fields{"Error": err.Error()}).Fatal("Error starting RPC service")

		fmt.Printf("Error starting RPC service: %q\n", err.Error())

		os.Exit(exitGeneralErr)
	}

	log.Debug("End startRPC()")
}

func initLogfile(logLevel, logFile string) {

	log.Debug("Begin log setup")

	switch logLevel {

	case levelPanic:
		log.SetLevel(log.PanicLevel)
	case levelFatal:
		log.SetLevel(log.FatalLevel)
	case levelError:
		log.SetLevel(log.ErrorLevel)
	default:
	case levelWarn:
		log.SetLevel(log.WarnLevel)
	case levelInfo:
		log.SetLevel(log.InfoLevel)
	case levelDebug:
		log.SetLevel(log.DebugLevel)
	}

	log.Debug("Begin parsing log level arguments")

	log.Debug("Begin parsing log file arguments and setup log file")

	if !strings.EqualFold(logFile, "") {

		log.Debug("Setting up log text formatter")
		tf := &log.TextFormatter{}
		tf.DisableColors = true
		tf.FullTimestamp = true
		log.SetFormatter(tf)
		log.WithField("TextFormatter", fmt.Sprintf("%+v", tf)).Debug("End set up log text formatter")

		log.AddHook(lfshook.NewHook(lfshook.PathMap{
			log.PanicLevel: logFile,
			log.FatalLevel: logFile,
			log.ErrorLevel: logFile,
			log.WarnLevel:  logFile,
			log.DebugLevel: logFile,
			log.InfoLevel:  logFile,
		}))

	} else {

		log.Debug("Will not be logging to file - logfile is empty.")
	}

	log.Debug("End log setup")
}

func checkAppInfo() {

	log.Debug("Begin checkAppInfo()")

	if flagVersion {

		versionInfo := getVersionInfo()

		fmt.Println(versionInfo)

		os.Exit(0)
	}

	log.Debug("End checkAppInfo()")
}

func getVersionInfo() string {

	return fmt.Sprintf("\n\n%s v%s (Built on %s) (%s)\n\n", applicationName, applicationVersion, applicationBuildDate, applicationPlatform)
}
