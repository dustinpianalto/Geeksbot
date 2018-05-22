import asyncpg
import asyncio


class DatabaseConnection:
    def __init__(self, host: str='localhost', port: int=5432, database: str='', username: str='', password: str=''):
        if username == '' or password == '' or database == '':
            raise RuntimeError('Username or Password are blank')
        self.kwargs = {'host': host, 'port': port, 'database': database, 'username': username, 'password': password}
        self._conn = asyncio.get_event_loop().run_until_complete(self.acquire())
        self.fetchval = self._conn.fetchval
        self.execute = self._conn.execute
        self.fetchall = self._conn.fetchall
        self.fetchone = self._conn.fetchone

    async def acquire(self):
        if not self._conn:
            self._conn = await asyncpg.create_pool(**self.kwargs)

    async def close(self):
        await self._conn.close()
        self._conn = None
