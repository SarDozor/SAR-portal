#для доставания из зипок файлы которые нам нужны для анализа
import os
import zipfile

#print(os.listdir('B1_series'))
allzipz = []
for s in os.listdir('B1_series'):
    allzipz.append(s)
allpathzip = []
for al in allzipz:
    f = f'B1_series/{al}'
    allpathzip.append(f)
#print(allpathzip)
'''with zipfile.ZipFile(allpathzip[0]) as firstzip:
    print(firstzip.namelist())'''

for al1 in allpathzip:
    with zipfile.ZipFile(al1) as firstzip:
        for zip in firstzip.namelist():
            if zip.endswith('_los_disp.tif') or zip.endswith('_corr.tif'):
                firstzip.extract(zip, 'B1_series')