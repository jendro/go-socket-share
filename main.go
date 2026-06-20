package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func newHub() *Hub {
	return &Hub{clients: make(map[*websocket.Conn]bool)}
}

func (h *Hub) broadcast(message string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for conn := range h.clients {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Printf("broadcast write error: %v", err)
			conn.Close()
			delete(h.clients, conn)
		}
	}
}

func (h *Hub) addClient(conn *websocket.Conn) {
	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()
}

func (h *Hub) removeClient(conn *websocket.Conn) {
	h.mu.Lock()
	if _, ok := h.clients[conn]; ok {
		conn.Close()
		delete(h.clients, conn)
	}
	h.mu.Unlock()
}

func main() {
	flag.Parse()
	hub := newHub()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, indexHTML)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("upgrade: %v", err)
			return
		}
		hub.addClient(conn)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("read: %v", err)
				hub.removeClient(conn)
				break
			}
			message := string(msg)
			if message == "" {
				continue
			}
			hub.broadcast(message)
		}
	})

	log.Printf("starting server at %s", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

const indexHTML = `<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>Shared JSON</title>
<style>
body { font-family: system-ui, sans-serif; margin: 0; padding: 0; background: #f7f7f7; }
.container { max-width: 760px; margin: 2rem auto; padding: 1rem; background: #fff; border-radius: 12px; box-shadow: 0 10px 30px rgba(0,0,0,0.08); }
h1 { margin-top: 0; }
textarea { width: 100%; min-height: 150px; padding: 12px; border: 1px solid #ccc; border-radius: 8px; font-family: monospace; font-size: 14px; resize: vertical; }
#messages { margin-top: 1rem; display: grid; gap: 10px; }
.message { padding: 1rem; border: 1px solid #e2e8f0; border-radius: 10px; background: #f9fafb; display: flex; justify-content: space-between; gap: 0.75rem; align-items: flex-start; }
.message pre { margin: 0; white-space: pre-wrap; word-break: break-word; flex: 1; }
button { background: #2563eb; color: white; border: none; border-radius: 8px; padding: 0.6rem 0.9rem; cursor: pointer; }
button:disabled { opacity: 0.5; cursor: default; }
.status { margin-top: 0.75rem; color: #555; }
</style>
</head>
<body>
<div class="container">
<h1>Shared JSON</h1>
<p>Tempel JSON ke textarea, tekan Enter untuk berbagi ke semua pengguna yang membuka halaman ini.</p>
<textarea id="input" placeholder="Tempel JSON lalu tekan Enter..."></textarea>
<div class="status" id="status">Menunggu koneksi WebSocket...</div>
<div id="messages"></div>
</div>
<script>
const input = document.getElementById('input');
const status = document.getElementById('status');
const messages = document.getElementById('messages');
let socket;

function setStatus(text) {
  status.textContent = text;
}

function connect() {
  socket = new WebSocket((location.protocol === 'https:' ? 'wss://' : 'ws://') + location.host + '/ws');
  socket.addEventListener('open', () => setStatus('Terhubung. Ketik JSON lalu Enter untuk membagikan.'));
  socket.addEventListener('close', () => setStatus('Terputus. Memuat ulang dalam 2 detik...') || setTimeout(connect, 2000));
  socket.addEventListener('error', () => setStatus('Gagal terhubung. Memuat ulang...') );
  socket.addEventListener('message', event => addMessage(event.data));
}

function addMessage(text) {
  const wrapper = document.createElement('div');
  wrapper.className = 'message';

  const pre = document.createElement('pre');
  pre.textContent = text;

  const copyButton = document.createElement('button');
  copyButton.textContent = 'Copy';
  copyButton.type = 'button';
  copyButton.addEventListener('click', async () => {
    try {
      await navigator.clipboard.writeText(text);
      copyButton.textContent = 'Copied';
      setTimeout(() => { copyButton.textContent = 'Copy'; }, 1200);
    } catch (err) {
      console.error(err);
      copyButton.textContent = 'Gagal';
      setTimeout(() => { copyButton.textContent = 'Copy'; }, 1500);
    }
  });

  wrapper.appendChild(pre);
  wrapper.appendChild(copyButton);
  messages.prepend(wrapper);
}

input.addEventListener('keydown', event => {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault();
    const text = input.value.trim();
    if (!text || socket.readyState !== WebSocket.OPEN) {
      return;
    }
    socket.send(text);
    input.value = '';
  }
});

connect();
</script>
</body>
</html>`
