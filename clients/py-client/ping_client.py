import time

from messages import PlayerLogin, Ping, Packet

from comm import connect, send_to_server, start_message_queue_thread, start_receiver_thread


USERNAME = 'player1'
PASSWORD = 'welcome'

PING_INTERVAL = 1


def main():

    connect()

    # Start the background threads needed before c
    start_message_queue_thread()

    start_receiver_thread()

    # Login
    print("Logging in...")
    packet = Packet(player_login=PlayerLogin(USERNAME, PASSWORD))
    send_to_server(bytes(packet))

    print("Sending pings, interval: %s sec" % PING_INTERVAL)
    while True:
        packet = Packet(ping=Ping(sent=int(time.time() * (10 ** 6))))
        send_to_server(packet)
        time.sleep(PING_INTERVAL)


if __name__ == "__main__":
    main()
