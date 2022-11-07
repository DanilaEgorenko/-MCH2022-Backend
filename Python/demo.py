import asyncio
from operator import itemgetter
from pprint import pprint
import joblib
import os
import numpy as np
from lib.StemmedTfidfVectorizer import StemmedTfidfVectorizer, tokenize
from lib.consts import cache_dir
import psycopg2

from lib.utils import load_deprecated
from lib.parser import merge_page_proccess_results, proccess_page, parse_site
from lib.consts import docs_extensions


# Запустите, чтобы увидеть как работает парсер сайтов.
# 
# Скрипт обходит сайт с url, заданным ниже, обращается к заранее обученной
# модели, чтобы определить предметную область сайта по тексту на нем и
# делает запись в БД.
# 
# последний пункт опциональный.
# Для выполнения последнего пункта требуется выполнить установку и настройку
# СУБД PostgreSQL и создать таблицу ParsedCompanies. См пункт "Установка".
# https://gitlab.com/Denactive/mch2022-backend
# 
# Впрочем, весь результат парсинга выводится в консоль, и строки под
# заголовком "Запись результатов в базу данных" (92-111) можно убрать.


url = 'https://www.goz.ru'


async def main():
  # ===========================================================================
  # 
  # настройка парсинга, парсинг и вывод результатов парсинга заданного url
  #
  # ===========================================================================

  # На эти хосты парсер заходить не будет
  deprecated = load_deprecated()

  # запуск парсинга с заданными настройками:
  # корневой url,
  # список игнорируемых хостов
  # список расширений файлов, считаемых документами: .odf, .odt, .xlsx и т.п.
  # нужно ли собирать документы. Если нет, на каждой листовой странице не будут выполняться HEAD-запросы
  # глубина обхода
  # лимит числа переходов по ссылками для _корневого_ url, чтобы не обходить сайты с кучей ссылок полностью
  pages, docs = await parse_site(
    url,
    deprecated = deprecated,
    docs_extensions = docs_extensions,
    search_docs = True,
    depth=3,
    max_pages_per_site=100,
  )
  pprint(pages)
  pprint(docs)

  result =  merge_page_proccess_results(list(proccess_page(page) for page in pages))
  print(result)


  # ===========================================================================
  # 
  # определение отрасли деятельности компании-владельца сайта с помощью TF-IDF
  #
  # ===========================================================================

  # загрузка модели
  le, clf = joblib.load(os.path.join(cache_dir, 'predict_industry_model.dat'))

  def get_encoded_industry_name(i):
    return le.inverse_transform([i])[0]

  def predict_industry(text):
    data = np.array([text])
    return get_encoded_industry_name(clf.predict(data))

  # определение отрасли по тексту на сайте
  industry = predict_industry(result['text'])
  print('\nпредсказанная отрасль:', industry )


  # ===========================================================================
  # 
  # Запись результатов в базу данных
  #
  # ===========================================================================

  docs = list(map(lambda x: x.url, docs))

  def get_query(industry, result, docs):
    q = f"INSERT INTO ParsedCompanies(category, phone, docks) VALUES('{industry}', ARRAY{result['phones'].split(' ')}, ARRAY{docs});"
    print('query', q)
    return q

  conn = psycopg2.connect(dbname='mch_db', user='user_mch', 
                        password='123456', host='localhost')

  conn.autocommit = True
  cursor = conn.cursor()

  cursor.execute(
    get_query(industry, result, docs)
  )

  cursor.close()
  conn.close()


# траблы с виндой (как обычно)
# os.name = 'posix' | 'nt' | 'Java' ~  Unix / Windows / Java
if os.name == 'nt':
  # https://stackoverflow.com/questions/45600579/asyncio-event-loop-is-closed-when-getting-loop
  asyncio.set_event_loop_policy(asyncio.WindowsSelectorEventLoopPolicy())
asyncio.run(main())
