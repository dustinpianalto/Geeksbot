import discord
import json
from src.imports import utils

owner_id = 351794468870946827


async def check_admin_role(bot, ctx, member):
    admin_roles = json.loads(await bot.db_concon.fetchval(f"select admin_roles from guild_config where guild_id = $1",
                                                          ctx.guild.id))
    for role in admin_roles:
        if discord.utils.get(ctx.guild.roles, id=admin_roles[role]) in member.roles:
            return True
    return member.id == ctx.guild.owner.id or member.id == owner_id


async def check_rcon_role(bot, ctx, member):
    rcon_admin_roles = json.loads(await bot.db_concon.fetchval("select rcon_admin_roles from guild_config "
                                                               "where guild_id = $1", ctx.guild.id))
    for role in rcon_admin_roles:
        if discord.utils.get(ctx.guild.roles, id=rcon_admin_roles[role]) in member.roles:
            return True
    return member.id == ctx.guild.owner.id or member.id == owner_id


async def is_admin(bot, ctx):
    admin_roles = json.loads(await bot.db_concon.fetchval("select admin_roles from guild_config where guild_id = $1",
                                                          ctx.guild.id))
    for role in admin_roles:
        if discord.utils.get(ctx.guild.roles, id=admin_roles[role]) in ctx.message.author.roles:
            return True
    return ctx.message.author.id == ctx.guild.owner.id or ctx.message.author.id == owner_id


async def is_guild_owner(ctx):
    if ctx.guild:
        return ctx.message.author.id == ctx.guild.owner.id or ctx.message.author.id == owner_id
    return False


async def is_rcon_admin(bot, ctx):
    rcon_admin_roles = json.loads(await bot.db_concon.fetchval("select rcon_admin_roles from guild_config "
                                                               "where guild_id = $1", ctx.guild.id))
    for role in rcon_admin_roles:
        if discord.utils.get(ctx.guild.roles, id=rcon_admin_roles[role]) in ctx.message.author.roles:
            return True
    return ctx.message.author.id == ctx.guild.owner.id or ctx.message.author.id == owner_id


def is_restricted_chan(ctx):
    if ctx.guild:
        if ctx.channel.overwrites_for(ctx.guild.default_role).read_messages is False:
            return True
        return False
    return False


async def is_spam(bot, ctx):
    max_rep = 5
    rep_time = 20
    spam_rep = 5
    spam_time = 3
    spam_check = 0
    msg_check = 0
    for msg in bot.recent_msgs[ctx.guild.id]:
        if msg['author'] == ctx.author:
            if msg['time'] > ctx.message.created_at.timestamp() - spam_time:
                spam_check += 1
            if msg['content'].lower() == ctx.content.lower() \
                    and msg['time'] > ctx.message.created_at.timestamp() - rep_time:
                msg_check += 1
        if spam_check == spam_rep - 1 or msg_check == max_rep - 1:
            await utils.mute(bot, ctx, admin=1, member=ctx.author.id)
            return True
    return False
