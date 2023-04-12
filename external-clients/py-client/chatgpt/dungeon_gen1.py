import random
import numpy as np
import matplotlib.pyplot as plt

WIDTH = 200
HEIGHT = 200
MIN_ROOM_SIZE = 10
MAX_ROOM_SIZE = 30
MIN_CORRIDOR_LENGTH = 10
MAX_CORRIDOR_LENGTH = 60


def intersects(r1, r2):
    return not (r1[0] + r1[2] <= r2[0] or r1[0] >= r2[0] + r2[2] or r1[1] + r1[3] <= r2[1] or r1[1] >= r2[1] + r2[3])


class Dungeon:
    def __init__(self, width, height):
        self.width = width
        self.height = height
        self.grid = np.zeros((height, width))

    def generate_rooms(self, n):
        rooms = []
        for i in range(n):
            w, h = random.randint(MIN_ROOM_SIZE, MAX_ROOM_SIZE), random.randint(MIN_ROOM_SIZE,
                                                                                MAX_ROOM_SIZE)
            x, y = random.randint(1, self.width - w - 1), random.randint(1, self.height - h - 1)
            room = (x, y, w, h)
            if not any(intersects(room, other) for other in rooms):
                rooms.append(room)
                self.grid[y:y + h, x:x + w] = 1
        return rooms

    def generate_corridors(self, rooms):
        while len(rooms) > 1:
            i, j = random.sample(range(len(rooms)), 2)
            x1, y1, w1, h1 = rooms[i]
            x2, y2, w2, h2 = rooms[j]
            cx1, cy1 = random.randint(x1, x1 + w1 - 1), random.randint(y1, y1 + h1 - 1)
            cx2, cy2 = random.randint(x2, x2 + w2 - 1), random.randint(y2, y2 + h2 - 1)
            if cx1 == cx2:
                self.grid[min(cy1, cy2):max(cy1, cy2) + 1, cx1] = 1
            else:
                hx1, hx2 = min(cx1, cx2), max(cx1, cx2)
                hy = random.randint(min(cy1, cy2), max(cy1, cy2))
                self.grid[hy, hx1:hx2 + 1] = 1
            rooms.pop(i)
            if i < j:
                rooms.pop(j - 1)
            else:
                rooms.pop(j)

    def print(self):
        for row in self.grid:
            line = ''
            for cell in row:
                line += '#' if cell == 0 else '.'
            print(line)

    def plot(self):
        plt.imshow(self.grid, cmap='binary')
        plt.xticks([])
        plt.yticks([])
        plt.show()

dungeon = Dungeon(WIDTH, HEIGHT)
rooms = dungeon.generate_rooms(40)
dungeon.generate_corridors(rooms)
dungeon.plot()
