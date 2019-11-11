import json
import os

os.chdir("../data")

if __name__ == "__main__":
    with open("language.json", "r") as f:
        translate_dict = json.load(f)
    with open("champions_graph.json", "r") as f:
        champions_graph = json.load(f)
    with open("champions.json", "r") as f:
        champions_list = json.load(f)

    # Generate nodes
    node_list = []
    for n, champion in enumerate(champions_list):
        node = dict()
        node['id'] = str(n)
        node["name"] = translate_dict[str(champion['name'])]
        node["cluster"] = translate_dict[champion["origin"][0]]
        node["value"] = champion["price"]*5
        node_list.append(node)

    # Generate edges
    edge_list = []
    for from_node, to_nodes in champions_graph.items():
        if to_nodes:
            for to_node in to_nodes:
                edge = dict()
                edge['source'] = str(from_node)
                edge['target'] = str(to_node)
                edge['sourceWeight'] = 10
                edge['targetWeight'] = 0
                edge_list.append(edge)

    total_dict = {
        "nodes": node_list,
        "edges": edge_list
    }
    with open("graph.json", 'w') as f:
        json.dump(total_dict, f, indent="\t", ensure_ascii=False)