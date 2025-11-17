local trxd = require("trxd")
local json = require("json")

function TestTypes(challenge)
    print(json.encode(challenge))
    challenge["DockerConfig"]["HashDomain"] = 4
    return challenge
end

trxd.on_event("challengeGet", "TestTypes")