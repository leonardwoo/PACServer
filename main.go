package main

import (
  "os"
  "fmt"
	"net/http"
	"log"
  "bytes"
  "bufio"
  "encoding/base64"
  "text/template"
)

const (
  SERVER_ADDR string := ""
  SERVER_PORT int := 8880

  REPLACE_ALL int := -1

  PROXY string := "PROXY 192.168.1.1:1080"
  RULESLISTADDR string := "https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt"
)

var RULES string := ""

type rulestmpl struct {
	/* proxy address */
	Proxy string

	/* gfwlist domains name */
	Domains map[string]int

	/* custom domains name */
	Custom map[string]int
}

func createRulesContent()  {
  //convert gfwlist to pac rules.

  if res, err := http.Get(RULESLISTADDR); err != nil {
		log.Fatal("GET error:", err)
    return
	}
  if decoded, err := base64.NewDecoder(res.Body); err != nil {
    log.Fatal("decode error:", err)
		return
	}
  var domains map[string]int = map[string]int{}
	reader := bufio.NewReader(decoder)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				log.Print(err)
			}
			break
		}

		s := parse(string(line))
		if s != "" {
			if _, ok := domains[s]; !ok {
				domains[s] = 1
			}
		}
	}
	res.Body.Close()

	if tmpl, err := template.ParseFiles("pac-rules.tmpl"); err != nil {
		log.Fatal(err)
	}

  bout := bytes.NewBuffer(make([]byte, 0))
  out := bufio.NewWriter(bout)
	if err = t.Execute(out, aa); err != nil {
		log.Fatal(err)
	}
  out.Flush()
  RULES = bout.toString()
}

func pacHandler(writer http.ResponseWriter, request *http.Request) {
  content string := ""
  if fin, err := os.Open("abs.js"); err != nil {
    log.Fatal("failed to load page", err)
    content = "failed to load page"
  } else {
    buf := make([]byte, 1024)
    for{
      n, _ := fin.Read(buf)
      if 0 == n { break }
      content += string(buf[:n])
    }

    strings.Replace(content, "${PROXY}", PROXY, REPLACE_ALL)

    createRulesContent()
    strings.Replace(content, "${RULES}", RULES, REPLACE_ALL)
  }
  fmt.Fprint(writer, content)
}

func homeHandler(writer http.ResponseWriter, request *http.Request) {
   fmt.Fprint(writer, "")
}

func main() {
  http.HandleFunc("/proxy.pac", pacHandler)
  http.HandleFunc("/", homeHandler)
  if err:= http.ListenAndServe(SERVER_ADDR + ":" + SERVER_PORT, nil); err != nil {
  	log.Fatal("failed to start server", err)
  }
}
