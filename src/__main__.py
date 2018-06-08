from typing import Dict

import discord
from discord.ext import commands
import logging
from datetime import datetime
import json
import aiohttp
from googleapiclient.discovery import build
from concurrent import futures
from src.shared_libs import database


log_format = '{asctime}.{msecs:03.0f}|{levelname:<8}|{name}::{message}'
date_format = '%Y.%m.%d %H.%M.%S'

log_dir = 'logs'

log_file = '{0}/geeksbot_{1}.log'.format(log_dir, datetime.now().strftime('%Y%m%d_%H%M%S%f'))

logging.basicConfig(level=logging.DEBUG, style='{', filename=log_file, datefmt=date_format, format=log_format)
console_handler = logging.StreamHandler()
console_handler.setLevel(logging.INFO)
formatter = logging.Formatter(log_format, style='{', datefmt=date_format)
console_handler.setFormatter(formatter)
logging.getLogger('').addHandler(console_handler)

config_dir = 'src/config/'
admin_id_file = 'admin_ids'
extension_dir = 'exts'
owner_id = 351794468870946827
bot_config_file = 'bot_config.json'
secrets_file = 'bot_secrets.json'
profane_words_file = 'profane_words'

emojis: Dict[str, str] = {
    'x': '❌',
    'y': '✅',
    'poop': '💩',
    'boom': '💥',
}


class Geeksbot(commands.Bot):
    def __init__(self, **kwargs):
        kwargs["command_prefix"] = self.get_custom_prefix
        self.description = 'I am Geeksbot Dev! Fear me!\nI might just break and take you with me :P'
        kwargs['description'] = self.description
        super().__init__(**kwargs)
        self.aio_session = aiohttp.ClientSession(loop=self.loop)
        with open(f'{config_dir}{bot_config_file}') as file:
            self.bot_config = json.load(file)
        with open(f'{config_dir}{secrets_file}') as file:
            self.bot_secrets = json.load(file)
        self.guild_config = {}
        self.infected = {}
        self.TOKEN = self.bot_secrets['token']
        self.embed_color = discord.Colour.from_rgb(49, 107, 111)
        del self.bot_secrets['token']
        self.db_con = database.DatabaseConnection(**self.bot_secrets['db_con'])
        self.default_prefix = 'g~'
        self.voice_chans = {}
        self.spam_list = {}
        self.owner_id = 351794468870946827
        self.__version__ = '0.5.0a'
        self.gcs_service = build('customsearch', 'v1', developerKey=self.bot_secrets['google_search_key'])
        self.tpe = futures.ThreadPoolExecutor()
        self.geo_api = '2d4e419c2be04c8abe91cb5dd1548c72'
        self.unicode_emojis: Dict[str, str] = {
                                        'x': '❌',
                                        'y': '✅',
                                        'poop': '💩',
                                        'boom': '💥',
                                        'left_fist': '🤛',
                                        'lock': '🔒',
                                        }
        self.book_emojis: Dict[str, str] = {
                                        'unlock': '🔓',
                                        'start': '⏮',
                                        'back': '◀',
                                        'hash': '#\N{COMBINING ENCLOSING KEYCAP}',
                                        'forward': '▶',
                                        'end': '⏭',
                                        'close': '🇽',
                                        }

    async def logout(self):
        await self.db_con.close()
        super().logout()

    @staticmethod
    async def get_custom_prefix(bot_inst, message):
        return await bot_inst.db_con.fetchval('select prefix from guild_config where guild_id = $1',
                                              message.guild.id) or bot_inst.default_prefix

    async def load_ext(self, ctx, mod=None):
        self.load_extension('src.{0}.{1}'.format(extension_dir, mod))
        if ctx is not None:
            await ctx.send('{0} loaded.'.format(mod))

    async def unload_ext(self, ctx, mod=None):
        self.unload_extension('src.{0}.{1}'.format(extension_dir, mod))
        if ctx is not None:
            await ctx.send('{0} unloaded.'.format(mod))

    async def close(self):
        await super().close()
        self.aio_session.close()  # aiohttp is drunk and can't decide if it's a coro or not


bot = Geeksbot(case_insensitive=True)


@bot.command(hidden=True)
@commands.is_owner()
async def load(ctx, mod=None):
    """Allows the owner to load extensions dynamically"""
    await bot.load_ext(ctx, mod)


@bot.command(hidden=True)
@commands.is_owner()
async def reload(ctx, mod=None):
    """Allows the owner to reload extensions dynamically"""
    if mod == 'all':
        load_list = bot.bot_config['load_list']
        for load_item in load_list:
            await bot.unload_ext(ctx, f'{load_item}')
            await bot.load_ext(ctx, f'{load_item}')
    else:
        await bot.unload_ext(ctx, mod)
        await bot.load_ext(ctx, mod)


@bot.command(hidden=True)
@commands.is_owner()
async def unload(ctx, mod):
    """Allows the owner to unload extensions dynamically"""
    await bot.unload_ext(ctx, mod)


@bot.event
async def on_message(ctx):
    if not ctx.author.bot:
        if ctx.guild:
            if int(await bot.db_con.fetchval("select channel_lockdown from guild_config where guild_id = $1",
                                             ctx.guild.id)):
                if ctx.channel.id in json.loads(await bot.db_con.fetchval("select allowed_channels from guild_config "
                                                                          "where guild_id = $1",
                                                                          ctx.guild.id)):
                    await bot.process_commands(ctx)
            elif ctx.channel.id == 418452585683484680:
                prefix = await bot.db_con.fetchval('select prefix from guild_config where guild_id = $1', ctx.guild.id)
                prefix = prefix[0] if prefix else bot.default_prefix
                ctx.content = f'{prefix}{ctx.content}'
                await bot.process_commands(ctx)
            else:
                await bot.process_commands(ctx)
        else:
            await bot.process_commands(ctx)


@bot.event
async def on_ready():
    bot.remove_command('help')
    bot.recent_msgs = {}
    logging.info('Logged in as {0.name}|{0.id}'.format(bot.user))
    load_list = bot.bot_config['load_list']
    for load_item in load_list:
        await bot.load_ext(None, f'{load_item}')
        logging.info('Extension Loaded: {0}'.format(load_item))
    with open(f'{config_dir}reboot', 'r') as f:
        reboot = f.readlines()
    if int(reboot[0]) == 1:
        await bot.get_channel(int(reboot[1])).send('Restart Finished.')
    with open(f'{config_dir}reboot', 'w') as f:
        f.write(f'0')
    logging.info('Done loading, Geeksbot is active.')

bot.run(bot.TOKEN)
