package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/vpaulo/seda/object"
)

// EventHandler is a function that handles UI events (button clicks, etc.)
type EventHandler func(callback object.Object) error

// keyboardButton extends widget.Button to support Enter key activation
type keyboardButton struct {
	widget.Button
}

// TypedKey handles keyboard events for the button
func (b *keyboardButton) TypedKey(key *fyne.KeyEvent) {
	// Activate button on Enter/Return key
	if key.Name == fyne.KeyReturn || key.Name == fyne.KeyEnter {
		b.OnTapped()
		return
	}
	// Let the base widget handle other keys (like Space)
	b.Button.TypedKey(key)
}

// KeyDown handles key down events (required for desktop.Keyable)
func (b *keyboardButton) KeyDown(key *fyne.KeyEvent) {
	// Activate button on Enter/Return key press
	if key.Name == fyne.KeyReturn || key.Name == fyne.KeyEnter {
		b.OnTapped()
	}
}

// KeyUp handles key up events (required for desktop.Keyable)
func (b *keyboardButton) KeyUp(key *fyne.KeyEvent) {
	// No action needed on key up
}

// Renderer converts Seda UI elements to Fyne widgets
type Renderer struct {
	app          fyne.App
	eventHandler EventHandler
}

// NewApplication creates a new Fyne Application
func NewApplication() fyne.App {
	return app.New()
}

// NewRenderer creates a new UI renderer
func NewRenderer(fyneApp fyne.App, eventHandler EventHandler) *Renderer {
	return &Renderer{
		app:          fyneApp,
		eventHandler: eventHandler,
	}
}

// BuildWidget recursively builds Fyne widgets from Seda UI elements
func (r *Renderer) BuildWidget(element *object.UIElement) (fyne.CanvasObject, error) {
	return r.buildWidget(element)
}

// buildWidget recursively builds Fyne widgets from Seda UI elements
func (r *Renderer) buildWidget(element *object.UIElement) (fyne.CanvasObject, error) {
	switch element.ElementType {
	case "Window":
		// Window is handled specially - just return its content
		return r.buildWindowContent(element)
	case "VBox":
		return r.buildVBox(element)
	case "HBox":
		return r.buildHBox(element)
	case "Text":
		return r.buildText(element)
	case "Button":
		return r.buildButton(element)
	default:
		return nil, fmt.Errorf("unknown UI element type: %s", element.ElementType)
	}
}

// buildWindowContent extracts the content from a Window element
func (r *Renderer) buildWindowContent(element *object.UIElement) (fyne.CanvasObject, error) {
	// If there's a single child, return it
	if len(element.Children) == 1 {
		return r.buildWidget(element.Children[0])
	}

	// If multiple children, wrap in VBox
	if len(element.Children) > 1 {
		return r.buildVBox(element)
	}

	// Empty window
	return widget.NewLabel(""), nil
}

// buildVBox creates a vertical box container
func (r *Renderer) buildVBox(element *object.UIElement) (fyne.CanvasObject, error) {
	var widgets []fyne.CanvasObject

	// Add all children vertically
	for _, child := range element.Children {
		widget, err := r.buildWidget(child)
		if err != nil {
			return nil, err
		}
		widgets = append(widgets, widget)
	}

	return container.NewVBox(widgets...), nil
}

// buildHBox creates a horizontal box container
func (r *Renderer) buildHBox(element *object.UIElement) (fyne.CanvasObject, error) {
	var widgets []fyne.CanvasObject

	// Add all children horizontally
	for _, child := range element.Children {
		widget, err := r.buildWidget(child)
		if err != nil {
			return nil, err
		}
		widgets = append(widgets, widget)
	}

	return container.NewHBox(widgets...), nil
}

// buildText creates a text label
func (r *Renderer) buildText(element *object.UIElement) (fyne.CanvasObject, error) {
	text := ""

	// Set text if provided
	if textObj, ok := element.Properties["text"]; ok {
		if textStr, ok := textObj.(*object.String); ok {
			text = textStr.Value
		}
	}

	label := widget.NewLabel(text)

	// Center align by default (matching original behavior)
	label.Alignment = fyne.TextAlignCenter

	return label, nil
}

// buildButton creates a button
func (r *Renderer) buildButton(element *object.UIElement) (fyne.CanvasObject, error) {
	buttonText := "Button"

	// Set button text if provided
	if textObj, ok := element.Properties["text"]; ok {
		if text, ok := textObj.(*object.String); ok {
			buttonText = text.Value
		}
	}

	// Create a keyboard-enabled button that responds to Enter key
	button := &keyboardButton{}
	button.Text = buttonText

	// Set onClick handler if provided
	if onClickObj, ok := element.Properties["onClick"]; ok {
		button.OnTapped = func() {
			// Invoke the event handler with the callback function
			if r.eventHandler != nil {
				if err := r.eventHandler(onClickObj); err != nil {
					// Log error but don't crash the UI
					fmt.Printf("Error in onClick handler: %v\n", err)
				}
			}
		}
	} else {
		button.OnTapped = func() {
			// No-op if no handler
		}
	}

	// Set button importance for visual styling
	button.Importance = widget.MediumImportance

	// Must call ExtendBaseWidget for custom widgets
	button.ExtendBaseWidget(button)

	return button, nil
}

// GetWindowTitle extracts the window title from the root element
func (r *Renderer) GetWindowTitle(element *object.UIElement) string {
	if element.ElementType == "Window" {
		if titleObj, ok := element.Properties["title"]; ok {
			if title, ok := titleObj.(*object.String); ok {
				return title.Value
			}
		}
	}
	return "Seda Application"
}
