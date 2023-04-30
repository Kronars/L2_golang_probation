package main

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	url := "https://google.com"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to fetch URL: %s\n", err)
		return
	}
	defer resp.Body.Close()

	// Создание директории для сохранения
	directory := strings.Split(url, "//")[1]
	if err := os.Mkdir(directory, os.ModePerm); err != nil {
		fmt.Printf("Failed to create directory: %s\n", err)
		return
	}

	// Сохранение
	htmlFile := filepath.Join(directory, "index.html")
	f, err := os.Create(htmlFile)
	if err != nil {
		fmt.Printf("Failed to create file: %s\n", err)
		return
	}
	defer f.Close()

	// Разбор ответа
	if _, err := io.Copy(f, resp.Body); err != nil {
		fmt.Printf("Failed to save file: %s\n", err)
		return
	}

	// Загрузка ресурсов страницы
	for _, link := range getLinks(htmlFile) {
		downloadFile(link, directory)
	}
}

func getLinks(htmlFile string) []string {
	f, err := os.Open(htmlFile)
	if err != nil {
		fmt.Printf("Failed to open file: %s\n", err)
		return nil
	}
	defer f.Close()

	var links []string
	var buf [4096]byte
	for {
		n, err := f.Read(buf[:])
		if err != nil && err != io.EOF {
			fmt.Printf("Failed to read file: %s\n", err)
			return nil
		}
		if n == 0 {
			break
		}
		links = append(links, getLinksFromText(string(buf[:n]))...)
	}
	return links
}

func getLinksFromText(text string) []string {
	var links []string
	for _, href := range strings.Split(text, "href=\"")[1:] {
		url := strings.Split(href, "\"")[0]
		if strings.HasPrefix(url, "http") || strings.HasPrefix(url, "https") {
			links = append(links, url)
		}
	}
	return links
}

func downloadFile(url, directory string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to fetch URL: %s\n", err)
		return
	}
	defer resp.Body.Close()

	filename := path.Base(url)
	filepath := filepath.Join(directory, filename)
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("Error creating file %s: %s\n", filepath, err)
		return
	}
	defer file.Close()

	// Запись документа в файл
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("Error writing file %s: %s\n", filepath, err)
		return
	}

	fmt.Printf("Downloaded %s to %s\n", url, filepath)
}
