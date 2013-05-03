from runner import ServerManager
import time

try:
    manager = ServerManager()
    manager.async_start_block_server(1234)
    manager.async_start_block_server(1235)
    manager.async_start_block_server(1236)
    manager.async_start_block_server(1237)
    while 1:
        time.sleep(0)
except:
    manager.kill_all()

