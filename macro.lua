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

macro("j", {
    hotkey = "KEY_J",
    action = function()
        if not win_class_active("spotify") then
          local res = os.execute(
             "dbus-send --print-reply --dest=org.mpris.MediaPlayer2.spotify /org/mpris/MediaPlayer2 org.mpris.MediaPlayer2.Player.OpenUri string:'spotify:track:6C2nVSSeXNqfoY8t6tliZ4'"
          )
          print(res)
        end
    end
})
