package main

import (
	"bufio"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"golang.org/x/term"
)

// ─── Встроенная карта ────────────────────────────────────────────
//
//go:embed map.png
var mapPNG []byte

// ─── Пароль ──────────────────────────────────────────────────────
//
// SHA-256 от пароля. Посчитать: echo -n "пароль" | sha256sum
const passwordHash = "8b4862430dfcab9bb37a104f1fcb6109c5f952e5e2af07991b4825057cced4d7"

const maxAttempts = 3

// ─── Цвета (ANSI) ────────────────────────────────────────────────
const (
	cReset  = "\033[0m"
	cDim    = "\033[2m"
	cGreen  = "\033[38;5;46m"
	cLime   = "\033[38;5;155m"
	cRed    = "\033[38;5;196m"
	cAmber  = "\033[38;5;214m"
	cGray   = "\033[38;5;245m"
	cBold   = "\033[1m"
)

// ─── Тексты ──────────────────────────────────────────────────────
const (
	appName    = "LIVECORP LOCATOR"
	appVersion = "v0.4"
	appPlace   = "Комплекс Омега · внутренний терминал"

	bannerTop    = "╔══════════════════════════════════════════════════╗"
	bannerBottom = "╚══════════════════════════════════════════════════╝"

	msgPrompt    = "Введите код доступа"
	msgGranted   = "ДОСТУП РАЗРЕШЁН"
	msgDenied    = "ДОСТУП ЗАПРЕЩЁН"
	msgAgent     = "Субъект: Агент 00"
	msgCoords    = "Координаты: 64.526374, 40.562241"
	msgStatus    = "Статус: жив. Пока."
	msgOpenFail  = "Не удалось открыть системный просмотрщик."
	msgSavedTo   = "Карта сохранена:"
	msgAttempts  = "Попытка %d из %d"
	msgMarked    = "Ваш терминал отмечен. Ожидайте."
	msgTheyCome  = "Устройство скопмпрометировано. Инициируется самоуничтожение."
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	printBanner()
	fmt.Println()

	attempts := 0
	for {
		fmt.Printf("%s%s%s %s›%s ", cDim, msgPrompt, cReset, cGreen, cReset)

		input, err := readSecret(reader)
		if err != nil {
			// не терминал (pipe) — читаем обычной строкой
			line, _ := reader.ReadString('\n')
			input = strings.TrimSpace(line)
		}
		input = strings.TrimSpace(input)
		fmt.Println()

		if hash(input) == passwordHash {
			grantAccess()
			return
		}

		attempts++
		fmt.Printf("  %s✗ %s%s %s(%s)%s\n\n",
			cRed, msgDenied, cReset,
			cGray, fmt.Sprintf(msgAttempts, attempts, maxAttempts), cReset)

		if attempts >= maxAttempts {
			denyAccess()
			return
		}
	}
}

// ─── Ввод ────────────────────────────────────────────────────────

// readSecret пытается прочитать пароль без эха. Если stdin — не терминал,
// возвращает ошибку, и вызывающий код читает строку обычным способом.
func readSecret(r *bufio.Reader) (string, error) {
	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return "", fmt.Errorf("not a terminal")
	}
	b, err := term.ReadPassword(fd)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ─── Логика ──────────────────────────────────────────────────────

func hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func grantAccess() {
	fmt.Printf("  %s✓ %s%s\n\n", cGreen, msgGranted, cReset)

	// Эффект «набора» — печатаем строки с задержкой
	typeLine(cGreen+msgAgent+cReset, 18*time.Millisecond)
	typeLine(cGreen+msgCoords+cReset, 14*time.Millisecond)
	typeLine(cAmber+msgStatus+cReset, 22*time.Millisecond)
	fmt.Println()

	time.Sleep(400 * time.Millisecond)

	// Пишем встроенный PNG во временный файл и открываем
	tmp, err := os.CreateTemp("", "locator-*.png")
	if err != nil {
		fmt.Printf("  %s!%s Не удалось создать временный файл: %v\n", cRed, cReset, err)
		return
	}
	defer tmp.Close()

	if _, err := tmp.Write(mapPNG); err != nil {
		fmt.Printf("  %s!%s Не удалось записать карту: %v\n", cRed, cReset, err)
		return
	}

	fmt.Printf("  %s%s%s %s\n\n", cDim, msgSavedTo, cReset, tmp.Name())

	if err := openFile(tmp.Name()); err != nil {
		fmt.Printf("  %s!%s %s\n", cRed, cReset, msgOpenFail)
		fmt.Printf("  %sФайл лежит здесь: %s%s\n", cGray, tmp.Name(), cReset)
	}
}

func denyAccess() {
	fmt.Println()
	typeLine(cRed+msgMarked+cReset, 40*time.Millisecond)
	time.Sleep(900 * time.Millisecond)
	typeLine(cRed+msgTheyCome+cReset, 60*time.Millisecond)
	time.Sleep(700 * time.Millisecond)
	fmt.Println()
}

// ─── Визуал ──────────────────────────────────────────────────────

// runeLen возвращает количество видимых рун в строке,
// игнорируя ANSI escape-последовательности.
func runeLen(s string) int {
	inEsc := false
	n := 0
	for _, r := range s {
		if r == '\033' {
			inEsc = true
			continue
		}
		if inEsc {
			if r == 'm' {
				inEsc = false
			}
			continue
		}
		n++
	}
	return n
}

// padRight дополняет строку пробелами до ширины w (по видимым рунам).
func padRight(s string, w int) string {
	n := runeLen(s)
	if n >= w {
		return s
	}
	return s + strings.Repeat(" ", w-n)
}

// frameLine собирает строку внутри рамки заданной внутренней ширины.
func frameLine(content string, inner int) string {
	return cGreen + "║" + cReset + padRight("  "+content, inner) + cGreen + "║" + cReset
}

func printBanner() {
	const inner = 50 // видимая ширина между ║ и ║

	top := cGreen + "╔" + strings.Repeat("═", inner) + "╗" + cReset
	bot := cGreen + "╚" + strings.Repeat("═", inner) + "╝" + cReset

	title := fmt.Sprintf("%s%s%s %s%s%s",
		cBold, appName, cReset, cDim, appVersion, cReset)
	place := fmt.Sprintf("%s%s%s", cGray, appPlace, cReset)

	lines := []string{
		top,
		frameLine(title, inner),
		frameLine(place, inner),
		bot,
	}

	for _, l := range lines {
		fmt.Println(l)
		time.Sleep(60 * time.Millisecond)
	}
}

// typeLine печатает строку по символам с заданной задержкой.
func typeLine(s string, d time.Duration) {
	for _, r := range s {
		fmt.Printf("%c", r)
		time.Sleep(d)
	}
	fmt.Println()
}

// openFile открывает файл системным просмотрщиком.
func openFile(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", path)
	default:
		// Linux/BSD: пробуем xdg-open, потом常见的 просмотрщики
		if _, err := exec.LookPath("xdg-open"); err == nil {
			cmd = exec.Command("xdg-open", path)
		} else if _, err := exec.LookPath("feh"); err == nil {
			cmd = exec.Command("feh", path)
		} else if _, err := exec.LookPath("eog"); err == nil {
			cmd = exec.Command("eog", path)
		} else {
			return fmt.Errorf("no image viewer found")
		}
	}
	return cmd.Start()
}

// на случай, если понадобится сохранить карту рядом с бинарником:
var _ = filepath.Join
