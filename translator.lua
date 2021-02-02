local enable = gui.new_checkbox('Enable', 'translator', false); enable:set_tooltip('Enables translator')
local translateRealtime = gui.new_checkbox('Enable realtime translation', 'translator_realtime', false); translateRealtime:set_tooltip('Translate realtime')
local apiKey = gui.new_textbox("API-Key", "translator_apikey"); apiKey:set_tooltip('Your google API-Key')
local api = gui.new_textbox("API", "translator_api_endpoint"); api:set_tooltip('The path to the backend you use')
local lang = gui.new_combobox("Target lang", "translator_langs", false, "en", "de", "ru")

function on_get_data(data, postInChat, playerName)
    if postInChat then
        engine_client.exec("say " .. data)
    else
        chat.write(playerName .. ": " .. data)
    end
end

function translate(message, postInChat, playerName)
    local postData = {
        username = info.ev0lve.username,
        apikey = apiKey:get_value()
    }
    http.post(api:get_value() .. "/api/translate/" .. string.lower(lang:get_value()) .. "/" .. message, postData, function (data)
        on_get_data(data, postInChat, playerName)
    end)
end

function on_player_say(event)
    if enable:get_value() then
        local input = event:get_string('text', "")

        local playerIndex = engine_client.get_player_for_userid(event:get_int('userid', 0))
        local player = entity_list.get_entity(playerIndex)

        local localPlayerIndex = engine_client.get_local_player()
        local localPlayer = entity_list.get_entity(localPlayerIndex)

        if translateRealtime:get_value() and localPlayerIndex ~= playerIndex then
            translate(input, false, player:get_player_info().name)
        end

        if localPlayerIndex == playerIndex and input:sub(1, 1) == "$" then

            print(input:sub(1, 7))

            if input:sub(1, 2) == "$t" then
                utils.run_delayed(1000, function ()
                    translate(input:sub(4, #input), true, nil)
                end)
            end

            if input:sub(1, 7) == "$apikey" then
                utils.run_delayed(1000, function ()
                    apiKey:set_value(input:sub(9, #input))
                end)
            end

            if input:sub(1, 4) == "$api" then
                utils.run_delayed(1000, function ()
                    api:set_value(input:sub(6, #input))
                end)
            end
        end
    end
end
