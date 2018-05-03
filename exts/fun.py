import discord
from discord.ext import commands
import logging
from datetime import datetime
import asyncio
import youtube_dl

config_dir = 'config/'
admin_id_file = 'admin_ids'
extension_dir = 'extensions'
owner_id = 351794468870946827
guild_config_dir = 'guild_config/'
rcon_config_file = 'server_rcon_config'
dododex_url = 'http://www.dododex.com'
embed_color = discord.Colour.from_rgb(49, 107, 111)
bot_config_file = 'bot_config'
default_guild_config_file = 'default_guild_config.json'
emoji_guild = 408524303164899338

events_log = logging.getLogger('events')

emojis = {
    'x': '❌',
    'y': '✅',
    'poop': '💩'
}


class Fun:
    def __init__(self, bot):
        self.bot = bot

    @commands.command()
    @commands.cooldown(1, 30, type=commands.BucketType.user)
    async def infect(self, ctx, member: discord.Member, emoji):
        if member.id == self.bot.user.id and ctx.author.id != owner_id:
            await ctx.send(f'You rolled a Critical Fail...\nInfection bounces off and rebounds on the attacker.')
            member = ctx.author
        if member in self.bot.infected:
            await ctx.send(f'{member.display_name} is already infected. '
                           f'Please wait until they are healed before infecting them again...')
        else:
            emoji = self.bot.get_emoji(int(emoji.split(':')[2].strip('>'))) if '<:' in emoji \
                                                                               or '<a:' in emoji else emoji
            self.bot.infected[member] = [emoji, datetime.now().timestamp()]
            await ctx.send(f"{member.display_name} has been infected with {emoji}")

    @commands.command()
    @commands.cooldown(1, 5, type=commands.BucketType.user)
    async def heal(self, ctx, member: discord.Member):
        if ctx.author == member and ctx.author.id != owner_id:
            await ctx.send('You can\'t heal yourself silly...')
        else:
            if member in self.bot.infected:
                del self.bot.infected[member]
                await ctx.send(f'{member.mention} You have been healed by {ctx.author.display_name}.')
            else:
                await ctx.send(f'{member.display_name} is not infected...')

    @commands.command()
    @commands.is_owner()
    async def print_infections(self, ctx):
        await ctx.author.send(f'```{self.bot.infected}```')

    @commands.command()
    @commands.cooldown(1, 5, type=commands.BucketType.user)
    async def slap(self, ctx, member: discord.Member):
        if member.id == self.bot.user.id and ctx.author.id != owner_id:
            await ctx.send(f'You rolled a Critical Fail...\nThe trout bounces off and rebounds on the attacker.')
            await ctx.send(f'{ctx.author.mention} '
                           f'You slap yourself in the face with a large trout <:trout:408543365085397013>')
        else:
            await ctx.send(f'{ctx.author.display_name} slaps '
                           f'{member.mention} around a bit with a large trout <:trout:408543365085397013>')

    @staticmethod
    def get_factorial(number):
        a = 1
        for i in range(1, int(number)):
            a = a * (i + 1)
        return a

    # @commands.command()
    # @commands.cooldown(1, 5, type=commands.BucketType.user)
    # async def fact(self, ctx, number:int):
    #     if number < 20001 and number > 0:
    #         n = 1990
    #         with ctx.channel.typing():
    #             a = await self.bot.loop.run_in_executor(None, self.get_factorial, number)
    #         if len(str(a)) > 6000:
    #             for b in [str(a)[i:i+n] for i in range(0, len(str(a)), n)]:
    #                 await ctx.author.send(f'```py\n{b}```')
    #             await ctx.send(f"{ctx.author.mention} Check your DMs.")
    #         else:
    #             for b in [str(a)[i:i+n] for i in range(0, len(str(a)), n)]:
    #                 await ctx.send(f'```py\n{b}```')
    #     else:
    #         await ctx.send("Invalid number. Please enter a number between 0 and 20,000")

    @commands.command(hidden=True)
    @commands.is_owner()
    async def play(self, ctx, url=None):
        if ctx.author.voice.channel.name in self.bot.voice_chans:
            if self.bot.voice_chans[ctx.author.voice.channel.name].is_playing():
                await ctx.send('There is currently a song playing. Please wait until it is done...')
                return
        else:
            self.bot.voice_chans[ctx.author.voice.channel.name] = await ctx.author.voice.channel.connect()
            asyncio.sleep(5)
        if url:
            opts = {"format": 'webm[abr>0]/bestaudio/best',
                    "ignoreerrors": True,
                    "default_search": "auto",
                    "source_address": "0.0.0.0",
                    'quiet': True}
            ydl = youtube_dl.YoutubeDL(opts)
            info = ydl.extract_info(url, download=False)
            self.bot.player = discord.FFmpegPCMAudio(info['url'])
        else:
            self.bot.player = discord.FFmpegPCMAudio('dead_puppies.mp3')
        self.bot.player = discord.PCMVolumeTransformer(self.bot.player, volume=0.3)
        self.bot.voice_chans[ctx.author.voice.channel.name].play(self.bot.player)

    @commands.command(hidden=True)
    @commands.is_owner()
    async def stop(self, ctx):
        if ctx.author.voice.channel.name in self.bot.voice_chans:
            if self.bot.voice_chans[ctx.author.voice.channel.name].is_playing():
                self.bot.voice_chans[ctx.author.voice.channel.name].stop()
            else:
                await ctx.send('Nothing is playing...')
        else:
            await ctx.send('Not connected to that voice channel.')

    @commands.command(hidden=True)
    @commands.is_owner()
    async def disconnect(self, ctx):
        if ctx.author.voice.channel.name in self.bot.voice_chans:
            await self.bot.voice_chans[ctx.author.voice.channel.name].disconnect()
            del self.bot.voice_chans[ctx.author.voice.channel.name]
        else:
            await ctx.send('Not connected to that voice channel.')

    @commands.command(hidden=True)
    @commands.is_owner()
    async def volume(self, ctx, volume: float):
        self.bot.player.volume = volume

    @commands.command(name='explode', aliases=['splode'])
    async def explode_user(self, ctx, member: discord.Member=None):
        if member is None:
            member = ctx.author

        trans = await self.bot.db_con.fetchval('select code from geeksbot_emojis where id = 405943174809255956')
        msg = await ctx.send(f'{member.mention}{trans*20}{self.bot.unicode_emojis["left_fist"]}')
        for i in range(20):
            await asyncio.sleep(0.1)
            await msg.edit(content=f'{member.mention}{trans*(20-i)}{self.bot.unicode_emojis["left_fist"]}')
        await asyncio.sleep(0.1)
        await msg.edit(content=f'{self.bot.unicode_emojis["boom"]}')
        await asyncio.sleep(0.1)
        await msg.edit(content=f'{self.bot.unicode_emojis["boom"]} <---- {member.mention} that was you...')


def setup(bot):
    bot.add_cog(Fun(bot))
