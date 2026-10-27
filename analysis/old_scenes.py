#здесь будет попытка найти ряд по брумадиньо
import asf_search as asf
from hyp3_sdk import HyP3

latitude = -20.119 #пока что примерные корды по брумадиньо(далее везде я буду писать B1 так как мне ужасно не нравится полное название)
longitude = -44.121
point1 = f"POINT ({longitude} {latitude})"
#print(point1)
srch = asf.geo_search(
    platform = asf.PLATFORM.SENTINEL1,
    processingLevel = asf.PRODUCT_TYPE.SLC,
    beamMode = asf.BEAMMODE.IW,
    intersectsWith = point1,
    start = '2018-01-01',
    end = '2019-01-25'
) #подключается к серверам asf и ищет снимки еще живого хвостохранилища B1 до катастрофы
#print(len(srch)) #за данный период выдало 65 снимков со спутника B1. рандомные брать нельзя, т.к. у всех разная траектория. их просто нельзя будет сравнить
ph1 = srch[0]
'''print(ph1.properties) #берем первый найденный снимок и смотрим какие у него свойства, и берем значения тех которые нужны
print("первый снимок")
print(f"дата: {ph1.properties['startTime']}, трек: {ph1.properties['pathNumber']}, кадр: {ph1.properties['frameNumber']}, направление: {ph1.properties['flightDirection']}, имя снимка: {ph1.properties['sceneName']}")
print("все снимки")'''
srch2 = []
for i in range(65):
    ph_any = srch[i]
    #print(f"дата: {ph_any.properties['startTime']}, трек: {ph_any.properties['pathNumber']}, кадр: {ph_any.properties['frameNumber']}, направление: {ph_any.properties['flightDirection']}, имя снимка: {ph_any.properties['sceneName']}")
    if ph_any.properties['pathNumber'] == 53:
        srch2.append(ph_any.properties['sceneName'])
srtd = []
srtd = sorted(srch2)
#print(len(srtd))
#print(srtd)

pairs = []
for i in range(0, len(srtd)-1):
    ph_any = srtd[i]
    ph_scnd = srtd[i+1]
    pair_i = []
    pair_i.append(ph_any)
    pair_i.append(ph_scnd)
    pairs.append(pair_i)
#print(len(pairs))
#print(pairs)
"""first = srtd[0]
print(first[17:25])"""
list1 = []
"""print(pairs[0])
print(pairs[1])"""
for pair in pairs:
    first, scnd = pair
    prepared_job1 = {
        'job_type': 'INSAR_GAMMA',
        'name': f'B1_{first[17:25]}__{scnd[17:25]}',
        'job_parameters': {
            'granules': pair,
            'looks': '10x2', #как оказалось кредитов много и можно увеличить разрешение
            'apply_water_mask': True,
            'include_look_vectors': True,
            'include_displacement_maps': True #пересчитывает развёрнутую фазу в смещение, в метрах
        }
    }
    list1.append(prepared_job1)
#print(list1)
#print(len(list1))
#возьмем тестовый снимок - первую пару
subm1t = False #выключатель, чтоб не отправить случайно
if subm1t:
    hyp3 = HyP3(prompt='password') #вход в аккаунт наса. пароль не видно, но он пишется
    
    save1 = hyp3.check_credits()
    print(f"credits before: {save1}")

    batches = hyp3.submit_prepared_jobs(prepared_jobs=list1) #отправка задачи на сервер nasa
    for batch in batches:
        print(f"ждем задачу: {batch}")
        batch = hyp3.watch(batch) #спрашивает у nasa закончилась ли обработка
        if batch.succeeded():
            downloaded_files = batch.download_files(location='B1_series') #создаю папку для Б1 рядов и туда скачиваю архивы
            print(f"скачанные файлы {downloaded_files}")
        else:
            print(f"ошибка: {batch.status_code}") 
    #весь ряд должен весить около 30 гб

    save2 = hyp3.check_credits()
    print(f"credits after: {save2}")
    print(f"raznica: {(save1 - save2)/len(list1)}")
#print(pairs[1], pairs[2], pairs[3])
fr, sr = pairs[0]
print(fr)
print(sr)