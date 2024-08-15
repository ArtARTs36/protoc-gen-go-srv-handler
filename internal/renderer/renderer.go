package renderer

import (
	"html/template"
	"io"

	"github.com/artarts36/protoc-gen-go-srv-handler/internal/entity"
	"github.com/artarts36/protoc-gen-go-srv-handler/internal/options"
	"github.com/artarts36/protoc-gen-go-srv-handler/templates"
)

type Renderer struct {
	templates *template.Template
}

func NewRenderer() (*Renderer, error) {
	rend := &Renderer{}

	tmpl, err := template.ParseFS(templates.FS, "*.template")
	if err != nil {
		return nil, err
	}

	rend.templates = tmpl

	return rend, nil
}

func (r *Renderer) RenderService(w io.Writer, srv *entity.Service) error {
	return r.templates.ExecuteTemplate(w, "service.template", map[string]interface{}{
		"Service": srv,
	})
}

func (r *Renderer) RenderServiceTest(w io.Writer, srv *entity.Service) error {
	return r.templates.ExecuteTemplate(w, "service_test.template", map[string]interface{}{
		"Service": srv,
	})
}

func (r *Renderer) RenderHandler(w io.Writer, hand *entity.Handler, params RenderHandlerParams) error {
	return r.templates.ExecuteTemplate(w, "handler.template", map[string]interface{}{
		"Service": hand.Service,
		"Handler": hand,
		"Params":  params,
	})
}

type RenderHandlerParams struct {
	RequestValidator options.RequestValidator
}

func (r *Renderer) RenderHandlerTest(w io.Writer, hand *entity.Handler, params RenderHandlerParams) error {
	return r.templates.ExecuteTemplate(w, "handler_test.template", map[string]interface{}{
		"Service": hand.Service,
		"Handler": hand,
		"Params":  params,
	})
}
