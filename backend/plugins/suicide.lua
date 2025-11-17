local trxd = require("trxd")

function Die(table)
    error ("bye!")
    return table
end

trxd.on_event("pluginsLoaded","Die")