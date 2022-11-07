import asyncio
from operator import itemgetter
from pprint import pprint
from lib.models import Document
import joblib
import os

from lib.utils import load_deprecated
from lib.parser import merge_page_proccess_results, proccess_page, parse_site
from lib.consts import docs_extensions, cache_dir, max_pages_per_site, dataset_dir

import numpy as np
import pandas as pd


tmp_dir = os.path.join(cache_dir, 'tmp_http_depth_3')
output_file = 'text_http.csv'

new_column_title = {
  'phones': 'Телефоны',
  # TODO: email
  'email': 'Email',
  'text': 'Текст',
  'dates': 'Даты',
  'docs': 'Документы'
}

async def proccess_url(url):
  deprecated = load_deprecated()

  pages, docs = await parse_site(
    url,
    deprecated = deprecated,
    docs_extensions = docs_extensions,
    search_docs = True,
    depth=2,
    max_pages_per_site=max_pages_per_site,
    # silent = True,
  )
  pprint(pages)
  pprint(docs)
  
  data = merge_page_proccess_results(list(proccess_page(page) for page in pages))
  data['docs'] = ' '.join(map(lambda x: x.url, docs))
  return data


async def main():
  # Открою датасет
  df = pd.read_csv(os.path.join(dataset_dir, 'cleared_indexed_79_rows.csv'), sep=",")
  # характеристики
  print("размер:", df.shape)
  print("\nколонки:\n", df.dtypes)

  for i, url in enumerate(df['Сайт']):
    try:
      proccess_data = await proccess_url('http://' + url)
      joblib.dump(proccess_data, os.path.join(tmp_dir, str(i) + '.tmp'))
      for key, chunk in proccess_data.items():
        df.loc[df.index[i], new_column_title[key]] = chunk
    except Exception as e:
      print('Error:', e)

  df.to_csv(output_file, index=False)

# траблы с виндой (как обычно)
# os.name = 'posix' | 'nt' | 'Java' ~  Unix / Windows / Java
if os.name == 'nt':
  # https://stackoverflow.com/questions/45600579/asyncio-event-loop-is-closed-when-getting-loop
  asyncio.set_event_loop_policy(asyncio.WindowsSelectorEventLoopPolicy())
asyncio.run(main())
