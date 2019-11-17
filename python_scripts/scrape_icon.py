import os
import json
import requests
from bs4 import BeautifulSoup

response = requests.get('https://na.leagueoflegends.com/en/news/game-updates/gameplay/teamfight-tactics-gameplay-guide?utm_source=web&utm_medium=web&utm_campaign=tft-microsite-2019#patch-champion-compendium')
soup = BeautifulSoup(response.text, 'lxml')

icons = soup.find_all('div', attrs={'class': "tft-icon-container"})
result = {}
for icon in icons:
    url = icon.find("img")['src']
    trait = icon.find('h4').getText()
    result[trait] = url
result['Hextech'] = "http://game.gtimg.cn/images/lol/act/a20190702loltftwf/icon-hks.png"

os.chdir('../data/trait_icons')
for name, url in result.items():
    response = requests.get(url)
    with open(name+'.png', 'wb') as f:
        f.write(response.content)
