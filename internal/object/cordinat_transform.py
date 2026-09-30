import rasterio
from rasterio.warp import transform

result = 'S1AA_20170118T232041_20170211T232040_VVP024_INT80_G_weF_DDAC'
latit = [-29.766708]
longit = [-69.976503]
vert_disp = f'{result}/{result}_vert_disp.tif'
with rasterio.open(vert_disp) as vd:
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
    all_table = vd.read(1) #весь слой данных. 1 в скобках тк всего слой 1
    #print(all_table)
    numb_table = all_table[numb_string, numb_column]
    print(numb_table)