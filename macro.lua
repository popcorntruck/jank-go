macro("notify", {
    hotkey = "KEY_E",
    action = function()
       send_notification("Test Macro Triggered")
    end
})

macro("print_and_notify", {
    hotkey = "KEY_Q",
    action = function()
       print("you can print stuff too!!")
       send_notification("The other macro ran")
    end
})