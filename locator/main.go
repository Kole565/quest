package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// SHA-256 от пароля.
const passwordHash = "c4b21ced8d989b4427f732c8f74b8db22689997fe186356ff38f890cb166fc04"

const (
	imageName   = "map.png"
	maxAttempts = 3
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║  Agent0152 LOCATOR v0.4                  ║")
	fmt.Println("║  Проект 'Протокол мертвой руки'          ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()

	attempts := 0
	for {
		fmt.Print("Введите код доступа: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if hash(input) == passwordHash {
			grantAccess()
			return
		}

		attempts++
		fmt.Printf("Доступ запрещён. Попытка %d из %d.\n\n", attempts, maxAttempts)

		if attempts >= maxAttempts {
			denyAccess()
			return
		}
	}
}

func hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func grantAccess() {
	fmt.Println()
	fmt.Println("Доступ разрешён.")
	fmt.Println("Субъект: Агент 00")
	fmt.Println("Координаты: 64.526374° N, 40.562241° E")
	fmt.Println("Статус: жив. Пока.")
	fmt.Println()

	// Ищем картинку рядом с бинарником
	exe, err := os.Executable()
	if err != nil {
		fmt.Println("Ошибка: не удалось определить путь к приложению.")
		return
	}
	dir := filepath.Dir(exe)
	imgPath := filepath.Join(dir, imageName)

	if _, err := os.Stat(imgPath); err != nil {
		fmt.Printf("Ошибка: файл %s не найден рядом с приложением.\n", imageName)
		return
	}

	if err := openImage(imgPath); err != nil {
		fmt.Printf("Не удалось открыть %s: %v\n", imageName, err)
		fmt.Printf("Файл лежит здесь: %s\n", imgPath)
	}
}

func denyAccess() {
	fmt.Println()
	fmt.Println("Три неудачные попытки.")
	fmt.Println("Ваш терминал отмечен. Ожидайте.")
	time.Sleep(2 * time.Second)
	fmt.Println("Ожидайте.")
	time.Sleep(2 * time.Second)
	fmt.Println("Пароль спрятан в книге.")
}

// openImage открывает png системным просмотрщиком.
func openImage(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", path)
	default: // linux, *bsd
		cmd = exec.Command("xdg-open", path)
	}

	return cmd.Start()
}
