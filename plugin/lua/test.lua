function test(params, ctx, bot)
    ctx.Reply('no straights')
end

function test_rate(params, ctx, box)
    ctx.Reply('1')
    ctx.Reply('2')
    ctx.Reply('3')
    ctx.Reply('4')
    ctx.Reply('5')
    ctx.Reply('6')
    ctx.Reply('7')
    ctx.Reply('8')
    ctx.Reply('9')
    ctx.Reply('10')
end

RegisterCommand("test", test)
RegisterCommand("test_rate", test_rate)
