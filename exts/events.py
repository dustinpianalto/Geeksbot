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
        await self.bot.db_con.execute('update messages set deleted_at = $1 where id = $2',
                                      datetime.utcnow(), msg_id)

    async def on_raw_bulk_message_delete(self, msg_ids, chan_id):
        sql = await self.bot.db_con.prepare('update messages set deleted_at = $1 where id = $2')
        for msg_id in msg_ids:
            await sql.execute(datetime.utcnow(), msg_id)

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
        await self.bot.db_con.execute(sql, *msg_data)
        if ctx.guild:
            if ctx.author != ctx.guild.me:
                if await self.bot.db_con.fetchval("select pg_filter from guild_config where guild_id = $1", ctx.guild.id):
                    profane = 0
                    for word in await self.bot.db_con.fetchval('select profane_words from guild_config '
                                                               'where guild_id = $1', ctx.guild.id):
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
        await self.bot.db_con.execute('update messages set reactions = $2 where id = $1',
                                      react.message.id, reacts)

    async def on_message_edit(self, before, ctx):
        previous_content = await self.bot.db_con.fetchval('select previous_content from messages where id = $1', ctx.id)
        if previous_content:
            previous_content.append(before.content)
        else:
            previous_content = [before.content]
        previous_embeds = await self.bot.db_con.fetchval('select previous_embeds from messages where id = $1', ctx.id)
        if previous_embeds:
            previous_embeds.append([json.dumps(e.to_dict()) for e in before.embeds])
        else:
            previous_embeds = [[json.dumps(e.to_dict()) for e in before.embeds]]
        sql = 'update messages set (edited_at, previous_content, previous_embeds, tts, type, content, embeds, ' \
              'channel, mention_everyone, mentions, channel_mentions, role_mentions, webhook, attachments, pinned, ' \
              'reactions, guild, created_at, system_content, author) = ' \
              '($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)' \
              'where id = $21'
        msg_data = [datetime.utcnow(), previous_content, previous_embeds, ctx.tts, str(ctx.type), ctx.content,
                    [json.dumps(e.to_dict()) for e in ctx.embeds], ctx.channel.id, ctx.mention_everyone,
                    [user.id for user in ctx.mentions], [channel.id for channel in ctx.channel_mentions],
                    [role.id for role in ctx.role_mentions], ctx.webhook_id,
                    [json.dumps({'id': a.id, 'size': a.size, 'height': a.height, 'width': a.width,
                                 'filename': a.filename, 'url': a.url}) for a in ctx.attachments], ctx.pinned,
                    [json.dumps({'emoji': r.emoji, 'count': r.count}) for r in ctx.reactions], ctx.guild.id,
                    ctx.created_at, ctx.system_content, ctx.author.id, ctx.id]
        await self.bot.db_con.execute(sql, *msg_data)

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
        await self.bot.db_con.execute("insert into guild_config(guild_id, guild_name, admin_roles, rcon_enabled, "
                                      "channel_lockdown, raid_status, pg_filter, patreon_enabled, referral_enabled) "
                                      "values ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
                                      default_config['guild_id'], default_config['name'],
                                      json.dumps(default_config['admin_roles']), default_config['rcon_enabled'],
                                      default_config['channel_lockdown'], default_config['raid_status'],
                                      default_config['pg_filter'], default_config['patreon_enabled'],
                                      default_config['referral_enabled'])
        events_log.info(f'Entry Created for {guild.name}')
        await guild.me.edit(nick='[g$] Geeksbot')

    async def on_guild_remove(self, guild):
        await self.bot.db_con.execute(f'delete from guild_config where guild_id = $1', guild.id)
        events_log.info(f'Left the {guild.name} guild.')

    async def on_member_join(self, member):
        events_log.info(f'Member joined: {member.name} {member.id} Guild: {member.guild.name} {member.guild.id}')
        join_chan = await self.bot.db_con.fetchval('select join_leave_chat from guild_config where guild_id = $1',
                                                   member.guild.id)
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
        mem_data = [member.id,
                    member.name,
                    member.discriminator,
                    member.bot
                    ]
        mem = await self.bot.db_con.fetchval('select guilds,nicks from user_data where id = $1', member.id)
        if mem:
            mem[1].append(json.dumps({member.guild.id: member.display_name}))
            mem[0].append(member.guild.id)
            mem_data.append(mem[1])
            mem_data.append(mem[0])
            self.bot.con.run('update user_data set (name, discriminator, bot, nicks, guilds) = '
                             '($2, $3, $4, $5, $6) where id = $1', *mem_data)
        else:
            mem_data.append([json.dumps({member.guild.id: member.display_name})])
            mem_data.append([member.guild.id])
            self.bot.con.run('insert into user_data (id, name, discriminator, bot, nicks, guilds) '
                             'values ($1, $2, $3, $4, $5, $6)', *mem_data)

    async def on_member_remove(self, member):
        leave_time = datetime.utcnow()
        events_log.info(f'Member left: {member.name} {member.id} Guild: {member.guild.name} {member.guild.id}')
        join_chan = await self.bot.db_con.fetchval('select join_leave_chat from guild_config where guild_id = $1',
                                                   member.guild.id)
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
        mem_data = [member.id,
                    member.name,
                    member.discriminator,
                    member.bot
                    ]
        mem = await self.bot.db_con.fetchrow('select guilds,nicks from user_data where id = $1', member.id)
        if mem:
            mem[0].remove(member.guild.id)
            mem[1].remove(json.dumps({member.guild.id: member.display_name}))
            mem_data.append(mem[1])
            mem_data.append(mem[0])
            self.bot.con.run('update user_data set (name, discriminator, bot, nicks, guilds) = '
                             '($2, $3, $4, $5, $6) where id = $1', *mem_data)


def setup(bot):
    bot.add_cog(BotEvents(bot))
