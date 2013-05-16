from runner import ServerManager
import time

NUMSERVERS = 10

def help():
    result =  " show_servers      | ss         Show available servers\n"
    result += " bring_up [port]   | bu [port]  Bring up a server\n"
    result += " kill [port]       | k [port]   Kill a server\n"
    result += " quit              | q          Quit\n"
    return result

if __name__ == "__main__":
    try:
        manager = ServerManager()
        print "Done with starting manager"
        ports = [1234]
        manager.async_start_block_server(1234, 1337)
        print "Spun up primary"

        for i in range(NUMSERVERS):
            manager.async_start_block_server(1235 + i)
            ports.append(1235 + i)

        while 1:
            inp = raw_input("(modata) ")
            if inp == "help" or inp == "h":
                print help()
            elif inp == "bring_up" or inp == "bu":
                port = int(raw_input("port? "))
                manager.async_start_block_server(port)
                ports.append(port)
            elif inp == "show_servers" or inp == "ss":
                print ports
            elif inp == "kill" or inp == "k":
                port = int(raw_input("port? "))
                manager.kill_block_server(port)
                ports.remove(port)
            elif inp == "q" or inp == "quit":
                break
            time.sleep(0)
    except:
        print "Killing all servers"
        manager.kill_all()

