# coding: utf-8
import string
import itertools


"""
Uncomplete:
    x3b.cb.gx

"""

validators = dict()


def Validate(tag, code):
    return validators[tag](code)


def bind(tag):
    def wrapper(func):
        # register tag
        validators[tag] = func

        def inner_wrapper(*args, **kwargs):
            num = func(*args, **kwargs)
            # log detail
            log_fmt = "Validate \n Tag: {tag}  Code: {args[0]} \n Num: {num}"
            print log_fmt.format(tag=tag, args=args, kwargs=kwargs, num=num)
            return num
        return inner_wrapper
    return wrapper


# 五星/直选/复式
@bind("x5.eq.batch")
def x5_eq_batch(code):
    num = 0
    if len(code) == 5:
        data = [i for i in code if len(i) > 0]
        if len(data) > 0:
            seed = [[i for i in n] for n in data]
            items = [item for item in itertools.product(*seed)]
            num = len(items)
    return num


# 五星/直选/单式
@bind("x5.eq.simple")
def x5_eq_simple(code):
    num = 0
    if len(code) > 0:
        for n in code:
            if len(n) != 5:
                return 0
            else:
                num += 1
    return num


# 五星/组选/组选120
@bind("x5.cb.g120")
def x5_cb_g120(code):
    num = 0
    if len(code) >= 5:
        seed = "".join(code)
        num = len([s for s in itertools.combinations(seed, 5)])
    return num


# 五星/组选/组选60
@bind("x5.cb.g60")
def x5_cb_g60(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == 2:
        result = []
        for p in data[0]:
            items = [i for i in itertools.combinations(data[1], 3)]
            for item in items:
                if p not in item:
                    result.append(item)
        num = len(result)
    return num


# 五星/组选/组选30
@bind("x5.cb.g30")
def x5_cb_g30(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == 2:
        result = []
        for p in data[1]:
            items = [i for i in itertools.combinations(data[0], 2)]
            for item in items:
                if p not in item:
                    result.append(item)
        num = len(result)
    return num


# 五星/组选/组选20
@bind("x5.cb.g20")
def x5_cb_g20(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == 2:
        result = []
        for p in data[0]:
            items = [i for i in itertools.combinations(data[1], 2)]
            for item in items:
                if p not in item:
                    result.append(item)
        num = len(result)
    return num


# 五星/组选/组选10
@bind("x5.cb.g10")
def x5_cb_g10(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == 2:
        result = []
        for p in data[0]:
            items = [i for i in itertools.combinations(data[1], 1)]
            for item in items:
                if p not in item:
                    result.append(item)
        num = len(result)
    return num


# 五星/组选/组选5
@bind("x5.cb.g5")
def x5_cb_g5(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == 2:
        result = []
        for p in data[0]:
            items = [i for i in itertools.combinations(data[1], 1)]
            for item in items:
                if p not in item:
                    result.append(item)
        num = len(result)
    return num


#五星/不定位/一码不定位
@bind("x5.unf.m1")
def x5_unf_m1(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == len(code):
        num = len(data)
    return num


#五星/不定位/二码不定位
@bind("x5.unf.m2")
def x5_unf_m2(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) >= 2:
        seed = "".join(data)
        items = [i for i in itertools.combinations(seed, 2)]
        num = len(items)
    return num


#五星/不定位/三码不定位
@bind("x5.unf.m3")
def x5_unf_m3(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) >= 3:
        seed = "".join(data)
        items = [i for i in itertools.combinations(seed, 3)]
        num = len(items)
    return num


# 五星/趣味/一帆风顺
@bind("x5.fun.f1")
def x5_fun_f1(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == len(code):
        num = len(data)
    return num


# 五星/趣味/好事成双
@bind("x5.fun.f2")
def x5_fun_f2(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == len(code):
        num = len(data)
    return num


# 五星/趣味/三星报喜
@bind("x5.fun.f3")
def x5_fun_f3(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == len(code):
        num = len(data)
    return num


# 五星/趣味/四季发财
@bind("x5.fun.f4")
def x5_fun_f4(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == len(code):
        num = len(data)
    return num


# 前四/直选/复式
@bind("x4f.eq.batch")
def x4f_eq_batch(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) > 0:
        seed = [[i for i in n] for n in data]
        items = [item for item in itertools.product(*seed)]
        num = len(items)
    return num


# 前四/直选/单式
@bind("x4f.eq.simple")
def x4f_eq_simple(code):
    num = 0
    if len(code) > 0:
        for n in code:
            if len(n) != 4:
                return 0
            else:
                num += 1
    return num


# 前四/组选/组选24
@bind("x4f.cb.g24")
def x4f_cb_g24(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 4:
        seed = "".join(code)
        items = [i for i in itertools.combinations(seed, 4)]
        num = len(items)
    return num


# 前四/组选/组选12
@bind("x4f.cb.g12")
def x4f_cb_g12(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == 2:
        result = []
        for p in data[0]:
            items = [i for i in itertools.combinations(data[1], 2)]
            for item in items:
                if p not in item:
                    result.append(item)
        num = len(result)
    return num


# 前四/组选/组选6
@bind("x4f.cb.g6")
def x4f_cb_g6(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) >= 2:
        seed = "".join(data)
        items = [i for i in itertools.combinations(seed, 2)]
        num = len(items)
    return num


# 前四/组选/组选4
@bind("x4f.cb.g4")
def x4f_cb_g4(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == 2:
        result = []
        for p in data[0]:
            items = [i for i in itertools.combinations(data[1], 1)]
            for item in items:
                if p not in item:
                    result.append(item)
        num = len(result)
    return num


# 前四/不定位/一码不定位
@bind("x4f.unf.m1")
def x4f_unf_m1(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == len(code):
        num = len(data)
    return num


# 前四/不定位/二码不定位
@bind("x4f.unf.m2")
def x4f_unf_m2(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) >= 2:
        seed = "".join(data)
        items = [i for i in itertools.combinations(seed, 2)]
        num = len(items)
    return num

# 后四/直选/复式
@bind("x4b.eq.batch")
def x4b_eq_batch(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) > 0:
        seed = [[i for i in n] for n in data]
        items = [item for item in itertools.product(*seed)]
        num = len(items)
    return num


# 后四/直选/单式
@bind("x4b.eq.simple")
def x4b_eq_simple(code):
    num = 0
    if len(code) > 0:
        for n in code:
            if len(n) != 4:
                return 0
            else:
                num += 1
    return num


# 后四/组选/组选24
@bind("x4b.cb.g24")
def x4b_cb_g24(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 4:
        seed = "".join(code)
        items = [i for i in itertools.combinations(seed, 4)]
        num = len(items)
    return num


# 后四/组选/组选12
@bind("x4b.cb.g12")
def x4b_cb_g12(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == 2:
        result = []
        for p in data[0]:
            items = [i for i in itertools.combinations(data[1], 2)]
            for item in items:
                if p not in item:
                    result.append(item)
        num = len(result)
    return num


# 后四/组选/组选6
@bind("x4b.cb.g6")
def x4b_cb_g6(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) >= 2:
        seed = "".join(data)
        items = [i for i in itertools.combinations(seed, 2)]
        num = len(items)
    return num


# 后四/组选/组选4
@bind("x4b.cb.g4")
def x4b_cb_g4(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == 2:
        result = []
        for p in data[0]:
            items = [i for i in itertools.combinations(data[1], 1)]
            for item in items:
                if p not in item:
                    result.append(item)
        num = len(result)
    return num


# 后四/不定位/一码不定位
@bind("x4b.unf.m1")
def x4b_unf_m1(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == len(code):
        num = len(data)
    return num


# 后四/不定位/二码不定位
@bind("x4b.unf.m2")
def x4b_unf_m2(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) >= 2:
        seed = "".join(data)
        items = [i for i in itertools.combinations(seed, 2)]
        num = len(items)
    return num


# 前三/直选/复式
@bind("x3f.eq.batch")
def x3f_eq_batch(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) > 0:
        seed = [[i for i in n] for n in code]
        items = [item for item in itertools.product(*seed)]
        num = len(items)
    return num


# 前三/直选/单式
@bind("x3f.eq.simple")
def x3f_eq_simple(code):
    num = 0
    if len(code) > 0:
        for n in code:
            if len(n) != 3:
                return 0
            else:
                num += 1
    return num


# 前三/直选/和值
@bind("x3f.eq.sum")
def x3f_eq_sum(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 1:
        result = {}
        seed = [string.digits] * 3
        items = [item for item in itertools.product(*seed)]
        for item in items:
            count = sum([int(i) for i in item])
            key = str(count)
            if key not in result.keys():
                result.setdefault(key, [])
            if item not in result[key]:
                result[key].append(item)
        num = sum([len(result[key]) for key in data])
    return num




# 前三/直选/跨度
@bind("x3f.eq.diff")
def x3f_eq_diff(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 1:
        result = {}
        seed = [string.digits] * 3
        items = [item for item in itertools.product(*seed)]
        for item in items:
            item_ = [int(i) for i in item]
            count = str(max(item_) - min(item_))
            key = str(count)
            if key not in result.keys():
                result.setdefault(key, [])
            if item not in result[key]:
                result[key].append(item)
        num = sum([len(result[key]) for key in data])
    return num


# 前三/组选/组选复式
@bind("x3f.cb.gx")
def x3f_cb_gx(code):
    num = 0
    return num


# 前三/组选/和值
@bind("x3f.cb.sum")
def x3f_cb_sum(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 1:
        result = {}
        seed = [string.digits] * 3
        items = [item for item in itertools.product(*seed) if len(set(item))>=2]
        for item in items:
            count = sum([int(i) for i in item])
            key = str(count)
            item_ = list(item)
            item_.sort()
            if key not in result.keys():
                result.setdefault(key, [])
            if item_ not in result[key]:
                result[key].append(item_)
        num = sum([len(result[key]) for key in data])
    return num


# 前三/组选/混合
@bind("x3f.cb.mix")
def x3f_cb_mix(code):
    num = 0
    data = [i for i in code if len(i) == 3]
    if len(data) == len(code):
        num = len(data)
    return num


# 前三/组选/包胆
@bind("x3f.cb.with")
def x3f_cb_with(code):
    num = 0
    data = code
    if len(data) == 1:
        seed = [string.digits] * 3
        items = [item for item in itertools.product(*seed) 
            if data[0] in item and len(set(item))==2]
        num = len(items)
    return num


# 前三/组选/组三复式
@bind("x3f.cb.g3")
def x3f_cb_g3(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 2:
        seed = "".join(code)
        items = [i for i in itertools.combinations(seed, 2)]
        num = len(items)
    return num


# 前三/组选/组六复式
@bind("x3f.cb.g6")
def x3f_cb_g6(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 3:
        seed = "".join(code)
        items = [i for i in itertools.combinations(seed, 3)]
        num = len(items)
    return num


# 前三/不定位/一码不定位
@bind("x3f.unf.m1")
def x3f_unf_m1(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == len(code):
        num = len(data)
    return num


# 前三/不定位/二码不定位
@bind("x3f.unf.m2")
def x3f_unf_m2(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) >= 2:
        seed = "".join(data)
        items = [i for i in itertools.combinations(seed, 2)]
        num = len(items)
    return num


# 中三/直选/复式
@bind("x3m.eq.batch")
def x3m_eq_batch(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) > 0:
        seed = [[i for i in n] for n in code]
        items = [item for item in itertools.product(*seed)]
        num = len(items)
    return num


# 中三/直选/单式
@bind("x3m.eq.simple")
def x3m_eq_simple(code):
    num = 0
    if len(code) > 0:
        for n in code:
            if len(n) != 3:
                return 0
            else:
                num += 1
    return num


# 中三/直选/和值
@bind("x3m.eq.sum")
def x3m_eq_sum(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 1:
        result = {}
        seed = [string.digits] * 3
        items = [item for item in itertools.product(*seed)]
        for item in items:
            count = sum([int(i) for i in item])
            key = str(count)
            if key not in result.keys():
                result.setdefault(key, [])
            if item not in result[key]:
                result[key].append(item)
        num = sum([len(result[key]) for key in data])
    return num




# 中三/直选/跨度
@bind("x3m.eq.diff")
def x3m_eq_diff(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 1:
        result = {}
        seed = [string.digits] * 3
        items = [item for item in itertools.product(*seed)]
        for item in items:
            item_ = [int(i) for i in item]
            count = str(max(item_) - min(item_))
            key = str(count)
            if key not in result.keys():
                result.setdefault(key, [])
            if item not in result[key]:
                result[key].append(item)
        num = sum([len(result[key]) for key in data])
    return num


# 中三/组选/组选复式
@bind("x3m.cb.gx")
def x3f_cb_gx(code):
    num = 0
    return num


# 中三/组选/和值
@bind("x3m.cb.sum")
def x3m_cb_sum(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 1:
        result = {}
        seed = [string.digits] * 3
        items = [item for item in itertools.product(*seed) if len(set(item))>=2]
        for item in items:
            count = sum([int(i) for i in item])
            key = str(count)
            item_ = list(item)
            item_.sort()
            if key not in result.keys():
                result.setdefault(key, [])
            if item_ not in result[key]:
                result[key].append(item_)
        num = sum([len(result[key]) for key in data])
    return num


# 中三/组选/混合
@bind("x3m.cb.mix")
def x3m_cb_mix(code):
    num = 0
    data = [i for i in code if len(i) == 3]
    if len(data) == len(code):
        num = len(data)
    return num


# 中三/组选/包胆
@bind("x3m.cb.with")
def x3m_cb_with(code):
    num = 0
    data = code
    if len(data) == 1:
        seed = [string.digits] * 3
        items = [item for item in itertools.product(*seed) 
            if data[0] in item and len(set(item))==2]
        num = len(items)
    return num


# 中三/组选/组三复式
@bind("x3m.cb.g3")
def x3m_cb_g3(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 2:
        seed = "".join(code)
        items = [i for i in itertools.combinations(seed, 2)]
        num = len(items)
    return num


# 中三/组选/组六复式
@bind("x3m.cb.g6")
def x3m_cb_g6(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 3:
        seed = "".join(code)
        items = [i for i in itertools.combinations(seed, 3)]
        num = len(items)
    return num


# 中三/不定位/一码不定位
@bind("x3m.unf.m1")
def x3m_unf_m1(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == len(code):
        num = len(data)
    return num


# 中三/不定位/二码不定位
@bind("x3m.unf.m2")
def x3m_unf_m2(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) >= 2:
        seed = "".join(data)
        items = [i for i in itertools.combinations(seed, 2)]
        num = len(items)
    return num


# 后三/直选/复式
@bind("x3b.eq.batch")
def x3b_eq_batch(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) > 0:
        seed = [[i for i in n] for n in code]
        items = [item for item in itertools.product(*seed)]
        num = len(items)
    return num


# 后三/直选/单式
@bind("x3b.eq.simple")
def x3b_eq_simple(code):
    num = 0
    if len(code) > 0:
        for n in code:
            if len(n) != 3:
                return 0
            else:
                num += 1
    return num


# 后三/直选/和值
@bind("x3b.eq.sum")
def x3b_eq_sum(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 1:
        result = {}
        seed = [string.digits] * 3
        items = [item for item in itertools.product(*seed)]
        for item in items:
            count = sum([int(i) for i in item])
            key = str(count)
            if key not in result.keys():
                result.setdefault(key, [])
            if item not in result[key]:
                result[key].append(item)
        num = sum([len(result[key]) for key in data])
    return num




# 后三/直选/跨度
@bind("x3b.eq.diff")
def x3b_eq_diff(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 1:
        result = {}
        seed = [string.digits] * 3
        items = [item for item in itertools.product(*seed)]
        for item in items:
            item_ = [int(i) for i in item]
            count = str(max(item_) - min(item_))
            key = str(count)
            if key not in result.keys():
                result.setdefault(key, [])
            if item not in result[key]:
                result[key].append(item)
        num = sum([len(result[key]) for key in data])
    return num


# 后三/组选/组选复式
@bind("x3b.cb.gx")
def x3b_cb_gx(code):
    num = 0
    
    return num


# 后三/组选/和值
@bind("x3b.cb.sum")
def x3b_cb_sum(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 1:
        result = {}
        seed = [string.digits] * 3
        items = [item for item in itertools.product(*seed) if len(set(item))>=2]
        for item in items:
            count = sum([int(i) for i in item])
            key = str(count)
            item_ = list(item)
            item_.sort()
            if key not in result.keys():
                result.setdefault(key, [])
            if item_ not in result[key]:
                result[key].append(item_)
        num = sum([len(result[key]) for key in data])
    return num


# 后三/组选/混合
@bind("x3b.cb.mix")
def x3b_cb_mix(code):
    num = 0
    data = [i for i in code if len(i) == 3]
    if len(data) == len(code):
        num = len(data)
    return num


# 后三/组选/包胆
@bind("x3b.cb.with")
def x3b_cb_with(code):
    num = 0
    data = code
    if len(data) == 1:
        seed = [string.digits] * 3
        items = [item for item in itertools.product(*seed) 
            if data[0] in item and len(set(item))==2]
        num = len(items)
    return num


# 后三/组选/组三复式
@bind("x3b.cb.g3")
def x3b_cb_g3(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 2:
        seed = "".join(code)
        items = [i for i in itertools.combinations(seed, 2)]
        num = len(items)
    return num


# 后三/组选/组六复式
@bind("x3b.cb.g6")
def x3b_cb_g6(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 3:
        seed = "".join(code)
        items = [i for i in itertools.combinations(seed, 3)]
        num = len(items)
    return num


# 后三/不定位/一码不定位
@bind("x3b.unf.m1")
def x3b_unf_m1(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) == len(code):
        num = len(data)
    return num


# 后三/不定位/二码不定位
@bind("x3b.unf.m2")
def x3b_unf_m2(code):
    num = 0
    data = [i for i in code if len(i)>0]
    if len(data) >= 2:
        seed = "".join(data)
        items = [i for i in itertools.combinations(seed, 2)]
        num = len(items)
    return num


# 前二/直选/复式
@bind("x2f.eq.batch")
def x2f_eq_batch(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) > 0:
        seed = [[i for i in n] for n in code]
        items = [item for item in itertools.product(*seed)]
        num = len(items)
    return num


# 前二/直选/单式
@bind("x2f.eq.simple")
def x2f_eq_simple(code):
    num = 0
    if len(code) > 0:
        for n in code:
            if len(n) != 2:
                return 0
            else:
                num += 1
    return num


# 前二/直选/和值
@bind("x2f.eq.sum")
def x2f_eq_sum(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 1:
        result = {}
        seed = [string.digits] * 2
        items = [item for item in itertools.product(*seed) if len(set(item))>=1]
        for item in items:
            count = sum([int(i) for i in item])
            key = str(count)
            if key not in result.keys():
                result.setdefault(key, [])
            if item not in result[key]:
                result[key].append(item)
        num = sum([len(result[key]) for key in data])
    return num


# 前二/直选/跨度
@bind("x2f.eq.diff")
def x2f_eq_diff(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 1:
        if len(data) >= 1:
            result = {}
            seed = [string.digits] * 2
            items = [item for item in itertools.product(*seed) if len(set(item))>=1]
            for item in items:
                item_ = [int(i) for i in item]
                diff = max(item_)-min(item_)
                key = str(diff)
                if key not in result.keys():
                    result.setdefault(key, [])
                if item not in result[key]:
                    result[key].append(item)
        num = sum([len(result[key]) for key in data])
    return num


# 前二/组选/复式
@bind("x2f.cb.batch")
def x2f_cb_batch(code):
    num = 0
    data = dict([(i, len(i) == 2) for i in code])
    if False not in data.values():
        num = len(data.keys())
    return num


# 前二/组选/单式
@bind("x2f.cb.simple")
def x2f_cb_simple(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 2:
        seed = "".join(data)
        items = [i for i in itertools.combinations(seed, 2)]
        num = len(items)
    return num


# 前二/组选/和值
@bind("x2f.cb.sum")
def x2f_cb_sum(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 1:
        result = {}
        seed = [string.digits] * 2
        items = [item for item in itertools.product(*seed) if len(set(item))==2]
        for item in items:
            count = sum([int(i) for i in item])
            key = str(count)
            item_ = list(item)
            item_.sort()
            if key not in result.keys():
                result.setdefault(key, [])
            if item_ not in result[key]:
                result[key].append(item_)
        num = sum([len(result[key]) for key in data])
    return num


# 前二/组选/包胆
@bind("x2f.cb.with")
def x2f_cb_with(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) == 1:
        seed = string.digits
        items = [item for item in itertools.combinations(seed, 2) if data[0] in item and len(set(item))==2]
        num = len(items)
    return num


# 后二/直选/复式
@bind("x2b.eq.batch")
def x2b_eq_batch(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) > 0:
        seed = [[i for i in n] for n in code]
        items = [item for item in itertools.product(*seed)]
        num = len(items)
    return num


# 后二/直选/单式
@bind("x2b.eq.simple")
def x2b_eq_simple(code):
    num = 0
    if len(code) > 0:
        for n in code:
            if len(n) != 2:
                return 0
            else:
                num += 1
    return num


# 后二/直选/和值
@bind("x2b.eq.sum")
def x2b_eq_sum(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 1:
        result = {}
        seed = [string.digits] * 2
        items = [item for item in itertools.product(*seed) if len(set(item))>=1]
        for item in items:
            count = sum([int(i) for i in item])
            key = str(count)
            if key not in result.keys():
                result.setdefault(key, [])
            if item not in result[key]:
                result[key].append(item)
        num = sum([len(result[key]) for key in data])
    return num


# 后二/直选/跨度
@bind("x2b.eq.diff")
def x2b_eq_diff(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 1:
        if len(data) >= 1:
            result = {}
            seed = [string.digits] * 2
            items = [item for item in itertools.product(*seed) if len(set(item))>=1]
            for item in items:
                item_ = [int(i) for i in item]
                diff = max(item_)-min(item_)
                key = str(diff)
                if key not in result.keys():
                    result.setdefault(key, [])
                if item not in result[key]:
                    result[key].append(item)
        num = sum([len(result[key]) for key in data])
    return num


# 后二/组选/复式
@bind("x2b.cb.batch")
def x2b_cb_batch(code):
    num = 0
    data = dict([(i, len(i) == 2) for i in code])
    if False not in data.values():
        num = len(data.keys())
    return num


# 后二/组选/单式
@bind("x2b.cb.simple")
def x2b_cb_simple(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 2:
        seed = "".join(data)
        items = [i for i in itertools.combinations(seed, 2)]
        num = len(items)
    return num


# 后二/组选/和值
@bind("x2b.cb.sum")
def x2b_cb_sum(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) >= 1:
        result = {}
        seed = [string.digits] * 2
        items = [item for item in itertools.product(*seed) if len(set(item))==2]
        for item in items:
            count = sum([int(i) for i in item])
            key = str(count)
            item_ = list(item)
            item_.sort()
            if key not in result.keys():
                result.setdefault(key, [])
            if item_ not in result[key]:
                result[key].append(item_)
        num = sum([len(result[key]) for key in data])
    return num


# 后二/组选/包胆
@bind("x2b.cb.with")
def x2b_cb_with(code):
    num = 0
    data = [i for i in code if len(i) > 0]
    if len(data) == 1:
        seed = string.digits
        items = [item for item in itertools.combinations(seed, 2) if data[0] in item and len(set(item))==2]
        num = len(items)
    return num


# 一星/定位胆/复式
@bind("x1.fix.batch")
def x1_fix_batch(code):
    num = 0
    data = [len(i) for i in code]
    num = reduce(lambda x, y: x + y, data)
    return num

'''
# 
@bind("")
def func(code):
    num = 0
    return num
'''

if __name__ == '__main__':
    # print validators["x5.eq.simple"](["01234", "56789"])
    # print Validate("x5.eq.simple", ["01234", "56789"])

    x5_eq_batch(["01", "34", "56", "78", "9"])
    x5_eq_simple(["01234", "56789", "56790"])
    