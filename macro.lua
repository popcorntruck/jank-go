local started = false

macro("notify", {
    hotkey = "KEY_E",
    action = function()
        if not started then
            started = true
           
           coroutine.wrap(function()
                while started do
                send_notification("click")
                   send_click()

                  sleep(1000)
                end
            end)()
        else
            started = false
            send_notification("Stop")
        end
        -- Assuming send_notification is a function that sends a notification
    end
})

local function create_spotify_macro(uri) 
    return function ()
        if win_class_active("spotify") or win_class_active("com.spotify.client") then
            os.execute(
                string.format("dbus-send --print-reply --dest=org.mpris.MediaPlayer2.spotify /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.OpenUri string:'%s'", uri)
            )
        end
    end
end

macro("t", {
    hotkey = "KEY_T",
    action = create_spotify_macro("spotify:track:0jNhSK5gotdRB1G4nMqEau")
})

macro("y", {
    hotkey = "KEY_Y",
    action = create_spotify_macro("spotify:artist:5K4W6rqBFWDnAN6FQUkS6x")
})

macro("j", {
    hotkey = "KEY_J",
    action = create_spotify_macro("spotify:track:6C2nVSSeXNqfoY8t6tliZ4")
})

