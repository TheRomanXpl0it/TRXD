local trxd = require("trxd")
local http = require("http")

function GetTRXD(_)
    local client = http.client()
    local request = http.request("GET", "https://training.theromanxpl0.it")
    local result, err = client:do_request(request)
    print("result:",result.body)
    return _
end

trxd.on_event("pluginsLoaded", "GetTRXD")