import cocos
from pyglet.window.mouse import buttons_string

from comm import send_to_server
from messages import MouseEvent


# Subclass a Layer and define the logic of you program here:
class MouseDisplay(cocos.layer.Layer):
    is_event_handler = True  #: enable director.window events

    def __init__(self):
        super(MouseDisplay, self).__init__()

        self.posx = 100
        self.posy = 240
        self.text = cocos.text.Label('No mouse events yet', font_size=18, x=self.posx, y=self.posy)
        self.add(self.text)

    def update_text(self, x, y):
        text = 'Mouse @ %d,%d' % (x, y)
        self.text.element.text = text
        self.text.element.x = self.posx
        self.text.element.y = self.posy

    def on_mouse_motion(self, x, y, dx, dy):
        """Called when the mouse moves over the app window with no button pressed

        (x, y) are the physical coordinates of the mouse
        (dx, dy) is the distance vector covered by the mouse pointer since the
          last call.
        """
        self.update_text(x, y)

        mouse_move_message = MouseEvent(x, y)
        send_to_server(mouse_move_message)

    def on_mouse_drag(self, x, y, dx, dy, buttons, modifiers):
        """Called when the mouse moves over the app window with some button(s) pressed

        (x, y) are the physical coordinates of the mouse
        (dx, dy) is the distance vector covered by the mouse pointer since the
          last call.
        'buttons' is a bitwise or of pyglet.window.mouse constants LEFT, MIDDLE, RIGHT
        'modifiers' is a bitwise or of pyglet.window.key modifier constants
           (values like 'SHIFT', 'OPTION', 'ALT')
        """
        self.update_text(x, y)

    def on_mouse_press(self, x, y, buttons, modifiers):
        """This function is called when any mouse button is pressed

        (x, y) are the physical coordinates of the mouse
        'buttons' is a bitwise or of pyglet.window.mouse constants LEFT, MIDDLE, RIGHT
        'modifiers' is a bitwise or of pyglet.window.key modifier constants
           (values like 'SHIFT', 'OPTION', 'ALT')
        """
        self.posx, self.posy = cocos.director.director.get_virtual_coordinates(x, y)
        self.update_text(x, y)

        bs = buttons_string(buttons)
        left_click = 'LEFT' in bs
        right_click = 'RIGHT' in bs

        mouse_move_message = MouseEvent(x, y, left_click, right_click)
        send_to_server(mouse_move_message)
