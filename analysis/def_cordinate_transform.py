import rasterio
from rasterio.warp import transform
import numpy as np
import os

def window_mean(table, numb_string, numb_column, ots):
    begin_string1 = numb_string - ots
    end_string1 = numb_string + ots + 1 
    begin_clmn1 = numb_column - ots
    end_clmn1 = numb_column + ots + 1
    srez1 = table[begin_string1:end_string1, begin_clmn1:end_clmn1]
    na_replaced1 = srez1[srez1 != 0]
    mean = np.mean(na_replaced1)
    return mean
def damb_move(result, longit, latit, ots, longit_sup, latit_sup, disp_type='vert_disp', base='.'):
    disp = f'{base}/{result}/{result}_{disp_type}.tif'
    coger = f'{base}/{result}/{result}_corr.tif'
    with rasterio.open(disp) as vd: #открытие файла вертикальных перемещений
        def_crs = 'EPSG:4326' #обычная система координат wgs84
        new_crs = vd.crs
        trns_x, trns_y = transform(def_crs, new_crs, [longit, longit_sup], [latit, latit_sup]) #трансформирует систему координат из обычной в utm
        x = trns_x[0]
        y = trns_y[0]
        x1 = trns_x[1]
        y1 = trns_y[1]
        numb_string, numb_column = vd.index(x, y)
        numb_string_sup, numb_column_sup = vd.index(x1, y1)
        vd_table = vd.read(1) #весь слой данных. 1 в скобках тк всего слой 1
        #numb_table = vd_table[numb_string, numb_column]
        begin_string = numb_string - ots
        end_string = numb_string + ots + 1 
        begin_clmn = numb_column - ots
        end_clmn = numb_column + ots + 1
        '''begin_string1 = numb_string - dam_ots
        end_string1 = numb_string + dam_ots + 1 
        begin_clmn1 = numb_column - dam_ots
        end_clmn1 = numb_column + dam_ots + 1'''
        srez = vd_table[begin_string:end_string, begin_clmn:end_clmn] #делает срез таблицы
        '''srez1 = vd_table[begin_string1:end_string1, begin_clmn1:end_clmn1]'''
        na_replaced = srez[srez != 0]
        '''na_replaced1 = srez1[srez1 != 0]'''
        median = np.median(na_replaced)
        '''mean = np.mean(na_replaced1)''' #добавив уже после ряда б1 среднее по квадрату пикселей, мы убедились, что шум создавался из за атмосферы. нужно искать опорную точку, которая не двигается
        moved = (window_mean(vd_table, numb_string, numb_column, 1) - median)*1000 #cмещение, будет выдано в мм
        moved_new = (window_mean(vd_table, numb_string, numb_column, 1) - window_mean(vd_table, numb_string_sup, numb_column_sup, 1))*1000
    with rasterio.open(coger) as cg: #открытие файла когерентности
        cg_table = cg.read(1)
        cg_result = cg_table[numb_string, numb_column] #выдает когерентность, т.е. проверку на то чисто ли измерение
        cg_result_new = cg_table[numb_string_sup, numb_column_sup]
    return moved, moved_new, cg_result, cg_result_new
'''displacement_mm, coherence = damb_move('S1BB_20190110T082905_20190122T082904_VVP012_INT80_G_weF_3C87', -44.120996, -20.119388, 40, disp_type='los_disp')
print(f"смещение в мм: {displacement_mm:.2f}, когерентность: {coherence:.2f}")'''
lon = -44.120996
lat = -20.119388
lon_sup = -44.134049
lat_sup = -20.121520
list2 = []
for n in os.listdir('B1_series'):
    if n.endswith('.zip'):
        pass
    else:
        list2.append(n)
list2 = sorted(list2)
summ = 0
for n in list2:
    file = f'B1_series/{n}/{n}_corr.tif'
    with rasterio.open(file) as f:
        trns_x, trns_y = transform('EPSG:4326', f.crs, [lon], [lat])
        f_ind_x, f_ind_y = f.index(trns_x[0], trns_y[0])
        f_table = f.read(1)
        ots = 40
        numb_table = f_table[f_ind_x, f_ind_y]
        begin_string = f_ind_x - ots
        end_string = f_ind_x + ots + 1 
        begin_clmn = f_ind_y - ots
        end_clmn = f_ind_y + ots + 1
        srez = f_table[begin_string:end_string, begin_clmn:end_clmn]
        summ = summ + srez
mean = summ/len(list2)
print(mean.shape)
print(np.median(mean[40:40]))
print(np.max(mean), np.min(mean))
'''dates = []
disps = []
cogers = []
disps_new = []
cogers_new = []
dp = 0
dp_new = 0
dps = []
dps_new = []
exmple = list2[0]
print(exmple[21:29])
for n in list2:
    displacement, displacement_new, coherence, coherence_new = damb_move(n, lon, lat, 80, lon_sup, lat_sup, disp_type='los_disp', base='B1_series')
    date2 = n[21:29]
    dp += displacement
    dp_new += displacement_new
    dates.append(date2)
    disps.append(displacement)
    cogers.append(coherence)
    dps.append(dp)
    disps_new.append(displacement_new)
    cogers_new.append(coherence_new)
    dps_new.append(dp_new)
    print("-----------------------------------------------------------------------------------------------------------------------------")
    print(f"дата: {date2}, смещение за этот интервал: {displacement:.2f}, когерентность: {coherence:.2f}, накопленное смещение: {dp:.2f}")
    print(f"НОВЫЕ: дата: {date2}, смещение за этот интервал: {displacement_new:.2f}, когерентность: {coherence_new:.2f}, накопленное смещение: {dp_new:.2f}")
print("END")'''