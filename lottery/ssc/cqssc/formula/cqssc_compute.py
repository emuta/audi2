# coding: utf-8
import string
import itertools


"""
Uncomplete:
    x3b.cb.gx

"""

computes = dict()


def Compute(tag, result, code):
    return computes[tag](result, code)


def bind(tag):
    def wrapper(func):
        # register tag
        computes[tag] = func

        def inner_wrapper(*args, **kwargs):
            num = func(*args, **kwargs)
            
            # log detail
            log_fmt = "Compute \n Tag: {tag} \n Result: {args[0]} \n Code: {args[1]} \n Num: {num}"
            print log_fmt.format(tag=tag, args=args, kwargs=kwargs, num=num)

            return num
        return inner_wrapper
    return wrapper


# 五星/直选/复式
@bind("x5.eq.batch")
def x5_eq_batch(result, code):
    num = 0
    result = tuple(result)
    data = [[i for i in n] for n in code]
    items = [i for i in itertools.product(*data)]
    num = items.count(result)
    return num


# 五星/直选/单式
@bind("x5.eq.simple")
def x5_eq_simple(result, code):
    num = 0
    result = "".join([i for i in result])
    num = code.count(result)
    return num


# 五星/组选/组选120
@bind("x5.cb.g120")
def x5_cb_g120(result, code):
    num = 0
    data = [i for i in code]
    items = filter(lambda x: x in data, result)
    if len(items) == 5:
        num = 1
    return num


# 五星/组选/组选60
@bind("x5.cb.g60")
def x5_cb_g60(result, code):
    num = 0
    items = [i for i in itertools.product([i for i in code[1]], repeat=3)]
    for p in code[0]:
        if result.count(p)==2:
            ex = [i for i in result if i!=p]
            for item in items:
                if ex == list(item):
                    num = 1
    return num


# 五星/组选/组选30
@bind("x5.cb.g30")
def x5_cb_g30(result, code):
    num = 0
    items = [i for i in itertools.combinations(code[0], 2)]
    for item in items:
        p1 = item[0]
        p2 = item[1]
        if result.count(p1) == 2 and result.count(p2) == 2:
            ex = [i for i in result if i not in item]
            for i in code[1]:
                if i in ex:
                    num = 1
    return num


# 五星/组选/组选20
@bind("x5.cb.g20")
def x5_cb_g20(result, code):
    num = 0
    items = [i for i in itertools.combinations(code[1], 2)]
    for p in code[0]:
        if result.count(p) == 3:
            ex = [i for i in result if i!= p]
            for item in items:
                if ex == list(item):
                    num = 1
    return num


# 五星/组选/组选10
@bind("x5.cb.g10")
def x5_cb_g10(result, code):
    num = 0
    for p in code[0]:
        if result.count(p) == 3:
            ex = [i for i in result if i!= p]
            for n in code[1]:
                if ex.count(n) == 2:
                    num = 1
    return num


# 五星/组选/组选5
@bind("x5.cb.g5")
def x5_cb_g5(result, code):
    num = 0
    for p in code[0]:
        if result.count(p) == 4:
            ex = [i for i in result if i!= p]
            for n in code[1]:
                if ex.count(n):
                    num = 1
    return num


#五星/不定位/一码不定位
@bind("x5.unf.m1")
def x5_unf_m1(result, code):
    num = 0
    for n in code:
        if result.count(n)>=1:
            num = 1
    return num


#五星/不定位/二码不定位
@bind("x5.unf.m2")
def x5_unf_m2(result, code):
    num = 0
    if len([i for i in code if i in result]) >= 2:
        num = 1
    return num


#五星/不定位/三码不定位
@bind("x5.unf.m3")
def x5_unf_m3(result, code):
    num = 0
    if len([i for i in code if i in result]) >= 3:
        num = 1
    return num


# 五星/趣味/一帆风顺
@bind("x5.fun.f1")
def x5_fun_f1(result, code):
    num = 0
    for n in code:
        if result.count(n) >= 1:
            num = 1
    return num


# 五星/趣味/好事成双
@bind("x5.fun.f2")
def x5_fun_f2(result, code):
    num = 0
    for n in code:
        if result.count(n) >= 2:
            num = 1
    return num


# 五星/趣味/三星报喜
@bind("x5.fun.f3")
def x5_fun_f3(result, code):
    num = 0
    for n in code:
        if result.count(n) >= 3:
            num = 1
    return num


# 五星/趣味/四季发财
@bind("x5.fun.f4")
def x5_fun_f4(result, code):
    num = 0
    for n in code:
        if result.count(n) >= 4:
            num = 1
    return num


# 前四/直选/复式
@bind("x4f.eq.batch")
def x4f_eq_batch(result, code):
    num = 0
    result = tuple(result[:4])
    seed = [[i for i in n] for n in code]
    items = [i for i in itertools.product(*seed)]
    num = items.count(result)
    return num


# 前四/直选/单式
@bind("x4f.eq.simple")
def x4f_eq_simple(result, code):
    num = 0
    result = result[:4]
    item = "".join([i for i in result])
    num = code.count(item)
    return num


# 前四/组选/组选24
@bind("x4f.cb.g24")
def x4f_cb_g24(result, code):
    num = 0
    result = tuple(result[:4])
    seed = [i for i in code]
    items = [i for i in itertools.product(seed, repeat=4)]
    num = items.count(result)
    return num


# 前四/组选/组选12
@bind("x4f.cb.g12")
def x4f_cb_g12(result, code):
    num = 0
    result = result[:4]
    for p in code[0]:
        if result.count(p)>=2:
            ex = [i for i in result if i!=p]
            seed = [i for i in code[1]]
            items = [i for i in itertools.product(seed, repeat=2)]
            for item in items:
                if list(item) == ex:
                    num = 1
    return num


# 前四/组选/组选6
@bind("x4f.cb.g6")
def x4f_cb_g6(result, code):
    num = 0
    result = tuple(set([i for i in result[:4]]))
    seed = "".join(code)
    items = [i for i in itertools.combinations(seed, 2)]
    for item in items:
        if item == result:
            num += 1
    return num


# 前四/组选/组选4
@bind("x4f.cb.g4")
def x4f_cb_g4(result, code):
    num = 0
    result = result[:4]
    for p in code[0]:
        if result.count(p)>=3:
            ex = [i for i in result if i!=p]
            for item in code[1]:
                if item in ex:
                    num = 1
    return num


# 前四/不定位/一码不定位
@bind("x4f.unf.m1")
def x4f_unf_m1(result, code):
    num = 0
    result = result[:4]
    if len([i for i in code if i in result]) >= 2:
        num =1
    return num


# 前四/不定位/二码不定位
@bind("x4f.unf.m2")
def x4f_unf_m2(result, code):
    num = 0
    result = result[:4]
    if len([i for i in code if i in result]) >= 3:
        num =1
    return num

# 后四/直选/复式
@bind("x4b.eq.batch")
def x4b_eq_batch(result, code):
    num = 0
    result = tuple(result[1:])
    seed = [[i for i in n] for n in code]
    items = [i for i in itertools.product(*seed)]
    num = items.count(result)
    return num


# 后四/直选/单式
@bind("x4b.eq.simple")
def x4b_eq_simple(result, code):
    num = 0
    result = result[1:]
    item = "".join([i for i in result])
    num = code.count(item)
    return num


# 后四/组选/组选24
@bind("x4b.cb.g24")
def x4b_cb_g24(result, code):
    num = 0
    result = tuple(result[1:])
    seed = [i for i in code]
    items = [i for i in itertools.product(seed, repeat=4)]
    num = items.count(result)
    return num


# 后四/组选/组选12
@bind("x4b.cb.g12")
def x4b_cb_g12(result, code):
    num = 0
    result = result[1:]
    for p in code[0]:
        if result.count(p)>=2:
            ex = [i for i in result if i!=p]
            seed = [i for i in code[1]]
            items = [i for i in itertools.product(seed, repeat=2)]
            for item in items:
                if list(item) == ex:
                    num = 1
    return num


# 后四/组选/组选6
@bind("x4b.cb.g6")
def x4b_cb_g6(result, code):
    num = 0
    result = tuple(set([i for i in result[1:]]))
    seed = "".join(code)
    items = [i for i in itertools.combinations(seed, 2)]
    for item in items:
        if item == result:
            num += 1
    return num


# 后四/组选/组选4
@bind("x4b.cb.g4")
def x4b_cb_g4(result, code):
    num = 0
    result = result[1:]
    for p in code[0]:
        if result.count(p)>=3:
            ex = [i for i in result if i!=p]
            for item in code[1]:
                if item in ex:
                    num = 1
    return num


# 后四/不定位/一码不定位
@bind("x4b.unf.m1")
def x4b_unf_m1(result, code):
    num = 0
    result = result[1:]
    if len([i for i in code if i in result]) >= 2:
        num =1
    return num


# 后四/不定位/二码不定位
@bind("x4b.unf.m2")
def x4b_unf_m2(result, code):
    num = 0
    result = result[1:]
    if len([i for i in code if i in result]) >= 3:
        num =1
    return num


# 前三/直选/复式
@bind("x3f.eq.batch")
def x3f_eq_batch(result, code):
    num = 0
    result = result[:3]
    item = "".join([i for i in result])
    num = code.count(item)
    return num


# 前三/直选/单式
@bind("x3f.eq.simple")
def x3f_eq_simple(result, code):
    num = 0
    result = result[:3]
    item = "".join([i for i in result])
    num = code.count(item)
    return num


# 前三/直选/和值
@bind("x3f.eq.sum")
def x3f_eq_sum(result, code):
    num = 0
    sum_ = sum([int(i) for i in result[:3]])
    for n in code:
        if int(n) == sum_:
            num = 1
    return num


# 前三/直选/跨度
@bind("x3f.eq.diff")
def x3f_eq_diff(result, code):
    num = 0
    result = [int(i) for i in result[:3]]
    diff = max(result) - min(result)
    for n in code:
        if int(n) == diff:
            num = 1
    return num


# 前三/组选/组选复式
@bind("x3f.cb.gx")
def x3f_cb_gx(result, code):
    num = 0
    result = result[0:3]
    result.sort()
    seed = [[i for i in n] for n in code]
    items = [i for i in itertools.product(*seed)]
    for item in items:
        m = list(item)
        m.sort()
        if m == result:
            num += 1
    return num


# 前三/组选/和值
@bind("x3f.cb.sum")
def x3f_cb_sum(result, code):
    num = 0
    result = list(set(result[:3]))
    if len(result) >= 2:
        sum_ = sum([int(i) for i in result])
        for n in code:
            if int(n) == sum_:
                num = 1
    return num


# 前三/组选/组三复式
@bind("x3f.cb.g3")
def x3f_cb_g3(result, code):
    num = 0
    result = list(set(result[:3]))
    if len(result) == 2:
        items = [i for i in code if i in result]
        if len(items) == 2:
            num = 1
    return num


# 前三/组选/组六复式
@bind("x3f.cb.g6")
def x3f_cb_g6(result, code):
    num = 0
    result = list(set(result[:3]))
    if len(result) == 3:
        items = [i for i in code if i in result]
        if len(items) == 3:
            num = 1
    return num


# 前三/组选/混合
@bind("x3f.cb.mix")
def x3f_cb_mix(result, code):
    num = 0
    result = result[:3]
    result.sort()
    code.sort()
    if code == result:
        num += 1
    return num


# 前三/组选/包胆
@bind("x3f.cb.with")
def x3f_cb_with(result, code):
    num = 0
    result = result[:3]
    if len(result) >= 2:
        for n in code:
            if result.count(n):
                num += 1
    return num


# 前三/不定位/一码不定位
@bind("x3f.unf.m1")
def x3f_unf_m1(result, code):
    num = 0
    result = result[:3]
    if len([i for i in code if i in result]) >= 1:
        num = 1
    return num


# 前三/不定位/二码不定位
@bind("x3f.unf.m2")
def x3f_unf_m2(result, code):
    num = 0
    result = result[:3]
    if len([i for i in code if i in result]) >= 2:
        num = 1
    return num


# 中三/直选/复式
@bind("x3m.eq.batch")
def x3m_eq_batch(result, code):
    num = 0
    result = result[1:-1]
    item = "".join([i for i in result])
    num = code.count(item)
    return num


# 中三/直选/单式
@bind("x3m.eq.simple")
def x3m_eq_simple(result, code):
    num = 0
    result = result[1:-1]
    item = "".join([i for i in result])
    num = code.count(item)
    return num


# 中三/直选/和值
@bind("x3m.eq.sum")
def x3m_eq_sum(result, code):
    num = 0
    sum_ = sum([int(i) for i in result[1:-1]])
    for n in code:
        if int(n) == sum_:
            num = 1
    return num


# 中三/直选/跨度
@bind("x3m.eq.diff")
def x3m_eq_diff(result, code):
    num = 0
    result = [int(i) for i in result[1:-1]]
    diff = max(result) - min(result)
    for n in code:
        if int(n) == diff:
            num = 1
    return num


# 中三/组选/组选复式
@bind("x3m.cb.gx")
def x3m_cb_gx(result, code):
    num = 0
    result = result[1:-1]
    result.sort()
    seed = [[i for i in n] for n in code]
    items = [i for i in itertools.product(*seed)]
    for item in items:
        m = list(item)
        m.sort()
        if m == result:
            num += 1
    return num


# 中三/组选/和值
@bind("x3m.cb.sum")
def x3m_cb_sum(result, code):
    num = 0
    result = list(set(result[1:-1]))
    if len(result) >= 2:
        sum_ = sum([int(i) for i in result])
        for n in code:
            if int(n) == sum_:
                num = 1
    return num


# 中三/组选/组三复式
@bind("x3m.cb.g3")
def x3m_cb_g3(result, code):
    num = 0
    result = list(set(result[1:-1]))
    if len(result) == 2:
        items = [i for i in code if i in result]
        if len(items) == 2:
            num = 1
    return num


# 中三/组选/组六复式
@bind("x3m.cb.g6")
def x3m_cb_g6(result, code):
    num = 0
    result = list(set(result[1:-1]))
    if len(result) == 3:
        items = [i for i in code if i in result]
        if len(items) == 3:
            num = 1
    return num


# 中三/组选/混合
@bind("x3m.cb.mix")
def x3m_cb_mix(result, code):
    num = 0
    result = result[1:-1]
    result.sort()
    code.sort()
    if code == result:
        num += 1
    return num


# 中三/组选/包胆
@bind("x3m.cb.with")
def x3m_cb_with(result, code):
    num = 0
    result = result[1:-1]
    if len(result) >= 2:
        for n in code:
            if result.count(n):
                num += 1
    return num


# 中三/不定位/一码不定位
@bind("x3m.unf.m1")
def x3m_unf_m1(result, code):
    num = 0
    result = result[1:-1]
    if len([i for i in code if i in result]) >= 1:
        num = 1
    return num


# 中三/不定位/二码不定位
@bind("x3m.unf.m2")
def x3m_unf_m2(result, code):
    num = 0
    result = result[1:-1]
    if len([i for i in code if i in result]) >= 2:
        num = 1
    return num


# 后三/直选/复式
@bind("x3b.eq.batch")
def x3b_eq_batch(result, code):
    num = 0
    result = tuple(result[2:])
    seed = [[i for i in n] for n in code]
    items = [i for i in itertools.product(*seed)]
    num = items.count(result)
    return num


# 后三/直选/单式
@bind("x3b.eq.simple")
def x3b_eq_simple(result, code):
    num = 0
    result = result[2:]
    item = "".join([i for i in result])
    num = code.count(item)
    return num


# 后三/直选/和值
@bind("x3b.eq.sum")
def x3b_eq_sum(result, code):
    num = 0
    sum_ = sum([int(i) for i in result[2:]])
    for n in code:
        if int(n) == sum_:
            num = 1
    return num



# 后三/直选/跨度
@bind("x3b.eq.diff")
def x3b_eq_diff(result, code):
    num = 0
    result = [int(i) for i in result[2:]]
    diff = max(result) - min(result)
    for n in code:
        if int(n) == diff:
            num = 1
    return num


# 后三/组选/组选复式
@bind("x3b.cb.gx")
def x3b_cb_gx(result, code):
    num = 0
    result = result[2:]
    result.sort()
    seed = [[i for i in n] for n in code]
    items = [i for i in itertools.product(*seed)]
    for item in items:
        m = list(item)
        m.sort()
        if m == result:
            num += 1
    return num


# 后三/组选/和值
@bind("x3b.cb.sum")
def x3b_cb_sum(result, code):
    num = 0
    result = list(set(result[2:]))
    if len(result) >= 2:
        sum_ = sum([int(i) for i in result])
        for n in code:
            if int(n) == sum_:
                num = 1
    return num


# 后三/组选/组三复式
@bind("x3b.cb.g3")
def x3b_cb_g3(result, code):
    num = 0
    result = list(set(result[2:]))
    if len(result) == 2:
        items = [i for i in code if i in result]
        if len(items) == 2:
            num = 1
    return num


# 后三/组选/组六复式
@bind("x3b.cb.g6")
def x3b_cb_g6(result, code):
    num = 0
    result = list(set(result[2:]))
    if len(result) == 3:
        items = [i for i in code if i in result]
        if len(items) == 3:
            num = 1
    return num


# 后三/组选/混合
@bind("x3b.cb.mix")
def x3b_cb_mix(result, code):
    num = 0
    result = result[2:]
    result.sort()
    code.sort()
    if code == result:
        num += 1
    return num


# 后三/组选/包胆
@bind("x3b.cb.with")
def x3b_cb_with(result, code):
    num = 0
    result = result[2:]
    if len(result) >= 2:
        for n in code:
            if result.count(n):
                num += 1
    return num


# 后三/不定位/一码不定位
@bind("x3b.unf.m1")
def x3b_unf_m1(result, code):
    num = 0
    result = result[2:]
    if len([i for i in code if i in result]) >= 1:
        num = 1
    return num


# 后三/不定位/二码不定位
@bind("x3b.unf.m2")
def x3b_unf_m2(result, code):
    num = 0
    result = result[2:]
    if len([i for i in code if i in result]) >= 2:
        num = 1
    return num


# 前二/直选/复式
@bind("x2f.eq.batch")
def x2f_eq_batch(result, code):
    num = 0
    result = tuple(result[:2])
    seed = [[i for i in n] for n in code]
    items = [i for i in itertools.product(*seed)]
    num = items.count(result)
    return num


# 前二/直选/单式
@bind("x2f.eq.simple")
def x2f_eq_simple(result, code):
    num = 0
    result = result[0:2]
    item = "".join([i for i in result])
    num = code.count(item)
    return num


# 前二/直选/和值
@bind("x2f.eq.sum")
def x2f_eq_sum(result, code):
    num = 0
    sum_ = sum([int(i) for i in result[:2]])
    for n in code:
        if int(n) == sum_:
            num = 1
    return num


# 前二/直选/跨度
@bind("x2f.eq.diff")
def x2f_eq_diff(result, code):
    num = 0
    result = [int(i) for i in result[:2]]
    diff = max(result) - min(result)
    for n in code:
        if int(n) == diff:
            num = 1
    return num


# 前二/组选/复式
@bind("x2f.cb.batch")
def x2f_cb_batch(result, code):
    num = 0
    result = [i for i in result[:2]]
    seed = "".join(code)
    items = [item for item in itertools.combinations(seed, 2)]
    for item in items:
        if list(item) == result:
            num = 1
    return num


# 前二/组选/单式
@bind("x2f.cb.simple")
def x2f_cb_simple(result, code):
    num = 0
    result = result[:2]
    for n in code:
        items = [i for i in n if i in result]
        if len(items) == 2:
            num += 1
    return num


# 前二/组选/和值
@bind("x2f.cb.sum")
def x2f_cb_sum(result, code):
    num = 0
    result = result[:2]
    if len(set(result)) == 2:
        sum_ = sum([int(i) for i in result])
        for n in code:
            if int(n) == sum_:
                num = 1
    return num


# 前二/组选/包胆
@bind("x2f.cb.with")
def x2f_cb_with(result, code):
    num = 0
    result = result[:2]
    if len(set(result)) == 2:
        for n in code:
            if n in result:
                num = 1
    return num


# 后二/直选/复式
@bind("x2b.eq.batch")
def x2b_eq_batch(result, code):
    num = 0
    result = tuple(result[-2:])
    seed = [[i for i in n] for n in code]
    items = [i for i in itertools.product(*seed)]
    num = items.count(result)
    return num


# 后二/直选/单式
@bind("x2b.eq.simple")
def x2b_eq_simple(result, code):
    num = 0
    result = result[0:2]
    item = "".join([i for i in result])
    num = code.count(item)
    return num


# 后二/直选/和值
@bind("x2b.eq.sum")
def x2b_eq_sum(result, code):
    num = 0
    sum_ = sum([int(i) for i in result[-2:]])
    for n in code:
        if int(n) == sum_:
            num = 1
    return num


# 后二/直选/跨度
@bind("x2b.eq.diff")
def x2b_eq_diff(result, code):
    num = 0
    result = [int(i) for i in result[-2:]]
    diff = max(result) - min(result)
    for n in code:
        if int(n) == diff:
            num = 1
    return num


# 后二/组选/复式
@bind("x2b.cb.batch")
def x2b_cb_batch(result, code):
    num = 0
    result = [i for i in result[-2:]]
    seed = "".join(code)
    items = [item for item in itertools.combinations(seed, 2)]
    for item in items:
        if list(item) == result:
            num = 1
    return num


# 后二/组选/单式
@bind("x2b.cb.simple")
def x2b_cb_simple(result, code):
    num = 0
    result = result[-2:]
    for n in code:
        items = [i for i in n if i in result]
        if len(items) == 2:
            num += 1
    return num


# 后二/组选/和值
@bind("x2b.cb.sum")
def x2b_cb_sum(result, code):
    num = 0
    result = result[-2:]
    if len(set(result)) == 2:
        sum_ = sum([int(i) for i in result])
        for n in code:
            if int(n) == sum_:
                num = 1
    return num


# 后二/组选/包胆
@bind("x2b.cb.with")
def x2b_cb_with(result, code):
    num = 0
    result = result[-2:]
    if len(set(result)) == 2:
        for n in code:
            if n in result:
                num = 1
    return num


# 一星/定位胆/复式
@bind("x1.fix.batch")
def x1_fix_batch(result, code):
    num = 0
    for i in range(len(result)):
        if result[i] in code[i]:
            num += 1
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

    x5_eq_batch(["0", "1", "2", "3", "4"], ["0123456789", "0123456789", "0123456789", "0123456789", "0123456789"])
    x5_eq_simple(["0", "1", "2", "3", "4"],["01234", "56789", "56790"])
    