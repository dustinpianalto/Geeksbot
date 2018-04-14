import discord
from discord.ext import commands
import json
from .imports import checks

config_dir = 'config'
extension_dir = 'extensions'
owner_id = 351794468870946827

class patreon():

    def __init__(self, bot):
        self.bot = bot

    @commands.command(aliases=['get_patreon','patreon'])
    @commands.cooldown(1, 5, type=commands.BucketType.user)
    async def get_patreon_links(self, ctx, target: discord.Member=None):
        'Prints Patreon information for creators on the server.'
        if self.bot.con.one('select patreon_enabled from guild_config where guild_id = %(id)s', {'id':ctx.guild.id}):
            patreon_info = self.bot.con.one('select patreon_message,patreon_links from guild_config where guild_id = %(id)s', {'id':ctx.guild.id})
            message = patreon_info[0].replace('\\n','\n')
            patreon_links = json.loads(patreon_info[1])
            for key in patreon_links:
                message = message + '\n{0}: {1}'.format(key, patreon_links[key])
            if target == None:
                await ctx.send(message)
            else:
                await ctx.send('{0}\n{1}'.format(target.mention, message))
        else:
            await ctx.send('Patreon links are not enabled on this guild.')

    @commands.command(aliases=['patreon_message'])
    async def set_patreon_message(self, ctx, message):
        if checks.is_admin(self.bot, ctx):
            patreon_message = self.bot.con.one('select patreon_message from guild_config where guild_id = %(id)s', {'id':ctx.guild.id})
            if message == patreon_message:
                await ctx.send('That is already the current message for this guild.')
            else:
                self.bot.con.run('update guild_config set patreon_message = %(message)s where guild_id = %(id)s', {'id':ctx.guild.id, 'message': message})
                await ctx.send(f'The patreon message for this guild has been set to:\n{message}')
        else:
            await ctx.send(f'You are not authorized to run this command.')

    @commands.command(aliases=['add_patreon', 'set_patreon'])
    async def add_patreon_info(self, ctx, name, url):
        if checks.is_admin(self.bot, ctx):
            patreon_info = self.bot.con.one('select patreon_links from guild_config where guild_id = %(id)s', {'id':ctx.guild.id})
            patreon_links = {}
            update = 0
            if patreon_info:
                patreon_links = json.loads(patreon_info)
                if name in patreon_links:
                    update = 1
            patreon_links[name] = url
            self.bot.con.run('update guild_config set patreon_links = %(links)s where guild_id = %(id)s', {'id':ctx.guild.id, 'links': json.dumps(patreon_links)})
            await ctx.send(f"The Patreon link for {name} has been {'updated to the new url.' if update else'added to the config for this guild.'}")
        else:
            await ctx.send(f'You are not authorized to run this command.')

    @commands.command(aliases=['remove_patreon'])
    async def remove_patreon_info(self, ctx, name):
        if checks.is_admin(self.bot, ctx):
            patreon_info = self.bot.con.one('select patreon_links from guild_config where guild_id = %(id)s', {'id':ctx.guild.id})
            patreon_links = {}
            if patreon_info:
                patreon_links = json.loads(patreon_info)
                if name in patreon_links:
                    del patreon_links[name]
                    self.bot.con.run('update guild_config set patreon_links = %(links)s where guild_id = %(id)s', {'id':ctx.guild.id, 'links': json.dumps(patreon_links)})
                    await ctx.send(f'The Patreon link for {name} has been removed from the config for this guild.')
                    return
                else:
                    await ctx.send(f'{name} is not in the Patreon config for this guild.')
            else:
                await ctx.send(f'There is no Patreon config for this guild.')
        else:
            await ctx.send(f'You are not authorized to run this command.')

    @commands.command()
    async def enable_patreon(self, ctx, state:bool=True):
        if checks.is_admin(self.bot, ctx):
            patreon_status = self.bot.con.one('select patreon_enabled from guild_config where guild_id = %(id)s', {'id':ctx.guild.id})
            if patreon_status and state:
                await ctx.send('Patreon is already enabled for this guild.')
            elif patreon_status and not state:
                self.bot.con.run('update guild_config set patreon_enabled = %(state)s where guild_id = %(id)s', {'id':ctx.guild.id, 'state': state})
                await ctx.send('Patreon has been disabled for this guild.')
            elif not patreon_status and state:
                self.bot.con.run('update guild_config set patreon_enabled = %(state)s where guild_id = %(id)s', {'id':ctx.guild.id, 'state': state})
                await ctx.send('Patreon has been enabled for this guild.')
            elif not patreon_status and not state:
                await ctx.send('Patreon is already disabled for this guild.')

    @commands.command(aliases=['referral','ref'])
    @commands.cooldown(1, 5, type=commands.BucketType.user)
    async def referral_links(self, ctx, target: discord.Member=None):
        'Prints G-Portal Referral Links.'
        if self.bot.con.one('select referral_enabled from guild_config where guild_id = %(id)s', {'id':ctx.guild.id}):
            referral_info = self.bot.con.one('select referral_message,referral_links from guild_config where guild_id = %(id)s', {'id':ctx.guild.id})
            message = referral_info[0]
            referral_links = json.loads(referral_info[1])
            for key in referral_links:
                message = message + '\n{0}: {1}'.format(key, referral_links[key])
            if target == None:
                await ctx.send(message)
            else:
                await ctx.send('{0}\n{1}'.format(target.mention, message))
        else:
            await ctx.send('Referrals are not enabled on this guild.')

def setup(bot):
    bot.add_cog(patreon(bot))
