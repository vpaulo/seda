## Simple Hello World UI example

component HelloWorld() ::
    Window {
        title: "Hello World",

        VBox {
            Text {
                text: "Welcome to Seda UI!"
            }
            Text {
                text: "Press Ctrl+C to exit"
            }
        }
    }
end

## Mount the component to render it
UI.mount(HelloWorld)
