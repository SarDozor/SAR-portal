import asf_search as asf

latitude = -29.766708
longitude = -69.976503
point = f"POINT ({longitude} {latitude})"
scenes = asf.geo_search(
    platform=asf.PLATFORM.SENTINEL1,
    processingLevel=asf.PRODUCT_TYPE.SLC,
    beamMode=asf.BEAMMODE.IW,
    intersectsWith=point,
    start='2025-01-01',
    end='2026-01-01'
)
print('всего найдено снимков:', len(scenes))
tracks = {}
for scene in scenes:
    path = scene.properties['pathNumber']
    frame = scene.properties['frameNumber']
    direction = scene.properties['flightDirection']
    group = (path, frame, direction)
    if group not in tracks:
        tracks[group] = 0
    tracks[group] += 1
print('трек, кадр, направление — количество снимков:')
for group in tracks:
    print(group, '—', tracks[group])
selected_scenes = []
for scene in scenes:
    if (
        scene.properties['pathNumber'] == 156
        and scene.properties['frameNumber'] == 688
        and scene.properties['flightDirection'] == 'DESCENDING'
    ):
        selected_scenes.append(scene.properties['sceneName'])
selected_scenes = sorted(selected_scenes)
print('\nвыбранные снимки:', len(selected_scenes))
for name in selected_scenes:
    print(name)