## Test UI module access and functions

## Define a simple component
component HelloWorld() ::
    Window {
        title: "Hello World",
        width: 300,
        height: 200
    }
end

println("Component defined: HelloWorld")

## Test UI module is accessible
var ui_check = isNull(UI)
println("UI module accessible: #{!ui_check}")

## Test UI.mount() with component
var result = UI.mount(HelloWorld)
println("Mount result: #{result}")

## Test defining a more complex component
component Counter(initial: Number) ::
    var count = initial

    Window {
        title: "Counter",
        VBox {
            Text {
                text: "Count: #{count}"
            }
            Button {
                text: "Click me"
            }
        }
    }
end

println("Component defined: Counter")

## Test mounting Counter component
var counter_result = UI.mount(Counter)
println("Counter mount result: #{counter_result}")
