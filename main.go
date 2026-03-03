package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"
	"net/url"
	"strconv"

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

// Maneja la conexión de cada cliente. Es como un router que solo decide qué función llamar según la ruta y el método HTTP.
func handleClient(conn net.Conn, db *sql.DB) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	conn.Read(buffer)
	request := string(buffer)

	lines := strings.Split(request, "\r\n")
	requestLine := strings.Split(lines[0], " ")

	method := requestLine[0]
	path := requestLine[1]

	if method == "GET" && path == "/" {
		serveHome(conn, db)
		return
	}

	if method == "GET" && path == "/create" {
		serveCreateForm(conn)
		return
	}

	if method == "POST" && path == "/create" {
		handleCreatePost(conn, request, db)
		return
	}

	if method == "POST" && strings.HasPrefix(path, "/update") {
	handleUpdate(conn, path, db)
	return
	}
}

// Muestra el formulario para agregar una nueva serie. Es una página HTML simple con un formulario que envía los datos al servidor.

func serveCreateForm(conn net.Conn) {

	html := `
	<html>
	<body>
	<h1>Agregar una nueva serie</h1>

	<form method="POST" action="/create">

		<label>Nombre:</label><br>
		<input type="text" name="series_name" required><br><br>

		<label>Capitulo Actual:</label><br>
		<input type="number" name="current_episode" min="1" value="1" required><br><br>

		<label>Total de Capitulos:</label><br>
		<input type="number" name="total_episodes" min="1" required><br><br>

		<button type="submit">Agregar</button>

	</form>

	<br>
	<a href="/">Volver a Inicio</a>

	</body>
	</html>
	`

	response := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n" + html
	conn.Write([]byte(response))
}

// Muestra la página principal con la lista de series. Consulta la base de datos para obtener las series y genera una tabla HTML para mostrarlas.
func serveHome(conn net.Conn, db *sql.DB) {

	rows, err := db.Query("SELECT id, name, current_episode, total_episodes FROM series")
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	html := `
	<html>
	<head>
	<title>Tracker de Series</title>

	<script>
	async function nextEpisode(id) {
		const url = "/update?id=" + id
		await fetch(url, { method: "POST" })
		location.reload()
	}
	</script>

	</head>
	<body>

	<h1>Tracker de Series</h1>

	<a href="/create">Agregar nueva serie</a>

	<table border="1">
	<tr>
	<th>#</th>
	<th>Name</th>
	<th>Current</th>
	<th>Total</th>
	<th>Actions</th>
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
			"<tr><td>%d</td><td>%s</td><td>%d</td><td>%d</td><td><button onclick='nextEpisode(%d)'>+1</button></td></tr>",
			id, name, current, total, id,
		)
	}

	html += `
	</table>
	</body>
	</html>
	`

	response := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n" + html
	conn.Write([]byte(response))
}

// Maneja la solicitud POST para crear una nueva serie. Lee el cuerpo de la solicitud, extrae los datos del formulario, los inserta en la base de datos y redirige al usuario a la página principal.

func handleCreatePost(conn net.Conn, request string, db *sql.DB) {

	parts := strings.SplitN(request, "\r\n\r\n", 2)
	headers := parts[0]
	body := parts[1]

	var contentLength int
	lines := strings.Split(headers, "\r\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "Content-Length:") {
			lengthStr := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
			contentLength, _ = strconv.Atoi(lengthStr)
		}
	}

	body = body[:contentLength]

	values, _ := url.ParseQuery(body)

	name := values.Get("series_name")
	current := values.Get("current_episode")
	total := values.Get("total_episodes")

	_, err := db.Exec(
	"INSERT INTO series (name, current_episode, total_episodes) VALUES (?, ?, ?)",
	name, current, total,
	)

	if err != nil {
		log.Println(err)
	}

	response := "HTTP/1.1 303 See Other\r\n" +
		"Location: /\r\n\r\n"

	conn.Write([]byte(response))
}

// Maneja la solicitud POST para actualizar el episodio actual de una serie. Extrae el ID de la serie de la URL, incrementa el episodio actual en la base de datos (si no ha llegado al total) y responde con un mensaje simple.

func handleUpdate(conn net.Conn, path string, db *sql.DB) {

	parts := strings.SplitN(path, "?", 2)

	var id string

	if len(parts) > 1 {
		params, _ := url.ParseQuery(parts[1])
		id = params.Get("id")
	}

	_, err := db.Exec(`
		UPDATE series
		SET current_episode = current_episode + 1
		WHERE id = ? AND current_episode < total_episodes
	`, id)

	if err != nil {
		log.Println(err)
	}

	response := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nok"
	conn.Write([]byte(response))
}