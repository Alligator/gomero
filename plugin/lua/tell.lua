tells = {}

function tell(params, ctx, bot)
    s = go.SplitN(params, " ", 2)
    if #s < 2 then
        return "what are you doing"
    end
    tells[s[1]] = s[2]
    return "I'll pass that along"
end

function tell_privmsg(params, ctx, bot)
    if tells[ctx.Nick] then
        ctx.Say(ctx.Nick..': '..tells[ctx.Nick])
        table.remove(tells, ctx.Nick)
    end
end

RegisterEvent("tell_privmsg", "PRIVMSG", tell_privmsg)
RegisterCommand("tell", tell)
