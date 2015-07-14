function irc_init(params, ctx, bot)
    for i = 1, #bot.Config.Channels do
        ctx.Raw("JOIN "..bot.Config.Channels[i])
    end
end

function irc_invite(params, ctx, bot)
    ctx.Raw("JOIN "..params)
end

RegisterEvent("irc_init", "376", irc_init)
RegisterEvent("irc_invite", "INVITE", irc_invite)
