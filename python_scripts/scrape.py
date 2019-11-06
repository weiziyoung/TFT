# -*- coding: utf-8 -*-
# @Time    : 2019/10/29 11:19 PM
# @Author  : weiziyang
# @FileName: scrape.py
# @Software: PyCharm
import os
import json
import time
import requests
from bs4 import BeautifulSoup

response = requests.get('https://na.leagueoflegends.com/en/news/game-updates/gameplay/teamfight-tactics-gameplay-guide?utm_source=web&utm_medium=web&utm_campaign=tft-microsite-2019#patch-champion-compendium')
soup = BeautifulSoup(response.text, 'lxml')
block = soup.find('div', attrs={'id': 'champion-compendium-slideshow'})
boxes = block.findAll('div', attrs={'class': 'content-border'})

champion_list = list()

for box in boxes:
    avatar_url = box.find('a', attrs={'class': 'reference-link'}).find('img')['src']
    name = box.find('h4', attrs={'class': 'change-title'}).getText()
    price = box.find('p', attrs={'class': 'summary price'}).getText().replace('gold', '')
    h4 = box.findAll('h4', attrs={'class': 'ability-title'})
    origin = []
    Class = []
    for each in h4[1:]:
        text = each.getText()
        if 'Origin' in text:
            origin.append(text.replace('Origin', '').strip().lower())
        elif 'Class' in text:
            Class.append(text.replace('Class', '').strip().lower())
    champion = {
        'name': name.strip(),
        'avatar': avatar_url,
        'price': int(price.strip()),
        'origin': origin,
        'class': Class
    }
    champion_list.append(champion)

with open('data.json', 'w') as f:
    json.dump(champion_list, f, indent='\t')

