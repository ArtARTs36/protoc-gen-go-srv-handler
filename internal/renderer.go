package internal

import (
	"html/template"
	"io"

	"github.com/artarts36/protoc-gen-go-srv-handler/templates"
)

type Renderer struct {
	templates struct {
		service *template.Template
	}
}

func NewRenderer() (*Renderer, error) {
	rend := &Renderer{}

	srvTmpl, err := template.ParseFS(templates.FS, "*.template")
	if err != nil {
		return nil, err
	}

	rend.templates.service = srvTmpl

	return rend, nil
}

func (r *Renderer) RenderService(w io.Writer, srv *Service) error {
	return r.templates.service.ExecuteTemplate(w, "service.template", map[string]interface{}{
		"Service": srv,
	})
}

func (r *Renderer) RenderServiceTest(w io.Writer, srv *Service) error {
	return r.templates.service.ExecuteTemplate(w, "service_test.template", map[string]interface{}{
		"Service": srv,
	})
}

func (r *Renderer) RenderHandler(w io.Writer, srv *Service, hand *Handler) error {
	return r.templates.service.ExecuteTemplate(w, "handler.template", map[string]interface{}{
		"Service": srv,
		"Handler": hand,
	})
}

func (r *Renderer) RenderHandlerTest(w io.Writer, hand *Handler) error {
	return r.templates.service.ExecuteTemplate(w, "handler_test.template", map[string]interface{}{
		"Service": hand.Service,
		"Handler": hand,
	})
}
