package main

import (
	ascii "ascii/app"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type ExecOutput struct {
	In  string
	Out string
}

func ValidAscii(s string) bool {
	for _, i := range []byte(s) {
		if i > 127 {
			return false
		}
	}
	return true
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	t, _ := template.ParseFiles("error/500.html")
	err := t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		switch r.Method {
		case "GET":
			t, err := template.ParseFiles("index.html")
			if err != nil {
				internalServerError(w, r)
			}
			t.Execute(w, nil)
		case "POST":
			r.ParseForm()
			if !ValidAscii(r.Form.Get("input")) {
				w.WriteHeader(http.StatusBadRequest)
				t, err := template.ParseFiles("error/400.html")
				if err != nil {
					internalServerError(w, r)
				}
				t.Execute(w, nil)
			} else {
				output, status := ascii.AsciiOutput(r.Form["input"][0], r.Form["font"][0])
				log.Printf("method: %v / font: %v / input: %v / statuscode: %v\n", r.Method, r.Form["font"][0], r.Form["input"][0], status)
				if status == 500 {
					internalServerError(w, r)
				} else {
					ex := ExecOutput{
						In:  r.Form["input"][0],
						Out: output,
					}
					t, err := template.ParseFiles("index.html")
					if err != nil {
						internalServerError(w, r)
						return
					}
					t.Execute(w, ex)
				}
			}
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		t, err := template.ParseFiles("error/404.html")
		if err != nil {
			internalServerError(w, r)
			return
		}
		t.Execute(w, nil)
	}
}

// func Openbrowser(zz string) {
// 	var err error
// 	switch runtime.GOOS {
// 	case "windows":
// 		err = exec.Command("rundll132", "url.dll,FileProtocolHandler", zz).Start()

// 	case "linux":
// 		err = exec.Command("xdg-open", zz).Start()
// 	}
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

func main() {
	log.Println("server is starting...")
	http.HandleFunc("/", Handler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("Server up at port 8080\nhttp status :", http.StatusOK)
	//Openbrowser("http.localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
