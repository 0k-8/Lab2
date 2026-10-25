package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const (
	Host = "localhost"
	Port = "9000"
)

func main() {
	addr := net.JoinHostPort(Host, Port)

	// 1. Створення TCP-слухача
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Не вдалося запустити TCP-сервер: %v", err)
	}
	defer listener.Close()

	fmt.Printf("TCP-сервер успішно запущено на %s\n", addr)

	// 2. Головний цикл прийняття з'єднань
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Помилка прийняття підключення: %v", err)
			continue
		}

		// Обробка кожного підключення в окремій горутині
		go handleConnection(conn)
	}
}

// handleConnection обробляє з'єднання з конкретним клієнтом
func handleConnection(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("[+] Нове підключення: %s\n", remoteAddr)

	// Надсилаємо привітальне повідомлення
	conn.Write([]byte("Ласкаво просимо до TCP-сервера команд!\nВведіть команду (ECHO, TIME, QUIT/CLOSE):\n"))

	reader := bufio.NewReader(conn)

	for {
		// Читання рядка до символу переносу '\n'
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("[-] Клієнт %s відключився\n", remoteAddr)
			return
		}

		// Очищення від символів \r та \n
		cmd := strings.TrimSpace(input)
		if cmd == "" {
			continue
		}

		// Розбираємо команду та аргументи
		parts := strings.SplitN(cmd, " ", 2)
		command := strings.ToUpper(parts[0])

		switch command {
		case "ECHO":
			message := "мала пой на сок "
			if len(parts) > 1 {
				message = parts[1]
			}
			response := fmt.Sprintf("[ECHO]: %s\n", message)
			conn.Write([]byte(response))

		case "TIME":
			currentTime := time.Now().Format(time.RFC3339)
			response := fmt.Sprintf("[TIME]: %s\n", currentTime)
			conn.Write([]byte(response))

		case "CLOSE", "QUIT":
			conn.Write([]byte("До побачення!\n"))
			fmt.Printf("[-] Клієнт %s завершив сеанс\n", remoteAddr)
			return

		default:
			errorMessage := fmt.Sprintf("Помилка: невідома команда '%s'. Доступні команди: ECHO <msg>, TIME, CLOSE, QUIT\n", command)
			conn.Write([]byte(errorMessage))
		}
	}
}
