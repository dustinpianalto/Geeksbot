    @checks.no_bots()
    @commands.cooldown(1,5,commands.BucketType.user)
    @commands.command()
    async def captcha(self, ctx, type, *, text):
        type = type.lower()
        if type not in "checked unchecked loading".split():
            raise commands.BadArgument(f"Invalid type {type!r}. Available "
                    "types: `unchecked`, `loading`, `checked`")
        font = ImageFont.truetype("Roboto-Regular.ttf", 14)
        async with ctx.typing():
            img = Image.open(f"blank-captcha-{type}.png")
            img.load()
            d = ImageDraw.Draw(img)
            fnc = functools.partial(d.text, (53,30), text, fill=(0,0,0,255),
                    font=font)
            await self.bot.loop.run_in_executor(None, fnc)
            img.save("captcha.png")
        await ctx.send(file=discord.File("captcha.png"))
        os.system("rm captcha.png")
        img.close()


import functools, youtube_dl
#bot.voice_chan = await ctx.author.voice.channel.connect()
bot.voice_chan.stop()
opts = {"format": 'webm[abr>0]/bestaudio/best',"ignoreerrors": True,"default_search": "auto","source_address": "0.0.0.0",'quiet': True}
ydl = youtube_dl.YoutubeDL(opts)
url = 'https://www.youtube.com/watch?v=hjbPszSt5Pc'
func = functools.partial(ydl.extract_info, url, download=False)
info = func()
#bot.voice_chan.play(discord.FFmpegPCMAudio('dead_puppies.mp3'))
bot.voice_chan.play(discord.FFmpegPCMAudio(info['url']))
#async while bot.voice_chan.is_playing():
#    pass
#await bot.voice_chan.disconnect()