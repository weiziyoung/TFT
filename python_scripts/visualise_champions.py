import matplotlib.pyplot as plt
import os
import numpy as np

os.chdir('../data')
width = 60*56
height = 60
avatars = [each for each in sorted(os.listdir('avatar'), key=lambda a:int(a[:2]))if not each.startswith('.')]
canvas = np.zeros((height, width, 3))

for n, avatar in enumerate(avatars):
    path = os.path.join('avatar', avatar)
    img = plt.imread(path)
    canvas[:, 60*n:60*(n+1), :] = img

plt.imsave('canvas.png', canvas)