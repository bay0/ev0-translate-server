local enable = gui.new_checkbox('Enable', 'translator', false); enable:set_tooltip('Enables translator')
local translateRealtime = gui.new_checkbox('Enable realtime translation', 'translator_realtime', false); translateRealtime:set_tooltip('Translate realtime')
local api = gui.new_textbox("API", "translator_api_endpoint")
local lang = gui.new_combobox("Target lang", "translator_langs", false, "en", "de", "ru")


function on_get_data(data, postInChat, playerName)
    if postInChat then
        engine_client.exec("say " .. data)
    else
        chat.write(playerName .. ": " .. data)
    end
end

function translate(message, postInChat, playerName)
    print(message)
    http.get(api:get_value() .. "translate/" .. string.lower(lang:get_value()) .. "/" .. message, function (data)
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

        if translateRealtime:get_value() and player ~= localPlayer then
            translate(input, false, player:get_player_info().name)
        end

        if localPlayer == player and input:sub(1, 1) == "$" then
            if input:sub(1, 2) == "$t" then
                utils.run_delayed(1000, function ()
                    translate(input:sub(4, #input), true, nil)
                end)
            end
        end
    end
end
