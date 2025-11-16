local trxd = require("trxd")

function ModifyFirstBlood (submission)
    print("was:",submission["first_blood"])
    submission["first_blood"] = not submission["first_blood"]
    print("now:",submission["first_blood"])
    return submission
end

trxd.on_event("submitFlag",ModifyFirstBlood)