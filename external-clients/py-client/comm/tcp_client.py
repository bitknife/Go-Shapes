import random
import socket
import threading


class MessageTCPClient(object):
    """
        Singleton representing/abstracting interactions with the server.

        Will _NOT_ connect upon creation.

        Does NOT assume/know of the specific message format. It just sends and receives
        a string of bytes, optionally prefixing the send operation with the size of the
        given buffer.

    """

    _instance = None

    # TODO: Supplied from server upon connect
    id = random.randint(0, 2**31)

    def __init__(self, hostname="127.0.0.1", port=7777, debug=False):
        self.server = hostname
        self.port = port
        self.addr = (self.server, self.port)
        self.debug = debug

        # SOCK_STREAM means TCP socket
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sock_lock = threading.RLock()

    @staticmethod
    def get_instance(**kwargs):
        if not MessageTCPClient._instance:
            MessageTCPClient._instance = MessageTCPClient(**kwargs)
        return MessageTCPClient._instance

    def connect(self):
        self.sock.connect(self.addr)

    def send(self, packet):
        with self.sock_lock:
            try:
                self.sock.sendall(packet)
            except socket.error as e:
                print(e)

    def receive(self, queue):
        size_byte = self.sock.recv(1)
        packet_size = int.from_bytes(size_byte, "big")
        packet_data = self.sock.recv(int(packet_size))

        # print("Received packet, len: %s" % (len(packet_data)))

        item = {
            'direction': 'CLIENT',
            'packet_data': packet_data
        }
        queue.put(item)
