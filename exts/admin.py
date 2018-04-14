import discord
from discord.ext import commands
import json
from srcds import rcon as rcon_con
import time, logging, math
from datetime import datetime, timedelta
import asyncio, inspect
import aiohttp, async_timeout
from bs4 import BeautifulSoup as bs
import traceback
import os, sys
from .imports import checks

config_dir = 'config/'
admin_id_file = 'admin_ids'
extension_dir = 'extensions'
owner_id = 351794468870946827
embed_color = discord.Colour.from_rgb(49,107,111)
bot_config_file = 'bot_config.json'
invite_match = '(https?://)?(www.)?discord(app.com/(invite|oauth2)|.gg|.io)/[\w\d_\-?=&/]+'

admin_log = logging.getLogger('admin')

class admin():
    def __init__(self, bot):
        self.bot = bot

    def _get_config_string(self, guild_config):
        config_str = ''
        for config in guild_config:
            if isinstance(guild_config[config], dict):
                config_str = f'{config_str}\n{" "*4}{config}'
                for item in guild_config[config]:
                    config_str = f'{config_str}\n{" "*8}{item}: {guild_config[config][item]}'
            elif isinstance(guild_config[config], list):
                config_str = f'{config_str}\n{" "*4}{config}'
                for item in guild_config[config]:
                    config_str = f'{config_str}\n{" "*8}{item}'
            else:
                config_str = f'{config_str}\n{" "*4}{config}: {guild_config[config]}'
        return config_str

    @commands.command(hidden=True)
    @commands.is_owner()
    async def reload_bot_config(self, ctx):
        with open(f'{config_dir}{bot_config_file}') as file:
            self.bot.bot_config = json.load(file)
        del self.bot.bot_config['token']
        del self.bot.bot_config['db_con']
        await ctx.send('Config reloaded.')

    @commands.command(hidden=True)
    @commands.is_owner()
    async def reboot(self, ctx):
        await ctx.send('Geeksbot is restarting.')
        with open(f'{config_dir}reboot', 'w') as f:
            f.write(f'1\n{ctx.channel.id}')
        os._exit(1)

    @commands.command(hidden=True)
    @commands.is_owner()
    async def get_bot_config(self, ctx):
        n = 2000
        config = [str(self.bot.bot_config)[i:i+n] for i in range(0, len(str(self.bot.bot_config)), n)]
        for conf in config:
            await ctx.message.author.send(conf)
        await ctx.send(f'{ctx.message.author.mention} check your DMs.')

    @commands.command(hidden=True)
    @commands.is_owner()
    async def update_emojis(self, ctx):
        emojis = self.bot.emojis
        for emoji in emojis:
            if emoji.animated:
                emoji_code = f'<a:{emoji.name}:{emoji.id}>'
            else:
                emoji_code = f'<:{emoji.name}:{emoji.id}>'
            if self.bot.con.all('select id from geeksbot_emojis where id = %(id)s', {'id':emoji.id}):
                self.bot.con.run("update geeksbot_emojis set id = %(id)s, name = %(name)s, code = %(emoji_code)s where name = %(name)s", {'name':emoji.name,'id':emoji.id,'emoji_code':emoji_code})
            else:
                self.bot.con.run("insert into geeksbot_emojis(id,name,code) values (%(id)s,%(name)s,%(emoji_code)s)", {'name':emoji.name,'id':emoji.id,'emoji_code':emoji_code})
        await ctx.message.add_reaction('✅')
        await ctx.send(f'Emojis have been updated in the database.')

    @commands.command(hidden=True)
    @commands.check(checks.is_guild_owner)
    async def get_guild_config(self, ctx):
        config = self.bot.con.one('select * from guild_config where guild_id = %(id)s', {'id':ctx.guild.id})
        configs = [str(config)[i:i+1990] for i in range(0, len(config), 1990)]
        await ctx.message.author.send(f'The current config for the {ctx.guild.name} guild is:\n')
        admin_log.info(configs)
        for config in configs:
            await ctx.message.author.send(f'```{config}```')
        await ctx.send(f'{ctx.message.author.mention} check your DMs.')

    @commands.group(case_insensitive=True)
    async def set(self, ctx):
        '''Run help set for more info'''
        pass

    @commands.group(case_insensitive=True)
    async def add(self, ctx):
        '''Run help set for more info'''
        pass

    @commands.group(case_insensitive=True)
    async def remove(self, ctx):
        '''Run help set for more info'''
        pass

    @set.command(name='admin_chan', aliases=['ac', 'admin_chat', 'admin chat'])
    async def _admin_channel(self, ctx, channel:discord.TextChannel=None):
        if ctx.guild:
            if checks.is_admin(self.bot, ctx):
                if channel != None:
                    self.bot.con.run('update guild_config set admin_chat = %(chan)s where guild_id = %(id)s', {'id':ctx.guild.id, 'chan': channel.id})
                    await ctx.send(f'{channel.name} is now set as the Admin Chat channel for this guild.')

    @set.command(name='channel_lockdown', aliases=['lockdown', 'restrict_access', 'cl'])
    async def _channel_lockdown(self, ctx, config='true'):
        if ctx.guild:
            if checks.is_admin(self.bot, ctx):
                if str(config).lower() == 'true':
                    if self.bot.con.one('select allowed_channels from guild_config where guild_id = %(id)s', {'id':ctx.guild.id}) == []:
                        await ctx.send('Please set at least one allowed channel before running this command.')
                    else:
                        self.bot.con.run('update guild_config set channel_lockdown = True where guild_id = %(id)s', {'id':ctx.guild.id})
                        await ctx.send('Channel Lockdown is now active.')
                elif str(config).lower() == 'false':
                    if self.bot.con.one('select channel_lockdown from guild_config where guild_id = %(id)s', {'id':ctx.guild.id}):
                        self.bot.con.run('update guild_config set channel_lockdown = False where guild_id = %(id)s', {'id':ctx.guild.id})
                        await ctx.send('Channel Lockdown has been deactivated.')
                    else:
                        await ctx.send('Channel Lockdown is already deactivated.')
            else:
                await ctx.send(f'You are not authorized to run this command.')
        else:
            await ctx.send('This command must be run from inside a guild.')


    @add.command(name='allowed_channels', aliases=['channel','ac'])
    async def _allowed_channels(self, ctx, *, channels):
        if ctx.guild:
            if checks.is_admin(self.bot, ctx):
                channels = channels.lower().replace(' ','').split(',')
                added = ''
                for channel in channels:
                    chnl = discord.utils.get(ctx.guild.channels, name=channel)
                    if chnl == None:
                        await ctx.send(f'{channel} is not a valid text channel in this guild.')
                    else:
                        admin_log.info('Chan found')
                        if self.bot.con.one('select allowed_channels from guild_config where guild_id = %(id)s', {'id':ctx.guild.id}):
                            if chnl.id in json.loads(self.bot.con.one('select allowed_channels from guild_config where guild_id = %(id)s', {'id':ctx.guild.id})):
                                admin_log.info('Chan found in config')
                                await ctx.send(f'{channel} is already in the list of allowed channels. Skipping...')
                            else:
                                admin_log.info('Chan not found in config')
                                allowed_channels = json.loads(self.bot.con.one('select allowed_channels from guild_config where guild_id = %(id)s', {'id':ctx.guild.id})).append(chnl.id)
                                self.bot.con.run('update guild_config set allowed_channels = %(channels)s where guild_id = %(id)s', {'id':ctx.guild.id, 'channels':allowed_channels})
                                added = f'{added}\n{channel}'
                        else:
                            admin_log.info('Chan not found in config')
                            allowed_channels = [chnl.id]
                            self.bot.con.run('update guild_config set allowed_channels = %(channels)s where guild_id = %(id)s', {'id':ctx.guild.id, 'channels':allowed_channels})
                            added = f'{added}\n{channel}'
                if added != '':
                    await ctx.send(f'The following channels have been added to the allowed channel list: {added}')
                await ctx.message.add_reaction('✅')
            else:
                await ctx.send(f'You are not authorized to run this command.')
        else:
            await ctx.send('This command must be run from inside a guild.')

    @commands.command()
    @commands.is_owner()
    async def view_code(self, ctx, code_name):
        await ctx.send(f"```py\n{inspect.getsource(self.bot.get_command(code_name).callback)}\n```")

    @add.command(aliases=['prefix','p'])
    @commands.cooldown(1, 5, type=commands.BucketType.guild)
    async def add_prefix(self, ctx, *, prefix=None):
        if ctx.guild:
            if checks.is_admin(self.bot, ctx):
                prefixes = self.bot.con.one('select prefix from guild_config where guild_id = %(id)s', {'id':ctx.guild.id})
                if prefix == None:
                    await ctx.send(prefixes)
                    return
                elif prefixes == None:
                    prefixes = prefix.replace(' ',',').split(',')
                else:
                    for p in prefix.replace(' ',',').split(','):
                        prefixes.append(p)
                if len(prefixes) > 10:
                    await ctx.send(f'Only 10 prefixes are allowed per guild.\nPlease remove some before adding more.')
                    prefixes = prefixes[:10]
                self.bot.con.run('update guild_config set prefix = %(prefixes)s where guild_id = %(id)s', {'id':ctx.guild.id, 'prefixes':prefixes})
                await ctx.guild.me.edit(nick=f'[{prefixes[0]}] Geeksbot')
                await ctx.send(f"Updated. You currently have {len(prefixes)} {'prefix' if len(prefixes) == 1 else 'prefixes'} in your config.\n{', '.join(prefixes)}")
            else:
                await ctx.send(f'You are not authorized to run this command.')
        else:
            await ctx.send(f'This command must be run from inside a guild.')

    @remove.command(aliases=['prefix','p'])
    @commands.cooldown(1, 5, type=commands.BucketType.guild)
    async def remove_prefix(self, ctx, *, prefix=None):
        if ctx.guild:
            if checks.is_admin(self.bot, ctx):
                prefixes = []
                prefixes = self.bot.con.one('select prefix from guild_config where guild_id = %(id)s', {'id':ctx.guild.id})
                found = 0
                if prefix == None:
                    await ctx.send(prefixes)
                    return
                elif prefixes == None or prefixes == []:
                    await ctx.send('There are no custom prefixes setup for this guild.')
                    return
                else:
                    prefix = prefix.replace(' ',',').split(',')
                    for p in prefix:
                        if p in prefixes:
                            prefixes.remove(p)
                            found = 1
                        else:
                            await ctx.send(f'The prefix {p} is not in the config for this guild.')
                if found:
                    self.bot.con.run('update guild_config set prefix = %(prefixes)s where guild_id = %(id)s', {'id':ctx.guild.id, 'prefixes':prefixes})
                    await ctx.guild.me.edit(nick=f'[{prefixes[0] if len(prefixes) != 0 else self.bot.default_prefix}] Geeksbot')
                    await ctx.send(f"Updated. You currently have {len(prefixes)} {'prefix' if len(prefixes) == 1 else 'prefixes'} in your config.\n{', '.join(prefixes)}")
            else:
                await ctx.send(f'You are not authorized to run this command.')
        else:
            await ctx.send(f'This command must be run from inside a guild.')

    @add.command(name='admin_role', aliases=['admin'])
    @commands.cooldown(1, 5, type=commands.BucketType.guild)
    @commands.check(checks.is_guild_owner)
    async def _add_admin_role(self, ctx, role=None):
        role = discord.utils.get(ctx.guild.roles, name=role)
        if role != None:
            roles = json.loads(self.bot.con.one('select admin_roles from guild_config where guild_id = %(id)s', {'id':ctx.guild.id}))
            if role.name in roles:
                await ctx.send(f'{role.name} is already registered as an admin role in this guild.')
            else:
                roles[role.name] = role.id
                self.bot.con.run('update guild_config set admin_roles = %(roles)s where guild_id = %(id)s', {'id':ctx.guild.id, 'roles':json.dumps(roles)})
                await ctx.send(f'{role.name} has been added to the list of admin roles for this guild.')
        else:
            await ctx.send('You must include a role with this command.')

    @remove.command(name='admin_role', aliases=['admin'])
    @commands.cooldown(1, 5, type=commands.BucketType.guild)
    @commands.check(checks.is_guild_owner)
    async def _remove_admin_role(self, ctx, role=None):
        role = discord.utils.get(ctx.guild.roles, name=role)
        if role != None:
            roles = json.loads(self.bot.con.one('select admin_roles from guild_config where guild_id = %(id)s', {'id':ctx.guild.id}))
            if role.name in roles:
                del roles[role.name]
                self.bot.con.run('update guild_config set admin_roles = %(roles)s where guild_id = %(id)s', {'id':ctx.guild.id, 'roles':roles})
                await ctx.send(f'{role.name} has been removed from the list of admin roles for this guild.')
            else:
                await ctx.send(f'{role.name} is not registered as an admin role in this guild.')
        else:
            await ctx.send('You must include a role with this command.')

def setup(bot):
    bot.add_cog(admin(bot))
