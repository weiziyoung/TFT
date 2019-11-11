import requests
import os
import json
import cv2

os.chdir('../data')
if __name__ == "__main__":
    with open('champions.json', 'r') as f:
        champions = json.load(f)
        for n, each in enumerate(champions):
            filename = 'avatar/'+each['name']+'.png'
            if not os.path.exists(filename):
                url = each['avatar']
                url = url.replace('resize=64', 'resize=120')
                response = requests.get(url)
                with open(filename, 'wb') as fp:
                    fp.write(response.content)
            image = cv2.imread(filename)
            image = cv2.resize(image, (25, 25))
            cv2.imwrite('avatar/'+str(n)+each['name']+'.png', image)
