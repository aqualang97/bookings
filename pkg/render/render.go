package render

import (
	"bytes"
	"github.com/aqualang97/bookings/pkg/config"
	"github.com/aqualang97/bookings/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

// NewTemplates set the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}
func AddDefaultData(td *models.TemplateDate) *models.TemplateDate {
	return td
}

// TemplateRender renders templates using http/template
func TemplateRender(w http.ResponseWriter, tmpl string, td *models.TemplateDate) {
	var tc map[string]*template.Template
	if app.UseCache {
		// create the template cache
		tc = app.TemplateCache

	} else {
		tc, _ = CreateTemplateCache()
	}
	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}
	buf := new(bytes.Buffer)

	td = AddDefaultData(td)
	_ = t.Execute(buf, td)

	//if err != nil {
	//	log.Println(err)
	//}
	// render the template

	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

// CreateTemplateCache is created cache for template
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files named *page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")

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

/*
// tc is variable that save template cache
var tc = make(map[string]*template.Template)

// TemplateRenderTest renders ...
func TemplateRender(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	// check to see if we already have the template in our cache
	_, inMap := tc[t]

	if !inMap {
		// need to create the template
		log.Println("creating template and adding to cache")
		err = createTemplateCache(t)
		if err != nil {
			log.Println(err)
		}
	} else {
		// we have the template in a cache
		log.Println("using cache template")
	}

	tmpl = tc[t]

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

// createTemplateCache is created cache for template
func createTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",
	}

	// parse the template
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	// add template to cache

	tc[t] = tmpl
	return nil
}

*/
