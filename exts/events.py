import discord
import logging
from datetime import datetime
import json
import re
from .imports import utils

config_dir = 'config/'
admin_id_file = 'admin_ids'
extension_dir = 'extensions'
owner_id = 351794468870946827
guild_config_dir = 'guild_config/'
rcon_config_file = 'server_rcon_config'
dododex_url = 'http://www.dododex.com'
embed_color = discord.Colour.from_rgb(49, 107, 111)
red_color = discord.Colour.from_rgb(142, 29, 31)
bot_config_file = 'bot_config'
default_guild_config_file = 'default_guild_config.json'

events_log = logging.getLogger('events')

emojis = {
    'x': '❌',
    'y': '✅',
    'poop': '💩',
    'crown': '👑',
    'eggplant': '🍆',
    'sob': '😭',
    'trident': '🔱'
}


class BotEvents:
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

    async def on_raw_message_delete(self, msg_id, chan_id):
        self.bot.con.run('update messages set deleted_at = %(time)s where id = %(id)s',
                         {'time': datetime.utcnow(), 'id': msg_id})

    async def on_raw_bulk_message_delete(self, msg_ids, chan_id):
        sql = ''
        for msg_id in msg_ids:
            sql += f';update messages set deleted_at = %(time)s where id = {msg_id}'
        self.bot.con.run(sql, {'time': datetime.utcnow()})

    async def on_message(self, ctx):
        try:
            if ctx.author in self.bot.infected:
                if datetime.now().timestamp() > self.bot.infected[ctx.author][1] + 300:
                    del self.bot.infected[ctx.author]
                    # await ctx.channel.send(f'{ctx.author.mention} You have been healed.')
                else:
                    await ctx.add_reaction(self.bot.infected[ctx.author][0])
        except Exception:
            pass
        sql = 'insert into messages (id, tts, type, content, embeds, channel, mention_everyone, mentions,\
        channel_mentions, role_mentions, webhook, attachments, pinned, reactions, guild, created_at,\
        system_content, author) \
        values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)'
        msg_data = [ctx.id, ctx.tts, str(ctx.type), ctx.content, [json.dumps(e.to_dict()) for e in ctx.embeds],
                    ctx.channel.id, ctx.mention_everyone, [user.id for user in ctx.mentions],
                    [channel.id for channel in ctx.channel_mentions], [role.id for role in ctx.role_mentions],
                    ctx.webhook_id, [json.dumps({'id': a.id, 'size': a.size, 'height': a.height, 'width': a.width,
                                     'filename': a.filename, 'url': a.url}) for a in ctx.attachments],
                    ctx.pinned, [json.dumps({'emoji': r.emoji, 'count': r.count}) for r in ctx.reactions],
                    ctx.guild.id, ctx.created_at, ctx.system_content, ctx.author.id]
        await self.bot.db_con.execute(sql, msg_data)
        if ctx.guild:
            if ctx.author != ctx.guild.me:
                if self.bot.con.one(f"select pg_filter from guild_config where guild_id = {ctx.guild.id}"):
                    profane = 0
                    for word in await self.bot.db_con.fetchval('select profane_words from guild_config '
                                                               'where guild_id = %(id)s',
                                                               {'id': ctx.guild.id}):
                        word = word.strip()
                        if word in ctx.content.lower():
                            events_log.info(f'Found non PG word {word}')
                            repl_str = '\*' * (len(word) - 1)
                            re_replace = re.compile(re.escape(word), re.IGNORECASE)
                            ctx.content = re_replace.sub(f'{word[:1]}{repl_str}', ctx.content)
                            profane = 1
                    if profane:
                        hook = await ctx.channel.create_webhook(name="PG Filter")
                        await ctx.delete()
                        if len(ctx.author.display_name) < 2:
                            username = f'឵{ctx.author.display_name}'
                        else:
                            username = ctx.author.display_name
                        await hook.send(ctx.content, username=username, avatar_url=ctx.author.avatar_url)
                        await hook.delete()

    async def on_reaction_add(self, react, user):
        if react.emoji == emojis['poop'] and react.message.author.id == 351794468870946827:
            await react.message.remove_reaction(emojis['poop'], user)
            await react.message.channel.send(f"You can't Poop on my Owner {user.mention} :P")
        if react.emoji == emojis['poop'] and react.message.author.id == 396588996706304010:
            await react.message.remove_reaction(emojis['poop'], user)
            await react.message.channel.send(f"You can't Poop on me {user.mention} :P")
        reactions = react.message.reactions
        reacts = [json.dumps({'emoji': r.emoji, 'count': r.count}) for r in reactions]
        self.bot.con.run('update messages set reactions = %(reacts)s where id = %(id)s',
                         {'id': react.message.id, 'reacts': reacts})

    async def on_message_edit(self, before, ctx):
        previous_content = self.bot.con.one('select previous_content from messages where id = %(id)s', {'id': ctx.id})
        if previous_content:
            previous_content.append(before.content)
        else:
            previous_content = [before.content]
        previous_embeds = self.bot.con.one('select previous_embeds from messages where id = %(id)s', {'id': ctx.id})
        if previous_embeds:
            previous_embeds.append([json.dumps(e.to_dict()) for e in before.embeds])
        else:
            previous_embeds = [[json.dumps(e.to_dict()) for e in before.embeds]]
        sql = 'update messages set (edited_at, previous_content, previous_embeds, tts, type, content,\
               embeds, channel, mention_everyone, mentions, channel_mentions, role_mentions, webhook,\
               attachments, pinned, reactions, guild, created_at, system_content, author) \
               = (%(edited_at)s, %(previous_content)s, %(previous_embeds)s, %(tts)s, %(type)s, %(content)s,\
               %(embeds)s, %(channel)s, %(mention_everyone)s, %(mentions)s, %(channel_mentions)s, %(role_mentions)s,\
               %(webhook)s, %(attachments)s, %(pinned)s, %(reactions)s, %(guild)s, %(created_at)s, %(system_content)s,\
               %(author)s) where id = %(id)s'
        msg_data = dict()
        msg_data['id'] = ctx.id
        msg_data['tts'] = ctx.tts
        msg_data['type'] = str(ctx.type)
        msg_data['content'] = ctx.content
        msg_data['embeds'] = [json.dumps(e.to_dict()) for e in ctx.embeds]
        msg_data['channel'] = ctx.channel.id
        msg_data['mention_everyone'] = ctx.mention_everyone
        msg_data['mentions'] = [user.id for user in ctx.mentions]
        msg_data['channel_mentions'] = [channel.id for channel in ctx.channel_mentions]
        msg_data['role_mentions'] = [role.id for role in ctx.role_mentions]
        msg_data['webhook'] = ctx.webhook_id
        msg_data['attachments'] = ctx.attachments
        msg_data['pinned'] = ctx.pinned
        msg_data['guild'] = ctx.guild.id
        msg_data['created_at'] = ctx.created_at
        msg_data['system_content'] = ctx.system_content
        msg_data['author'] = ctx.author.id
        msg_data['reactions'] = [json.dumps({'emoji': r.emoji, 'count': r.count}) for r in ctx.reactions]
        msg_data['previous_content'] = previous_content
        msg_data['previous_embeds'] = previous_embeds
        msg_data['edited_at'] = datetime.utcnow()
        self.bot.con.run(sql, msg_data)

    @staticmethod
    async def on_command_error(ctx, error):
        if ctx.channel.id == 418452585683484680 and type(error) == discord.ext.commands.errors.CommandNotFound:
            return
        for page in utils.paginate(error):
            await ctx.send(page)

    async def on_guild_join(self, guild):
        with open(f"{config_dir}{default_guild_config_file}", 'r') as file:
            default_config = json.loads(file.read())
        admin_role = guild.role_hierarchy[0]
        default_config['admin_roles'] = {admin_role.name: admin_role.id}
        default_config['name'] = guild.name.replace("'", "\\'")
        default_config['guild_id'] = guild.id
        events_log.info(default_config)
        self.bot.con.run("insert into guild_config(guild_id, guild_name, admin_roles, rcon_enabled, channel_lockdown,\
                          raid_status, pg_filter, patreon_enabled, referral_enabled)\
                          values (%(guild_id)s, %(name)s, %(admin_roles)s, %(rcon_enabled)s, %(channel_lockdown)s,\
                          %(raid_status)s, %(pg_filter)s, %(patreon_enabled)s, %(referral_enabled)s)",
                         {'guild_id': default_config['guild_id'],
                          'name': default_config['name'],
                          'admin_roles': json.dumps(default_config['admin_roles']),
                          'rcon_enabled': default_config['rcon_enabled'],
                          'channel_lockdown': default_config['channel_lockdown'],
                          'raid_status': default_config['raid_status'],
                          'pg_filter': default_config['pg_filter'],
                          'patreon_enabled': default_config['patreon_enabled'],
                          'referral_enabled': default_config['referral_enabled']
                          })
        events_log.info(f'Entry Created for {guild.name}')
        await guild.me.edit(nick='[g$] Geeksbot')

    async def on_guild_remove(self, guild):
        self.bot.con.run(f'delete from guild_config where guild_id = %(id)s', {'id': guild.id})
        events_log.info(f'Left the {guild.name} guild.')

    async def on_member_join(self, member):
        events_log.info(f'Member joined: {member.name} {member.id} Guild: {member.guild.name} {member.guild.id}')
        join_chan = self.bot.con.one('select join_leave_chat from guild_config where guild_id = %(id)s',
                                     {'id': member.guild.id})
        if join_chan:
            em = discord.Embed(style='rich',
                               color=embed_color
                               )
            em.set_thumbnail(url=member.avatar_url)
            em.add_field(name=f'Welcome {member.name}#{member.discriminator}', value=member.id, inline=False)
            em.add_field(name='User created on:', value=member.created_at.strftime('%Y-%m-%d at %H:%M:%S GMT'),
                         inline=True)
            em.add_field(name='Bot:', value=str(member.bot))
            em.set_footer(text=f"{member.guild.name} | {member.joined_at.strftime('%Y-%m-%d at %H:%M:%S GMT')}",
                          icon_url=member.guild.icon_url)
            await discord.utils.get(member.guild.channels, id=join_chan).send(embed=em)
        mem_data = {'id': member.id,
                    'name': member.name,
                    'discriminator': member.discriminator,
                    'bot': member.bot
                    }
        mem = self.bot.con.one('select guilds,nicks from user_data where id = %(id)s', {'id': member.id})
        if mem:
            mem[1].append(json.dumps({member.guild.id: member.display_name}))
            mem[0].append(member.guild.id)
            mem_data['nicks'] = mem[1]
            mem_data['guilds'] = mem[0]
            self.bot.con.run('update user_data set (name, discriminator, bot, nicks, guilds) =\
                              (%(name)s, %(discriminator)s, %(bot)s, %(nicks)s, %(guilds)s) where\
                               id = %(id)s', mem_data)
        else:
            mem_data['nicks'] = [json.dumps({member.guild.id: member.display_name})]
            mem_data['guilds'] = [member.guild.id]
            self.bot.con.run('insert into user_data (id, name, discriminator, bot, nicks, guilds) values\
                              (%(id)s, %(name)s, %(discriminator)s, %(bot)s, %(nicks)s, %(guilds)s)', mem_data)

    async def on_member_remove(self, member):
        leave_time = datetime.utcnow()
        events_log.info(f'Member left: {member.name} {member.id} Guild: {member.guild.name} {member.guild.id}')
        join_chan = self.bot.con.one('select join_leave_chat from guild_config where guild_id = %(id)s',
                                     {'id': member.guild.id})
        if join_chan:
            em = discord.Embed(style='rich',
                               color=red_color
                               )
            em.set_thumbnail(url=member.avatar_url)
            em.add_field(name=f'RIP {member.name}#{member.discriminator}', value=member.id, inline=False)
            join_time = member.joined_at
            em.add_field(name='Joined on:', value=join_time.strftime('%Y-%m-%d at %H:%M:%S GMT'), inline=True)
            em.add_field(name='Bot:', value=str(member.bot), inline=True)
            em.add_field(name='Left on:', value=leave_time.strftime('%Y-%m-%d at %H:%M:%S GMT'), inline=False)
            total_time = leave_time - join_time
            days, remainder = divmod(total_time.total_seconds(), 86400)
            hours, remainder = divmod(remainder, 3600)
            minutes, seconds = divmod(remainder, 60)
            days_str = str(int(days)) + ' days, ' if days != 0 else ''
            hours_str = str(int(hours)) + ' hours, ' if hours != 0 else ''
            minutes_str = str(int(minutes)) + ' minutes and ' if minutes != 0 else ''
            time_str = f"{days_str}{hours_str}{minutes_str}{int(seconds)} seconds"
            em.add_field(name='Total time in Guild:', value=time_str, inline=False)
            em.set_footer(text=f"{member.guild.name} | {datetime.utcnow().strftime('%Y-%m-%d at %H:%M:%S GMT')}",
                          icon_url=member.guild.icon_url)
            await discord.utils.get(member.guild.channels, id=join_chan).send(embed=em)


def setup(bot):
    bot.add_cog(BotEvents(bot))
