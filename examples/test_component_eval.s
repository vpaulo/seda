## Test component evaluation

component SimpleButton(label: String) ::
    var clicks = 0

    Window {
        title: "Button Test"
    }
end

## Test that component is defined
println("Component defined successfully")
