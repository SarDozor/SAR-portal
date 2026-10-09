import rasterio
from rasterio.warp import transform
import numpy as np

result = 'S1AA_20170118T232041_20170211T232040_VVP024_INT80_G_weF_DDAC'
latit = [-29.766708]
longit = [-69.976503]
vert_disp = f'{result}/{result}_vert_disp.tif'
with rasterio.open(vert_disp) as vd: #открытие файла вертикальных перемещений
    print(vd.crs) #выдает код системы координат utm
    def_crs = 'EPSG:4326' #обычная система координат wgs84
    new_crs = vd.crs
    trns_x, trns_y = transform(def_crs, new_crs, longit, latit) #трансформирует систему координат из обычной в utm
    print(trns_x, trns_y)
    x = trns_x[0]
    y = trns_y[0]
    numb_string, numb_column = vd.index(x, y)
    print(f"№ столбца: {numb_column} | № строки: {numb_string}")
    print(f"кол-во столбцов: {vd.width} | кол-во строк: {vd.height}")
    vd_table = vd.read(1) #весь слой данных. 1 в скобках тк всего слой 1
    #print(vd_table)
    numb_table = vd_table[numb_string, numb_column]
    print(f"смещение {numb_table}") #выдает 5 см смещения за 24 дня, но не факт что реальное оседание
    ots = 40 #отступ пикселей от дамбы для квадрата нахождения опорной точки
    begin_string = numb_string - ots
    end_string = numb_string + ots + 1 
    begin_clmn = numb_column - ots
    end_clmn = numb_column + ots + 1 
    srez = vd_table[begin_string:end_string, begin_clmn:end_clmn] #делает срез таблицы 
    print(f" размер таблицы среза: {srez.shape}") #выдает 81 на 81 что правильно, срез получился 
    print(f"srez in the center: {srez[ots, ots]}")
    na_replaced = srez[srez != 0]
    median = np.median(na_replaced)
    print(f"median: {median}") # медиана окрестностей дамбы = -5,26 см, а сама дамба = -5,25 -> разница всего 0,13 мм, что почти равно нулю, следовательно 5 см смещения это не оседание дамбы, а вклад атмосферы, получается, что дамба не двигалась за эти 24 дня
    print(f"deafult - median: {numb_table - median}") # разница всего 0,13 мм, скорее всего обычная погрешность. дамба абсолютно не двигалась за эти 24 дня. можем сделать вывод, что сырым данным InSAR нельзя верить напрямую
    print(f"default table size: {vd_table.size}")
    print(f"filtred data table size: {na_replaced.size}")
coger = f'{result}/{result}_corr.tif'
with rasterio.open(coger) as cg: #открытие файла когерентности
    print(f"кол-во столбцов: {cg.width} | кол-во строк: {cg.height}") # -> сетка сходтися с _vert_disp.tif
    cg_table = cg.read(1)
    print(cg_table)
    cg_result = cg_table[numb_string, numb_column]
    print(f"когерентность: {cg_result}") # когерентность 0.8, что значит что поверхность между взятыми снимками не изменнилась и лишнего шума нет -> измерение чистое
