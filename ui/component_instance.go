package ui

import (
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	"github.com/vpaulo/seda/object"
)

// ComponentInstance represents a running component with reactive state
type ComponentInstance struct {
	component      *object.UIComponent
	args           []object.Object
	Env            *object.Environment
	renderer       *Renderer
	app            fyne.App
	window         fyne.Window
	rootWidget     fyne.CanvasObject
	mutex          sync.Mutex
	initialized    bool // Track if component statements have been evaluated
}

// NewComponentInstance creates a new component instance
func NewComponentInstance(component *object.UIComponent, args []object.Object) *ComponentInstance {
	return &ComponentInstance{
		component: component,
		args:      args,
		Env:       object.NewEnclosedEnvironment(component.Env),
	}
}

// SetApp sets the Fyne application
func (ci *ComponentInstance) SetApp(app fyne.App) {
	ci.app = app
}

// SetWindow sets the Fyne window
func (ci *ComponentInstance) SetWindow(window fyne.Window) {
	ci.window = window
}

// SetRenderer sets the renderer
func (ci *ComponentInstance) SetRenderer(renderer *Renderer) {
	ci.renderer = renderer
}

// RenderComponent evaluates the component and builds the UI tree
func (ci *ComponentInstance) RenderComponent() error {
	ci.mutex.Lock()
	defer ci.mutex.Unlock()

	// Only evaluate component body statements (var declarations) on first render
	// On rerenders, we want to keep the existing variable values (reactive state)
	if !ci.initialized {
		for _, stmt := range ci.component.Body.Statements {
			// Use the Eval function from evaluator package
			result := evalFunc(stmt, ci.Env)
			if isErrorFunc(result) {
				return &ComponentError{Message: result.(*object.Error).Message}
			}
		}
		ci.initialized = true
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

	// Set the content on the window
	if ci.window != nil {
		ci.window.SetContent(widget)
	}

	return nil
}

// Rerender re-renders the component from event handlers
func (ci *ComponentInstance) Rerender() {
	// Fyne requires UI updates to happen on the main goroutine
	// But since button clicks already happen on the UI thread, we can call directly
	if err := ci.RenderComponent(); err != nil {
		// Log error but don't crash
		fmt.Printf("Error during rerender: %v\n", err)
		return
	}

	// The new content is already set by RenderComponent via SetContent()
	// No additional refresh needed - SetContent handles it
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
