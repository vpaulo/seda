package ui

import (
	"sync"

	"github.com/rivo/tview"
	"github.com/vpaulo/seda/object"
)

// Application is an alias for tview.Application
type Application = tview.Application

// ComponentInstance represents a running component with reactive state
type ComponentInstance struct {
	component *object.UIComponent
	args      []object.Object
	Env       *object.Environment
	renderer  *Renderer
	app       *Application
	rootWidget tview.Primitive
	mutex     sync.Mutex
}

// NewComponentInstance creates a new component instance
func NewComponentInstance(component *object.UIComponent, args []object.Object) *ComponentInstance {
	return &ComponentInstance{
		component: component,
		args:      args,
		Env:       object.NewEnclosedEnvironment(component.Env),
	}
}

// SetApp sets the tview application
func (ci *ComponentInstance) SetApp(app *Application) {
	ci.app = app
}

// SetRenderer sets the renderer
func (ci *ComponentInstance) SetRenderer(renderer *Renderer) {
	ci.renderer = renderer
}

// RenderComponent evaluates the component and builds the UI tree
func (ci *ComponentInstance) RenderComponent() error {
	ci.mutex.Lock()
	defer ci.mutex.Unlock()

	// Re-evaluate component body statements (var declarations)
	for _, stmt := range ci.component.Body.Statements {
		// Use the Eval function from evaluator package
		// Note: We'll need to expose this function
		result := evalFunc(stmt, ci.Env)
		if isErrorFunc(result) {
			return &ComponentError{Message: result.(*object.Error).Message}
		}
	}

	// Evaluate the root UI element
	if ci.component.Body.Root == nil {
		return &ComponentError{Message: "component has no UI tree"}
	}

	rootObj := evalUIElementFunc(ci.component.Body.Root, ci.Env)
	if isErrorFunc(rootObj) {
		return &ComponentError{Message: rootObj.(*object.Error).Message}
	}

	rootElement, ok := rootObj.(*object.UIElement)
	if !ok {
		return &ComponentError{Message: "component root is not a UI element"}
	}

	// Build the widget tree
	widget, err := ci.renderer.BuildWidget(rootElement)
	if err != nil {
		return err
	}

	ci.rootWidget = widget

	// Set the root widget on the app
	if ci.app != nil {
		ci.app.SetRoot(widget, true)
	}

	return nil
}

// Rerender re-renders the component from event handlers
func (ci *ComponentInstance) Rerender() {
	// This will be called from event handlers (different goroutine)
	// Use QueueUpdateDraw to safely update UI from any goroutine
	if ci.app != nil {
		ci.app.QueueUpdateDraw(func() {
			ci.RenderComponent()
		})
	}
}

// ComponentError represents a component rendering error
type ComponentError struct {
	Message string
}

func (e *ComponentError) Error() string {
	return e.Message
}

// These functions will be set by the evaluator package
var evalFunc func(interface{}, *object.Environment) object.Object
var evalUIElementFunc func(interface{}, *object.Environment) object.Object
var isErrorFunc func(object.Object) bool

// SetEvalFunc sets the Eval function from evaluator
func SetEvalFunc(fn func(interface{}, *object.Environment) object.Object) {
	evalFunc = fn
}

// SetEvalUIElementFunc sets the eval_ui_element function from evaluator
func SetEvalUIElementFunc(fn func(interface{}, *object.Environment) object.Object) {
	evalUIElementFunc = fn
}

// SetIsErrorFunc sets the is_error function from evaluator
func SetIsErrorFunc(fn func(object.Object) bool) {
	isErrorFunc = fn
}
