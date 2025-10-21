## Static Counter Example (no reactivity yet)
## This demonstrates the UI structure for a counter
## Reactivity will be added in Phase 4

component Counter(initial: Number) ::
    var count = initial

    Window {
        title: "Counter App",

        VBox {
            Text {
                text: "Count: #{count}"
            }
            Text {
                text: "(Static - no reactivity yet)"
            }
            Button {
                text: "Increment",
                onClick: fn() ::
                    ## This will be wired up to actually work in Phase 4
                    println("Button clicked! (handler not wired yet)")
                end
            }
        }
    }
end

## Mount the component
UI.mount(Counter, 0)
