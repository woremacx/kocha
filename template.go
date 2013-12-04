package kocha

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

var (
	TemplateFuncs = template.FuncMap{
		"eq": func(a, b interface{}) bool {
			// TODO: remove in Go 1.2
			//       see http://tip.golang.org/pkg/text/template/#hdr-Functions
			return a == b
		},
		"ne": func(a, b interface{}) bool {
			// TODO: remove in Go 1.2
			//       see http://tip.golang.org/pkg/text/template/#hdr-Functions
			return a != b
		},
		"in": func(a, b interface{}) bool {
			v := reflect.ValueOf(a)
			switch v.Kind() {
			case reflect.Slice, reflect.Array, reflect.String:
				if v.IsNil() {
					return false
				}
				for i := 0; i < v.Len(); i++ {
					if v.Index(i).Interface() == b {
						return true
					}
				}
			default:
				panic(fmt.Errorf("invalid type %v: valid types are slice, array and string", v.Type().Name()))
			}
			return false
		},
		"url": Reverse,
		"nl2br": func(text string) template.HTML {
			return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
		},
		"raw": func(text string) template.HTML {
			return template.HTML(text)
		},
		"date": func(date time.Time, layout string) string {
			return date.Format(layout)
		},
	}
)

type TemplateSet map[string]AppTemplateSet
type AppTemplateSet map[string]LayoutTemplateSet
type LayoutTemplateSet map[string]FileExtTemplateSet
type FileExtTemplateSet map[string]*template.Template

// Get gets a parsed template.
func (t TemplateSet) Get(appName, layoutName, name, format string) *template.Template {
	return t[appName][layoutName][format][ToSnakeCase(name)]
}

func (t TemplateSet) Ident(appName, layoutName, name, format string) string {
	return fmt.Sprintf("%s:%s %s.%s", appName, layoutName, ToSnakeCase(name), format)
}

// TemplateSetFromPaths returns TemplateSet constructed from templateSetPaths.
func TemplateSetFromPaths(templateSetPaths map[string][]string) TemplateSet {
	layoutPaths := make(map[string]map[string]map[string]string)
	templatePaths := make(map[string]map[string]map[string]string)
	templateSet := make(TemplateSet)
	for appName, paths := range templateSetPaths {
		layoutPaths[appName] = make(map[string]map[string]string)
		templatePaths[appName] = make(map[string]map[string]string)
		for _, rootPath := range paths {
			layoutDir := filepath.Join(rootPath, "layouts")
			if err := collectLayoutPaths(layoutPaths[appName], layoutDir); err != nil {
				panic(err)
			}
			if err := collectTemplatePaths(templatePaths[appName], rootPath, layoutDir); err != nil {
				panic(err)
			}
		}
		templateSet[appName] = make(AppTemplateSet)
	}
	for appName, templates := range templatePaths {
		if err := buildSingleAppTemplateSet(templateSet[appName], templates); err != nil {
			panic(err)
		}
	}
	for appName, layouts := range layoutPaths {
		if err := buildLayoutAppTemplateSet(templateSet[appName], layouts, templatePaths[appName]); err != nil {
			panic(err)
		}
	}
	return templateSet
}

func collectLayoutPaths(layoutPaths map[string]map[string]string, layoutDir string) error {
	return filepath.Walk(layoutDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		baseName, err := filepath.Rel(layoutDir, path)
		if err != nil {
			return err
		}
		name, ext := SplitExt(baseName)
		if _, exists := layoutPaths[name]; !exists {
			layoutPaths[name] = make(map[string]string)
		}
		if layoutPath, exists := layoutPaths[name][ext]; exists {
			return fmt.Errorf("duplicate name of layout file:\n  1. %s\n  2. %s\n", layoutPath, path)
		}
		layoutPaths[name][ext] = path
		return nil
	})
}

func collectTemplatePaths(templatePaths map[string]map[string]string, templateDir, excludeDir string) error {
	return filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if path == excludeDir {
				return filepath.SkipDir
			}
			return nil
		}
		baseName, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}
		name, ext := SplitExt(baseName)
		if _, exists := templatePaths[ext]; !exists {
			templatePaths[ext] = make(map[string]string)
		}
		if templatePath, exists := templatePaths[ext][name]; exists {
			return fmt.Errorf("duplicate name of template file:\n  1. %s\n  2. %s\n", templatePath, path)
		}
		templatePaths[ext][name] = path
		return nil
	})
}

func buildSingleAppTemplateSet(appTemplateSet AppTemplateSet, templates map[string]map[string]string) error {
	layoutTemplateSet := make(LayoutTemplateSet)
	for ext, templateInfos := range templates {
		layoutTemplateSet[ext] = make(FileExtTemplateSet)
		for name, path := range templateInfos {
			templateBytes, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			t := template.Must(template.New(name).Funcs(TemplateFuncs).Parse(string(templateBytes)))
			layoutTemplateSet[ext][name] = t
		}
	}
	appTemplateSet[""] = layoutTemplateSet
	return nil
}

func buildLayoutAppTemplateSet(appTemplateSet AppTemplateSet, layouts map[string]map[string]string, templates map[string]map[string]string) error {
	for layoutName, layoutInfos := range layouts {
		layoutTemplateSet := make(LayoutTemplateSet)
		for ext, layoutPath := range layoutInfos {
			layoutTemplateSet[ext] = make(FileExtTemplateSet)
			layoutBytes, err := ioutil.ReadFile(layoutPath)
			if err != nil {
				return err
			}
			for name, path := range templates[ext] {
				// do not use the layoutTemplate.Clone() in order to retrieve layout as string by `kocha build`
				layout := template.Must(template.New("layout").Funcs(TemplateFuncs).Parse(string(layoutBytes)))
				t := template.Must(layout.ParseFiles(path))
				layoutTemplateSet[ext][name] = t
			}
		}
		appTemplateSet[layoutName] = layoutTemplateSet
	}
	return nil
}
