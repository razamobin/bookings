package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/razamobin/bookings/pkg/config"
	"github.com/razamobin/bookings/pkg/models"
)

var app *config.AppConfig
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	// add default data that goes to every page (as needed)
	td.StringMap["defaultdata"] = "it's the default bro"
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	var myCache map[string]*template.Template
	if app.UseCache {
		myCache = app.TemplateCache
	} else {
		var err error
		myCache, err = CreateTemplateCache()
		if err != nil {
			log.Fatal(err)
		}
	}

	t, exists := myCache[tmpl]
	if !exists {
		log.Fatal("template cache didn't have",tmpl)
	}

	td = AddDefaultData(td)

	buf := new(bytes.Buffer)
	err := t.Execute(buf, td)
	if err != nil {
		log.Fatal(err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

	//parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl, "./templates/base.layout.tmpl")
	//err = parsedTemplate.Execute(w, nil)
	//if err != nil {
	//	log.Println("error parsing template", err)
	//	return
	//}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	log.Println("creating template cache lol")

	// pre-create entire cache at once, by studying templates folder

	// *.page.tmpl first
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	// range through those pages
	for _, page := range pages {
		// get filename
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		// look for layouts now
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}