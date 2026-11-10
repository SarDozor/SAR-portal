import os
os.environ.pop('PROJ_LIB', None)
os.environ.pop('PROJ_DATA', None)
import rasterio
from rasterio.warp import transform
import numpy as np


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
#print(mean.shape)
#print(np.median(mean[40, 40]))
#print(np.max(mean), np.min(mean))
mask = mean > 0.8
mask[30:51, 30:51] = False
foldername = 'S1BB_20180208T082857_20180220T082857_VVP012_INT40_G_weF_3093'
dem = f'{foldername}/{foldername}/{foldername}_dem.tif'
dopusk = [100, 75, 50, 30]
#for z in dopusk:
with rasterio.open(dem) as d:
#print(d.shape)
    trns_x, trns_y = transform('EPSG:4326', d.crs, [lon], [lat])
    f_ind_x, f_ind_y = d.index(trns_x[0], trns_y[0])
    f_table = d.read(1)
    ots = 40
    numb_table = f_table[f_ind_x, f_ind_y]
    begin_string = f_ind_x - ots
    end_string = f_ind_x + ots + 1 
    begin_clmn = f_ind_y - ots
    end_clmn = f_ind_y + ots + 1
    dem_srez = f_table[begin_string:end_string, begin_clmn:end_clmn]
    mask2 = dem_srez - dem_srez[40, 40]
    mask2 = np.abs(mask2) <= 100
    mask_z = mask & mask2
    print(np.sum(mask_z))
coords = np.argwhere(mask)
random21 = np.random.choice(len(coords), size=21, replace=False)
index21 = coords[random21]
new_mask = np.zeros_like(mask)
new_mask[index21[:, 0], index21[:, 1]] = True
#print(np.sum(new_mask))
dp = 0
intervals = []
disp_windows = []
results = []
height_intervals = []
for n in list2:
    file = f'B1_series/{n}/{n}_los_disp.tif'
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
        ref_vals = srez[(new_mask) & (srez != 0)]
        disp_windows.append(srez.copy())
        opora = np.median(ref_vals)
        damba = window_mean(f_table, f_ind_x, f_ind_y, 1)
        dps = (damba - opora) * 1000
        dp += dps
        intervals.append(dps)
        data = n[21:29]
        #print(f"DATA: {data}, displacement za interval: {dps:.2f}, nakoplennoe: {dp:.2f}")
    #print(f'{np.std(intervals):.2f}')
#print(len(disp_windows))
#print(disp_windows[5].shape)
print("исходных пикселей:", np.sum(mask))
print("вырезок:", len(disp_windows))
np.random.seed(42)
for i in range(100):
    coords = np.argwhere(mask)
    random21 = np.random.choice(len(coords), size=120, replace=False)
    if i == 0:
        print("первые 10 выбранных индексов:", random21[:10])
    index21 = coords[random21]
    new_mask = np.zeros_like(mask)
    new_mask[index21[:, 0], index21[:, 1]] = True
    intervals = []
    for srez in disp_windows:
        ref_vals = srez[(new_mask) & (srez != 0)]
        opora = np.median(ref_vals)
        damba = window_mean(srez, 40, 40, 1)
        dps = (damba - opora) * 1000
        intervals.append(dps)
    results.append(np.std(intervals))
print(results)
height_intervals = []
for srez in disp_windows:
    ref_vals = srez[mask_z & (srez != 0)]
    opora = np.median(ref_vals)
    damba = window_mean(srez, 40, 40, 1)
    height_intervals.append((damba - opora) * 1000)
height_std = np.std(height_intervals)
print("разброс с фильтром высот, мм:", height_std)
def window_grid(path):
    with rasterio.open(path) as src:
        xs, ys = transform('EPSG:4326', src.crs, [lon], [lat])
        row, col = src.index(xs[0], ys[0])
        points = np.array([
            src.xy(row - 40, col - 40),
            src.xy(row - 40, col + 40),
            src.xy(row + 40, col - 40),
        ]) #центры трёх угловых пикселей вырезки 81×81
        return src.crs, points
dem_crs, dem_points = window_grid(dem)
matched = 0
checked = 0
for n in list2:
    for kind in ('corr', 'los_disp'):
        path = f'B1_series/{n}/{n}_{kind}.tif'
        crs, points = window_grid(path)
        checked += 1

        same_grid = (
            crs == dem_crs
            and np.allclose(points, dem_points, rtol=0, atol=1e-6)
        )

        if same_grid:
            matched += 1
        else:
            print("несовпадение сетки:", n, kind)
            if crs == dem_crs:
                print("разности координат углов:", points - dem_points)
            else:
                print("разные CRS:", crs, dem_crs)
print(f"совпало с сеткой DEM: {matched}/{checked}")



'''#check = f'B1_series/{list2[0]}/{list2[0]}_los_disp.tif'
with rasterio.open(check) as c:
    print(c.shape)
dates = []
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