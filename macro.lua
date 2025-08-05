macro("notify", {
    hotkey = "KEY_E",
    action = function()
       send_notification("Test Macro Triggered")
    end
})

macro("j", {
    hotkey = "KEY_J",
    action = function()

        -- Now send the OpenUri command
        local res = os.execute(
           "dbus-send --print-reply --dest=org.mpris.MediaPlayer2.spotify /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.OpenUri string:'spotify:track:6C2nVSSeXNqfoY8t6tliZ4'"
        )
        print(res)
    end
})