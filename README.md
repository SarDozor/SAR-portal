## Повторный расчёт смещений для B1

1. Положите ZIP-архивы HyP3 в папку `B1_series/` в корне проекта.
2. Проверьте, что файл DEM лежит по пути, указанному в переменной `dem`
   в `analysis/def_cordinate_transform.py`.
3. Откройте терминал в корне проекта и выполните:

   ```powershell
   python analysis/unzip_b1.py
   python analysis/def_cordinate_transform.py
   ```

Первый скрипт распакует растры из ZIP. Второй рассчитает смещения
и создаст `b1_measurements_long.csv` для импорта. Для текущего набора B1
ожидается 96 строк данных и 12 столбцов.
