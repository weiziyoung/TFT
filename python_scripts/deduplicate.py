import json


# jaccard 相似度
def diff_set(set1 :set, set2 :set) -> float:
    return len(set1 & set2) / len(set1 | set2)


def setify(combo:dict) -> set:
    trait_detail = combo['trait_detail']
    result_set = set()
    for trait, num in trait_detail.items():
        result_set.add((trait, num,))
    return result_set


if __name__ == "__main__":
    num = 9
    with open(f'../data/output/total_strength/champions_comb{num}.json', 'r') as f:
        combos = json.load(f)

    trait_list = [setify(combos[0])]
    for combo in combos[1:]:
        combo_set = setify(combo)
        for trait_set in trait_list:
            # 如果相似度大于80%就是相同体系羁绊
            sim = diff_set(trait_set, combo_set)
            if sim >= 0.7:
                break
        else:
            trait_list.append(combo_set)
    result = [str(list(each)) for each in trait_list]
    with open(f"../data/final_result/total_strength/final_result{num}.json", 'w') as f:
        json.dump(result, f, indent='\t', ensure_ascii=False)
