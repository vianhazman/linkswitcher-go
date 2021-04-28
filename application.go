package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var choices map[string]string
var values []string
var random bool
var manual string

func pickRandom(w http.ResponseWriter, r *http.Request) {
	var selected string
	if random {
		selected = values[rand.Intn(len(values))]
	} else {
		selected = choices[manual]
	}
	fmt.Println("Random: " + strconv.FormatBool(random) + " - " + selected)
	fmt.Println(selected)
	http.Redirect(w, r, selected, 303)
}

func switchMode(w http.ResponseWriter, r *http.Request) {
	if random {
		random = false
	} else {
		random = true
	}
	fmt.Fprintf(w, "Mode changed! Random: "+strconv.FormatBool(random))
}

func changeCurrent(w http.ResponseWriter, r *http.Request) {
	if random {
		fmt.Fprintf(w, `Change mode to manual first`)
	} else {
		manual = strings.ToUpper(strings.TrimPrefix(r.URL.Path, "/change/"))
		fmt.Fprintf(w, manual)
	}
}

func checkMode(w http.ResponseWriter, r *http.Request) {
	if random {
		fmt.Fprintf(w, `Current mode is random`)

	} else {
		fmt.Fprintf(w, `Current mode is manual with SET: `+manual)

	}
}

func main() {
	getenvironment := func(data []string, getkeyval func(item string) (key, val string)) map[string]string {
		items := make(map[string]string)
		for _, item := range data {
			key, val := getkeyval(item)
			items[key] = val
		}
		return items
	}
	choices := getenvironment(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = splits[1]
		return
	})
	random = true
	manual = "A"

	values = make([]string, 0, len(choices))
	for _, v := range choices {
		values = append(values, v)
	}

	fmt.Println(choices)
	http.HandleFunc("/", pickRandom)
	http.HandleFunc("/mode", switchMode)
	http.HandleFunc("/check", checkMode)
	http.HandleFunc("/change/", changeCurrent)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
