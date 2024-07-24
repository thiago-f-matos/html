package html

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// Titulo recebe uma lista de URLs e retorna um canal de strings contendo os títulos das páginas.
func Titulo(urls ...string) <-chan string {
	c := make(chan string)
	for _, url := range urls {
		go func(url string) {
			resp, err := http.Get(url)
			if err != nil {
				c <- fmt.Sprintf("erro ao acessar %s: %v", url, err)
				return
			}
			defer resp.Body.Close()

			html, err := io.ReadAll(resp.Body)
			if err != nil {
				c <- fmt.Sprintf("erro ao ler corpo da resposta de %s: %v", url, err)
				return
			}

			// Compila a expressão regular para extrair o título da página
			r, err := regexp.Compile(`<title>(.*?)</title>`)
			if err != nil {
				c <- fmt.Sprintf("erro ao compilar regexp para %s: %v", url, err)
				return
			}

			// Encontra o título na página HTML
			matches := r.FindStringSubmatch(string(html))
			if len(matches) > 1 {
				c <- strings.TrimSpace(matches[1])
			} else {
				c <- fmt.Sprintf("título não encontrado em %s", url)
			}
		}(url)
	}
	return c
}
