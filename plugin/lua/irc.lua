function irc_init(params, ctx, bot)
    return "JOIN "..bot.Config.Channels[1]
end

RegisterEvent("irc_init", "376", irc_init)
