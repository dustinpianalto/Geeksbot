import discord
from discord.ext import commands
import logging
from src.imports.utils import Paginator, run_command
import asyncio

owner_id = 351794468870946827
embed_color = discord.Colour.from_rgb(49, 107, 111)

git_log = logging.getLogger('git')


class Git:
    def __init__(self, bot):
        self.bot = bot

    @commands.group(case_insensitive=True, invoke_without_command=True)
    async def git(self, ctx):
        """Run help git for more info"""
        em = discord.Embed(style='rich',
                           title=f'Here is where you can find my code',
                           url='https://github.com/dustinpianalto/Geeksbot/tree/development',
                           description='I am the development branch of Geeksbot. You can find the master branch here:\n'
                                       'https://github.com/dustinpianalto/Geeksbot/',
                           color=embed_color)
        em.set_thumbnail(url=f'{ctx.guild.me.avatar_url}')
        await ctx.send(embed=em)

    @git.command()
    @commands.is_owner()
    async def pull(self, ctx):
        pag = Paginator(max_line_length=60, max_lines=30, max_chars=1014)
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
        pag.add(result)
        for page in pag.pages():
            em.add_field(name='￲', value=f'{page}')
        await ctx.send(embed=em)

    @git.command()
    @commands.is_owner()
    async def status(self, ctx):
        pag = Paginator(max_line_length=60, max_lines=30, max_chars=1014)
        em = discord.Embed(style='rich',
                           title=f'Git Pull',
                           color=embed_color)
        em.set_thumbnail(url=f'{ctx.guild.me.avatar_url}')
        result = await asyncio.wait_for(self.bot.loop.create_task(run_command('git status')), 10)
        pag.add(result)
        for page in pag.pages():
            em.add_field(name='￲', value=f'{page}')
        await ctx.send(embed=em)


def setup(bot):
    bot.add_cog(Git(bot))