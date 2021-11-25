package router

import (
	"github.com/Postcord/objects"
	"github.com/Postcord/rest"
)

// Defines the form item.
type formItem struct {
	// TODO
}

// Form is used to define a forms contents.
type Form struct {
	// Title is used to define the title of this form.
	Title string `json:"title"`

	// Path is used to define the path which this is mounted on.
	Path string `json:"path"`

	// Generator is used to define the generator which is used to generate the form.
	// This is useful if you pass through parameters to the form from another context.
	Generator FormGenerator `json:"generator"`

	// Handler is used to define the handler which is called when this form is submitted.
	Handler FormHandler `json:"handler"`
}

// FormGenerationCtx is used to define the context of a form in a generation state.
type FormGenerationCtx struct {
	// Defines the interaction which started this.
	*objects.Interaction

	// RESTClient is used to define the REST client.
	RESTClient *rest.Client `json:"rest_client"`

	// Defines the form contents.
	contents []*formItem
}

// FormSubmitCtx is used to define the context of a form in a submission state.
type FormSubmitCtx struct {
	// TODO
}

// FormGenerator is the function type which is used for form generation. The "fnParams" are passed from the
// parent .Form function that called this.
type FormGenerator func(ctx *FormGenerationCtx, fnParams ...interface{}) error

// FormHandler is used to handle a form submission. Params are used to define the parameters in the path.
type FormHandler func(ctx *FormSubmitCtx, params map[string]string) error
