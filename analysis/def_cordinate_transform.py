import rasterio
from rasterio.warp import transform
import numpy as np
import os

def damb_move(result, longit, latit, ots, disp_type='vert_disp', base='.'):
    disp = f'{base}/{result}/{result}_{disp_type}.tif'
    coger = f'{base}/{result}/{result}_corr.tif'
    with rasterio.open(disp) as vd: #открытие файла вертикальных перемещений
        def_crs = 'EPSG:4326' #обычная система координат wgs84
        new_crs = vd.crs
        trns_x, trns_y = transform(def_crs, new_crs, [longit], [latit]) #трансформирует систему координат из обычной в utm
        x = trns_x[0]
        y = trns_y[0]
        numb_string, numb_column = vd.index(x, y)
        vd_table = vd.read(1) #весь слой данных. 1 в скобках тк всего слой 1
        numb_table = vd_table[numb_string, numb_column]
        begin_string = numb_string - ots
        end_string = numb_string + ots + 1 
        begin_clmn = numb_column - ots
        end_clmn = numb_column + ots + 1 
        srez = vd_table[begin_string:end_string, begin_clmn:end_clmn] #делает срез таблицы 
        na_replaced = srez[srez != 0]
        median = np.median(na_replaced)
        moved = (numb_table - median)*1000 #cмещение, будет выдано в мм
    with rasterio.open(coger) as cg: #открытие файла когерентности
        cg_table = cg.read(1)
        cg_result = cg_table[numb_string, numb_column] #выдает когерентность, т.е. проверку на то чисто ли измерение
    return moved, cg_result
'''displacement_mm, coherence = damb_move('S1BB_20190110T082905_20190122T082904_VVP012_INT80_G_weF_3C87', -44.120996, -20.119388, 40, disp_type='los_disp')
print(f"смещение в мм: {displacement_mm:.2f}, когерентность: {coherence:.2f}")'''
lon = -44.120996
lat = -20.119388
list2 = []
for n in os.listdir('B1_series'):
    if n.endswith('.zip'):
        pass
    else:
        list2.append(n)
list2 = sorted(list2)
dates = []
disps = []
cogers = []
dp = 0
dps = []
'''exmple = list2[0]
print(exmple[21:29])'''
for n in list2:
    displacement, coherence = damb_move(n, lon, lat, 80, disp_type='los_disp', base='B1_series')
    date2 = n[21:29]
    dp += displacement
    dates.append(date2)
    disps.append(displacement)
    cogers.append(coherence)
    dps.append(dp)
    print(f"дата: {date2}, смещение за этот интервал: {displacement:.2f}, когерентность: {coherence:.2f}, накопленное смещение: {dp:.2f}")