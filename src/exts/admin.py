import discord
from discord.ext import commands
import json
import logging
import inspect
import os
from src.imports import checks, utils

config_dir = 'src/config/'
admin_id_file = 'admin_ids'
owner_id = 351794468870946827
embed_color = discord.Colour.from_rgb(49, 107, 111)
bot_config_file = 'bot_config.json'
invite_match = '(https?://)?(www.)?discord(app.com/(invite|oauth2)|.gg|.io)/[\w\d_\-?=&/]+'

admin_log = logging.getLogger('admin')


class Admin:
    def __init__(self, bot):
        self.bot = bot

    @staticmethod
    def _get_config_string(guild_config):
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
        # noinspection PyProtectedMember
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
            if await self.bot.db_con.fetch('select id from geeksbot_emojis where id = $1', emoji.id):
                await self.bot.db_con.execute("update geeksbot_emojis set id = $2, name = $1, code = $3 "
                                              "where name = $1", emoji.name, emoji.id, emoji_code)
            else:
                await self.bot.db_con.execute("insert into geeksbot_emojis(id,name,code) values ($2,$1,$3)",
                                              emoji.name, emoji.id, emoji_code)
        await ctx.message.add_reaction('✅')
        await ctx.send(f'Emojis have been updated in the database.')

    @commands.command(hidden=True)
    @commands.check(checks.is_guild_owner)
    async def get_guild_config(self, ctx):
        config = await self.bot.db_con.fetchrow('select * from guild_config where guild_id = $1', ctx.guild.id)
        configs = [str(config)[i:i+1990] for i in range(0, len(config), 1990)]
        await ctx.message.author.send(f'The current config for the {ctx.guild.name} guild is:\n')
        admin_log.info(configs)
        for config in configs:
            await ctx.message.author.send(f'```{config}```')
        await ctx.send(f'{ctx.message.author.mention} check your DMs.')

    @commands.group(case_insensitive=True)
    async def set(self, ctx):
        """Run help set for more info"""
        pass

    @commands.group(case_insensitive=True)
    async def add(self, ctx):
        """Run help set for more info"""
        pass

    @commands.group(case_insensitive=True)
    async def remove(self, ctx):
        """Run help set for more info"""
        pass

    @set.command(name='admin_chan', aliases=['ac', 'admin_chat', 'admin chat'])
    async def _admin_channel(self, ctx, channel: discord.TextChannel=None):
        """Sets the channel for admin specific notifications"""
        if ctx.guild:
            if await checks.is_admin(self.bot, ctx):
                if channel is not None:
                    await self.bot.db_con.execute('update guild_config set admin_chat = $2 where guild_id = $1',
                                                  ctx.guild.id, channel.id)
                    await ctx.send(f'{channel.name} is now set as the Admin Chat channel for this guild.')

    @set.command(name='channel_lockdown', aliases=['lockdown', 'restrict_access', 'cl'])
    async def _channel_lockdown(self, ctx, config='true'):
        """Toggles the channel lockdown restricting Geeksbot to only access channels defined in allowed_channels
        If you run this before configuring allowed_channels it will tell you to run that command first."""
        if ctx.guild:
            if await checks.is_admin(self.bot, ctx):
                if str(config).lower() == 'true':
                    if await self.bot.db_con.fetchval('select allowed_channels from guild_config '
                                                      'where guild_id = $1', ctx.guild.id) is []:
                        await ctx.send('Please set at least one allowed channel before running this command.')
                    else:
                        await self.bot.db_con.execute('update guild_config set channel_lockdown = True '
                                                      'where guild_id = $1', ctx.guild.id)
                        await ctx.send('Channel Lockdown is now active.')
                elif str(config).lower() == 'false':
                    if await self.bot.db_con.fetchval('select channel_lockdown from guild_config where guild_id = $1',
                                                      ctx.guild.id):
                        await self.bot.db_con.execute('update guild_config set channel_lockdown = False '
                                                      'where guild_id = $1', ctx.guild.id)
                        await ctx.send('Channel Lockdown has been deactivated.')
                    else:
                        await ctx.send('Channel Lockdown is already deactivated.')
            else:
                await ctx.send(f'You are not authorized to run this command.')
        else:
            await ctx.send('This command must be run from inside a guild.')

    @add.command(name='allowed_channels', aliases=['channel', 'ac'])
    async def _allowed_channels(self, ctx, *, channels):
        """Allows Admin to restrict what channels Geeksbot is allowed to access
        This only takes effect if channel_lockdown is enabled.
        If one of the channels passed is not found then it is ignored."""
        if ctx.guild:
            if await checks.is_admin(self.bot, ctx):
                channels = channels.lower().replace(' ', '').split(',')
                added = ''
                admin_log.info(channels)
                allowed_channels = await self.bot.db_con.fetchval('select allowed_channels from guild_config '
                                                                  'where guild_id = $1', ctx.guild.id)
                if allowed_channels == 'null':
                    allowed_channels = None

                channels = [discord.utils.get(ctx.guild.channels, name=channel)
                            for channel in channels if channel is not None]

                if allowed_channels and channels:
                    allowed_channels = [int(channel) for channel in json.loads(allowed_channels)]
                    channels = [channel.id for channel in channels if channel.id not in allowed_channels]
                    allowed_channels += channels
                    await self.bot.db_con.execute('update guild_config set allowed_channels = $2 where guild_id = $1',
                                                  ctx.guild.id, json.dumps(allowed_channels))
                elif channels:
                    admin_log.info('Config is empty')
                    allowed_channels = [channel.id for channel in channels]
                    await self.bot.db_con.execute('update guild_config set allowed_channels = $2 '
                                                  'where guild_id = $1', ctx.guild.id,
                                                  json.dumps(allowed_channels))
                else:
                    await ctx.send('None of those are valid text channels for this guild.')
                    return

                if channels:
                    channel_str = '\n'.join([str(channel.name) for channel in channels])
                    await ctx.send('The following channels have been added to the allowed channel list: '
                                   f'{channel_str}')
                await ctx.message.add_reaction('✅')
            else:
                await ctx.send(f'You are not authorized to run this command.')
        else:
            await ctx.send('This command must be run from inside a guild.')

    @commands.command()
    @commands.is_owner()
    async def view_code(self, ctx, code_name):
        pag = utils.Paginator(self.bot, prefix='```py', suffix='```')
        pag.add(inspect.getsource(self.bot.get_command(code_name).callback))
        for page in pag.pages():
            await ctx.send(page)

    @add.command(aliases=['prefix', 'p'])
    @commands.cooldown(1, 5, type=commands.BucketType.guild)
    async def add_prefix(self, ctx, *, prefix=None):
        """Adds a guild specific prefix to the guild config
        Note: This overwrites the default of g$. If you would
        like to keep using g$ you will need to add it to the
        Guild config as well."""
        if ctx.guild:
            if await checks.is_admin(self.bot, ctx):
                prefixes = await self.bot.db_con.fetchval('select prefix from guild_config where guild_id = $1',
                                                          ctx.guild.id)
                if prefix is None:
                    await ctx.send(prefixes)
                    return
                elif prefixes is None:
                    prefixes = prefix.replace(' ', ',').split(',')
                else:
                    for p in prefix.replace(' ', ',').split(','):
                        prefixes.append(p)
                if len(prefixes) > 10:
                    await ctx.send(f'Only 10 prefixes are allowed per guild.\nPlease remove some before adding more.')
                    prefixes = prefixes[:10]
                await self.bot.db_con.execute('update guild_config set prefix = $2 where guild_id = $1',
                                              ctx.guild.id, prefixes)
                await ctx.guild.me.edit(nick=f'[{prefixes[0]}] Geeksbot')
                await ctx.send(f"Updated. You currently have {len(prefixes)} "
                               f"{'prefix' if len(prefixes) == 1 else 'prefixes'} "
                               f"in your config.\n{', '.join(prefixes)}")
            else:
                await ctx.send(f'You are not authorized to run this command.')
        else:
            await ctx.send(f'This command must be run from inside a guild.')

    @remove.command(aliases=['prefix', 'p'])
    @commands.cooldown(1, 5, type=commands.BucketType.guild)
    async def remove_prefix(self, ctx, *, prefix=None):
        """Removes a guild specific prefix from the guild config
        If the last prefix is removed then Geeksbot will default
        Back to g$"""
        if ctx.guild:
            if await checks.is_admin(self.bot, ctx):
                prefixes = await self.bot.db_con.fetchval('select prefix from guild_config where guild_id = $1',
                                                          ctx.guild.id)
                found = 0
                if prefix is None:
                    await ctx.send(prefixes)
                    return
                elif prefixes is None or prefixes == []:
                    await ctx.send('There are no custom prefixes setup for this guild.')
                    return
                else:
                    prefix = prefix.replace(' ', ',').split(',')
                    for p in prefix:
                        if p in prefixes:
                            prefixes.remove(p)
                            found = 1
                        else:
                            await ctx.send(f'The prefix {p} is not in the config for this guild.')
                if found:
                    await self.bot.db_con.execute('update guild_config set prefix = $2 where guild_id = $1',
                                                  ctx.guild.id, prefixes)
                    await ctx.guild.me.edit(nick=f'[{prefixes[0] if len(prefixes) != 0 else self.bot.default_prefix}] '
                                                 f'Geeksbot')
                    await ctx.send(f"Updated. You currently have {len(prefixes)} "
                                   f"{'prefix' if len(prefixes) == 1 else 'prefixes'} "
                                   f"in your config.\n{', '.join(prefixes)}")
            else:
                await ctx.send(f'You are not authorized to run this command.')
        else:
            await ctx.send(f'This command must be run from inside a guild.')

    @add.command(name='admin_role', aliases=['admin'])
    @commands.cooldown(1, 5, type=commands.BucketType.guild)
    @commands.check(checks.is_guild_owner)
    async def _add_admin_role(self, ctx, role=None):
        """The Guild owner can add a role to the admin list
        Allowing members of that role to run admin commands
        on the current guild."""
        role = discord.utils.get(ctx.guild.roles, name=role)
        if role is not None:
            roles = json.loads(await self.bot.db_con.fetchval('select admin_roles from guild_config '
                                                              'where guild_id = $1', ctx.guild.id))
            if role.name in roles:
                await ctx.send(f'{role.name} is already registered as an admin role in this guild.')
            else:
                roles[role.name] = role.id
                await self.bot.db_con.execute('update guild_config set admin_roles = $2 where guild_id = $1',
                                              ctx.guild.id, json.dumps(roles))
                await ctx.send(f'{role.name} has been added to the list of admin roles for this guild.')
        else:
            await ctx.send('You must include a role with this command.')

    @remove.command(name='admin_role', aliases=['admin'])
    @commands.cooldown(1, 5, type=commands.BucketType.guild)
    @commands.check(checks.is_guild_owner)
    async def _remove_admin_role(self, ctx, role=None):
        """The Guild owner can remove a role from the admin list"""
        role = discord.utils.get(ctx.guild.roles, name=role)
        if role is not None:
            roles = json.loads(await self.bot.db_con.fetchval('select admin_roles from guild_config '
                                                              'where guild_id = $1', ctx.guild.id))
            if role.name in roles:
                del roles[role.name]
                await self.bot.db_con.execute('update guild_config set admin_roles = $2 where guild_id = $1',
                                              ctx.guild.id, json.dumps(roles))
                await ctx.send(f'{role.name} has been removed from the list of admin roles for this guild.')
            else:
                await ctx.send(f'{role.name} is not registered as an admin role in this guild.')
        else:
            await ctx.send('You must include a role with this command.')


def setup(bot):
    bot.add_cog(Admin(bot))
