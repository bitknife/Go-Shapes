import unittest

from comm import get_type_id
import messages


class TestPBMessages(unittest.TestCase):

    def test_should_serialize_login_message(self):

        login = messages.PlayerLogin(username='Hans', password='foobar')

        data = bytes(login)

        self.assertEqual(b'\n\x04Hans\x12\x06foobar', data)

    def test_should_serialize_mouse_message(self):

        mouse_move = messages.MouseEvent(x=1, y=1)

        data = bytes(mouse_move)

        # The bytes seen here, is the raw bytes of the Protobuf
        # this will be sent over the wire
        self.assertEqual(b'\x08\x01\x10\x01', data)

    def test_should_prefix_length_to_mouse_message(self):
        mouse_move = messages.MouseEvent(x=1, y=1)

        data = bytes(mouse_move)
        lb = bytes([len(data)])
        packet = lb+data

        # A packet is the data prefixed with the data_len (which is 4):
        self.assertEqual(b'\x04\x08\x01\x10\x01', bytes(packet))


class TestPacket(unittest.TestCase):

    def test_should_get_correct_type_ids(self):

        type_id = get_type_id(messages.PlayerLogin('a', 'b'))
        self.assertEqual(messages.MType.PLAYER_LOGIN, type_id)

        type_id = get_type_id(messages.MouseEvent(x=1, y=1))
        self.assertEqual(messages.MType.MOUSE_EVENT, type_id)
