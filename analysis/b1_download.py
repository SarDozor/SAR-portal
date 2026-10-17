#файл для скачивание зипок обработки отдельно. в данном случае по б1
from hyp3_sdk import HyP3

hyp3 = HyP3(prompt='password')
#print(hyp3.find_jobs(job_type='INSAR_GAMMA')) #показывает статусы работ: успешно, провалено и тп
b2 = []
b1 = hyp3.find_jobs(job_type='INSAR_GAMMA')
#print(vars(b1[0])) #vars дало много инфы: стоимость обработки в 10х2 стоит 15 кредитов, один архив снимков в этом разрешении весит около 0.9 ГБ, обработка шла около 35 минут, файлы для скачивания доступны всего две недели
sobm = False
if sobm:
    for b in b1:
        if b.name.startswith('B1_'):
            b2.append(b)
    for b in b2:
        print(b.name)
        b.download_files(location='B1_series')
    print("скачано")