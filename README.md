> 本文希望读者玩过云顶之弈，不懂编程的可以直接拉到最下面去看结论，懂编程的希望你了解递归、分治、图、堆这些基本概念，并掌握Python或者Go语言。
代码已公开在github上:[https://github.com/weiziyoung/TFT](https://github.com/weiziyoung/TFT) ，转载请注明来源。

今天是11月11日，首先恭喜FPX一顿摧枯拉朽横扫G2， 拿下S赛冠军！证明了LPL是世界第一赛区，也让电竞作为一种赛事在这代人心中铭记。本届S赛结束，也就意味着，S8告一段落，S9即将上线。而云顶之弈作为今年刚出的新模式，在上周11月6日也发布了S2元素崛起版本，一时间各种打法也是层出不穷，小编我也是一名忠实的云顶之弈玩家，但目前还没有玩过S2版本，主要想把这篇文章先写好分享给想读的人。

其实早在今年暑假刚出这个新模式，大家都还不会玩，还在摸索各种阵容的时候，我就在思考一件事——如何通过编程的手段搜索到6人口、7人口、8人口甚至9人口时凑到的最多的羁绊？这种想法来源于一个惨痛的经历，就是我第一次玩的时候，大概只凑出来了一个3贵族2骑士羁绊，就草草第七名带走了...当时就觉得这个游戏太难了，这么多卡片怎么能全记住？除了英雄之外，还有装备合成也跟英雄联盟差的很远，但玩个两三局，大概就明白：

这个游戏想吃鸡有三个核心——**羁绊**、**英雄等级**、**装备**， 三个核心有两个占优势，基本可以达到前四，三个都占优势，就稳定吃鸡了。这里我们主要讨论的就是去搜索**羁绊**，从而在这个维度上不吃亏。而装备这块比较靠脸，所以不做讨论，英雄等级这块其实可以根据每张卡在每个阶段出现的概率来估算出来这个阵容成型的难易程度，但是在本片博客里不做讨论，这里只讨论一个问题，就是**羁绊**。

# 文章大纲
- 云顶之弈游戏简介
- 基本算法思路
- 准备实验数据
- 排列组合的原理和实现
- 用图降低搜索复杂度
- 评估函数的设计和实现
- 最小堆维护Top 100阵容
- 结果展示
- 分析与总结

# 云顶之弈游戏简介
一般读到这里的读者应该都玩过云顶之弈，但为了照顾有些只打匹配排位从不下棋的同学，这里还是简单介绍一下这个游戏机制。
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x37e60f05a6fa8aimage.png)

- 方框`1`所在的是小小英雄，就是召唤师啦，好看用的。
- 方框`2`是你目前的队伍，这个队伍可以由不同英雄组成，但是队伍规模取决于你等级的高低。
- 方框`3`是你候选区的英雄，放一些你暂时不想上场的英雄，当然这个区域大多数是用来合成英雄的，3个相同1星英雄可以合成成2星，3个相同2星英雄可以合成为3星。当然我们这里不讨论如何优化英雄等级的话题。
- 方框`4`是发牌员给你发的牌，还有你目前有多少钱，每回合发牌员会给你发5张牌，你需要用金币去购买，这里只需要记住一点，星级越高的英雄越难抽到，并且也越强。
- 方框`5`就是我们的核心——羁绊了，它是根据场上的英雄的种族和职业所确定的，比如目前场上小炮和男枪可以组成一个枪手的buff，这个Buff可以使得枪手在攻击时造成2个目标伤害，而劫的话自己是个忍者，所以可以组成一个忍者buff，它可以提升自己暴击几率和攻击速度。每个羁绊都有自己的效果，同时，羁绊也有自己的等级，比如当你只有2个枪手的时候，你的枪手能够同时造成2个敌人的伤害，而4个枪手的时候，你可以再攻击时同时造成3个目标的伤害；同时羁绊也有范围，有的羁绊只对单个人有效，比如双帝国、三贵族、单忍者，大多数羁绊对同种族的有效，比如狂野、枪手、剑士，少数羁绊对队伍里所有英雄都有效，比如骑士、法师。

具体的S1版本英雄羁绊图如下(有一些后期英雄没加上去，比如潘森、卡萨、海克斯):
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x37e611184d2ee6image.png)

总共是56只英雄，大多英雄拥有一个种族，一个职业，船长的职业既是剑士也是枪手，纳尔的种族既是约德尔人也是狂野。一般来说，这个游戏在七人口时阵容成型，这个阶段基本能看出谁胜谁负，所以我们的目的就是选7个英雄，组成羁绊上最强的阵容。

# 基本算法思路

就像之前所说的，我们的目的是在56个英雄里选n个英雄，然后从里面选出羁绊最强的前K个。这句话可以拆分为这三个问题:
1. 首先，如何让计算机去自动把所有组合的可能性一个不拉地遍历出来？不重复也不漏检？
2. 其次，给定一个阵容，如何去评判羁绊的强度？
3. 第三，怎么去保存前K个羁绊最强的结果？

对于第一个问题来说，很多编程语言都有combination的拓展库，方便程序员求出一个列表的元素所有的组合可能性。但是这个是个好的方案嘛？真的可行嘛？如果不可行，怎么去优化？

对于第二个问题来说，我们在评估一个东西，或者说量化一个东西的时候，应该采用哪些指标？羁绊多是不是意味着羁绊就强？如果不是的话，是否需要引入主观性的一些指标，比如单个羁绊对英雄的增益程度？另外这个羁绊好成型嘛？是不是容易在组建的半路上暴毙？这些都是需要注意的问题。

对于第三个问题来说，看起来很容易，但排序真的可行吗？由于我们搜索的结果多达几百万个的阵容组合，全部排序后再取前K个现实嘛？

# 准备实验数据
本次主要使用语言为Go，并且用Python做一些脚本辅助我们做一些分析，之所以采用Go来写核心代码，是因为这种上百万轮次的搜索，Go往往比Python能快出一个数量级，同时Go工程化之类的也做的更好一些，语法也不至于像C++和Java那样繁琐。

程序 = 算法 + 数据。数据是一切的基石，要实现我们这次的目标，我们至少需要拥有两个数据:英雄数据、羁绊数据。在国外英雄联盟官网上，我们可以找到这个页面:[TFT GAMEPLAY GUIDE](https://na.leagueoflegends.com/en/news/game-updates/gameplay/teamfight-tactics-gameplay-guide?utm_source=web&utm_medium=web&utm_campaign=tft-microsite-2019)，接下来只要用Python 的BeautifulSoup包吧页面解析出来就可以了，大概20行代码就可以搞定了，由于思路比较简单，这里就不放代码了，给个链接自己看：[python_scripts/scrape.py](https://github.com/weiziyoung/TFT/blob/master/python_scripts/scrape.py)。

如下所示，这里我们需要记录英雄的元数据包括:名字、头像、费用、种族和职业，总共56个英雄，这里不展示了。需要的自己去取:[data/champions.json](https://github.com/weiziyoung/TFT/blob/master/data/champions.json)
```json
	{
		"name": "Varus",
		"avatar": "https://am-a.akamaihd.net/image?f=https://news-a.akamaihd.net/public/images/articles/2019/june/tftcompendium/Champions/Varus.png&resize=64:",
		"price": 2,
		"origin": ["demon"],
		"class": ["ranger"]
	},
```

另外是羁绊数据，这个数据可以从英雄数据里面整理出来，同时也要我们自己手填一些数据，以**恶魔**为例:
```json
	{
		"name": "demon",
		"bonus_num": [2,4,6],
		"scope": [2,2,2],
		"champions": [
		"Varus","Elise","Morgana","Evelynn",
		"Aatrox","Brand","Swain"]
	},
```
恶魔羁绊需要在2只时触发，且在4,6时羁绊进阶，那`bonus_num`就是`[2,4,6]`，而恶魔羁绊无论多少级，都是只有同种族的受益，所以范围序号是`2`，具体范围序号含义我们定义如下
1. `1`代表只有一个英雄能吃到这个羁绊buff的效果，典型的比如3贵族、2帝国。
2. `2`代表持有该羁绊的能够吃到这个buff效果，大多数羁绊都属于这个效果，比如恶魔、冰川、狂野、变形者、刺客、枪手、剑士、4帝国等等。
3. `3`代表队伍全部都可以吃到这个buff，比如6贵族、骑士、法师这些。
4. `4`代表一个特殊的羁绊范围，就是护卫了，护卫是除了护卫本身，其周围的人都能吃到buff。

`champions`就是持有这个羁绊的所有英雄了，全部羁绊数据在这里:[data/traits.json](https://github.com/weiziyoung/TFT/blob/master/data/traits.json)。这就是我们现在所能拿到的所有客观数据，不掺杂任何拍脑袋给的主观权重。实际上在评估时，这种数据越多越好，主观性太强的指标例如英雄强度、羁绊强度这种，公说公有理，婆说婆有理，很难有客观的结论，尽量少引入到评价体系中。

# 排列组合的原理和实现
现在我们有所有英雄了，作为召唤师，我觉得很有必要把它们一字排开欣赏欣赏...毕竟S2就看不到他们的绝大多数了。
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x37e6ce133434acimage.png)
所以我们的任务就是从55个英雄里面挑出8个英雄，让他们的羁绊数量最多。所以这是一个排列组合里的组合问题，可以根据公式求出组合数量:
$$ C_{n}^{m}=\frac{n!}{m!(n-m)!} $$
其中`n`等于55，`m`等于8，也就是八人口时，需要搜索`231917400`个不重复的可能性。
## 如何实现组合呢
最经典的思路就是分治了，看个简单的问题，比如对`[a,b,c,d,e]`求个数为3的所有组合。那么，我们首先会先把`a`取出来，问题简化成了对`[b,c,d,e]`求个数为2的所有组合。其次，我们把`b`取出来，问题简化成了对`[c,d,e]`求个数为1的所有组合，这时候问题就简单了.示意图如下：
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x37e6d4b7680260image.png)
红框表示你现在已经选择的字母，红框下面的数字代表需要继续进行组合的元素，到三层结束。
Python实现代码，非常短小精干，需要仔细品味和研读，理解递归、分治的优雅:
```Python
def combine(data, step, selected_data, target_num):
    if len(selected_data) == target_num:   # 递归的结束条件:已选择的元素数量等于目标数量
        print(selected_data)
        return
    if step == len(data):               # 游标到头了还没找出来，就结束吧
        return
    selected_data.append(data[step])             # 选择当前元素把她加到已选择的队伍里
    combine(data, step + 1, selected_data, target_num)  # 将游标推进，进入递归去找下一层
    selected_data.pop()                         # 把选择过的元素踢出去
    combine(data, step + 1, selected_data, target_num) #在不选择刚才被踢出去的元素情况下继续递归
if __name__ == '__main__':
    data = ['a','b','c','d', 'e']
    combine(data, 0, [], 3)
```
理解了上面这个代码，换个变量名，加入evaluate函数，就可以用于搜索我们的全羁绊了。
```Python
def combine(champions, step ,combo, max_num):
    if len(combo) == max_num: # 如果队伍到了最大的人口，就进行评估
        evaluate(combo)
        return
    if step == len(combo):  
        return
    combo.append(champions[step]) # 把游标所指定的英雄加到队伍里面去
    combine(champions, step + 1, combo, max_num)  # 游标往前进，继续抓壮丁
    combo.pop()  # 把刚才指定的英雄踢出去
    combine(champions, step+1, combo, max_num)   # 再继续往前进抓壮丁

def evaluate(combo):
    # 这里写给定一个阵容，怎么去评估它的强度，应该返回一个数值，或者是多个维度的评分结构体。
    # 往后再议
    pass

def init_champions():
    # 这里从json里读数据，代码略
    pass

if __name__ == "__main__":
    champions = init_champions()  # 把英雄数据导入进去，每个英雄应该是个结构体，或者是个字典。
    combine(champions, 0, [], 7) 
```
跑了一下，自行感受一下Python🐌蜗牛一般的速度吧:
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x37e77878ffd19cimage.png)
平均每秒遍历`36979`个结点，搜索6人口的最优羁绊竟然要花14分钟，作为一个堆效率有追求的程序员，怎么能够容忍这种事情出现？？我只想对这个结果说:
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x37e779287edda3image.png)
所以接下来就没有Python代码了,同样的算法用Go跑的话，速度是每秒大约20w个结点， 大概是Python的7倍左右，如果用C++来写会更快，但如果让我用C++来写可能要明年你们才能看到我这篇文章了，所以程序员要在开发效率和运行速度中取得一个平衡:
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x37e77b30a5a8ffimage.png)

# 用图降低搜索复杂度
## 穷举法的弊端
由之前的公式：

$$ C_{n}^{m}=\frac{n!}{m!(n-m)!} $$

我们可以算出，八人口需要搜索`231917400`个结点，用Python搜索大概需要1.7个小时左右，用Golang搜索大概需要20分钟，速度还是很不够看，从语言上已经优化不了了，那就从算法上进行优化。

结合这个游戏，仔细思考一下我们是否真的需要对56个英雄都组合一遍呢？这么看不够直观，我举个非常简单的栗子
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x37e78fc64dcd0aimage.png)
给定图上的五只英雄：蜘蛛、盖伦、浪人、维鲁斯、豹女、寒冰，选出三个英雄，目标是让他们组成的羁绊数量最大，用大脑去看，那结果一定是“蜘蛛、维鲁斯、寒冰”，但是，我们模拟之前穷举法的过程，首先选出蜘蛛，其次选择第二位的盖伦，如果真的有人会在拿到蜘蛛的情况下去第二位去选择盖伦凑羁绊，大概会让人觉得：
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x37e790c69b29a5image.png)
## 基于羁绊的思路
正常的人拿完蜘蛛，下一步一定是拿维鲁斯或者豹女，拿维鲁斯因为刚好可以凑一个恶魔，维鲁斯又是一个比较强的打工二费卡，何乐而不为？拿豹女是因为后面可能可以凑3换型，能凑出3换型，前期坦度是妥妥的，所以我们在拿到蜘蛛的情况下，不可能去考虑下一步拿盖伦和狼人，在下一步拿到维鲁斯的情况下，去考虑豹女和寒冰，（思考一下为什么要考虑豹女？），这样我们就达到了最多羁绊:双恶魔加双寒冰。综上，我们简化搜索的主要逻辑就是**每次只选择与他能产生羁绊的对象**，基于这个想法，我们的搜索就变成：
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x597283d640b96image.png)

而图就是用来描述每个对象之间关系的一种数据结构，在这里，图用来描述英雄之间的羁绊关系，而图的表示方法有两种：邻接矩阵法和邻接表法，两者的取舍取决于图的稀疏程度。将上面官方给的羁绊-英雄图转个方式就得到了英雄-羁绊邻接矩阵图（57*57的矩阵，有相同羁绊则输出1， 没有则输出0）由图中可以看出，该矩阵为稀疏矩阵，所以我们后面用邻接表法来表示该矩阵):
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x37eb7e5e8dd158image.png)
另外，所有的英雄都和机器人、浪人、忍者有羁绊，因为队伍里只要添上它们任何中的一个，都可以为羁绊数+1，符合我们的优化预期。亚索在这里不是孤儿了。
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x37eb85fea153ddimage.png)

那么怎么利用这个信息去优化我们的算法呢？这需要进一步地去理解“组合”搜索究竟做了什么？是否可以用图的方式来进行组合搜索？答案是肯定的，以刚才组合`a,b,c,d,e`选出3个进行组合为例，换个思路来想这个事，实际上他们彼此之间也可以用有向图来表示:
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x37eb83b1ecd59aimage.png)
所以之前那个组合示意图，也可以这么理解：
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x37eb8407a5c20aimage.png)
综上所述，对于组合而言，我们只要把每个结点指向起后面的所有结点，然后用普通的图搜索，就可以得到组合的结果。
而利用羁绊图，我们可以不用把每个结点指向后面的所有结点，相反，我们只要把每个结点指向后面所有能跟当前组合产生羁绊的结点就可以了，注意！不能只考虑和当前结点产生羁绊，而要考虑队伍里所有英雄所拥有的所有结点，否则会漏搜索！我们优化的初衷是，保证搜索结果不变的情况下，减少不必要的搜索，而不能漏搜索。
因此核心搜索代码如下:
```go
type Graph map[int][]int

// GenerateGraph 生成羁绊图
func GenerateGraph(championList models.ChampionList) Graph{
	graph := make(Graph)
	positionMap := make(map[string]int)
	for index, champion := range championList {
		positionMap[champion.Name] = index
	}
	for no, champion := range championList {
		// children 排序
		children := make([]int, 0, 30)
		// 加入相同职业的英雄
		classes := champion.Class
		for _, class := range classes {
			sameClassChampions := globals.TraitDict[class].Champions
			for _, champion := range sameClassChampions {
				index := positionMap[champion]
				if index > no{
					children = append(children, index)
				}
			}
		}
		// 加入相同种族的英雄
		origins := champion.Origin
		for _, origin := range origins {
			sameOriginChampions := globals.TraitDict[origin].Champions
			for _, champion := range sameOriginChampions {
				index := positionMap[champion]
				if index > no {
					children = append(children, index)
				}
			}
		}
		// 加入1羁绊的英雄
		for _, championName := range globals.OneTraitChampionNameList {
			index := positionMap[championName]
			if index > no {
				children = append(children, index)
			}
		}
		// 对index从小到大排序
		sort.Ints(children)
		children = utils.Deduplicate(children)
		graph[no] = children
	}
	return graph
}


// Traverse 图遍历，
// championList, 英雄列表，固定不变。 graph 羁绊图，也是固定不变。node 为当前的结点, selected 为已选择的英雄， oldChildren是父节点的children
func Traverse(championList models.ChampionList, graph Graph, node int, selected []int, oldChildren []int) {
	selected = append(selected, node)
	if len(selected) == lim {
		combo := make(models.ChampionList, lim)
		for index, no := range selected {
			unit := championList[no]
			combo[index] = unit
		}
		metric := evaluate.Evaluate(combo)
		heap.Push(&Result, metric)

		// 超过最大就pop
		if len(Result) == globals.Global.MaximumHeap {
			heap.Remove(&Result, 0)
		}
		return
	}
	newChildren := graph[node]
	children := append(oldChildren, newChildren...)
	sort.Ints(children)
	children = utils.DeduplicateAndFilter(children, node)
	copyChildren := make([]int, len(children), 50)
	copy(copyChildren, children)
	for _, child := range children {
		copySelected := make([]int, len(selected), lim)
		copy(copySelected, selected)
		Traverse(championList, graph, child, copySelected, copyChildren)
	}
}
// TraitBasedGraphSearch 基于羁绊图的图搜索
func TraitBasedGraphSearch(championList models.ChampionList, teamSize int) models.ComboMetricHeap {
	graph := GenerateGraph(championList)
	lim = teamSize

	heap.Init(&Result)
	startPoints := getSlice(0, len(championList)-teamSize + 1)
	for _,startNode := range startPoints{
		Traverse(championList, graph, startNode, make([]int, 0, teamSize), make([]int, 0, 57))
	}
	return Result
}
```
用这种方法所产生的有向图如下图所示（这里顺手安利一个网络图可视化的js库[antv-G6](https://www.yuque.com/antv/g6/ie7zi7)），大幅度简化了初始的搜索图（自行想象一下所有结点连接所有后续结点密密麻麻的效果图）。

![image10.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images/screencapture-localhost-63342-TFT-visualisation-champion-graph-html-2019-11-10-16_16_59-Recovered.jpg)
实际上，我认为这种启发式搜索，有点A star搜索的意思在里面，核心思想就是讲后续children进行排序，将预期离目标结果近的放在前面。这里做的极端了一些，我们把没有产生羁绊的后续结点全部咔嚓了，但实际上这并不会造成漏检（读者可以自己实验一下）

最后，比较一下基于羁绊图的结点搜索数量和不基于羁绊图的结点搜索数量，横坐标是人口，纵坐标是结点数量，注意一下纵坐标的跨度，是指数级别的。
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x37eb8c389e91dbimage.png)
所以到这里，这篇博客的核心部分就讲完了，基本思想就是利用现有的知识(英雄之间产生的羁绊）来大幅度简化搜索。

# 评估函数的设计与实现
之前我们一直都没有实现评估函数，其实这个评估函数的设计是非常灵活的，也是玩家可以加入自己玩游戏的经验的一部分。这里我们用4个指标来描述阵容强度：
```go
type ComboMetric struct {
	// 英雄组合
	Combo []string `json:"combo"`
	// 队伍总羁绊数量 = sigma{羁绊} * 羁绊等级
	TraitNum int `json:"trait_num"`
	// 具体羁绊
	TraitDetail map[string]int `json:"trait_detail"`
	// 总英雄收益羁绊数量 = sigma{羁绊} 羁绊范围 * 羁绊等级
	TotalTraitNum int `json:"total_trait_num"`
	// 当前阵容羁绊强度 = sigma{羁绊} 羁绊范围 *  羁绊强度
	TotalTraitStrength float32 `json:"total_trait_strength"`
	// 当前阵容强度 = sigma{英雄} 英雄强度 * 羁绊强度
	TotalStrength float64 `json:"total_strength"`
}
```

- 队伍总羁绊数量: 这个是最好理解的，你可以理解为你左侧边栏有多少个羁绊，也就是这个部分,谁不喜欢亮刷刷的一排羁绊呢？看的就很舒服。注意，像6恶魔这种算3个羁绊，而不能只算1个羁绊，6贵族算2个羁绊。这也是我们最开始的motivation，就是寻找怎么能让左边的羁绊灯亮的最多。
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x37eb8df9467873image.png)

- 英雄总收益羁绊数量： 这个也是好理解的，灯亮的多并不代表强，我的经验告诉我，往往吃鸡的阵容，灯亮的往往并不多，有时候甚至就三四个，因此需要引入其他衡量标准。因为不同羁绊有不同的收益范围，所以这个指标就是计算的就是每个羁绊羁绊收益范围乘以它等级的总和。6贵族羁绊之所以挺强，强的不在于它单个属性有多强，而在于它产生了单个buff到群体buff的一个质变，骑士Buff好用也是这个道理，为什么大家都喜欢用骑士过渡，甚至到后面主流吃鸡阵容就包括骑士+枪呢？本质上就是因为骑士能够提供全队伍的收益，而不是只针对本种族的收益。

- 当前阵容羁绊强度： 这个指标开始就加入人为指标了，也就是其输出取决于玩家对羁绊的理解，这个指标引入了羁绊强度这个概念，这个参数是指：当英雄拥有该羁绊时，能够比不拥有羁绊时强多少倍，比如在这里我设置贵族buff可以让英雄强1.8倍，双恶魔buff能让英雄强1.3倍，龙buff能够直接增强2倍...具体可以看我[data/traits.json](https://github.com/weiziyoung/TFT/blob/master/data/traits.json)文件。

- 当前阵容整体强度: 这个跟上一版差别就在于考虑了英雄的等级，比如同样是双骑士，你拿个盖伦加诺手，肯定比不上你拿个波比加猪妹。这里为了简化情景，所以设定，2星英雄比1星英雄强1.25倍，3星又比2星强1.25倍...以此类推，最后5星英雄大约比1星英雄强2.5倍，如果你觉得这个数值低了，可以自己在配置文件里面调整。

最后我们的evaluate评估函数如下,注意一个问题，就是忍者buff的奇异设定，游戏规定，忍者Buff只在1和4的时触发，在2,3时会熄灭，这不同于其他任何一个羁绊规则，所以要拎出来单独处理一下:
```go
// Evaluate 评估当前组合的羁绊数量、单位收益羁绊总数、羁绊强度
func Evaluate(combo []models.ChampionDict) models.ComboMetric {
	var traitDetail = make(map[string]int)

	comboName := make([]string, 0, len(combo))
	traitNum := 0
	totalTraitNum := 0
	totalTraitStrength := float32(0.0)

	// 初始化英雄强度向量
	unitsStrength := make([]float64, len(combo), len(combo))
	traitChampionsDict := make(map[string][]int)
	for index, unit := range combo {
		comboName = append(comboName, unit.Name)
		unitStrength := math.Pow(globals.Global.GainLevel, float64(unit.Price-1))
		unitsStrength[index] = unitStrength
		for _, origin := range unit.Origin {
			traitChampionsDict[origin] = append(traitChampionsDict[origin], index)
		}
		for _, class := range unit.Class {
			traitChampionsDict[class] = append(traitChampionsDict[class], index)
		}
	}

	for trait, champions := range traitChampionsDict {
		num := len(champions)
		bonusRequirement := globals.TraitDict[trait].BonusNum
		var bonusLevel = len(bonusRequirement)
		for index, requirement := range bonusRequirement {
			if requirement > num {
				bonusLevel = index
				break
			}
		}

		// 忍者只有在1只和4只时触发，其他不触发
		if trait == "ninja" && 1 < num && num < 4 {
			bonusLevel = 0
		}
		if bonusLevel > 0 {
			traitDetail[trait] = bonusRequirement[bonusLevel-1]
			bonusScope := globals.TraitDict[trait].Scope[bonusLevel-1]
			traitNum += bonusLevel
			bonusStrength := globals.TraitDict[trait].Strength[bonusLevel-1]
			benefitedNum := 0
			switch bonusScope {
			case 1:
				{
					benefitedNum = 1 // 单体Buff，例如 机器人、浪人、三贵族、双帝国
					for _, champion := range champions {
						unitsStrength[champion] *= float64(bonusStrength)
					}
				}
			case 2:
				{
					benefitedNum = num // 对同一种族的Buff，大多数羁绊都是这种
					for _, champion := range champions {
						unitsStrength[champion] *= float64(bonusStrength)
					}
				}
			case 3:
				{
					benefitedNum = len(combo) // 群体Buff，如骑士、六贵族、四帝国
					for index, _ := range unitsStrength {
						unitsStrength[index] *= float64(bonusStrength)
					}
				}
			case 4:
				{
					benefitedNum = len(combo) - 2 // 护卫Buff，比较特殊，除护卫本身外，其他均能吃到buff
					for index, _ := range unitsStrength {
						isGuard := false
						for _, champion := range champions {
							if index == champion {
								isGuard = true
								break
							}
						}
						if !isGuard {
							unitsStrength[index] *= float64(bonusStrength)
						}
					}
				}
			}
			totalTraitNum += bonusLevel * benefitedNum
			totalTraitStrength += float32(benefitedNum) * bonusStrength
		}
	}
	metric := models.ComboMetric{
		Combo:              comboName,
		TraitNum:           traitNum,
		TotalTraitNum:      totalTraitNum,
		TraitDetail:        traitDetail,
		TotalTraitStrength: totalTraitStrength,
		TotalStrength:      utils.Sum(unitsStrength),
	}
	return metric
}

```

# 最小堆维护Top 100阵容
之前也提到了，我们每次搜索都是对上千万乃至上亿的叶子结点进行评估，那么如何取出评估结点的前100名呢？我们会想到把结果存起来，然后排序，但这么做可行嘛？

想一下我们十人口进行搜索，总共搜索了25844630个结点，假设每存一个metric需要消耗1kb，那最后把它们全部存下来，大约需要24G，记住这是存在内存里的哦，而不是在硬盘上的噢，正常PC的内存条能有16G很不错了吧，更何况还要跑个操作系统在上面，所以这个方案一定是不行的，那有什么更好的方案呢？

这就需要联系我上个月写的博客，[详解数据结构——堆](https://zhuanlan.zhihu.com/p/85518062),这篇博文里我们讲到利用堆，我们只需要在内存里开辟堆长度大小的空间即可，比如我们想保留前100个结果，那我们只要开辟100k的内存即可，而每次插入删除，都是`log n`的复杂度，非常快。

而保留前K个结果，需要使用的是最小堆，golang里集成了堆的数据结构，只需要重写它的一些接口就可以用了，所以我们的ComboMetric完整版实现就是这样，具体用起来就是每次都push，满了就把堆顶pop出来即可，最后剩下来的就是前K个结果，把它们最后排个序即可:
```go
package models

type ComboMetric struct {
	// 英雄组合
	Combo []string `json:"combo"`
	// 队伍总羁绊数量 = sigma{羁绊} * 羁绊等级
	TraitNum int `json:"trait_num"`
	// 具体羁绊
	TraitDetail map[string]int `json:"trait_detail"`
	// 总英雄收益羁绊数量 = sigma{羁绊} 羁绊范围 * 羁绊等级
	TotalTraitNum int `json:"total_trait_num"`
	// 当前阵容羁绊强度 = sigma{羁绊} 羁绊范围 *  羁绊强度
	TotalTraitStrength float32 `json:"total_trait_strength"`
	// 当前阵容强度 = sigma{英雄} 英雄强度 * 羁绊强度
	TotalStrength float64 `json:"total_strength"`
}

// 定义一个最小堆，保留前K个羁绊
type ComboMetricHeap []ComboMetric

func (h ComboMetricHeap) Len() int {
	return len(h)
}

func (h ComboMetricHeap) Less(i,j int) bool {
	return h[i].TotalStrength < h[j].TotalStrength
}

func (h ComboMetricHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *ComboMetricHeap) Push(x interface{}) {
	*h = append(*h, x.(ComboMetric))
}

func (h *ComboMetricHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

```

# 结果展示
这篇博客最最重要的环节要来了，我们需要检验，计算机搜索出来的最强阵容，是否和S1版本的吃鸡阵容是相符的。全部结果文件在[result](https://github.com/weiziyoung/TFT/tree/master/data/final_result)里，读者也可以自己把代码下下来编译跑一下。

另外因为很多阵容之间的区别仅仅是换了一个相同羁绊的英雄，或者改了一个小羁绊，所以我们这里对搜索结果做了一个很简单的去重融合，当两个阵容羁绊相似度过高时进行合并，相似度可以用[Jaccard similarity coefficient](https://en.wikipedia.org/wiki/Jaccard_index) 来计算集合之间的相似度，如果相似度大于0.7，则认为属于同一套阵容：
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x37ebee8b58e0a2image.png)

## 羁绊数最多的阵容
首先我们看不可能有错的一个指标——羁绊数。直观来说，就是搜索出让左边的羁绊灯亮最多的阵容（这种阵容不一定强)
- 六人口
```json
    {
        "combo": ["艾希","狗熊","机器人","劫","螳螂","卡萨"],
        "trait_num": 7,
        "trait_detail": {
            "冰川": 2,
            "刺客": 3,
            "忍者": 1,
            "斗士": 2,
            "机器人": 1,
            "游侠": 2,
            "虚空": 2
        },
        "total_trait_num": 12,
        "total_trait_strength": 16.4,
        "total_strength": 19.052499984405003
    },
```
总羁绊数达到了7个羁绊，注意这是6人口，正常咱们玩自走棋，6人口大约是4个羁绊数左右，毕竟阵容还没成型，但是实际上6人口在不用铲子的情况下最多可以有7个羁绊。

- 七人口
```json
{
        "combo": [
            "狗熊","猪妹","机器人","慎","船长","卡密尔","金克斯"
        ],
        "trait_num": 7,
        "trait_detail": {
            "冰川": 2,
            "剑士": 3,
            "忍者": 1,
            "斗士": 2,
            "机器人": 1,
            "枪手": 2,
            "海克斯": 2
        },
        "total_trait_num": 17,
        "total_trait_strength": 22.6,
        "total_strength": 21.726248967722068
    },
```
七人口最大羁绊数竟然还是7。不过不同于6人口只有一种组合能达到7羁绊，七人口前100个中基本都达到了7羁绊。

- 八人口
```json
    {
        "combo": [
            "艾希", "狗熊","机器人","劫","螳螂","挖掘机",
            "大虫子","卡萨"
        ],
        "trait_num": 9,
        "trait_detail": {
            "冰川": 2,
            "刺客": 3,
            "忍者": 1,
            "斗士": 4,
            "机器人": 1,
            "游侠": 2,
            "虚空": 4
        },
        "total_trait_num": 25,
        "total_trait_strength": 24.4,
        "total_strength": 27.058749668872913
    }
```
总共是9个羁绊，看着阵容好像是虚空斗刺哈哈哈，但虚空斗刺没有艾希。这套阵容强度看上去还是可以的。

- 九人口
```json
    {
        "combo": [
            "狗熊","猪妹","机器人","盖伦","薇恩","天使","劫","螳螂","卡萨"
        ],
        "trait_num": 9,
        "trait_detail": {
            "冰川": 2,
            "刺客": 3,
            "忍者": 1,
            "斗士": 2,
            "机器人": 1,
            "游侠": 2,
            "虚空": 2,
            "贵族": 3,
            "骑士": 2
        },
        "total_trait_num": 22,
        "total_trait_strength": 28.55,
        "total_strength": 31.621585006726214
    }
```
总之我没看过亮9栈灯的阵容，看样子挺花哨的，但这个阵容其实不妥的。羁绊只是吃鸡的一小部分，实际上更多的需要依靠英雄等级、装备、输出和坦克的组合。

- 十人口
```json
    {
        "combo": [
            "维鲁斯","乌鸦","亚索","机器人","诺手","天使",
            "阿卡丽","螳螂","挖掘机","卡萨"
        ],
        "trait_num": 10,
        "trait_detail": {
            "刺客": 3,"帝国": 2,"忍者": 1,"恶魔": 2,
            "斗士": 2,"机器人": 1,"浪人": 1,"游侠": 2,
            "虚空": 2,"骑士": 2
        },
        "total_trait_num": 24,
        "total_trait_strength": 31,
        "total_strength": 37.67076561712962
    }
```
亮了10栈灯，这种阵容基本看看就好，不可能成型并且吃鸡的，因为这是个有5个5费卡的阵容。

## 强度最高的阵容
正如之前说的，羁绊多阵容并不一定强，所以一定要结合英雄等级、羁绊强度、羁绊范围这些来算，这里英雄等级的增益和羁绊强度都是具有主观判断在里面的，而且算上这些指标实际上也是不够的，看下计算出的阵容就知道了:

- 六人口
```json
    {
        "combo": [
            "潘森","布隆","丽桑卓","狗熊","冰鸟","凯南"
        ],
        "trait_num": 5,
        "trait_detail": {
            "元素师": 3,
            "冰川": 4,
            "忍者": 1,
            "护卫": 2
        },
        "total_trait_num": 14,
        "total_trait_strength": 15.2,
        "total_strength": 30.272461525164545
    },
```
这看上去是一个冰川元素阵容，游戏刚出的时候，这套阵容还是很容易吃鸡的，主要就是利用丽桑卓和冰鸟都是冰川+元素，导致这套阵容又有控制又有坦度，在以前谁都不会玩这个游戏的年代很容易吃鸡，小编我第一次吃鸡用的就是冰川元素流。但冰川元素逐渐没落了，原因就是后来大家都会玩这个游戏了，导致游戏节奏加快，而这个阵容一个最大的缺点就是成型有点困难，猪妹和冰鸟都不是那么容易抽到的，前期靠布隆一个坦度点是肯定不够的。

- 七人口
```json
    {
        "combo": [
            "莫甘娜","龙王","潘森","日女","天使","铁男","死歌"
        ],
        "trait_num": 5,
        "trait_detail": {
            "幽灵": 2,"护卫": 2,"法师": 3,"骑士": 2,"龙": 2
        },
        "total_trait_num": 22,
        "total_trait_strength": 29.9,
        "total_strength": 40.17980836913922
    },
```
这个看上去是护卫龙，但又不太像，因为护卫龙好像没有人配法师的，但这不是最重要的，最重要的是，这套阵容太不容易成型了！！因为我们的评价指标里没有考虑羁绊的成型难易度，导致它更偏好等级高的英雄，强度看上去还可以，有输出有坦克，但有谁7人口能凑出来3个五星，2个四星呢？

- 八人口
```json
    {
        "combo": [
            "龙王","潘森","布隆","丽桑卓","冰鸟",
            "凯南","露露","小法"
        ],
        "trait_num": 7,
        "trait_detail": {
            "元素师": 3,"冰川": 2,"忍者": 1,
            "护卫": 2,"法师": 3,"约德尔": 3,
            "龙": 2
        },
        "total_trait_num": 24,
        "total_trait_strength": 34.100002,
        "total_strength": 50.979002334643155
    }
```
跟上面有点像(其实我不太清楚为什么七八人口都是护卫龙)，这套阵容其实是缺乏坦度的hhh还不容易成型。所以我们的评估指标还是有问题哈哈哈，看到这套阵容人傻了。

- 九人口
```json
{
        "combo": [
            "潘森","亚索","剑姬","盖伦",
            "薇恩","卢锡安","日女","天使","船长"
        ],
        "trait_num": 7,
        "trait_detail": {
            "剑士": 3,"护卫": 2,"枪手": 2,
            "浪人": 1,"贵族": 6,"骑士": 2
        },
        "total_trait_num": 47,
        "total_trait_strength": 54.249996,
        "total_strength": 61.73055001568699
    }
```
这套阵容我还是用过的，能不能吃鸡要看装备，亚索能2星并且吃到装备基本能吃鸡，吃不到装备就很缺乏输出，据说也可以把装备给船长养船长这个点，不过没试过。九人口贵族崛起大概是因为贵族的全范围buff比较给力。

# 分析与总结

## 贡献
直到云顶之弈S1结束，网上并没有一篇用图搜索来组建羁绊阵容的文章，这篇文章就当是弥补这一块的空白吧，它从另一个角度去为我们推荐了阵容。核心思想就是利用英雄之间的相互羁绊来简化暴力搜索。

## 缺陷

实际上我觉得在评估阵容强度的时候，模型还是过于粗糙的，具体表现如下:

1. **首先忽视了坦度和输出的配合这个维度**。导致有些推荐阵容全是坦克没有输出，有些阵容只有输出没有坦克。


2. **其次忽视了羁绊之间的克制关系**。可以看到七八人口的时候，计算出来的都是以护卫龙为核心的阵容，因为护卫羁绊提供的收益范围很大，但前提条件是你把英雄都集中放护卫周围，但这种方法实际上是被海克斯完克的，所以在实际时间上，护卫buff的收益并没有这里计算中的那么大。

3. **忽略了阵过渡的平滑程度**。这是这里存在的最大问题，由于我们在评价阵容的时候，给高等级英雄倾向了一些权重，导致阵容中会有数量较多的高费英雄，实际上不考虑阵容成型难易程度的推荐就是在耍流氓。比如潘森刚出来的时候，很多人推荐贵族护卫龙，实际应用上效果并不好。

4. **没有考虑英雄升星的难易程度**。这个实际上跟上面是一种问题，我在搜索结果里找赌刺的阵容，直接被排名拍到了40多名，但赌刺绝对是6人口的T1阵容，这里面的原因就是刺客的卡费普遍是低的，导致在这套算法里赚不到便宜，但其实低费卡更容易到三星，而三星低费卡的强度是高于高费卡的，尤其是像三星劫这样的英雄。
![image.png](https://wzy-zone.oss-cn-shanghai.aliyuncs.com/article_images%2F0x37ebf0f7478498image.png)

5. **没有考虑金铲铲**。因为简化问题，这里没有考虑金铲铲，如果考虑金铲铲的话，搜索空间将会变得极其庞大，相当于为每个英雄都给配剑士、刺客、骑士、冰川、约德尔、恶魔这些羁绊。这些加上去以后，复杂度也就跟全搜索差不多了。

## 踩坑记录

1. Golang append函数，函数原型如下:
```go
func append(slice []Type, elems ...Type) []Type
```
从原型上看是传入一个切片，和若干需要加入的元素，返回一个切片。但实际上传入的slice切片在运行的过程中会被修改，返回的那个切片实际上就是你传入的slice切片。所以在使用golang里面的append函数的时候，记得把接受变量设置成你传入的第一个slice变量，或者使用前对slice进行copy。

2. 保留前K大个数实际上要用小顶堆，而不是想当然地使用大顶堆。

3. 在考虑当前英雄的后续结点的时候，不能只考虑当前英雄的羁绊，而要考虑队伍里所有英雄的羁绊，否则会漏检。