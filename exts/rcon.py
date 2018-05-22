import discord
from discord.ext import commands
import json
from srcds import rcon as rcon_con
import time
import logging
from datetime import datetime
import asyncio
import traceback
from .imports import checks

config_dir = 'config'
admin_id_file = 'admin_ids'
extension_dir = 'extensions'
owner_id = 351794468870946827
guild_config_file = 'guild_config'
rcon_config_file = 'server_rcon_config'
guild_config_dir = 'guild_config/'


rcon_log = logging.getLogger('rcon')

game_commands = ['admin', ]
game_prefix = '$'


class Rcon:

    def __init__(self, bot):
        self.bot = bot

    @staticmethod
    def _listplayers(con_info):
        con = rcon_con.RconConnection(con_info['ip'], con_info['port'], con_info['password'])
        asyncio.sleep(5)
        response = con.exec_command('listplayers')
        rcon_log.info(response)
        while response == b'Keep Alive\x00\x00':
            asyncio.sleep(5)
            response = con.exec_command('listplayers')
            rcon_log.info(response)
        return response.strip(b'\n \x00\x00').decode('ascii').strip()

    @staticmethod
    def _saveworld(con_info):
        con = rcon_con.RconConnection(con_info['ip'], con_info['port'], con_info['password'])
        asyncio.sleep(5)
        response = con.exec_command('saveworld')
        rcon_log.info(response)
        while response == b'Keep Alive\x00\x00':
            asyncio.sleep(5)
            response = con.exec_command('saveworld')
            rcon_log.info(response)
        return response.strip(b'\n \x00\x00').decode('ascii').strip()

    @staticmethod
    def _whitelist(con_info, steam_ids):
        messages = []
        con = rcon_con.RconConnection(con_info['ip'], con_info['port'], con_info['password'])
        asyncio.sleep(5)
        for steam_id in steam_ids:
            response = con.exec_command('AllowPlayerToJoinNoCheck {0}'.format(steam_id))
            rcon_log.info(response)
            while response == b'Keep Alive\x00\x00':
                asyncio.sleep(5)
                response = con.exec_command('AllowPlayerToJoinNoCheck {0}'.format(steam_id))
                rcon_log.info(response)
            messages.append(response.strip(b'\n \x00\x00').decode('ascii').strip())
        return messages

    @staticmethod
    def _broadcast(con_info, message):
        messages = []
        con = rcon_con.RconConnection(con_info['ip'], con_info['port'], con_info['password'])
        asyncio.sleep(5)
        response = con.exec_command('broadcast {0}'.format(message))
        rcon_log.info(response)
        while response == b'Keep Alive\x00\x00':
            asyncio.sleep(5)
            response = con.exec_command('broadcast {0}'.format(message))
            rcon_log.info(response)
        messages.append(response.strip(b'\n \x00\x00').decode('ascii').strip())
        return messages

    @staticmethod
    def _get_current_chat(con):
        response = con.exec_command('getchat')
        rcon_log.debug(response)
        return response.strip(b'\n \x00\x00').decode('ascii').strip()

    def server_chat_background_process(self, guild_id, con):
        return_messages = []
        try:
            message = 'Server received, But no response!!'
            time_now = datetime.now().timestamp()
            while 'Server received, But no response!!' in message and time_now + 20 > datetime.now().timestamp():
                message = self._get_current_chat(con)
                rcon_log.debug(message)
                time.sleep(1)
            if 'Server received, But no response!!' not in message:
                for msg in message.split('\n'):
                    msg = '{0} ||| {1}'.format(datetime.now().strftime('%Y-%m-%d %H:%M:%S'), msg.strip())
                    return_messages.append(msg)
        except rcon_con.RconError:
            rcon_log.error('RCON Error {0}\n{1}'.format(guild_id, traceback.format_exc()))
            return_messages.append('RCON Error')
        except Exception as e:
            rcon_log.error('Exception {0}\n{1}'.format(guild_id, traceback.format_exc()))
            return_messages.append('Exception')
            return_messages.append(e)
        rcon_log.debug(return_messages)
        return return_messages

    @staticmethod
    def admin(ctx, msg, rcon_server, admin_roles):
        player = msg.split(' ||| ')[1].split(' (')[0]
        con = rcon_con.RconConnection(rcon_server['ip'],
                                      rcon_server['port'],
                                      rcon_server['password'],
                                      True)
        con.exec_command('ServerChatToPlayer "{0}" GeeksBot: Admin Geeks have been notified you need assistance. '
                         'Please be patient.'.format(player))
        # noinspection PyProtectedMember
        con._sock.close()
        for role in admin_roles:
            msg = '{0} {1}'.format(msg, discord.utils.get(ctx.guild.roles, id=admin_roles[role]).mention)
        return msg

    @commands.command()
    @commands.guild_only()
    async def monitor_chat(self, ctx, *, server=None):
        """Begins monitoring the specified ARK server for chat messages and other events.
        The specified server must already be in the current guild\'s configuration.
        To add and remove ARK servers from the guild see add_rcon_server and remove_rcon_server.
        The server argument is not case sensitive and if the server name has two
            words it can be in one of the following forms:
        first last
        first_last
        "first last"
        To view all the valid ARK servers for this guild see list_ark_servers."""
        if checks.is_rcon_admin(self.bot, ctx):
            if server is not None:
                rcon_connections = json.loads(await self.bot.db_con.fetchval('select rcon_connections '
                                                                             'from guild_config where guild_id = $1',
                                                                             ctx.guild.id))
                server = server.replace('_', ' ').title()
                if server in rcon_connections:
                    rcon_connections[server]["monitoring_chat"] = 1
                    await self.bot.db_con.execute('update guild_config set rcon_connections = $2 where guild_id = $1',
                                                  ctx.guild.id, json.dumps(rcon_connections))
                    channel = self.bot.get_channel(rcon_connections[server]['game_chat_chan_id'])
                    await channel.send('Started monitoring on the {0} server.'.format(server))
                    await ctx.message.add_reaction('✅')
                    rcon_log.debug('Started monitoring on the {0} server.'.format(server))
                    while rcon_connections[server]["monitoring_chat"] == 1:
                        try:
                            con = rcon_con.RconConnection(rcon_connections[server]['ip'],
                                                          rcon_connections[server]['port'],
                                                          rcon_connections[server]['password'],
                                                          True)
                            messages = await self.bot.loop.run_in_executor(None, self.server_chat_background_process,
                                                                           ctx.guild.id, con)
                            # noinspection PyProtectedMember
                            con._sock.close()
                        except TimeoutError:
                            rcon_log.error(traceback.format_exc())
                            await channel.send('TimeoutError')
                            await asyncio.sleep(30)
                        else:
                            rcon_log.debug('Got chat from {0}.'.format(server))
                            for message in messages:
                                rcon_log.info(message)
                                message = '```{0}```'.format(message)
                                for command in game_commands:
                                    prefix_command = '{0}{1}'.format(game_prefix, command)
                                    if prefix_command in message:
                                        try:
                                            func = getattr(self, command)
                                        except AttributeError:
                                            rcon_log.warning('Function not found "{0}"'.format(command))
                                        else:
                                            rcon_log.info(f'Sending to {command}')
                                            message = func(ctx, message, rcon_connections['server'])
                                await channel.send('{0}'.format(message))
                        await asyncio.sleep(1)
                        rcon_connections = json.loads(await self.bot.db_con.fetchval('select rcon_connections '
                                                                                     'from guild_config '
                                                                                     'where guild_id = $1',
                                                                                     ctx.guild.id))
                    await channel.send('Monitoring Stopped')
                else:
                    await ctx.send(f'Server not found: {server}')
            else:
                await ctx.send(f'You must include a server in this command.')
        else:
            await ctx.send(f'You are not authorized to run this command.')

    @commands.command()
    @commands.guild_only()
    async def end_monitor_chat(self, ctx, *, server=None):
        """Ends chat monitoring on the specified server.
        Context is the same as monitor_chat"""
        if checks.is_rcon_admin(self.bot, ctx):
            if server is not None:
                rcon_connections = json.loads(await self.bot.db_con.fetchval('select rcon_connections '
                                                                             'from guild_config where guild_id = $1',
                                                                             ctx.guild.id))
                server = server.replace('_', ' ').title()
                if server in rcon_connections:
                    rcon_connections[server]["monitoring_chat"] = 0
                    await self.bot.db_con.execute('update guild_config set rcon_connections = $2 where guild_id = $1',
                                                  ctx.guild.id, json.dumps(rcon_connections))
                else:
                    await ctx.send(f'Server not found: {server}')
            else:
                await ctx.send(f'You must include a server in this command.')
        else:
            await ctx.send(f'You are not authorized to run this command.')

    @commands.command()
    @commands.guild_only()
    async def listplayers(self, ctx, *, server=None):
        """Lists the players currently connected to the specified ARK server.
        The specified server must already be in the current guild\'s configuration.
        To add and remove ARK servers from the guild see add_rcon_server and remove_rcon_server.
        The server argument is not case sensitive and if the server name has two
            words it can be in one of the following forms:
        first last
        first_last
        "first last"
        To view all the valid ARK servers for this guild see list_ark_servers."""
        if checks.is_rcon_admin(self.bot, ctx):
            rcon_connections = json.loads(await self.bot.db_con.fetchval('select rcon_connections from guild_config '
                                                                         'where guild_id = $1', ctx.guild.id))
            if server is not None:
                server = server.replace('_', ' ').title()
                if server in rcon_connections:
                    connection_info = rcon_connections[server]
                    msg = await ctx.send('Getting Data for the {0} server'.format(server.title()))
                    await ctx.channel.trigger_typing()
                    message = self._listplayers(connection_info)
                    await ctx.channel.trigger_typing()
                    await msg.delete()
                    await ctx.send('Players currently on the {0} server:\n{1}'.format(server.title(), message))
                else:
                    await ctx.send('That server is not in my configuration.\nPlease add it via !add_rcon_server "{0}" '
                                   '"ip" port "password" if you would like to get info from it.'.format(server))
            else:
                for server in rcon_connections:
                    msg = await ctx.send('Getting Data for the {0} server'.format(server.title()))
                    # noinspection PyBroadException
                    try:
                        connection_info = rcon_connections[server]
                        async with ctx.channel.typing():
                            message = self._listplayers(connection_info)
                    except Exception as e:
                        await msg.delete()
                        await ctx.send(f'Player listing failed on the {server.title()}\n{e}')
                    else:
                        await msg.delete()
                        await ctx.send('Players currently on the {0} server:\n{1}'.format(server.title(), message))
        else:
            await ctx.send(f'You are not authorized to run this command.')

    @commands.command()
    @commands.guild_only()
    async def add_rcon_server(self, ctx, server, ip, port, password):
        """Adds the specified server to the current guild\'s rcon config.
        All strings (<server>, <ip>, <password>) must be contained inside double quotes."""
        if checks.is_rcon_admin(self.bot, ctx):
            server = server.title()
            rcon_connections = json.loads(await self.bot.db_con.fetchval('select rcon_connections from guild_config '
                                                                         'where guild_id = $1', ctx.guild.id))
            if server not in rcon_connections:
                rcon_connections[server] = {
                    'ip': ip,
                    'port': port,
                    'password': password,
                    'name': server.lower().replace(' ', '_'),
                    'game_chat_chan_id': 0,
                    'msg_chan_id': 0,
                    'monitoring_chat': 0
                }
                await self.bot.db_con.execute('update guild_config set rcon_connections = $2 where guild_id = $1',
                                              ctx.guild.id, json.dumps(rcon_connections))
                await ctx.send('{0} server has been added to my configuration.'.format(server))
            else:
                await ctx.send('This server name is already in my configuration. Please choose another.')
        else:
            await ctx.send(f'You are not authorized to run this command.')
        await ctx.message.delete()
        await ctx.send('Command deleted to prevent password leak.')

    @commands.command()
    @commands.guild_only()
    async def remove_rcon_server(self, ctx, server):
        """removes the specified server from the current guild\'s rcon config.
        All strings <server> must be contained inside double quotes."""
        if checks.is_rcon_admin(self.bot, ctx):
            server = server.title()
            rcon_connections = json.loads(await self.bot.db_con.fetchval('select rcon_connections from guild_config '
                                                                         'where guild_id = $1', ctx.guild.id))
            if server in rcon_connections:
                del rcon_connections[server]
                await self.bot.db_con.execute('update guild_config set rcon_connections = $2 where guild_id = $1',
                                              ctx.guild.id, json.dumps(rcon_connections))
                await ctx.send('{0} has been removed from my configuration.'.format(server))
            else:
                await ctx.send('{0} is not in my configuration.'.format(server))
        else:
            await ctx.send(f'You are not authorized to run this command.')

    @commands.command()
    @commands.guild_only()
    async def add_whitelist(self, ctx, *, steam_ids=None):
        """Adds the included Steam 64 IDs to the running whitelist on all the ARK servers in the current guild\'s rcon config.
        Steam 64 IDs should be a comma seperated list of IDs.
        Example: 76561198024193239,76561198024193239,76561198024193239"""
        if checks.is_rcon_admin(self.bot, ctx):
            if steam_ids is not None:
                rcon_connections = json.loads(await self.bot.db_con.fetchval('select rcon_connections '
                                                                             'from guild_config where guild_id = $1',
                                                                             ctx.guild.id))
                error = 0
                error_msg = ''
                success_msg = 'Adding to the running whitelist on all servers.'
                steam_ids = steam_ids.replace(', ', ',').replace(' ', ',').split(',')
                for (i, steam_id) in enumerate(steam_ids):
                    try:
                        steam_id = int(steam_id)
                    except ValueError:
                        error = 1
                        error_msg = '{0}\n__**ERROR:**__ {1} is not a valid Steam64 ID'.format(error_msg, steam_id)
                    else:
                        steam_ids[i] = steam_id
                if error == 0:
                    msg = await ctx.send(success_msg)
                    for server in rcon_connections:
                        try:
                            success_msg = '{0}\n\n{1}:'.format(success_msg, server.title())
                            await msg.edit(content=success_msg.strip())
                            messages = await self.bot.loop.run_in_executor(None,
                                                                           self._whitelist,
                                                                           rcon_connections[server],
                                                                           steam_ids)
                        except Exception as e:
                            success_msg = '{0}\n{1}'.format(success_msg, e)
                            await msg.edit(content=success_msg.strip())
                        else:
                            for message in messages:
                                success_msg = '{0}\n{1}'.format(success_msg, message.strip())
                            await msg.edit(content=success_msg.strip())
                    await msg.add_reaction('✅')
                else:
                    await ctx.send(error_msg)
            else:
                await ctx.send('I need a list of steam IDs to add to the whitelist.')
        else:
            await ctx.send(f'You are not authorized to run this command.')

    @commands.command()
    @commands.guild_only()
    async def saveworld(self, ctx, *, server=None):
        """Runs SaveWorld on the specified ARK server.
        If a server is not specified it will default to running saveworld on all servers in the guild\'s config.
        Will print out "World Saved" for each server when the command completes successfully."""
        if checks.is_rcon_admin(self.bot, ctx):
            rcon_connections = json.loads(await self.bot.db_con.fetchval('select rcon_connections from guild_config '
                                                                         'where guild_id = $1', ctx.guild.id))
            success_msg = 'Running saveworld'
            if server is None:
                success_msg += ' on all the servers:'
                msg = await ctx.send(success_msg)
                for server in rcon_connections:
                    try:
                        success_msg = '{0}\n\n{1}:'.format(success_msg, server.title())
                        await msg.edit(content=success_msg.strip())
                        message = await self.bot.loop.run_in_executor(None, self._saveworld, rcon_connections[server])
                    except Exception as e:
                        success_msg = '{0}\n{1}'.format(success_msg, e)
                        await msg.edit(content=success_msg.strip())
                    else:
                        success_msg = '{0}\n{1}'.format(success_msg, message.strip())
                        await msg.edit(content=success_msg.strip())
                await msg.add_reaction('✅')
            elif server.replace('_', ' ').title() in rcon_connections:
                success_msg = '{0} {1}:'.format(success_msg, server.replace('_', ' ').title())
                msg = await ctx.send(success_msg)
                message = await self.bot.loop.run_in_executor(None, self._saveworld, rcon_connections[server.title()])
                success_msg = '{0}\n{1}'.format(success_msg, message.strip())
                await msg.edit(content=success_msg.strip())
                await msg.add_reaction('✅')
            else:
                await ctx.send(f'{server.title()} is not currently in the configuration for this guild.')
        else:
            await ctx.send(f'You are not authorized to run this command.')

    @commands.group(case_insensitive=True)
    async def broadcast(self, ctx):
        """Run help broadcast for more info"""
        pass

    @broadcast.command(name='all')
    @commands.guild_only()
    async def broadcast_all(self, ctx, *, message=None):
        """Sends a broadcast message to all servers in the guild config.
        The message will be prefixed with the Discord name of the person running the command.
        Will print "Success" for each server once the broadcast is sent."""
        if checks.is_rcon_admin(self.bot, ctx):
            rcon_connections = json.loads(await self.bot.db_con.fetchval('select rcon_connections from guild_config '
                                                                         'where guild_id = $1', ctx.guild.id))
            if message is not None:
                message = f'{ctx.author.display_name}: {message}'
                success_msg = f'Broadcasting "{message}" to all servers.'
                msg = await ctx.send(success_msg)
                for server in rcon_connections:
                    try:
                        success_msg = '{0}\n\n{1}:'.format(success_msg, server.title())
                        await msg.edit(content=success_msg.strip())
                        messages = await self.bot.loop.run_in_executor(None,
                                                                       self._broadcast,
                                                                       rcon_connections[server],
                                                                       message)
                    except Exception as e:
                        success_msg = '{0}\n{1}'.format(success_msg, e)
                        await msg.edit(content=success_msg.strip())
                    else:
                        for mesg in messages:
                            if mesg == 'Server received, But no response!!':
                                mesg = 'Success'
                            success_msg = '{0}\n{1}'.format(success_msg, mesg.strip())
                        await msg.edit(content=success_msg.strip())
                await msg.add_reaction('✅')
            else:
                await ctx.send('You must include a message with this command.')
        else:
            await ctx.send(f'You are not authorized to run this command.')

    @broadcast.command(name='server')
    @commands.guild_only()
    async def broadcast_server(self, ctx, server, *, message=None):
        """Sends a broadcast message to the specified server that is in the guild's config.
        The message will be prefixed with the Discord name of the person running the command.
        If <server> has more than one word in it's name it will either need to be surrounded
        by double quotes or the words separated by _"""
        if checks.is_rcon_admin(self.bot, ctx):
            rcon_connections = json.loads(await self.bot.db_con.fetchval('select rcon_connections from guild_config '
                                                                         'where guild_id = $1', ctx.guild.id))
            if server is not None:
                server = server.replace('_', ' ').title()
                if message is not None:
                    message = f'{ctx.author.display_name}: {message}'
                    success_msg = f'Broadcasting "{message}" to {server}.'
                    msg = await ctx.send(success_msg)
                    if server in rcon_connections:
                        messages = await self.bot.loop.run_in_executor(None,
                                                                       self._broadcast,
                                                                       rcon_connections[server],
                                                                       message)
                        for mesg in messages:
                            if mesg != 'Server received, But no response!!':
                                success_msg = '{0}\n{1}'.format(success_msg, mesg.strip())
                                await msg.edit(content=success_msg.strip())
                        await msg.add_reaction('✅')
                    else:
                        await ctx.send(f'{server} is not in the config for this guild')
                else:
                    await ctx.send('You must include a message with this command.')
            else:
                await ctx.send('You must include a server with this command')
        else:
            await ctx.send(f'You are not authorized to run this command.')

    @commands.command()
    @commands.guild_only()
    async def create_server_chat_chans(self, ctx):
        """Creates a category and server chat channels in the current guild.
        The category will be named "Server Chats" and read_messages is disabled for the guild default role (everyone)
        This can be overridden by modifying the category's permissions.
        Inside this category a channel will be created for each server in the current guild's rcon config
        and these channel's permissions will be synced to the category.
        These channels will be added to the guild's rcon config and are where the
        server chat messages will be sent when monitor_chat is run."""
        if checks.is_rcon_admin(self.bot, ctx):
            rcon_connections = json.loads(await self.bot.db_con.fetchval('select rcon_connections from guild_config '
                                                                         'where guild_id = $1', ctx.guild.id))
            edited = 0
            category = discord.utils.get(ctx.guild.categories, name='Server Chats')
            if category is None:
                overrides = {ctx.guild.default_role: discord.PermissionOverwrite(read_messages=False)}
                category = await ctx.guild.create_category('Server Chats', overwrites=overrides)
            channels = ctx.guild.channels
            cat_chans = []
            for channel in channels:
                if channel.category_id == category.id:
                    cat_chans.append(channel)
            for server in rcon_connections:
                exists = 0
                if cat_chans:
                    for channel in cat_chans:
                        if rcon_connections[server]['game_chat_chan_id'] == channel.id:
                            exists = 1
                if exists == 0:
                    print('Creating {}'.format(server))
                    chan = await ctx.guild.create_text_channel(rcon_connections[server]['name'], category=category)
                    rcon_connections[server]['game_chat_chan_id'] = chan.id
                    edited = 1
            if edited == 1:
                await self.bot.db_con.execute('update guild_config set rcon_connections = $2 where guild_id = $1',
                                              ctx.guild.id, json.dumps(rcon_connections))
            await ctx.message.add_reaction('✅')
        else:
            await ctx.send(f'You are not authorized to run this command.')

    @commands.command(aliases=['servers', 'list_servers'])
    @commands.guild_only()
    @commands.check(checks.is_restricted_chan)
    async def list_ark_servers(self, ctx):
        """Returns a list of all the ARK servers in the current guild\'s config."""
        servers = json.loads(await self.bot.db_con.fetchval('select rcon_connections from guild_config '
                                                            'where guild_id = $1', ctx.guild.id))
        em = discord.Embed(style='rich',
                           title=f'__**There are currently {len(servers)} ARK servers in my config:**__',
                           color=discord.Colour.green()
                           )
        if ctx.guild.icon:
            em.set_thumbnail(url=f'{ctx.guild.icon_url}')
        for server in servers:
            description = f"""
            ឵          **IP:** {servers[server]['ip']}:{servers[server]['port']}
            ឵          **Steam Connect:** [steam://connect/{servers[server]['ip']}:{servers[server]['port']}]\
            (steam://connect/{servers[server]['ip']}:{servers[server]['port']})"""
            em.add_field(name=f'__***{server}***__', value=description, inline=False)
        await ctx.send(embed=em)


def setup(bot):
    bot.add_cog(Rcon(bot))
