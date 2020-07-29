package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
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
	choices = map[string]string{
		"A": "https://docs.google.com/forms/d/1_T6Dyg8PGbiJKwGFBT-RDFRnlkcPDU5wFDDbOaLA1aU",
		"B": "https://docs.google.com/forms/d/1x89NQ06-Q2DCis8FDHQN7W6C3CNicDp1oSpvrXbTD1w",
		"C": "https://docs.google.com/forms/d/1eOZSuWHZ51ytM_ZLKQ7QkXYSrK_e-gM5xYc17oJvrpQ",
	}
	random = true
	manual = "A"

	values = make([]string, 0, len(choices))
	for k, v := range choices {
		fmt.Println(k)
		values = append(values, v)
	}

	fmt.Println(choices)
	http.HandleFunc("/", pickRandom)
	http.HandleFunc("/mode", switchMode)
	http.HandleFunc("/check", checkMode)
	http.HandleFunc("/change/", changeCurrent)

	log.Fatal(http.ListenAndServe(":5000", nil))
}
