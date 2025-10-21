## Test UI element evaluation with property binding

component Counter(initial: Number) ::
    var count = initial
    var title = "Counter App"

    Window {
        title: title,
        width: 400,
        height: 300,

        VBox {
            Text {
                text: "Count: #{count}",
                fontSize: 24
            }
            Button {
                text: "Increment",
                onClick: fn() :: count = count + 1 end
            }
        }
    }
end

println("Component Counter with nested UI elements defined successfully")

component SimpleWindow() ::
    Window {
        title: "Simple Window",
        width: 200
    }
end

println("Component SimpleWindow defined successfully")

## Test component with interpolated string properties
component Greeting(name: String) ::
    Window {
        title: "Hello #{name}!",
        Text {
            text: "Welcome, #{name}!"
        }
    }
end

println("Component Greeting with interpolated properties defined successfully")
