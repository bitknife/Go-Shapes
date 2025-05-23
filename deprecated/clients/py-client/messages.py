# Generated by the protocol buffer compiler.  DO NOT EDIT!
# sources: messages.proto
# plugin: python-betterproto
from dataclasses import dataclass
from typing import Dict

import betterproto


class GameObjectKind(betterproto.Enum):
    """
    GameObjectEvent:Sent by server to clients.Manages lifecycle and basic
    properties of all GameObjects.Other more specific events for certain
    classes sets morespecific attributes.Client:- SPAWNS or REMOVES
    representations.
    """

    DOT = 0
    BOX = 1
    DIAMOND = 2


@dataclass
class PlayerLogin(betterproto.Message):
    username: str = betterproto.string_field(1)
    password: str = betterproto.string_field(2)


@dataclass
class PlayerLogout(betterproto.Message):
    username: str = betterproto.string_field(1)


@dataclass
class Ping(betterproto.Message):
    sent: int = betterproto.uint64_field(1)
    bounced: int = betterproto.uint64_field(2)
    received: int = betterproto.uint64_field(3)


@dataclass
class MouseInput(betterproto.Message):
    # For communicating exact mose movements, dont use in games, better abstract
    # control type
    mouse_x: int = betterproto.int32_field(1)
    mouse_y: int = betterproto.int32_field(2)
    right_click: bool = betterproto.bool_field(3)
    left_click: bool = betterproto.bool_field(4)


@dataclass
class PlayerMovement(betterproto.Message):
    deg: int = betterproto.int32_field(1)
    vel: int = betterproto.int32_field(2)
    acc: int = betterproto.int32_field(3)


@dataclass
class PlayerAction(betterproto.Message):
    # Not sure how to model this properly, how generic it should be etc.
    primary: bool = betterproto.bool_field(1)
    secondary: bool = betterproto.bool_field(2)
    str_attrs: Dict[str, str] = betterproto.map_field(
        10, betterproto.TYPE_STRING, betterproto.TYPE_STRING
    )


@dataclass
class GameObject(betterproto.Message):
    id: str = betterproto.string_field(1)
    action: str = betterproto.string_field(2)
    # Server frame this happened in
    tick: int = betterproto.int64_field(3)
    # What class (on client typically) to spawn.
    kind: "GameObjectKind" = betterproto.enum_field(4)
    # Translation within current segment
    x: int = betterproto.int32_field(5)
    y: int = betterproto.int32_field(6)
    z: int = betterproto.int32_field(7)
    # Width and height
    w: int = betterproto.int32_field(8)
    h: int = betterproto.int32_field(9)
    # Rotation in degrees
    r: int = betterproto.int32_field(10)
    # Initial values for this specific class to be set during this action
    str_attrs: Dict[str, str] = betterproto.map_field(
        11, betterproto.TYPE_STRING, betterproto.TYPE_STRING
    )
    int_attrs: Dict[str, int] = betterproto.map_field(
        12, betterproto.TYPE_STRING, betterproto.TYPE_INT32
    )
    fl_attrs: Dict[str, float] = betterproto.map_field(
        13, betterproto.TYPE_STRING, betterproto.TYPE_FLOAT
    )


@dataclass
class SetGameObjectAttributes(betterproto.Message):
    """https://protobuf.dev/programming-guides/proto3/#maps"""

    id: str = betterproto.string_field(1)
    attributes: Dict[str, str] = betterproto.map_field(
        2, betterproto.TYPE_STRING, betterproto.TYPE_STRING
    )


@dataclass
class Packet(betterproto.Message):
    player_login: "PlayerLogin" = betterproto.message_field(1, group="payload")
    player_logout: "PlayerLogout" = betterproto.message_field(2, group="payload")
    ping: "Ping" = betterproto.message_field(3, group="payload")
    mouse_input: "MouseInput" = betterproto.message_field(4, group="payload")
    player_movement: "PlayerMovement" = betterproto.message_field(5, group="payload")
    player_action: "PlayerAction" = betterproto.message_field(6, group="payload")
    game_object: "GameObject" = betterproto.message_field(10, group="payload")
    set_game_object_attributes: "SetGameObjectAttributes" = betterproto.message_field(
        11, group="payload"
    )
