package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/vpaulo/seda/object"
)

// EventHandler is a function that handles UI events (button clicks, etc.)
type EventHandler func(callback object.Object) error

// Renderer converts Seda UI elements to tview widgets
type Renderer struct {
	app          *tview.Application
	eventHandler EventHandler
}

// NewApplication creates a new tview Application
func NewApplication() *tview.Application {
	return tview.NewApplication()
}

// NewRenderer creates a new UI renderer
func NewRenderer(eventHandler EventHandler) *Renderer {
	return &Renderer{
		eventHandler: eventHandler,
	}
}

// Render renders a UI element tree and starts the application
func (r *Renderer) Render(element *object.UIElement) error {
	// Build the tview widget tree
	widget, err := r.buildWidget(element)
	if err != nil {
		return err
	}

	// Set the root widget and run the application
	r.app.SetRoot(widget, true)
	return r.app.Run()
}

// BuildWidget recursively builds tview widgets from Seda UI elements (exported)
func (r *Renderer) BuildWidget(element *object.UIElement) (tview.Primitive, error) {
	return r.buildWidget(element)
}

// buildWidget recursively builds tview widgets from Seda UI elements
func (r *Renderer) buildWidget(element *object.UIElement) (tview.Primitive, error) {
	switch element.ElementType {
	case "Window":
		return r.buildWindow(element)
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

// buildWindow creates a window (Box with border and title)
func (r *Renderer) buildWindow(element *object.UIElement) (tview.Primitive, error) {
	box := tview.NewBox()
	box.SetBorder(true)

	// Set title if provided
	if titleObj, ok := element.Properties["title"]; ok {
		if title, ok := titleObj.(*object.String); ok {
			box.SetTitle(title.Value)
		}
	}

	// If there's a single child, render it inside the window
	if len(element.Children) == 1 {
		child, err := r.buildWidget(element.Children[0])
		if err != nil {
			return nil, err
		}

		// Create a flex container to hold the child
		flex := tview.NewFlex()
		flex.SetBorder(true)

		// Copy title from box
		if titleObj, ok := element.Properties["title"]; ok {
			if title, ok := titleObj.(*object.String); ok {
				flex.SetTitle(title.Value)
			}
		}

		flex.AddItem(child, 0, 1, true)
		return flex, nil
	}

	return box, nil
}

// buildVBox creates a vertical flex container
func (r *Renderer) buildVBox(element *object.UIElement) (tview.Primitive, error) {
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)

	// Add all children vertically
	for _, child := range element.Children {
		widget, err := r.buildWidget(child)
		if err != nil {
			return nil, err
		}
		flex.AddItem(widget, 0, 1, true)
	}

	return flex, nil
}

// buildHBox creates a horizontal flex container
func (r *Renderer) buildHBox(element *object.UIElement) (tview.Primitive, error) {
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexColumn)

	// Add all children horizontally
	for _, child := range element.Children {
		widget, err := r.buildWidget(child)
		if err != nil {
			return nil, err
		}
		flex.AddItem(widget, 0, 1, true)
	}

	return flex, nil
}

// buildText creates a text view
func (r *Renderer) buildText(element *object.UIElement) (tview.Primitive, error) {
	textView := tview.NewTextView()
	textView.SetDynamicColors(true)
	textView.SetTextAlign(tview.AlignCenter)

	// Set text if provided
	if textObj, ok := element.Properties["text"]; ok {
		if text, ok := textObj.(*object.String); ok {
			fmt.Fprintf(textView, text.Value)
		}
	}

	// Set font size (approximate with padding)
	if fontSizeObj, ok := element.Properties["fontSize"]; ok {
		if fontSize, ok := fontSizeObj.(*object.Number); ok {
			// Use font size to determine padding (rough approximation)
			if fontSize.Value > 20 {
				textView.SetTextAlign(tview.AlignCenter)
			}
		}
	}

	return textView, nil
}

// buildButton creates a button
func (r *Renderer) buildButton(element *object.UIElement) (tview.Primitive, error) {
	button := tview.NewButton("Button")

	// Set button text if provided
	if textObj, ok := element.Properties["text"]; ok {
		if text, ok := textObj.(*object.String); ok {
			button.SetLabel(text.Value)
		}
	}

	// Set onClick handler if provided
	if onClickObj, ok := element.Properties["onClick"]; ok {
		// Wire up the click handler
		button.SetSelectedFunc(func() {
			// Invoke the event handler with the callback function
			if r.eventHandler != nil {
				if err := r.eventHandler(onClickObj); err != nil {
					// Log error but don't crash the UI
					// TODO: Better error handling
				}
			}
		})
	}

	// Style the button
	button.SetBackgroundColor(tcell.ColorBlue)
	button.SetLabelColor(tcell.ColorWhite)

	return button, nil
}
