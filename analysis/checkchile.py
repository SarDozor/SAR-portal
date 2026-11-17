import os
os.environ.pop('PROJ_LIB', None)
os.environ.pop('PROJ_DATA', None)
from zipfile import ZipFile
import numpy as np
import rasterio
from rasterio.warp import transform

product = 'S1AA_20250127T100316_20250208T100315_VVP012_INT40_G_weF_E083'
archive_path = f'{product}.zip'
lon = -69.976503
lat = -29.766708

with ZipFile(archive_path) as archive:
    for name in archive.namelist():
        if name.endswith(('_corr.tif', '_los_disp.tif', '_dem.tif')):
            print('извлекаем:', name)
            archive.extract(name, 'chile_test')
windows = {}
grids = {}
for kind in ('corr', 'los_disp', 'dem'):
    path = f'chile_test/{product}/{product}_{kind}.tif'
    with rasterio.open(path) as src:
        xs, ys = transform('EPSG:4326', src.crs, [lon], [lat])
        row, col = src.index(xs[0], ys[0])
        print(f'\ntipe of file: {kind}')
        print('coords:', src.crs)
        print('raster size:', src.shape)
        print('string and column of dot:', row, col)
        if row < 40 or col < 40 or row + 40 >= src.height or col + 40 >= src.width:
            print('точка за пределами растра или окно 81×81 не помещается.')
            continue
        table = src.read(1, masked=True)
        window = table[row - 40:row + 41, col - 40:col + 41]
        windows[kind] = window.copy()
        grids[kind] = (src.crs, src.transform, src.shape, row, col)
        valid = ~np.ma.getmaskarray(window) & np.isfinite(window.data)
        print('размер вырезки:', window.shape)
        print('значение в центре:', window[40, 40])
        print('корректных пикселей в окне:', np.sum(valid))
        if kind == 'corr':
            candidates = valid & (window.data > 0.8)
            candidates[30:51, 30:51] = False
            print('пикселей с когерентностью > 0.8 вне центра 21×21:', np.sum(candidates))
if len(windows) == 3:
    if grids['corr'] == grids['los_disp'] == grids['dem']:
        print('\nсетки трёх файлов совпадают')
        corr = windows['corr']
        disp = windows['los_disp']
        dem = windows['dem']
        valid_corr = ~np.ma.getmaskarray(corr) & np.isfinite(corr.data)
        valid_disp = ~np.ma.getmaskarray(disp) & np.isfinite(disp.data)
        valid_dem = ~np.ma.getmaskarray(dem) & np.isfinite(dem.data)
        base_mask = valid_corr & valid_disp & (corr.data > 0.8)
        base_mask[30:51, 30:51] = False
        height_mask = (
            base_mask
            & valid_dem
            & (np.abs(dem.data - dem.data[40, 40]) <= 100)
        )
        dam_window = disp[39:42, 39:42]
        dam_values = dam_window.compressed()
        dam_values = dam_values[np.isfinite(dam_values)]
        if len(dam_values) == 0:
            print('нет корректных значений в окне дамбы 3×3')
        elif not valid_dem[40, 40]:
            print('нет корректной высоты в центре DEM')
        elif np.sum(base_mask) == 0 or np.sum(height_mask) == 0:
            print('для одного из способов нет подходящих опорных пикселей')
        else:
            dam_mean = np.mean(dam_values)

            base_reference = np.median(disp.data[base_mask])
            height_reference = np.median(disp.data[height_mask])

            base_move = (dam_mean - base_reference) * 1000
            height_move = (dam_mean - height_reference) * 1000

            print('пикселей базовой опоры:', np.sum(base_mask))
            print('пикселей опоры по высоте:', np.sum(height_mask))
            print(f'смещение с базовой опорой: {base_move:.3f} мм')
            print(f'смещение с опорой по высоте: {height_move:.3f} мм')
    else:
        print('\nсетки файлов отличаются — маски нельзя напрямую совмещать')
else:
    print('\nне удалось получить все три вырезки')