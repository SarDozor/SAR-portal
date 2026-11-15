import os
os.environ.pop('PROJ_LIB', None)
os.environ.pop('PROJ_DATA', None)
import rasterio
from rasterio.warp import transform
import numpy as np
import matplotlib.pyplot as plt
from datetime import datetime
import pandas as pd
#from pathlib import Path
from zipfile import ZipFile

object_code = 'B1'
series_folder = 'B1_series'
lon = -44.120996
lat = -20.119388
foldername = 'S1BB_20180208T082857_20180220T082857_VVP012_INT40_G_weF_3093'
dem = f'{foldername}/{foldername}/{foldername}_dem.tif'
output_prefix = 'b1'

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
lon_sup = -44.134049
lat_sup = -20.121520
list2 = []
for n in os.listdir(series_folder):
    if n.endswith('.zip'):
        pass
    else:
        list2.append(n)
list2 = sorted(list2)
summ = 0
dam_coherence_means = []
dam_coherence_counts = []
for n in list2:
    file = f'{series_folder}/{n}/{n}_corr.tif'
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
        dam_corr = srez[39:42, 39:42]
        valid_corr = dam_corr[np.isfinite(dam_corr) & (dam_corr != 0)]
        dam_coherence_counts.append(valid_corr.size)
        if valid_corr.size > 0:
            dam_coherence_means.append(np.mean(valid_corr))
        else:
            dam_coherence_means.append(np.nan)
        summ = summ + srez
mean = summ/len(list2)
#print(mean.shape)
#print(np.median(mean[40, 40]))
#print(np.max(mean), np.min(mean))
mask = mean > 0.8
mask[30:51, 30:51] = False
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
disp_windows = []
results = []
base_intervals = []
large_window_intervals = []
for n in list2:
    file = f'{series_folder}/{n}/{n}_los_disp.tif'
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
        disp_windows.append(srez.copy())
        damba = window_mean(srez, 40, 40, 1)
        base_vals = srez[mask & np.isfinite(srez) & (srez != 0)]
        base_ref = np.median(base_vals)
        base_intervals.append((damba - base_ref) * 1000)
        large_srez = f_table[f_ind_x - 80:f_ind_x + 81, f_ind_y - 80:f_ind_y + 81]
        if large_srez.shape != (161, 161):
            raise ValueError(f"Большое окно обрезано в файле {n}")
        large_vals = large_srez[np.isfinite(large_srez) & (large_srez != 0)]
        if large_vals.size == 0:
            raise ValueError(f"Большое окно не содержит данных: {n}")
        large_ref = np.median(large_vals)
        large_window_intervals.append((damba - large_ref) * 1000)
print("исходных пикселей:", np.sum(mask))
print("вырезок:", len(disp_windows))
np.random.seed(42)
for i in range(1000):
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
height_intervals = []
for srez in disp_windows:
    ref_vals = srez[mask_z & (srez != 0)]
    opora = np.median(ref_vals)
    damba = window_mean(srez, 40, 40, 1)
    height_intervals.append((damba - opora) * 1000)
height_std = np.std(height_intervals)
print("\nсравнение способов опоры:")
print(f"исходная маска: {np.std(base_intervals):.3f} мм")
print(f"фильтр ±100 м: {height_std:.3f} мм")
print(f"большое окно 161×161: {np.std(large_window_intervals):.3f} мм")
print("\n")
random_stds = np.array(results)
if not np.all(np.isfinite(random_stds)):
    raise ValueError("случайных результатах есть nan или бесконечность")
if not np.isfinite(height_std):
    raise ValueError("разброс фильтра высот некорректен")
not_worse = np.count_nonzero(random_stds <= height_std)
print("всего экспериментов:", len(random_stds))
print(f"минимум случайных результатов: {np.min(random_stds):.3f} мм")
print(f"медиана: {np.median(random_stds):.3f} мм")
print(f"максимум: {np.max(random_stds):.3f} мм")
print(f"результат фильтра высот: {height_std:.3f} мм")
print("случайных результатов не хуже фильтра:", not_worse)
print("разброс с фильтром высот, мм:", height_std)
def window_grid(path):
    with rasterio.open(path) as src:
        xs, ys = transform('EPSG:4326', src.crs, [lon], [lat])
        row, col = src.index(xs[0], ys[0])
        points = np.array([
            src.xy(row - 40, col - 40),
            src.xy(row - 40, col + 40),
            src.xy(row + 40, col - 40),]) #центры трёх угловых пикселей вырезки 81×81
        return src.crs, points
dem_crs, dem_points = window_grid(dem)
matched = 0
checked = 0
for n in list2:
    for kind in ('corr', 'los_disp'):
        path = f'{series_folder}/{n}/{n}_{kind}.tif'
        crs, points = window_grid(path)
        checked += 1
        same_grid = (crs == dem_crs and np.allclose(points, dem_points, rtol=0, atol=1e-6))
        if same_grid:
            matched += 1
        else:
            print("несовпадение сетки:", n, kind)
            if crs == dem_crs:
                print("разности координат углов:", points - dem_points)
            else:
                print("разные CRS:", crs, dem_crs)
print(f"совпало с сеткой DEM: {matched}/{checked}")
fig, ax = plt.subplots(figsize=(9, 5))
ax.hist(
    random_stds,
    bins=25,
    color="steelblue",
    edgecolor="white",
    label="случайный отбор: 120 пикселей",
)
ax.axvline(
    height_std,
    color="darkred",
    linewidth=2,
    label=f"фильтр ±100 м: {height_std:.3f} мм",
)
ax.set_xlabel("cтандартное отклонение интервальных смещений, мм")
ax.set_ylabel("количество случайных выборок")
ax.set_title("сравнение отбора по высоте с 1000 случайными выборками")
ax.legend()
fig.tight_layout()
fig.savefig("reference_selection_histogram.png", dpi=300)
for kind in ('corr', 'los_disp'):
    nodata_values = set()
    for n in list2:
        path = f'{series_folder}/{n}/{n}_{kind}.tif'
        with rasterio.open(path) as src:
            nodata_values.add(str(src.nodata))
    print(f"nodata для {kind}:", nodata_values)
with rasterio.open(dem) as src:
    print("nodata для dem:", src.nodata)
base_counts = []
height_counts = []
dam_counts = []
for srez in disp_windows:
    valid = np.isfinite(srez) & (srez != 0)
    base_counts.append(np.count_nonzero(mask & valid))
    height_counts.append(np.count_nonzero(mask_z & valid))
    dam_counts.append(np.count_nonzero(valid[39:42, 39:42]))
print("исходная опора: мин, макс:", min(base_counts), max(base_counts))
print("опора по высоте: мин, макс:", min(height_counts), max(height_counts))
print("дамба 3×3: мин, макс:", min(dam_counts), max(dam_counts))
dates = []
for n in list2:
    date_text = n[21:29]
    date = datetime.strptime(date_text, "%Y%m%d")
    dates.append(date)
fig_time, ax_time = plt.subplots(figsize=(11, 5))
ax_time.plot(dates, base_intervals, marker="o", label="исходная маска")
ax_time.plot(dates, height_intervals, marker="o", label="фильтр ±100 м")
ax_time.plot(dates, large_window_intervals, marker="o", label="окно 161×161")
ax_time.axhline(0, color="gray", linewidth=1)
ax_time.set_xlabel("дата второго снимка пары")
ax_time.set_ylabel("интервальное смещение LOS, мм")
ax_time.set_title("смещения дамбы при трёх способах выбора опоры")
ax_time.legend()
fig_time.autofmt_xdate()
fig_time.tight_layout()
fig_time.savefig("reference_methods_timeseries.png", dpi=300)
first_date = datetime.strptime(list2[0][5:13], "%Y%m%d")
cumulative_dates = [first_date] + dates
base_cumulative = np.cumsum([0] + base_intervals)
height_cumulative = np.cumsum([0] + height_intervals)
large_cumulative = np.cumsum([0] + large_window_intervals)
fig_sum, ax_sum = plt.subplots(figsize=(11, 5))
ax_sum.plot(cumulative_dates, base_cumulative, marker="o",
            label="исходная маска")
ax_sum.plot(cumulative_dates, height_cumulative, marker="o",
            label="фильтр ±100 м")
ax_sum.plot(cumulative_dates, large_cumulative, marker="o",
            label="окно 161×161")
ax_sum.axhline(0, color="gray", linewidth=1)
ax_sum.set_xlabel("дата")
ax_sum.set_ylabel("накопленное относительное смещение LOS, мм")
ax_sum.set_title("накопленные смещения относительно 03.01.2018")
ax_sum.legend()
fig_sum.autofmt_xdate()
fig_sum.tight_layout()
fig_sum.savefig("reference_methods_cumulative.png", dpi=300)
print("итоговое накопленное смещение:")
print(f"исходная маска: {base_cumulative[-1]:.3f} мм")
print(f"фильтр ±100 м: {height_cumulative[-1]:.3f} мм")
print(f"большое окно: {large_cumulative[-1]:.3f} мм")
difference = np.array(height_intervals) - np.array(base_intervals)
idx = np.argmax(np.abs(difference))
print("пара с максимальным расхождением:", list2[idx])
print(f"исходная маска: {base_intervals[idx]:.3f} мм")
print(f"фильтр высот: {height_intervals[idx]:.3f} мм")
print(f"большое окно: {large_window_intervals[idx]:.3f} мм")
print(f"разность фильтр − исходная маска: {difference[idx]:.3f} мм")
start_dates = []
end_dates = []
for n in list2:
    start = datetime.strptime(n[5:13], "%Y%m%d")
    end = datetime.strptime(n[21:29], "%Y%m%d")

    start_dates.append(start.strftime("%Y-%m-%d"))
    end_dates.append(end.strftime("%Y-%m-%d"))
export_table = pd.DataFrame({
    "object_code": object_code,
    "pair_name": list2,
    "reference_date": start_dates,
    "secondary_date": end_dates,
    "latitude": lat,
    "longitude": lon,
    "base_interval_mm": base_intervals,
    "height_interval_mm": height_intervals,
    "large_window_interval_mm": large_window_intervals,
    "base_cumulative_mm": base_cumulative[1:],
    "height_cumulative_mm": height_cumulative[1:],
    "large_window_cumulative_mm": large_cumulative[1:],
    "dam_coherence_mean": dam_coherence_means,
    "dam_coherence_valid_pixels": dam_coherence_counts,
})
export_table.to_csv(
    f'{output_prefix}_measurements.csv',
    index=False
)
print(export_table.head(3))
print("table size:", export_table.shape)
measurement_tables = []
methods = [
    ("base", "base_interval_mm"),
    ("height_100m", "height_interval_mm"),
    ("large_window_161", "large_window_interval_mm"),
]
for method_name, displacement_column in methods:
    method_table = export_table[[
        "object_code",
        "pair_name",
        "reference_date",
        "secondary_date",
        "latitude",
        "longitude",
        "dam_coherence_mean",
        "dam_coherence_valid_pixels",
    ]].copy()
    method_table["reference_method"] = method_name
    method_table["displacement_mm"] = export_table[displacement_column]
    measurement_tables.append(method_table)
measurements_long = pd.concat(measurement_tables, ignore_index=True)
granules_by_pair = {}
for pair_name in list2:
    archive_path = f"{series_folder}/{pair_name}.zip"
    member = f"{pair_name}/{pair_name}.txt"
    with ZipFile(archive_path) as archive:
        contents = archive.read(member).decode("utf-8-sig")
    fields = {}
    for line in contents.splitlines():
        name, _, value = line.partition(":")
        name = name.strip()
        if name in ("Reference Granule", "Secondary Granule"):
            fields[name] = value.strip()
    granules_by_pair[pair_name] = (
        fields["Reference Granule"],
        fields["Secondary Granule"],
    )
measurements_long["reference_granule"] = [
    granules_by_pair[name][0] for name in measurements_long["pair_name"]
]
measurements_long["secondary_granule"] = [
    granules_by_pair[name][1] for name in measurements_long["pair_name"]
]
measurements_long.to_csv(
    f'{output_prefix}_measurements_long.csv',
    index=False
)
print("размер таблицы для передачи:", measurements_long.shape)
print(measurements_long.head(3))
#plt.show()
print("first pair:", list2[0])
print("last pair:", list2[-1])
previous_end = None
breaks = 0
durations = []
for n in list2:
    start = datetime.strptime(n[5:13], "%Y%m%d")
    end = datetime.strptime(n[21:29], "%Y%m%d")
    durations.append((end - start).days)
    if previous_end is not None and start != previous_end:
        print("стык не совпадает:", previous_end.date(), start.date())
        breaks += 1
    previous_end = end
print("длительности пар, дней:", sorted(set(durations)))
print("несовпадающих стыков:", breaks)











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