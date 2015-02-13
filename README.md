# Websocket Handler for Go

A simple websocket server for go that supports channels.

Usage looks like this:

```go
package main

import (
  "github.com/joncalhoun/webjack"
  "log"
  "net/http"
)

func main() {
  server := webjack.NewServer()
  http.Handle("/ws", server.GetHandler())
  http.Handle("/", http.FileServer(http.Dir("public")))
  log.Fatal(http.ListenAndServe(":3000", nil))
}
```

And then in your javascript:

```javascript
<script>
  var connection = new WebSocket("ws://localhost:3000/websockets?name=" + encodeURIComponent("your-channel-name"));
  connection.onmessage = function(e) {
    // Do stuff with msg
    var json = JSON.parse(e.data);
    console.log(json)
  }
</script>
```

If no channel is provided, the default channel of `""` *(empty string)* is used.
