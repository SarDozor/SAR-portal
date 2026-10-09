import pandas as pd
import openpyxl


#здесь переносим дб с данными в csv формат

df = pd.read_excel("data/dambs.xlsx", header=2)
df = df.loc[:, ['tsf', 'country', 'latitude', 'longitude', 'status', 'raise_type', 'current_maximum_height']]
df['longitude'] = df['longitude'].astype(str) #на этапе с файлом по брумадиньо оказалось, что весь столбец с долготами - строчный, поэтому нужно исключить все ошибки и заново создать csv
df['longitude'] = df['longitude'].str.strip()
df['longitude'] = pd.to_numeric(df['longitude'], errors = 'coerce')
df['latitude'] = df['latitude'].astype(str) #на всякий для широт. сначала превращаю в строки, потом убираю случайные пробелы из за которых не получается перевести строки в числа при помощи .strip. потом в числа
df['latitude'] = df['latitude'].str.strip()
df['latitude'] = pd.to_numeric(df['latitude'], errors = 'coerce')
df = df.dropna(subset=['latitude', 'longitude'])
print(df.dtypes) #проверка изменений

print(df.shape)
print(df.columns)
print(df.head())

df.to_csv('data/clean_dambs.csv', index=False) 


df = pd.read_csv("data/clean_dambs.csv")
df['tsf'] = df['tsf'].str.replace('\n', ' ')
df = df.drop_duplicates(subset=['latitude', 'longitude'])
#df_upstream = df[df['raise_type'] == 'Upstream']
df_upstream = df[df['raise_type'].str.contains('Upstream', na=False)]
df_upstream = df_upstream.sort_values(by= 'current_maximum_height', ascending=False)


print(df_upstream.shape)
print(df_upstream.columns)
print(df_upstream.head(8))


df_testdambs = df_upstream.head(8)
#print(df_testdambs)

df_safe = df[~df['raise_type'].str.contains('In-pit', na=False)]
df_safe = df_safe[~df_safe['raise_type'].str.contains('Upstream', na=False)]
df_safe = df_safe[df_safe['status'] == 'Active']
df_safe = df_safe.sort_values(by= 'current_maximum_height', ascending=False)


print(df_safe.shape)
print(df_safe.columns)
print(df_safe.head(4))


df_testdambs2 = df_safe[df_safe['current_maximum_height'] >0].head(4)
df_result = pd.concat([df_testdambs, df_testdambs2])
print(df_result)
df_result.to_csv('data/test_dambs.csv', index=False)
