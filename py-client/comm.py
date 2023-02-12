import queue
import messages
import tcp_client
import threading

comm_queue = queue.Queue()


def connect():

    # Lets connect to server first to see that stuff works
    gc = tcp_client.MessageTCPClient.get_instance(debug=True)

    gc.connect()


def send_to_server(game_message):
    #
    # Thread safe method called from whoever wants to send a packet
    # by passing it onto a queue that is
    #
    item = {
        'direction': 'SERVER',
        'game_message': game_message
    }
    comm_queue.put(item)


def _build_packet(game_message):
    """
    This is close to the wire, called (only?) from a low layer (TCP/UDP) to
    build a tight packet:

    [1 byte packet length] + [1 byte packet type] + [N bytes packet data]

    :param game_message:
    :return:
    """
    gm_id = get_type_id(game_message)
    msg_bytes = bytes(game_message)
    header = bytes([len(msg_bytes), gm_id])
    return header + msg_bytes


def get_type_id(game_message):

    match(type(game_message)):
        case messages.PlayerLogin:
            return messages.MType.PLAYER_LOGIN

        case messages.MouseEvent:
            return messages.MType.MOUSE_EVENT

        case _:
            raise Exception("Unsupported message type: %s" % type(game_message))


def process_queue():
    #
    # Will block until an item arrives
    #
    gc = tcp_client.MessageTCPClient.get_instance(debug=True)

    while True:
        item = comm_queue.get()

        if item['direction'] == 'SERVER':
            packet = _build_packet(item['game_message'])
            gc.send(packet)
        elif item['direction'] == 'CLIENT':
            # From server, coming here
            pass


def receive_packets():
    gc = tcp_client.MessageTCPClient.get_instance(debug=True)

    while True:
        gc.receive(comm_queue)


def start_message_queue_thread():
    t = threading.Thread(target=process_queue)
    t.start()
    return t


def start_receiver_thread():
    t = threading.Thread(target=receive_packets)
    t.start()
    return t
