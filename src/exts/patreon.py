import discord
from discord.ext import commands
import json
from src.imports import checks

config_dir = 'config'
extension_dir = 'extensions'
owner_id = 351794468870946827


class Patreon:

    def __init__(self, bot):
        self.bot = bot

    @commands.command(aliases=['get_patreon', 'patreon'])
    @commands.cooldown(1, 5, type=commands.BucketType.user)
    async def get_patreon_links(self, ctx, target: discord.Member=None):
        """Prints Patreon information for creators on the server."""
        if await self.bot.db_con.fetchval('select patreon_enabled from guild_config where guild_id = $1', ctx.guild.id):
            patreon_info = await self.bot.db_con.fetchrow('select patreon_message,patreon_links from guild_config '
                                                          'where guild_id = $1', ctx.guild.id)
            print(patreon_info)
            message = patreon_info.patreon_mesage.replace('\\n', '\n')
            patreon_links = json.loads(patreon_info.patreon_links)
            for key in patreon_links:
                message = message + '\n{0}: {1}'.format(key, patreon_links[key])
            if target is None:
                await ctx.send(message)
            else:
                await ctx.send('{0}\n{1}'.format(target.mention, message))
        else:
            await ctx.send('Patreon links are not enabled on this guild.')

    @commands.command(aliases=['patreon_message'])
    async def set_patreon_message(self, ctx, message):
        if await checks.is_admin(self.bot, ctx):
            patreon_message = await self.bot.db_con.fetchval('select patreon_message from guild_config '
                                                             'where guild_id = $1', ctx.guild.id)
            if message == patreon_message:
                await ctx.send('That is already the current message for this guild.')
            else:
                await self.bot.db_con.execute('update guild_config set patreon_message = $2 where guild_id = $1',
                                              ctx.guild.id, message)
                await ctx.send(f'The patreon message for this guild has been set to:\n{message}')
        else:
            await ctx.send(f'You are not authorized to run this command.')

    @commands.command(aliases=['add_patreon', 'set_patreon'])
    async def add_patreon_info(self, ctx, name, url):
        if await checks.is_admin(self.bot, ctx):
            patreon_info = await self.bot.db_con.fetchval('select patreon_links from guild_config where guild_id = $1',
                                                          ctx.guild.id)
            patreon_links = {}
            update = 0
            if patreon_info:
                patreon_links = json.loads(patreon_info)
                if name in patreon_links:
                    update = 1
            patreon_links[name] = url
            await self.bot.db_con.execute('update guild_config set patreon_links = $2 where guild_id = $1',
                                          ctx.guild.id, json.dumps(patreon_links))
            await ctx.send(f"The Patreon link for {name} has been "
                           f"{'updated to the new url.' if update else'added to the config for this guild.'}")
        else:
            await ctx.send(f'You are not authorized to run this command.')

    @commands.command(aliases=['remove_patreon'])
    async def remove_patreon_info(self, ctx, name):
        if await checks.is_admin(self.bot, ctx):
            patreon_info = await self.bot.db_con.fetchval('select patreon_links from guild_config where guild_id = $1',
                                                          ctx.guild.id)
            if patreon_info:
                patreon_links = json.loads(patreon_info)
                if name in patreon_links:
                    del patreon_links[name]
                    await self.bot.db_con.execute('update guild_config set patreon_links = $2 where guild_id = $1',
                                                  ctx.guild.id, json.dumps(patreon_links))
                    await ctx.send(f'The Patreon link for {name} has been removed from the config for this guild.')
                    return
                else:
                    await ctx.send(f'{name} is not in the Patreon config for this guild.')
            else:
                await ctx.send(f'There is no Patreon config for this guild.')
        else:
            await ctx.send(f'You are not authorized to run this command.')

    @commands.command()
    async def enable_patreon(self, ctx, state: bool=True):
        if await checks.is_admin(self.bot, ctx):
            patreon_status = await self.bot.db_con.fetchval('select patreon_enabled from guild_config '
                                                            'where guild_id = $1', ctx.guild.id)
            if patreon_status and state:
                await ctx.send('Patreon is already enabled for this guild.')
            elif patreon_status and not state:
                await self.bot.db_con.execute('update guild_config set patreon_enabled = $2 where guild_id = $1',
                                              ctx.guild.id, state)
                await ctx.send('Patreon has been disabled for this guild.')
            elif not patreon_status and state:
                await self.bot.db_con.execute('update guild_config set patreon_enabled = $2 where guild_id = $1',
                                              ctx.guild.id, state)
                await ctx.send('Patreon has been enabled for this guild.')
            elif not patreon_status and not state:
                await ctx.send('Patreon is already disabled for this guild.')

    @commands.command(aliases=['referral', 'ref'])
    @commands.cooldown(1, 5, type=commands.BucketType.user)
    async def referral_links(self, ctx, target: discord.Member=None):
        """Prints G-Portal Referral Links."""
        if await self.bot.db_con.fetchval('select referral_enabled from guild_config where guild_id = $1',
                                          ctx.guild.id):
            referral_info = await self.bot.db_con.fetchval('select referral_message,referral_links from guild_config '
                                                           'where guild_id = $1', ctx.guild.id)
            message = referral_info[0]
            referral_links = json.loads(referral_info[1])
            for key in referral_links:
                message = message + '\n{0}: {1}'.format(key, referral_links[key])
            if target is None:
                await ctx.send(message)
            else:
                await ctx.send('{0}\n{1}'.format(target.mention, message))
        else:
            await ctx.send('Referrals are not enabled on this guild.')


def setup(bot):
    bot.add_cog(Patreon(bot))
