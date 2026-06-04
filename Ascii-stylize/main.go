package main

import (
	"html/template"
	"net/http"
	"stylize/ascii"
)

type PageData struct {
	Art    string
	Text   string
	Banner string
}

func GenerateArt(text, banner string) (string, error) {
	// banner is "standard" / "shadow" / "thinkertoy" → build the file path
	bannerFile := "banners/" + banner + ".txt"

	lines, err := ascii.ReadBanner(bannerFile)
	if err != nil {
		return "", err
	}
	asciiMap := ascii.BuildAsciiMap(lines)
	return ascii.PrintAscii(text, asciiMap), nil
}

func handler(w http.ResponseWriter, r *http.Request) {
    tmpl := template.Must(template.ParseFiles("index.html"))
    data := PageData{}

    if r.Method == http.MethodPost {
        text := r.FormValue("text")
        banner := r.FormValue("banner")

        valid := map[string]bool{"standard": true, "shadow": true, "thinkertoy": true}
        if !valid[banner] {
            banner = "standard"
        }

        art, err := GenerateArt(text, banner)
        if err != nil {
            http.Error(w, "Could not generate art: "+err.Error(), http.StatusInternalServerError)
            return
        }
        data.Art = art       // store everything on the SAME object...
        data.Text = text
        data.Banner = banner
    }

    tmpl.Execute(w, data)    // ...and render THAT object
}

func main() {
	http.HandleFunc("/", handler)

	// Serve CSS/static files — the gotcha I warned you about
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.ListenAndServe(":8080", nil)
}
