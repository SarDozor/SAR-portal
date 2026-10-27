import os
os.environ.pop('PROJ_LIB', None)
os.environ.pop('PROJ_DATA', None)
from hyp3_sdk import HyP3
import rasterio
from rasterio.warp import transform

checker = False
if checker:
    hype = HyP3(prompt='password')
    first = 'S1B_IW_SLC__1SDV_20180208T082857_20180208T082924_009529_011271_2D26'
    second = 'S1B_IW_SLC__1SDV_20180220T082857_20180220T082924_009704_01182A_9D0F'
    prp_job = {
        'job_type': 'INSAR_GAMMA',
        'name': 'b1_dem',
        'job_parameters': {
            'granules': [first, second],
            'looks': '10x2', #как оказалось кредитов много и можно увеличить разрешение
            'apply_water_mask': True,
            'include_dem': True
        }
    }

    batch = hype.submit_prepared_jobs(prepared_jobs=[prp_job]) #отправка задачи на сервер nasa
    job = batch[0]
    print(f"задача отправлена {job}")
    job = hype.watch(job) #спрашивает у nasa закончилась ли обработка
    if job.succeeded():
        downloaded_files = job.download_files()
        print(f"скачанные файлы {downloaded_files}")
    else:
        print(f"ошибка: {job.status_code}")

foldername = 'S1BB_20180208T082857_20180220T082857_VVP012_INT40_G_weF_3093'
file = f'{foldername}/{foldername}/{foldername}_dem.tif'
lon = -44.120996
lat = -20.119388
lon1 = -44.134049
lat1 = -20.121520
with rasterio.open(file) as f:
    print(f.crs)
    print(f.shape)
    x, y = transform('EPSG:4326', f.crs, [lon], [lat])
    row, col = f.index(x[0], y[0])
    print(x, y)
    print(row, col)
    f_table = f.read(1)
    print(f_table[row][col])
print('-----------------------')
with rasterio.open(file) as f:
    print(f.crs)
    print(f.shape)
    x, y = transform('EPSG:4326', f.crs, [lon1], [lat1])
    row, col = f.index(x[0], y[0])
    print(x, y)
    print(row, col)
    f_table = f.read(1)
    print(f_table[row][col])