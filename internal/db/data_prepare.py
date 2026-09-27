import pandas as pd
import openpyxl

"""
#здесь переносим Бд с данными в csv формат

df = pd.read_excel("dambs.xlsx", header=2)
df = df.loc[:, ['tsf', 'country', 'latitude', 'longitude', 'status', 'raise_type', 'current_maximum_height']]
df = df.dropna(subset=['latitude', 'longitude'])

print(df.shape)
print(df.columns)
print(df.head())

df.to_csv('clean_dambs.csv', index=False) 
"""

df = pd.read_csv("clean_dambs.csv")
df['tsf'] = df['tsf'].str.replace('\n', ' ')
df = df.drop_duplicates(subset=['latitude', 'longitude'])
#df_upstream = df[df['raise_type'] == 'Upstream']
df_upstream = df[df['raise_type'].str.contains('Upstream', na=False)]
df_upstream = df_upstream.sort_values(by= 'current_maximum_height', ascending=False)

"""
print(df_upstream.shape)
print(df_upstream.columns)
print(df_upstream.head(8))
"""

df_testdambs = df_upstream.head(8)
#print(df_testdambs)

df_safe = df[~df['raise_type'].str.contains('In-pit', na=False)]
df_safe = df_safe[~df_safe['raise_type'].str.contains('Upstream', na=False)]
df_safe = df_safe[df_safe['status'] == 'Active']
df_safe = df_safe.sort_values(by= 'current_maximum_height', ascending=False)

"""
print(df_safe.shape)
print(df_safe.columns)
print(df_safe.head(4))
"""

df_testdambs2 = df_safe[df_safe['current_maximum_height'] >0].head(4)
df_result = pd.concat([df_testdambs, df_testdambs2])
print(df_result)
df_result.to_csv('test_dambs.csv', index=False)
