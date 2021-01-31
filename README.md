# ev0-translate-server

API for the ev0 lua to translate messages. ([preview](https://streamable.com/e/hlojer))

## Build
```bash
go build .
```

## Binaries
Get binaries from the build server [![](https://github.com/bay0/ev0-translate-server/workflows/build/badge.svg)](https://github.com/bay0/ev0-translate-server/actions).

## Usage
Create the config.json
```json
{
    "apiKey": "",
    "port": 8080
}
```
## Dependencies
[Cloud Translation](https://cloud.google.com/translate)

Lua to use with [ev0lve.xyz](https://ev0lve.xyz/)

```lua
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
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
