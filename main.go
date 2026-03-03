package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "file:series.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("Servidor en http://localhost:8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleClient(conn, db)
	}
}

func handleClient(conn net.Conn, db *sql.DB) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	conn.Read(buffer)
	request := string(buffer)

	if !strings.HasPrefix(request, "GET / ") {
		return
	}

		rows, err := db.Query("SELECT id, name, current_episode, total_episodes FROM series")
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	
	html := `
	<html>
	<head>
	<title>My Series Tracker</title>
	<style>
	body { 
		font-family: Arial; 
		background-color: #f4f4f4; 
	}
	table { border-collapse: collapse; width: 60%; margin: auto; background: white; }
	th, td { border: 1px solid #ccc; padding: 10px; text-align: center; }
	th { background-color: #333; color: white; }
	h1 { text-align: center; }
	button { display: block; margin: 15px auto; padding: 10px; }
	</style>
	</head>
	<body>

	<h1>My Series Tracker</h1>

	<button onclick="toggleChaos()">MODO CAOS</button>

	<table>
	<tr>
	<th>#</th>
	<th>Name</th>
	<th>Current</th>
	<th>Total</th>
	</tr>
	`

	
	for rows.Next() {
		var id int
		var name string
		var current int
		var total int

		err := rows.Scan(&id, &name, &current, &total)
		if err != nil {
			log.Println(err)
			continue
		}

		html += fmt.Sprintf(
			"<tr><td>%d</td><td>%s</td><td>%d</td><td>%d</td></tr>",
			id, name, current, total,
		)
	}

	html += `
	</table>

	<script>
	let chaosInterval;

	function toggleChaos() {
		if (chaosInterval) {
			clearInterval(chaosInterval);
			chaosInterval = null;
			document.body.style.backgroundColor = "#f4f4f4";
			return;
		}

		chaosInterval = setInterval(() => {
			const randomColor = "#" + Math.floor(Math.random()*16777215).toString(16);
			document.body.style.backgroundColor = randomColor;
		}, 100); // cambia cada 100ms (ultra rápido)
	}
	</script>

	</body>
	</html>
	`

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/html\r\n" +
		"\r\n" +
		html

	conn.Write([]byte(response))
}