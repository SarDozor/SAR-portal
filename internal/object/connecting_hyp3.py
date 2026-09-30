from hyp3_sdk import HyP3

hyp3 = HyP3(prompt='password') #вход в аккаунт наса. пароль не видно, но он пишется
first_granule = 'S1A_IW_SLC__1SSV_20170118T232041_20170118T232108_014892_0184B7_13CE'
second_granule = 'S1A_IW_SLC__1SSV_20170211T232040_20170211T232107_015242_018F81_1B1D'
#print(hyp3.check_credits()) #хорошая строка для проверки токенов(кредитов) в этом месяце у моего аккаунта nasa
#job_name = 'dangerous_damb'
prepared_job = {
    'job_type': 'INSAR_GAMMA',
    'name': 'dangerous_damb',
    'job_parameters': {
        'granules': [first_granule, second_granule],
        'looks': '20x4',
        'apply_water_mask': True,
        'include_look_vectors': True,
        'include_displacement_maps': True #модуль переводит пересчёт геометрии съемки в сантиметры
    }
} #настройка интерферометрии dinsar
batch = hyp3.submit_prepared_jobs(prepared_jobs=[prepared_job]) #отправка задачи на сервер nasa
job = batch[0]
print(f"задача отправлена {job}")
job = hyp3.watch(job) #спрашивает у nasa закончилась ли обработка
if job.succeeded():
    downloaded_files = job.download_files()
    print(f"скачанные файлы {downloaded_files}")
else:
    print(f"ошибка: {job.status_code}")
