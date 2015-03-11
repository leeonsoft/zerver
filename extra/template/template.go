package template

import (
	htmpl "html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cosiner/zerver"
)

func _tmpSetTmplDelims(left, right string) {
	zerver.TmpHSet("tmpl", "delims_left", left)
	zerver.TmpHSet("tmpl", "delims_right", right)
}

func _tmpGetTmplDelims() (string, string) {
	left := zerver.TmpHGet("tmpl", "delims_left")
	right := zerver.TmpHGet("tmpl", "delims_right")
	if left == nil || right == nil {
		return "{{", "}}"
	}
	return left.(string), right.(string)
}

func _tmpAddTmpl(name string) {
	zerver.TmpHSet("tmpl", "tmpls", append(_tmpTmpls(), name))
}

func _tmpTmpls() []string {
	var tmpls []string
	if ts := zerver.TmpHGet("tmpl", "tmpls"); ts != nil {
		tmpls = ts.([]string)
	}
	return tmpls
}

type (
	// TemplateEngine is a template engine which support set delimiters, add file suffix name,
	// add template file and dirs, add template functions, compile template, and
	// render template
	TemplateEngine interface {
		SetTemplateDelims(left, right string)
		AddTemplateSuffix(s []string)
		AddTemplates(path ...string) error
		AddTemplateFunc(name string, fn interface{})
		AddTemplateFuncs(funcs map[string]interface{})
		CompileTemplates() error
		RenderTemplate(wr io.Writer, name string, value interface{}) error
	}

	// template implements TemplateEngine interface use standard html/template package
	template struct {
		tmpl *htmpl.Template
	}
)

var (
	// GlobalTmplFuncs is the default template functions
	GlobalTmplFuncs = map[string]interface{}{
	// "I18N": I18N,
	}
	// tmplSuffixes is all template file's suffix
	tmplSuffixes = map[string]bool{"tmpl": true, "html": true}
)

// NewTemplateEngine create a new template engine
func NewTemplateEngine() TemplateEngine {
	return new(template)
}

// isTemplate check whether a file name is recognized template file
func (*template) isTemplate(name string) (is bool) {
	index := strings.LastIndex(name, ".")
	if is = (index >= 0); is {
		is = tmplSuffixes[name[index+1:]]
	}
	return
}

// AddTemplateSuffix add suffix for template
func (*template) AddTemplateSuffix(suffixes []string) {
	for _, suffix := range suffixes {
		if suffix != "" {
			if suffix[0] == '.' {
				suffix = suffix[1:]
			}
			tmplSuffixes[suffix] = true
		}
	}
}

// SetTemplateDelims set default template delimeters
func (*template) SetTemplateDelims(left, right string) {
	_tmpSetTmplDelims(left, right)
}

// AddTemplates add templates to server, all templates will be
// parsed on server start
func (t *template) AddTemplates(names ...string) (err error) {
	addTmpl := func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && t.isTemplate(path) {
			_tmpAddTmpl(path)
		}
		return err
	}
	for _, name := range names {
		if err = filepath.Walk(name, addTmpl); err != nil {
			break
		}
	}
	return
}

// CompileTemplates compile all added templates
func (t *template) CompileTemplates() (err error) {
	if tmpls := _tmpTmpls(); len(tmpls) != 0 {
		var tmpl *htmpl.Template
		tmpl, err = htmpl.New("tmpl").
			Delims(_tmpGetTmplDelims()).
			Funcs(GlobalTmplFuncs).
			ParseFiles(tmpls...)
		if err == nil {
			t.tmpl = tmpl
		}
	}
	return
}

// AddTemplateFunc register a function used in templates
func (*template) AddTemplateFunc(name string, fn interface{}) {
	GlobalTmplFuncs[name] = fn
}

// AddTemplateFuncs register some functions used in templates
func (*template) AddTemplateFuncs(funcs map[string]interface{}) {
	for name, fn := range funcs {
		GlobalTmplFuncs[name] = fn
	}
}

// RenderTemplate render a template with given name use given
// value to given writer
func (t *template) RenderTemplate(wr io.Writer, name string, val interface{}) error {
	return t.tmpl.ExecuteTemplate(wr, name, val)
}
