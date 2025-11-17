local trxd = require("trxd")
local json = require("json")

function ModifyChallenges(challenges)
    -- print(json.encode(challenges))

    -- iterate over slice from Go
    for _, challenge in ipairs(challenges) do
        if challenge.Category == "Web" then
            challenge.Name = "Palle"
        end
    end

    return challenges
end

trxd.on_event("challengesGet","ModifyChallenges")