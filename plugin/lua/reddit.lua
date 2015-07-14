function alligator(params, ctx, bot)
    math.randomseed(os.time())
    math.random()
    js = go.GetJSON("http://www.reddit.com/r/britishproblems.json")
    posts = js["data"]["children"]
    return posts[math.random(#posts)]["data"]["title"]
end

RegisterCommand("alligator", alligator)
