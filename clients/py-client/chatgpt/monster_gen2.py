import random
import math


class Monster:
    def __init__(self, monster_type, life, defense, attack, speed, experience):
        self.monster_type = monster_type
        self.life = life
        self.defense = defense
        self.attack = attack
        self.speed = speed
        self.experience = experience

    def __str__(self):
        return f"{self.life:.1f} | {self.defense:.1f} | {self.attack:.1f} | {self.speed:.1f} | {self.experience:.0f} | {self.monster_type}"

# Generate list of 30 unique insect names with random size values between 1 and 10
insect_names = []
insect_sizes = []
while len(insect_names) < 15:
    insect_name = f"{random.choice(['butterfly', 'beetle', 'ant', 'spider', 'wasp', 'moth', 'mosquito', 'grasshopper', 'dragonfly', 'caterpillar', 'centipede', 'cricket', 'snail', 'fly', 'termite', 'ladybug', 'bee', 'cockroach', 'scorpion', 'cicada'])}"
    if insect_name not in insect_names:
        insect_names.append(insect_name)
        insect_sizes.append(random.randint(1, 10))

# Generate list of 30 unique material names with random hardness values between 1 and 10
material_names = []
material_hardnesses = []
while len(material_names) < 20:
    material_name = f"{random.choice(['adamantite', 'kalendrite', 'mythril', 'orichalcum', 'rune', 'sardonyx', 'tanzanite', 'zircon', 'carbide', 'cobalt', 'ceramic', 'corundum', 'titanium', 'chromium', 'nickel', 'aluminum', 'iron', 'copper', 'silver', 'gold'])}"
    if material_name not in material_names:
        material_names.append(material_name)
        material_hardnesses.append(random.randint(1, 10))

# Generate list of 30 unique monsters
monsters = []
while len(monsters) < 30:
    # Randomly choose an insect and a material
    insect_idx = random.randint(0, len(insect_names)-1)
    material_idx = random.randint(0, len(material_names)-1)
    insect_name = insect_names[insect_idx]
    material_name = material_names[material_idx]
    size = insect_sizes[insect_idx]
    hardness = material_hardnesses[material_idx]

    # Calculate monster properties
    life = size * 10
    defense = size
    attack = 1 + (hardness / 10) * 9
    speed = 1 + ((10 - hardness) / 10) * 9
    experience = round(life * math.sqrt(attack) * math.sqrt(defense))

    # Create and append monster to list
    monster_type = f"{material_name}_{insect_name}"
    monster = Monster(monster_type, life, defense, attack, speed, experience)
    monsters.append(monster)

# Sort monsters by experience highest first
monsters.sort(key=lambda x: x.experience, reverse=True)

# Print monsters as a nicely formatted table
print("Life | Defense | Attack | Speed | Experience | Monster Type")
print("-" * 55)
for monster in monsters:
    print(monster)
