local enable = gui.new_checkbox('Enable', 'translator', false); enable:set_tooltip('Enables translator')
local api = gui.new_textbox("API", "translator_api_endpoint")
local lang = gui.new_combobox("Target lang", "translator_langs", false, "en", "de", "ru")


function on_get_data(data)
    print(data)
    engine_client.exec("say " .. data)
end

function translate(message)
    print(message)
    http.get(api:get_value() .. "translate/" .. string.lower(lang:get_value()) .. "/" .. message, on_get_data)
end

function on_player_say(event)
    if enable:get_value() then
        local input = event:get_string('text', "")
        if input:sub(1, 1) == "$" then
            local userID = event:get_int('userid', 0)
            print("UserID: " .. userID .. " Input: " .. input)
            if input:sub(1, 2) == "$t" then
                print("UserID: " .. userID .. " accessed command $t")
                utils.run_delayed(1000, function ()
                    translate(input:sub(4, #input))
                end)
            end
        end
    end
end