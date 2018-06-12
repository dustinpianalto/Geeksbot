import discord
from discord.ext import commands
import logging
from src.imports.utils import Paginator, run_command, Book
import asyncio

owner_id = 351794468870946827
embed_color = discord.Colour.from_rgb(49, 107, 111)

git_log = logging.getLogger('git')


class Git:
    def __init__(self, bot):
        self.bot = bot

    @commands.group(case_insensitive=True, invoke_without_command=True)
    async def git(self, ctx):
        """Shows my Git link"""
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
        """Pulls updates from GitHub rebasing branch."""
        pag = Paginator(self.bot, max_line_length=44, embed=True)
        pag.set_embed_meta(title='Git Pull',
                           color=self.bot.embed_color,
                           thumbnail=f'{ctx.guild.me.avatar_url}')
        pag.add('\uFFF6' + await asyncio.wait_for(self.bot.loop.create_task(run_command('git fetch --all')), 120))
        pag.add(await asyncio.wait_for(self.bot.loop.create_task(run_command('git reset --hard '
                                                                             'origin/$(git '
                                                                             'rev-parse --symbolic-full-name'
                                                                             ' --abbrev-ref HEAD)')), 120))
        pag.add('\uFFF7\n\uFFF8')
        pag.add(await asyncio.wait_for(self.bot.loop.create_task(run_command('git show --stat | '
                                                                             'sed "s/.*@.*[.].*/ /g"')), 10))
        book = Book(pag, (None, ctx.channel, self.bot, ctx.message))
        await book.create_book()

    @git.command()
    @commands.is_owner()
    async def status(self, ctx):
        """Gets status of current branch."""
        pag = Paginator(self.bot, max_line_length=44, max_lines=30, embed=True)
        pag.set_embed_meta(title='Git Status',
                           color=self.bot.embed_color,
                           thumbnail=f'{ctx.guild.me.avatar_url}')
        result = await asyncio.wait_for(self.bot.loop.create_task(run_command('git status')), 10)
        pag.add(result)
        book = Book(pag, (None, ctx.channel, self.bot, ctx.message))
        await book.create_book()



def setup(bot):
    bot.add_cog(Git(bot))
