import unittest

from comm import get_type_id
import messages


class TestPBMessages(unittest.TestCase):

    def test_should_serialize_login_message(self):

        player_login = messages.PlayerLogin(username='Hans', password='foobar')

        data = bytes(player_login)

        self.assertEqual(b'\n\x04Hans\x12\x06foobar', data)

    def test_should_serialize_packet(self):

        player_login = messages.PlayerLogin(username='Hans', password='foobar')
        packet = messages.Packet(player_login=player_login)

        data = bytes(packet)

        # The bytes seen here, is the raw bytes of the Protobuf
        # this will be sent over the wire
        self.assertEqual(b'\x12\x0e\n\x04Hans\x12\x06foobar', data)

        new_packet = messages.Packet.FromString(data)
        new_player_login = new_packet.player_login
        self.assertEqual('Hans', new_player_login.username)
        self.assertEqual('foobar', new_player_login.password)
