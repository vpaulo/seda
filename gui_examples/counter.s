## Interactive Counter Example
## This demonstrates reactive state with event handlers

component Counter(initial: Number) ::
    var count = initial

    Window {
        title: "Interactive Counter",

        VBox {
            Text {
                text: "Count: #{count}"
            }
            Button {
                text: "Increment",
                onClick: fn() ::
                    count = count + 1
                end
            }
            Button {
                text: "Decrement",
                onClick: fn() ::
                    count = count - 1
                end
            }
            Button {
                text: "Reset",
                onClick: fn() ::
                    count = initial
                end
            }
            Text {
                text: "Press Tab to switch buttons, Enter to click"
            }
        }
    }
end

## Mount the component
UI.mount(Counter, 0)
