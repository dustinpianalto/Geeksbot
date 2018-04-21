import discord
from discord.ext import commands
import logging
from .imports.utils import paginate, run_command
import asyncio

owner_id = 351794468870946827
embed_color = discord.Colour.from_rgb(49, 107, 111)

git_log = logging.getLogger('git')


class Git:
    def __init__(self, bot):
        self.bot = bot

    @commands.group(case_insensitive=True)
    async def git(self, ctx):
        """Run help git for more info"""
        pass

    @git.command()
    @commands.is_owner()
    async def pull(self, ctx):
        em = discord.Embed(style='rich',
                           title=f'Git Pull',
                           color=embed_color)
        em.set_thumbnail(url=f'{ctx.guild.me.avatar_url}')
        result = await asyncio.wait_for(self.bot.loop.create_task(run_command('git fetch --all')), 120) + '\n'
        result += await asyncio.wait_for(self.bot.loop.create_task(run_command('git reset --hard '
                                                                               'origin/$(git '
                                                                               'rev-parse --symbolic-full-name'
                                                                               ' --abbrev-ref HEAD)')), 120) + '\n\n'
        result += await asyncio.wait_for(self.bot.loop.create_task(run_command('git show --stat | '
                                                                               'sed "s/.*@.*[.].*/ /g"')), 10)
        results = paginate(result, maxlen=1014)
        for page in results[:5]:
            em.add_field(name='￲', value=f'{page}')
        await ctx.send(embed=em)

    @git.command()
    @commands.is_owner()
    async def status(self, ctx):
        em = discord.Embed(style='rich',
                           title=f'Git Pull',
                           color=embed_color)
        em.set_thumbnail(url=f'{ctx.guild.me.avatar_url}')
        result = await asyncio.wait_for(self.bot.loop.create_task(run_command('git status')), 10)
        results = paginate(result, maxlen=1014)
        for page in results[:5]:
            em.add_field(name='￲', value=f'{page}')
        await ctx.send(embed=em)


def setup(bot):
    bot.add_cog(Git(bot))
