package testserver

import (
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// addr = flag.String("addr", "127.0.0.1:2003", "http service address")
	// addr = "127.0.0.1:2003"

	// use default options
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// template frontDemo
	homeTemplate = template.Must(template.New("").Parse(`
	<!DOCTYPE html>
	<head>
	<meta charset="utf-8">
	<script>
	window.addEventListener("load", function(evt) {
		
		var output = document.getElementById("output");
		var input = document.getElementById("input");
		var ws;

		var print = function(message) {
			var d = document.createElement("div");
			d.innerHTML = message;
			output.appendChild(d);
		};

		document.getElementById("open").onclick = function(evt) {
			if (ws) {
				return false;
			}
			var web_url=document.getElementById("web_url").value
			ws = new WebSocket(web_url);
			ws.onopen = function(evt) {
				print("OPEN");
			}
			ws.onclose = function(evt) {
				print("CLOSE");
				ws = null;
			}
			ws.onmessage = function(evt) {
				print("RESPONSE: " + evt.data);
			}
			ws.onerror = function(evt) {
				print("ERROR: " + evt.data);
			}
			return false;
		};

		document.getElementById("send").onclick = function(evt) {
			if (!ws) {
				return false
			}
			print("SEND: " + input.value);
			ws.send(input.value);
			return false
		};

		document.getElementById("close").onclick = function(evt) {
			if (!ws) {
				return false;
			}
			ws.close();
			return false;
		};
	});
	</script>
	</head>
	<body>
	<table>
	<tr><td valign="top" width="60%">
	<p>Click "Open" to create a connection to the server,
	"send" to send a message to the server and "Close" to close the connection.
	You can change the message and send multiple times.
	<p>
	<form>
	<button id="open">Open</button>
	<button id="close">Close</button>
	<p><input id="web_url" type="text" value="{{.}}">
	<p><input id="input" type="text" value="Hello world!">
	<button id="send">Send</button>
	</form>
	</td><td valign="top" width="40%">
	<div id="output"></div>
	</td></tr></table>
	</body>
	</html>
	`))
)

type RealServer struct {
	Addr string
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()

	//
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read : ", err)
			break
		}
		log.Printf("recv: %s\n", message)

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write : ", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func (r *RealServer) Run() {
	// flag.Parse()
	// log.SetFlags(0)
	mux := http.NewServeMux()
	mux.HandleFunc("/echo", echo)
	mux.HandleFunc("/", home)
	server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: 3 * time.Second,
		Handler:      mux,
	}
	// log.Println("Starting websocket server at " + *addr)
	// log.Fatal(http.ListenAndServe(*addr, nil))
	log.Println("Starting websocket server at " + r.Addr)
	go func() {
		log.Fatal(server.ListenAndServe())
	}()
}

// func main() {

// }
