# ev0-translate-server

API for the ev0 lua to translate messages. ([preview](https://streamable.com/e/hlojer))

## Build
```bash
go build .
```

## Binaries
Get binaries from the build server [![](https://github.com/bay0/ev0-translate-server/workflows/build/badge.svg)](https://github.com/bay0/ev0-translate-server/actions).

## Usage
Create the environment variables
```
{
    APP_PORT= //port the api runs on e.x 8080
    APP_GOOGLE_CLOUD_API_KEY= //google cloud api key e.x WojsIIWikm1mRPtQ5rSLmH5drgTvstgo5363FxK
}
```
## Dependencies
[Cloud Translation](https://cloud.google.com/translate)

Lua to use with [ev0lve.xyz](https://ev0lve.xyz/)

```lua
local enable = gui.new_checkbox('Enable', 'translator', false); enable:set_tooltip('Enables translator')
local translateRealtime = gui.new_checkbox('Enable realtime translation(Crashes)', 'translator_realtime', false); translateRealtime:set_tooltip('Translate realtime')
local apiKey = gui.new_textbox("API-Key", "translator_apikey"); apiKey:set_tooltip('Your google API-Key')
local api = gui.new_textbox("API", "translator_api_endpoint"); api:set_tooltip('The path to the backend you use')

local function get_val(t)
    local va = {}
    for i,v in pairs(t) do
        table.insert(va, i)
    end
    return unpack(va)
end

local languages = {
    ["Chinese (Simplified)"] = "zh",
    ["Japanese"] = "ja",
    ["Korean"] = "ko",
    ["Arabic"] = "ar",
    ["Swedish"] = "sv",
    ["Italian"] = "it",
    ["Norwegian"] = "no",
    ["Romanian"] = "ro",
    ["Portuguese"] = "pt",
    ["Polish"] = "pl",
    ["Turkish"] = "tr",
    ["French"] = "fr",
    ["Russian"] = "ru",
    ["English"] = "en",
    ["German"] = "de"
}

local lang = gui.new_combobox("Target lang", "translator_langs", false, get_val(languages))


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
    http.post(api:get_value() .. "/api/translate/" .. string.lower(languages[lang:get_value()]) .. "/" .. message, postData, function (data)
        on_get_data(data, postInChat, playerName)
    end)
end

function on_player_say(event)
    if enable:get_value() and engine_client.is_ingame() then
        local input = event:get_string('text', "")

        local playerIndex = engine_client.get_player_for_userid(event:get_int('userid', 0))
        local player = entity_list.get_entity(playerIndex)

        local localPlayerIndex = engine_client.get_local_player()
        local localPlayer = entity_list.get_entity(localPlayerIndex)

        if translateRealtime:get_value() and localPlayerIndex ~= playerIndex then
            translate(input, false, player:get_player_info().name)
        end

        if localPlayerIndex == playerIndex and input:sub(1, 1) == "$" then

            if input:sub(1, 2) == "$t" then
                utils.run_delayed(1000, function ()
                    translate(input:sub(4, #input), true, nil)
                end)
            end

            if input:sub(1, 7) == "$apikey" then
                apiKey:set_value(input:sub(9, #input))
            end

            if input:sub(1, 4) == "$api" then
                api:set_value(input:sub(6, #input))
            end
        end
    end
end
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
