from hyp3_sdk import HyP3

first = 'S1A_IW_SLC__1SDV_20250103T100317_20250103T100344_057278_070BF6_E58D'
second = 'S1A_IW_SLC__1SDV_20250115T100316_20250115T100343_057453_0712E0_5818'
hype = HyP3(prompt='password')

prp_job = {
    'job_type': 'INSAR_GAMMA',
    'name': 'chile_test_20250103_20250115',
    'job_parameters': {
        'granules': [first, second],
        'looks': '10x2',
        'apply_water_mask': True,
        'include_dem': True,
        'include_displacement_maps': True
    }
}
batch = hype.submit_prepared_jobs(prepared_jobs=[prp_job])
job = batch[0]

print(f'Задача отправлена: {job}')

job = hype.watch(job)

if job.succeeded():
    downloaded_files = job.download_files()
    print(f'скачанные файлы: {downloaded_files}')
else:
    print(f'ошибка обработки: {job.status_code}')