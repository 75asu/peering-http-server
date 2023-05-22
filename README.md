## Server Code Documentation

### Overview

This code implements a server that responds to GET requests on the path "/ping" by returning a "pong" message along with the hostname of the server. The server also periodically reads a configuration file that contains a list of host-port combinations and sends GET requests to each of these combinations, logging the response from each server.


### Main Function

The `main` function is the entry point of the server application. It does the following:

1. Defines the path "/ping" and its corresponding handler function, which returns a "pong" response along with the hostname of the server.
2. Starts the server on port 8080 using `http.ListenAndServe`.
3. Creates a channel `configUpdateCh` to receive configuration update signals.
4. Starts a goroutine to periodically check for configuration updates by calling the `watchConfigUpdates` function.
5. Enters an infinite loop where it reads the configuration, sends GET requests based on the current configuration, and waits for a configuration update signal.

### sendPingRequests Function

The `sendPingRequests` function is responsible for sending GET requests to the host-port combinations specified in the configuration. It performs the following steps:

1. Splits the host-port configuration into individual lines.
2. Iterates over each line, trimming any leading or trailing whitespace.
3. If the line is not empty, splits it into host and port components.
4. Sends a GET request to the specified host and port by constructing the URL.
5. Reads the response body and logs the response along with the hostname of the server that handled the request.

### readConfig Function

The `readConfig` function reads the host-port configuration from the file specified by `configFilePath`. It performs the following steps:

1. Reads the contents of the configuration file using `ioutil.ReadFile`.
2. Converts the read data to a string.
3. Returns the host-port configuration string.

### watchConfigUpdates Function

The `watchConfigUpdates` function watches for changes in the configuration file and sends an update signal through the `configUpdateCh` channel when a modification is detected. It operates as follows:

1. Starts a goroutine to periodically check the last modified time of the configuration file.
2. Compares the last modified time with the previously recorded time.
3. If the file has been modified, sends an update signal through the `configUpdateCh` channel.
4. Sleeps for a short interval before checking again.

### Note

Make sure to replace `configFilePath` with the actual path to your configuration file before running the server.

# peering-http-server
# peering-http-server
