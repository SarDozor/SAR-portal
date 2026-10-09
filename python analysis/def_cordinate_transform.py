import rasterio
from rasterio.warp import transform
import numpy as np

def damb_move(result, longit, latit, ots):
    vert_disp = f'{result}/{result}_vert_disp.tif'
    coger = f'{result}/{result}_corr.tif'
    with rasterio.open(vert_disp) as vd: #открытие файла вертикальных перемещений
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
displacement_mm, coherence = damb_move('S1AA_20170118T232041_20170211T232040_VVP024_INT80_G_weF_DDAC', -69.976503, -29.766708, 40)
print(f"смещение в мм: {displacement_mm:.2f}, когерентность: {coherence:.2f}")
