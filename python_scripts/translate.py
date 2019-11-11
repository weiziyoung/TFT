import json

# with open('../data/traits.json', 'r') as f:
#     champion_list = json.load(f)
#
# name_list = []
# for each in champion_list:
#     name_list.append({each['name']: ""})
#
# with open('../data/language.json', 'w') as f:
#     json.dump(name_list, f, indent='\t', ensure_ascii=False)
dic = {}

with open("../data/language.json", 'r') as f:
    l = json.load(f)

for each in l:
    key = list(each.keys())[0]
    value = list(each.values())[0]
    dic[key] = value

with open("../data/language_copy.json", "w") as f:
    json.dump(dic, f, indent='\t', ensure_ascii=False)