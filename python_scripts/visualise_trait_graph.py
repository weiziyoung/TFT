import os
import json
import numpy as np
import cv2
import matplotlib.pyplot as plt

os.chdir('..')
with open('data/champions.json') as f:
    champions = json.load(f)
trait_graph = [["" for _ in range(len(champions))] for each in range(len(champions))]

prior_trait = {'ninja', 'exile', 'robot'}
for y, champion_a in enumerate(champions):
    origins = champion_a['origin']
    classes = champion_a['class']
    champion_a_traits = set(origins + classes)
    for x, champion_b in enumerate(champions):
        if x == y:
            trait_graph[y][x] = "same"
            continue
        origins = champion_b['origin']
        classes = champion_b['class']
        champion_b_traits = set(origins + classes)
        # 两个英雄共同的羁绊
        common_traits = (champion_a_traits & champion_b_traits)
        # 是否存在忍者、浪人、机器人
        if common_traits:
            extra_traits = set()
        else:
            extra_traits = (champion_b_traits | champion_a_traits) & prior_trait
        all_traits = common_traits | extra_traits
        if all_traits:
            trait_graph[y][x] = str(','.join(list(all_traits)))

champions_num = len(champions)
size = 40
width = size + champions_num * size * 2 + size
height = size + champions_num * size + size
canvas = np.zeros((height, width, 3))*255
start_point = (size, size)

# 首先画英雄
file_names = sorted([each for each in os.listdir('data/avatar') if not each.startswith('.')], key=lambda a: int(a[:2]))
file_names = [os.path.join('data', 'avatar', each) for each in file_names]
for n, file_name in enumerate(file_names):
    img = cv2.cvtColor(cv2.imread(file_name), cv2.COLOR_BGR2RGB)
    resized_img = cv2.resize(img, (size, size))
    # 先画横向的
    y = 0
    x = size + size * 2 * n
    canvas[y:y+size, x:x+size, :] = resized_img
    canvas[:, x+size:x+size+1, :] = np.ones((height, 1, 3))*255
    # 再画纵向的
    x = 0
    y = size + size * n
    canvas[y:y+size, x:x+size, :] = resized_img
    canvas[y+size:y+size+1, :, :] = np.ones((1, width, 3))*255

# 载入羁绊
trait_file_names = [each for each in os.listdir('data/trait_icons') if not each.startswith('.')]
name = [each.split('.')[0].lower() for each in trait_file_names]
trait_icon_dict = {}
for n, trait_file in enumerate(trait_file_names):
    img = cv2.cvtColor(cv2.imread(os.path.join('data', 'trait_icons', trait_file)), cv2.COLOR_BGR2RGB)
    resized_img = cv2.resize(img, (size, size))
    trait_icon_dict[name[n]] = resized_img
trait_icon_dict['Hextech'] = trait_icon_dict['hextech']
# 再画羁绊
for row in range(len(trait_graph)):
    for col in range((len(trait_graph))):
        content = trait_graph[row][col]
        if content and content != 'same':
            traits = content.split(',')
            for index, trait in enumerate(traits):
                icon = trait_icon_dict[trait]
                y = size + row * size
                if index == 0:
                    x = size + 2 * col * size
                if index == 1:
                    x = size + 2 * col * size + size
                canvas[y:y+size, x:x+size] = trait_icon_dict[trait]

cv2.imwrite("data/trait_graph.png", canvas)