# RPC-AGENT

To get started:
1. Get the repo
2. In the applications/rpc-agent directory type 'go get github.com/wptechinnovation/wpw-sdk-go'
3. Then you can run the ./build-all.sh to build all the rpc-agents!

An application that hosts the WorldpayWithin SDK behind an Apache Thrift RPC listener. This agent and the SDK are written in Go but Thrift allows for other languages to interact with the SDK via it's RPC mechanism. 

You can test multiple devices on the same dev machine, by running multiple instances of the RPC-agent. Please be aware that in order to do this, each RPC-agent must be attached to a different IP address, for each of the client / device applications.

## Quick usage

Start the application and specify a port using the `-port` flag. e.g. `rpc-agent -port=9090` or `.rpc-agent -port 9090`

## Usage in more detail

### How to install the RPC agent
* Change directory to `cd $GOPATH/src/github.com/wptechninnovation/worldpay-within-sdk/applications/rpc-agent`
* Type `go install`
* This should build, package up, and install the binaries for the rpc-agent into your bin directory `$GOPATH/bin`
* If there are any errors around missing packages do additional `go get <package-repo-path>`
* If there are any compile errors, it is likely you are running a version of go that is too old (we have seen this most commonly on Ubuntu Linux)

### How to run the RPC agent
* Change to the bin directory `cd $GOPATH/bin`
* Type the following command to run the RPC agent and see the command line flags that you can pass; `./rpc-agent -help`
* You can manually set the parameters, to get everything running quickly you just need to set the prot e.g. `./rpc-agent -port 9090`

### Logging output of the RPC agent
* In the config file, specify the log name
* The log will output to the directory of the executable at the moment
* Use tail `-f <nameoflog>.log` to get the log file to output

### How to get the log file to output to the browser on a socket connection

There is the capability to output to a web browser, so you can see the logs on the device, this is not to be used in the first release. Core and RPC agent are two seperate logs, the sdk should log independent of he application. So we are planning on having two separate logs one for the SDK and one for the app.

### Current parameters

Below lists the current parameters as of 16th August 2016

```
Usage of ./rpc-agent:
  -buffer int
    	Buffer size. (default 8192)
  -buffered
    	Buffered transmission - bool.
  -framed
    	Framed transmission - bool.
  -host string
    	Listening host. (default "127.0.0.1")
  -logfile string
    	Log file, if set, outputs to file, if not, not logfile.
  -loglevel string
    	Log level (default "warn")
  -port int
    	Port to listen on. Required.
  -protocol string
    	Transport protocol. (default "binary")
  -secure
    	Secured transport - bool.
```

