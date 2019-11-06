# -*- coding: utf-8 -*-
# @Time    : 2019/10/31 12:58 AM
# @Author  : weiziyang
# @FileName: traits.py
# @Software: PyCharm
import json

ONLY_ONE = 1
SAME_ORIGIN = 2
ALL = 3
NINJIA = 12
BOOST = 23
SUPERBOOST = 13
GUARD = 4


if __name__ == "__main__":
    with open('champions_data.json', 'r') as f:
        champion_list = json.load(f)
        origin_dict = dict()
        class_dict = dict()
        for champion in champion_list:
            origins = champion['origin']
            for origin in origins:
                if origin_dict.get(origin):
                    origin_dict[origin]['champions'].append(champion['name'])
                else:
                    origin_dict[origin] = {
                        'bonus_num': [],
                        'scope': 2,
                        'strength': 1,
                        'champions': [champion['name']]
                    }
            classes = champion['class']
            for _class in classes:
                if class_dict.get(_class):
                    class_dict[_class]['champions'].append(champion['name'])
                else:
                    class_dict[_class] = {
                        'bonus_num': [],
                        'scope': 2,
                        'strength': 1,
                        'champions': [champion['name']]
                    }

    trait_dict = {}
    trait_dict.update(origin_dict)
    trait_dict.update(class_dict)

    with open('traits_data.json', 'w') as f:
        json.dump(trait_dict, f, indent='\t')
