import os
os.environ.pop('PROJ_LIB', None)
os.environ.pop('PROJ_DATA', None)
import numpy as np
import rasterio
from rasterio.warp import transform
import matplotlib.pyplot as plt

products = [
    'S1AA_20250103T100317_20250115T100316_VVP012_INT40_G_weF_BD3A',
    'S1AA_20250115T100316_20250127T100316_VVP012_INT40_G_weF_8927',
    'S1AA_20250127T100316_20250208T100315_VVP012_INT40_G_weF_E083',
]
lon = -69.976503
lat = -29.766708

first_grid = None
base_masks = []
height_masks = []
dam_masks = []
displacements = []
for product in products:
    windows = {}
    for kind in ('corr', 'los_disp', 'dem'):
        path = f'chile_test/{product}/{product}_{kind}.tif'
        with rasterio.open(path) as src:
            xs, ys = transform('EPSG:4326', src.crs, [lon], [lat])
            row, col = src.index(xs[0], ys[0])
            corners = np.array([
                src.xy(row - 40, col - 40),
                src.xy(row - 40, col + 40),
                src.xy(row + 40, col - 40),
            ])
            if first_grid is None:
                first_grid = (src.crs, corners)
            first_crs, first_corners = first_grid
            if src.crs != first_crs or not np.allclose(
                corners, first_corners, rtol=0, atol=1e-6
            ):
                raise ValueError(f'сетка вырезки отличается: {product}, {kind}')
            table = src.read(1, masked=True)
            window = table[row - 40:row + 41, col - 40:col + 41].copy()
            if window.shape != (81, 81):
                raise ValueError(f'неполная вырезка: {product}, {kind}')
            windows[kind] = window
    corr = windows['corr']
    disp = windows['los_disp']
    dem = windows['dem']
    valid_corr = ~np.ma.getmaskarray(corr) & np.isfinite(corr.data)
    valid_disp = ~np.ma.getmaskarray(disp) & np.isfinite(disp.data)
    valid_dem = ~np.ma.getmaskarray(dem) & np.isfinite(dem.data)
    if not valid_dem[40, 40]:
        raise ValueError(f'нет высоты в центре: {product}')
    base_mask = valid_corr & valid_disp & (corr.data > 0.8)
    base_mask[30:51, 30:51] = False
    height_mask = (
        base_mask & valid_dem
        & (np.abs(dem.data - dem.data[40, 40]) <= 100)
    )
    base_masks.append(base_mask)
    height_masks.append(height_mask)
    dam_masks.append(valid_disp[39:42, 39:42])
    displacements.append(disp.data)
common_base = base_masks[0] & base_masks[1] & base_masks[2]
common_height = height_masks[0] & height_masks[1] & height_masks[2]
common_dam = dam_masks[0] & dam_masks[1] & dam_masks[2]
print('сетки всех вырезок совпадают')
print('общих пикселей базовой опоры:', np.sum(common_base))
print('общих пикселей опоры по высоте:', np.sum(common_height))
print('общих пикселей в участке дамбы 3×3:', np.sum(common_dam))
if not common_base.any() or not common_height.any() or not common_dam.any():
    raise ValueError('лдна из общих масок пуста')
for product, disp in zip(products, displacements):
    dam_mean = np.mean(disp[39:42, 39:42][common_dam])
    base_mm = (dam_mean - np.median(disp[common_base])) * 1000
    height_mm = (dam_mean - np.median(disp[common_height])) * 1000
    print(f'{product[5:13]} → {product[21:29]}')
    print(f'базовая: {base_mm:+.3f} мм, по высоте: {height_mm:+.3f} мм')
second_disp = displacements[1]
excluded = common_base & ~common_height
print(f'медиана всей базовой опоры: {np.median(second_disp[common_base]) * 1000:.3f} мм')
print(f'медиана оставшейся опоры: {np.median(second_disp[common_height]) * 1000:.3f} мм')
print(f'медиана исключённых пикселей: {np.median(second_disp[excluded]) * 1000:.3f} мм')
fig, (ax_disp, ax_ref) = plt.subplots(1, 2, figsize=(11, 5))
image = ax_disp.imshow(second_disp * 1000, cmap='viridis')
fig.colorbar(image, ax=ax_disp, label='LOS, мм')
ax_disp.set_title('смещения в растре, 15–27 января')
ax_disp.scatter(40, 40, marker='+', color='white', s=100)
for mask, label, color in [
    (excluded, 'исключённые', 'tab:orange'),
    (common_height, 'оставшиеся', 'tab:blue'),
]:
    rows, cols = np.where(mask)
    ax_ref.scatter(cols, rows, s=4, color=color, label=label)
ax_ref.scatter(40, 40, marker='+', color='black', s=100, label='точка дамбы')
ax_ref.set_title('какие пиксели оставил фильтр')
ax_ref.set_xlim(-0.5, 80.5)
ax_ref.set_ylim(80.5, -0.5)
ax_ref.set_aspect('equal')
#ax_ref.legend()
for ax in (ax_disp, ax_ref):
    ax.set_xlabel('col')
    ax.set_ylabel('row')
fig.tight_layout()
base_without_right = common_base.copy()
base_without_right[:, 60:] = False
dam_mm = np.mean(second_disp[39:42, 39:42][common_dam]) * 1000
reference_mm = np.median(second_disp[base_without_right]) * 1000
print('пикселей в опоре без правой полосы:', np.sum(base_without_right))
print(f'медиана опоры без правой полосы: {reference_mm:.3f} мм')
print(f'смещение дамбы с этой опорой: {dam_mm - reference_mm:+.3f} мм')
for product, disp in zip(products, displacements):
    dam_mean = np.mean(disp[39:42, 39:42][common_dam])
    full_ref = np.median(disp[common_base])
    trimmed_ref = np.median(disp[base_without_right])
    full_mm = (dam_mean - full_ref) * 1000
    trimmed_mm = (dam_mean - trimmed_ref) * 1000
    print(f'{product[5:13]} → {product[21:29]}')
    print(f'полная опора: {full_mm:+.3f} мм')
    print(f'без правой полосы: {trimmed_mm:+.3f} мм')
    print(f'изменение оценки: {trimmed_mm - full_mm:+.3f} мм')
fig.savefig('chile_reference_jan15_jan27.png', dpi=180)
#plt.show()