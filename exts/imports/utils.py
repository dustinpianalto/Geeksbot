from io import StringIO
import sys
import asyncio
import discord
from discord.ext.commands.formatter import Paginator
from . import checks
import re
import typing


class Capturing(list):
    def __enter__(self):
        self._stdout = sys.stdout
        sys.stdout = self._stringio = StringIO()
        return self

    def __exit__(self, *args):
        self.extend(self._stringio.getvalue().splitlines())
        del self._stringio    # free up some memory
        sys.stdout = self._stdout


async def mute(bot, ctx, admin=0, member_id=None):
    mute_role = bot.db_con.fetchval(f'select muted_role from guild_config where guild_id = $1', ctx.guild.id)
    if mute_role:
        if admin or await checks.is_admin(bot, ctx):
            if ctx.guild.me.guild_permissions.manage_roles:
                if member_id:
                    ctx.guild.get_member(member_id).edit(roles=[discord.utils.get(ctx.guild.roles, id=mute_role)])


def to_list_of_str(items, out: list=list(), level=1, recurse=0):
    # noinspection PyShadowingNames
    def rec_loop(item, key, out, level):
        quote = '"'
        if type(item) == list:
            out.append(f'{"    "*level}{quote+key+quote+": " if key else ""}[')
            new_level = level + 1
            out = to_list_of_str(item, out, new_level, 1)
            out.append(f'{"    "*level}]')
        elif type(item) == dict:
            out.append(f'{"    "*level}{quote+key+quote+": " if key else ""}{{')
            new_level = level + 1
            out = to_list_of_str(item, out, new_level, 1)
            out.append(f'{"    "*level}}}')
        else:
            out.append(f'{"    "*level}{quote+key+quote+": " if key else ""}{repr(item)},')

    if type(items) == list:
        if not recurse:
            out = list()
            out.append('[')
        for item in items:
            rec_loop(item, None, out, level)
        if not recurse:
            out.append(']')
    elif type(items) == dict:
        if not recurse:
            out = list()
            out.append('{')
        for key in items:
            rec_loop(items[key], key, out, level)
        if not recurse:
            out.append('}')

    return out


def replace_text_ignorecase(in_str: str, old: str, new: str='') -> str:
    re_replace = re.compile(re.escape(old), re.IGNORECASE)
    return re_replace.sub(f'{new}', in_str)


def paginate(text, maxlen=1990):
    paginator = Paginator(prefix='```py', max_size=maxlen+10)
    if type(text) == list:
        data = to_list_of_str(text)
    elif type(text) == dict:
        data = to_list_of_str(text)
    else:
        data = str(text).split('\n')
    for line in data:
        if len(line) > maxlen:
            n = maxlen
            for l in [line[i:i+n] for i in range(0, len(line), n)]:
                paginator.add_line(l)
        else:
            paginator.add_line(line)
    return paginator.pages


async def run_command(args):
    # Create subprocess
    process = await asyncio.create_subprocess_shell(
        args,
        # stdout must a pipe to be accessible as process.stdout
        stdout=asyncio.subprocess.PIPE)
    # Wait for the subprocess to finish
    stdout, stderr = await process.communicate()
    # Return stdout
    return stdout.decode().strip()


# TODO Add Paginator
class Paginator:
    def __init__(self,
                 max_chars: int=1990,
                 max_lines: int=20,
                 prefix: str='```md',
                 suffix: str='```',
                 page_break: str=''):
        assert 0 < max_lines <= max_chars

        self._parts = list()
        self._prefix = prefix
        self._suffix = suffix
        self._max_chars = max_chars if max_chars + len(prefix) + len(suffix) + 2 <= 2000 \
            else 2000 - len(prefix) - len(suffix) - 2
        self._max_lines = max_lines - (prefix + suffix).count('\n') + 1
        self._page_break = page_break

