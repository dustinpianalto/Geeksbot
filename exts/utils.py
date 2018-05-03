import discord
from discord.ext import commands
import json
import logging
import math
import psutil
from datetime import datetime, timedelta
import asyncio
import async_timeout
from .imports import checks
import pytz
import gspread
from oauth2client.service_account import ServiceAccountCredentials
import matplotlib as mpl
mpl.use('Agg')
import matplotlib.pyplot as plt
from mpl_toolkits.basemap import Basemap
from io import BytesIO

config_dir = 'config/'
admin_id_file = 'admin_ids'
extension_dir = 'extensions'
owner_id = 351794468870946827
embed_color = discord.Colour.from_rgb(49, 107, 111)
bot_config_file = 'bot_config.json'
invite_match = '(https?://)?(www.)?discord(app.com/(invite|oauth2)|.gg|.io)/[\w\d_\-?=&/]+'

utils_log = logging.getLogger('utils')
clock_emojis = ['🕛', '🕐', '🕑', '🕒', '🕓', '🕔', '🕕', '🕖', '🕗', '🕘', '🕙', '🕚']


class Utils:
    def __init__(self, bot):
        self.bot = bot

    async def _4_hour_ping(self, channel, message, wait_time):
        channel = self.bot.get_channel(channel)
        time_now = datetime.utcnow()
        await channel.send(message)
        while True:
            if time_now < datetime.utcnow() - timedelta(seconds=wait_time+10):
                await channel.send(message)
                time_now = datetime.utcnow()
            await asyncio.sleep(wait_time/100)

    async def background_ping(self, ctx, i):
        def check(message):
            return message.author == ctx.guild.me and 'ping_test' in message.content
        msg = await self.bot.wait_for('message', timeout=5, check=check)
        self.bot.ping_times[i]['rec'] = msg

    @commands.command()
    async def channel_ping(self, ctx, wait_time: float=10, message: str='=bump', channel: int=265828729970753537):
        await ctx.send('Starting Background Process.')
        self.bot.loop.create_task(self._4_hour_ping(channel, message, wait_time))

    @commands.command()
    @commands.is_owner()
    async def sysinfo(self, ctx):
        await ctx.send(f'''
        ```ml
        CPU Percentages: {psutil.cpu_percent(percpu=True)}
        Memory Usage: {psutil.virtual_memory().percent}%
        Disc Usage: {psutil.disk_usage("/").percent}%
        ```''')

    @commands.command(hidden=True)
    async def role(self, ctx, role: str):
        if ctx.guild.id == 396156980974059531 and role != 'Admin' and role != 'Admin Geeks':
            role = discord.utils.get(ctx.guild.roles, name=role)
            if role is not None:
                await ctx.message.author.add_roles(role)
                await ctx.send("Roles Updated")
            else:
                await ctx.send('Unknown Role')
        else:
            await ctx.send("You are not authorized to send this command.")

    @commands.command(aliases=['oauth', 'link'])
    @commands.cooldown(1, 5, type=commands.BucketType.user)
    async def invite(self, ctx, guy: discord.User=None):
        """Shows you the bot's invite link.
           If you pass in an ID of another bot, it gives you the invite link to that bot.
        """
        guy = guy or self.bot.user
        url = discord.utils.oauth_url(guy.id)
        await ctx.send(f'**{url}**')

    @staticmethod
    def create_date_string(time, time_now):
        diff = (time_now - time)
        date_str = time.strftime('%Y-%m-%d %H:%M:%S')
        return f"{diff.days} {'day' if diff.days == 1 else 'days'} " \
               f"{diff.seconds // 3600} {'hour' if diff.seconds // 3600 == 1 else 'hours'} " \
               f"{diff.seconds % 3600 // 60} {'minute' if diff.seconds % 3600 // 60 == 1 else 'minutes'} " \
               f"{diff.seconds % 3600 % 60} {'second' if diff.seconds % 3600 % 60 == 1 else 'seconds'} ago.\n{date_str}"

    @commands.command()
    @commands.cooldown(1, 5, type=commands.BucketType.user)
    async def me(self, ctx):
        """Prints out your user information."""
        em = discord.Embed(style='rich',
                           title=f'{ctx.author.name}#{ctx.author.discriminator} ({ctx.author.display_name})',
                           description=f'({ctx.author.id})',
                           color=embed_color)
        em.set_thumbnail(url=f'{ctx.author.avatar_url}')
        em.add_field(name=f'Highest Role:',
                     value=f'{ctx.author.top_role}',
                     inline=True)
        em.add_field(name=f'Bot:',
                     value=f'{ctx.author.bot}',
                     inline=True)
        em.add_field(name=f'Joined Guild:',
                     value=f'{self.create_date_string(ctx.author.joined_at, ctx.message.created_at)}',
                     inline=False)
        em.add_field(name=f'Joined Discord:',
                     value=f'{self.create_date_string(ctx.author.created_at, ctx.message.created_at)}',
                     inline=False)
        em.add_field(name=f'Current Status:',
                     value=f'{ctx.author.status}',
                     inline=True)
        em.add_field(name=f"Currently{' '+ctx.author.activity.type.name.title() if ctx.author.activity else ''}:",
                     value=f"{ctx.author.activity.name if ctx.author.activity else 'Not doing anything important.'}",
                     inline=True)
        count = 0
        async for message in ctx.channel.history(after=(ctx.message.created_at - timedelta(hours=1))):
            if message.author == ctx.author:
                count += 1
        em.add_field(name=f'Activity:',
                     value=f'You have sent {count} '
                           f'{"message" if count == 1 else "messages"} in the last hour to this channel.',
                     inline=False)
        await ctx.send(embed=em)

    @commands.command()
    @commands.cooldown(1, 5, type=commands.BucketType.user)
    async def user(self, ctx, member: discord.Member):
        """Prints User information.
        <member> should be in the form @Dusty.P#0001"""
        em = discord.Embed(style='rich',
                           title=f'{member.name}#{member.discriminator} ({member.display_name})',
                           description=f'({member.id})',
                           color=embed_color)
        em.set_thumbnail(url=f'{member.avatar_url}')
        em.add_field(name=f'Highest Role:',
                     value=f'{member.top_role}',
                     inline=True)
        em.add_field(name=f'Bot:',
                     value=f'{member.bot}',
                     inline=True)
        em.add_field(name=f'Joined Guild:',
                     value=f'{self.create_date_string(member.joined_at,ctx.message.created_at)}',
                     inline=False)
        em.add_field(name=f'Joined Discord:',
                     value=f'{self.create_date_string(member.created_at,ctx.message.created_at)}',
                     inline=False)
        em.add_field(name=f'Current Status:',
                     value=f'{member.status}',
                     inline=True)
        em.add_field(name=f"Currently{' '+member.activity.type.name.title() if member.activity else ''}:",
                     value=f"{member.activity.name if member.activity else 'Not doing anything important.'}",
                     inline=True)
        count = 0
        async for message in ctx.channel.history(after=(ctx.message.created_at - timedelta(hours=1))):
            if message.author == member:
                count += 1
        em.add_field(name=f'Activity:',
                     value=f'{member.display_name} has sent {count} '
                           f'{"message" if count == 1 else "messages"} in the last hour to this channel.',
                     inline=False)
        await ctx.send(embed=em)

    @commands.command()
    @commands.cooldown(1, 5, type=commands.BucketType.user)
    async def ping(self, ctx, mode='normal', count: int=2):
        """Check the Bot\'s connection to Discord"""
        em = discord.Embed(style='rich',
                           title=f'Pong 🏓',
                           color=discord.Colour.green()
                           )
        msg = await ctx.send(embed=em)
        time1 = ctx.message.created_at
        time = (msg.created_at - time1).total_seconds() * 1000
        em.description = f'''Response Time: **{math.ceil(time)}ms**
        Discord Latency: **{math.ceil(self.bot.latency*1000)}ms**'''
        await msg.edit(embed=em)

        if mode == 'comp':
            try:
                count = int(count)
            except ValueError:
                await ctx.send('Not a valid count. Must be a whole number.')
            else:
                if count > 24:
                    await ctx.send('24 Pings is the max allowed. Setting count to 24.', delete_after=5)
                    count = 24
                self.bot.ping_times = []
                times = []
                for i in range(count):
                    self.bot.ping_times.append({})
                    self.bot.loop.create_task(self.background_ping(ctx, i))
                    await asyncio.sleep(0.1)
                    self.bot.ping_times[i]['snd'] = await ctx.send('ping_test')
                    now = datetime.utcnow()
                    while 'rec' not in self.bot.ping_times[i]:
                        now = datetime.utcnow()
                        if now.timestamp() > self.bot.ping_times[i]['snd'].created_at.timestamp()+5:
                            break
                    if 'rec' in self.bot.ping_times[i]:
                        time = now - self.bot.ping_times[i]['snd'].created_at
                        time = time.total_seconds()
                        times.append(time)
                        value = f"Message Sent:" \
                                f"{datetime.strftime(self.bot.ping_times[i]['snd'].created_at, '%H:%M:%S.%f')}" \
                                f"Response Received: {datetime.strftime(now, '%H:%M:%S.%f')}" \
                                f"Total Time: {math.ceil(time * 1000)}ms"
                        await self.bot.ping_times[i]['rec'].delete()
                        em.add_field(name=f'Ping Test {i}', value=value, inline=True)
                    else:
                        em.add_field(name=f'Ping Test {i}', value='Timeout...', inline=True)
                total_time = 0
                print(times)
                for time in times:
                    total_time += time * 1000
                em.add_field(value=f'Total Time for Comprehensive test: {math.ceil(total_time)}ms',
                             name=f'Average: **{round(total_time/count,1)}ms**',
                             inline=False)
        await msg.edit(embed=em)

    @commands.group(case_insensitive=True)
    async def admin(self, ctx):
        """Run help admin for more info"""
        pass

    @admin.command(name='new', aliases=['nr'])
    @commands.cooldown(1, 30, type=commands.BucketType.user)
    async def new_admin_request(self, ctx, *, request_msg=None):
        """Submit a new request for admin assistance.
        The admin will be notified when your request is made and it will be added to the request list for this guild.
        """
        if ctx.guild:
            if request_msg is not None:
                if len(request_msg) < 1000:
                    self.bot.con.run('insert into admin_requests (issuing_member_id, guild_orig, request_text,'
                                     'request_time) values (%(member_id)s, %(guild_id)s, %(text)s, %(time)s)',
                                     {'member_id': ctx.author.id, 'guild_id': ctx.guild.id, 'text': request_msg,
                                      'time': ctx.message.created_at})
                    channel = self.bot.con.one(f'select admin_chat from guild_config where guild_id = {ctx.guild.id}')
                    if channel:
                        chan = discord.utils.get(ctx.guild.channels, id=channel)
                        msg = ''
                        admin_roles = []
                        roles = self.bot.con.one(f'select admin_roles,rcon_admin_roles from guild_config where '
                                                 f'guild_id = %(id)s', {'id': ctx.guild.id})
                        request_id = self.bot.con.one(f'select id from admin_requests where '
                                                      f'issuing_member_id = %(member_id)s and request_time = %(time)s',
                                                      {'member_id': ctx.author.id, 'time': ctx.message.created_at})
                        for item in roles:
                            i = json.loads(item)
                            for j in i:
                                if i[j] not in admin_roles:
                                    admin_roles.append(i[j])
                        for role in admin_roles:
                            msg = '{0} {1}'.format(msg, discord.utils.get(ctx.guild.roles, id=role).mention)
                        msg += f"New Request ID: {request_id} " \
                               f"{ctx.author.mention} has requested assistance: \n" \
                               f"```{request_msg}``` \n" \
                               f"Requested on: {datetime.strftime(ctx.message.created_at, '%Y-%m-%d at %H:%M:%S')} GMT"
                        await chan.send(msg)
                    await ctx.send('The Admin have received your request.')
                else:
                    await ctx.send('Request is too long, please keep your message to less than 1000 characters.')
            else:
                await ctx.send('Please include a message containing information about your request.')
        else:
            await ctx.send('This command must be run from inside a guild.')

    @admin.command(name='list', aliases=['lr'])
    @commands.cooldown(1, 5, type=commands.BucketType.user)
    async def list_admin_requests(self, ctx, assigned_to: discord.Member=None):
        """Returns a list of all active Admin help requests for this guild

        If a user runs this command it will return all the requests that they have submitted and are still open.
            - The [assigned_to] argument is ignored but will still give an error if an incorrect value is entered.

        If an admin runs this command it will return all the open requests for this guild.
            - If the [assigned_to] argument is included it will instead return all open requests that
                are assigned to the specified admin.
        """
        em = discord.Embed(style='rich',
                           title=f'Admin Help Requests',
                           color=discord.Colour.green()
                           )
        if checks.is_admin(self.bot, ctx) or checks.is_rcon_admin(self.bot, ctx):
            if assigned_to is None:
                requests = self.bot.con.all(f'select * from admin_requests where guild_orig = %(guild_id)s '
                                            f'and completed_time is null', {'guild_id': ctx.guild.id})
                em.title = f'Admin help requests for {ctx.guild.name}'
                if requests:
                    for request in requests:
                        member = discord.utils.get(ctx.guild.members, id=request[1])
                        admin = discord.utils.get(ctx.guild.members, id=request[2])
                        title = f"{'Request ID':^12}{'Requested By':^20}{'Assigned to':^20}\n" + \
                                f"{request[0]:^12}{member.display_name if member else 'None':^20}" + \
                                f"{admin.display_name if admin else 'None':^20}"
                        em.add_field(name='￲',
                                     value=f"```{title} \n\n"
                                     f"{request[4]}\n\n"
                                     f"Requested on: {datetime.strftime(request[5], '%Y-%m-%d at %H:%M:%S')} GMT```",
                                     inline=False)
                else:
                    em.add_field(name='There are no pending requests for this guild.', value='￰', inline=False)
            else:
                if checks.check_admin_role(self.bot, ctx, assigned_to)\
                        or checks.check_rcon_role(self.bot, ctx, assigned_to):
                    requests = self.bot.con.all('select * from admin_requests where assigned_to = %(admin_id)s '
                                                'and guild_orig = %(guild_id)s and completed_time is null',
                                                {'admin_id': assigned_to.id, 'guild_id': ctx.guild.id})
                    em.title = f'Admin help requests assigned to {assigned_to.display_name} in {ctx.guild.name}'
                    if requests:
                        for request in requests:
                            member = discord.utils.get(ctx.guild.members, id=request[1])
                            em.add_field(name=f"Request ID: {request[0]} Requested By:"
                                              f"{member.display_name if member else 'None'}",
                                         value=f"{request[4]} \n"
                                         f"Requested on: {datetime.strftime(request[5], '%Y-%m-%d at %H:%M:%S')} GMT\n"
                                         "￰",
                                         inline=False)
                    else:
                        em.add_field(name=f'There are no pending requests for {assigned_to.display_name} on this guild.',
                                     value='￰',
                                     inline=False)
                else:
                    em.title = f'{assigned_to.display_name} is not an admin in this guild.'
        else:
            requests = self.bot.con.all('select * from admin_requests where issuing_member_id = %(member_id)s '
                                        'and guild_orig = %(guild_id)s and completed_time is null',
                                        {'member_id': ctx.author.id, 'guild_id': ctx.guild.id})
            em.title = f'Admin help requests for {ctx.author.display_name}'
            if requests:
                for request in requests:
                    admin = discord.utils.get(ctx.guild.members, id=request[2])
                    em.add_field(name=f"Request ID: {request[0]}"
                                      f"{' Assigned to: ' + admin.display_name if admin else ''}",
                                 value=f"{request[4]}\n"
                                 f"Requested on: {datetime.strftime(request[5], '%Y-%m-%d at %H:%M:%S')} GMT\n"
                                 "￰",
                                 inline=False)
            else:
                em.add_field(name='You have no pending Admin Help requests.',
                             value='To submit a request please use `admin new <message>`',
                             inline=False)
        await ctx.send(embed=em)

    @admin.command(name='close')
    async def close_request(self, ctx, *, request_ids=None):
        """Allows Admin to close admin help tickets.
        [request_id] must be a valid integer pointing to an open Request ID
        """
        if checks.is_admin(self.bot, ctx) or checks.is_rcon_admin(self.bot, ctx):
            if request_ids:
                request_ids = request_ids.replace(' ', '').split(',')
                for request_id in request_ids:
                    try:
                        request_id = int(request_id)
                    except ValueError:
                        await ctx.send(f'{request_id} is not a valid request id.')
                    else:
                        request = self.bot.con.one(f'select * from admin_requests where id = %(request_id)s',
                                                   {'request_id': request_id})
                        if request:
                            if request[3] == ctx.guild.id:
                                if request[6] is None:
                                    self.bot.con.run('update admin_requests set completed_time = %(time_now)s where '
                                                     'id = %(request_id)s',
                                                     {'time_now': ctx.message.created_at, 'request_id': request_id})
                                    await ctx.send(f'Request {request_id} by '
                                                   f'{ctx.guild.get_member(request[1]).display_name}'
                                                   f' has been marked complete.')
                                else:
                                    await ctx.send(f'Request {request_id} is already marked complete.')
                            else:
                                await ctx.send(f'Request {request_id} is not registered to this guild.')
                        else:
                            await ctx.send(f'{request_id} is not a valid request id.')
            else:
                await ctx.send('You must include at least one request id to close.')
        else:
            await ctx.send(f'You are not authorized to run this command.')

    @commands.command(name='weather', aliases=['wu'])
    @commands.cooldown(5, 15, type=commands.BucketType.default)
    async def get_weather(self, ctx, *, location='palmer ak'):
        """Gets the weather data for the location provided,
        If no location is included then it will get the weather for the Bot's home location.
        """
        try:
            url = f'http://autocomplete.wunderground.com/aq?query={location}&format=JSON'
            with async_timeout.timeout(10):
                async with self.bot.aio_session.get(url) as response:
                    data = await response.json()
            link = data['RESULTS'][0]['l']
            url = f'http://api.wunderground.com/api/88e14343b2dd6d8e/geolookup/conditions/{link}.json'
            with async_timeout.timeout(10):
                async with self.bot.aio_session.get(url) as response:
                    data = json.loads(await response.text())
            utils_log.info(data)
            em = discord.Embed()
            em.title = f"Weather for {data['current_observation']['display_location']['full']}"
            em.url = data['current_observation']['forecast_url']
            em.description = data['current_observation']['observation_time']
            em.set_thumbnail(url=data['current_observation']['icon_url'])
            em.set_footer(text='Data provided by wunderground.com',
                          icon_url=data['current_observation']['image']['url'])
            value_str = f'''```
{'Temp:':<20}{data['current_observation']['temperature_string']:<22}
{'Feels Like:':<20}{data['current_observation']['feelslike_string']:<22}
{'Relative Humidity:':<20}{data['current_observation']['relative_humidity']:<22}
{'Wind:':<5}{data['current_observation']['wind_string']:^44}
```'''
            em.add_field(name=f"Current Conditions: {data['current_observation']['weather']}",
                         value=value_str,
                         inline=False)
            await ctx.send(embed=em)
        except IndexError:
            await ctx.send('Can\'t find that location, please try again.')

    @commands.command(name='localtime', aliases=['time', 'lt'])
    @commands.cooldown(1, 3, type=commands.BucketType.user)
    async def get_localtime(self, ctx, timezone: str='Anchorage'):
        em = discord.Embed()
        try:
            tz = pytz.timezone(timezone)
            localtime = datetime.now(tz=tz)
            em.title = f'{clock_emojis[(localtime.hour % 12)]} {tz}'
            em.description = localtime.strftime('%c')
            em.colour = embed_color
            await ctx.send(embed=em)
        except pytz.exceptions.UnknownTimeZoneError:
            for tz in pytz.all_timezones:
                if timezone.lower() in tz.lower():
                    localtime = datetime.now(tz=pytz.timezone(tz))
                    em.title = f'{clock_emojis[(localtime.hour % 12)]} {tz}'
                    em.description = localtime.strftime('%c')
                    em.colour = embed_color
                    await ctx.send(embed=em)
                    return
            em.title = 'Unknown Timezone.'
            em.colour = discord.Colour.red()
            await ctx.send(embed=em)

    @commands.command(name='purge', aliases=['clean', 'erase'])
    @commands.cooldown(1, 3, type=commands.BucketType.user)
    async def purge_messages(self, ctx, number: int=20, member: discord.Member=None):
        def is_me(message):
            if message.author == self.bot.user:
                return True
            prefixes = self.bot.con.one('select prefix from guild_config where guild_id = %(id)s', {'id': ctx.guild.id})
            if prefixes:
                for prefix in prefixes:
                    if message.content.startswith(prefix):
                        return True
                return False
            return message.content.startswith(self.bot.default_prefix)

        def is_member(message):
            return message.author == member

        def is_author(message):
            return message.author == ctx.author

        if checks.is_admin(self.bot, ctx):
            if member:
                deleted = await ctx.channel.purge(limit=number, check=is_member)
                if member != ctx.author:
                    await ctx.message.delete()
            else:
                deleted = await ctx.channel.purge(limit=number, check=is_me)
        else:
            deleted = await ctx.channel.purge(limit=number, check=is_author)
        em = discord.Embed(title='❌ Purge', colour=discord.Colour.red())
        em.description = f'Deleted {len(deleted)} messages.'
        await ctx.send(embed=em, delete_after=5)

    @commands.command(name='purge_all', aliases=['cls', 'clear'])
    @commands.cooldown(1, 3, type=commands.BucketType.user)
    async def purge_all(self, ctx, number: int=20, contents: str='all'):
        if checks.is_admin(self.bot, ctx):
            if contents != 'all':
                deleted = await ctx.channel.purge(limit=number, check=lambda message: message.content == contents)
            else:
                deleted = await ctx.channel.purge(limit=number)
            em = discord.Embed(title='❌ Purge', colour=discord.Colour.red())
            em.description = f'Deleted {len(deleted)} messages.'
            if contents != ctx.message.content and contents != 'all':
                await ctx.message.delete()
            await ctx.send(embed=em, delete_after=5)

    @commands.command(name='google', aliases=['g', 'search'])
    async def google_search(self, ctx, *, search):
        res = self.bot.gcs_service.cse().list(q=search, cx=self.bot.bot_secrets['cx']).execute()
        results = res['items'][:4]
        em = discord.Embed()
        em.title = f'Google Search'
        em.description = f'Top 4 results for "{search}"'
        em.colour = embed_color
        # TODO Fix layout of Results
        for result in results:
            em.add_field(name=f'{result["title"]}', value=f'{result["snippet"]}\n{result["link"]}')
        await ctx.send(embed=em)

    @commands.command(hidden=True, name='sheets')
    async def google_sheets(self, ctx, member: discord.Member):
        if checks.is_admin(self.bot, ctx):
            scope = ['https://spreadsheets.google.com/feeds',
                     'https://www.googleapis.com/auth/drive']
            credentials = ServiceAccountCredentials.from_json_keyfile_name('config/google_client_secret.json', scope)
            gc = gspread.authorize(credentials)
            sh = gc.open_by_key(self.bot.bot_secrets['sheet'])
            ws = sh.worksheet('Current Whitelist')
            names = ws.col_values('3')
            steam = ws.col_values('6')
            tier = ws.col_values('5')
            patron = ws.col_values('4')
            em = discord.Embed()
            em.title = f'User Data from Whitelist Sheet'
            em.colour = embed_color
            for i, name in enumerate(names):
                if member.name.lower() in name.lower()\
                        or member.display_name.lower() in name.lower()\
                        or name.lower() in member.name.lower()\
                        or name.lower() in member.display_name.lower():
                    em.add_field(name=name,
                                 value=f'Steam ID: {steam[i]}\nPatreon Level: {tier[i]}\nPatron of: {patron[i]}')
            await ctx.send(embed=em)

    @commands.command(name='iss')
    async def iss_loc(self, ctx):
        def gen_image(iss_loc):
            lat = iss_loc['latitude']
            lon = iss_loc['longitude']
            plt.figure(figsize=(8, 8))
            m = Basemap(projection='ortho', resolution=None, lat_0=lat, lon_0=lon)
            m.bluemarble(scale=0.5)
            x, y = m(lon, lat)
            plt.plot(x, y, 'ok', markersize=10, color='red')
            plt.text(x, y, '  ISS', fontsize=20, color='red')

            img = BytesIO()
            plt.savefig(img, format='png', transparent=True)
            img.seek(0)
            return img

        async with ctx.typing():
            async with self.bot.aio_session.get('https://api.wheretheiss.at/v1/satellites/25544') as response:
                loc = await response.json()
            async with self.bot.loop.run_in_executor(self.bot.tpe, gen_image, loc) as output:
                await ctx.send(file=discord.File(output, 'output.png'))

# TODO Create Help command


def setup(bot):
    bot.add_cog(Utils(bot))
