import random
from math import sqrt


class Monster:
    insect_names = ['ant', 'aphid', 'bee', 'beetle', 'butterfly', 'caterpillar', 'cockroach', 'cricket', 'dragonfly', 'firefly', 'fly', 'grasshopper', 'ladybug', 'mosquito', 'moth', 'spider', 'termite', 'wasp', 'weevil', 'worm']
    materials = ['iron', 'bronze', 'silver', 'gold', 'platinum', 'diamond', 'emerald', 'ruby', 'sapphire', 'obsidian']

    def __init__(self):
        self.monster_type = self.generate_monster_type()
        self.life = self.generate_life()
        self.defense = self.generate_defense()
        self.attack = self.generate_attack()
        self.speed = self.generate_speed()
        self.experience = self.calculate_experience()

    def generate_monster_type(self):
        insect_name = random.choice(self.insect_names)
        material = random.choice(self.materials)
        return f"{material}_{insect_name}"

    def generate_life(self):
        index = self.insect_names.index(self.monster_type.split('_')[1])
        return min((index + 1) * 5, 10)

    def generate_defense(self):
        index = self.insect_names.index(self.monster_type.split('_')[1])
        return min((index + 1) * 0.5, 10)

    def generate_attack(self):
        index = self.materials.index(self.monster_type.split('_')[0])
        return min((index + 1) * 1.0, 10)

    def generate_speed(self):
        index = self.materials.index(self.monster_type.split('_')[0])
        return min((index + 1) * 0.5, 10)

    def calculate_experience(self):
        return round(self.life * sqrt(self.attack) * sqrt(self.defense))

    def __str__(self):
        return f"{self.monster_type}\t{self.life:.1f}\t{self.defense:.1f}\t{self.attack:.1f}\t{self.speed:.1f}\t{self.experience}"


monsters = [Monster() for _ in range(30)]
monsters.sort(key=lambda x: x.experience, reverse=True)

print("Monster Type\tLife\tDefense\tAttack\tSpeed\tExperience")
print("--------------------------------------------------------------")
for monster in monsters:
    print(monster)
