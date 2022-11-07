import os

cache_dir = 'cache'
dataset_dir = 'datasets'

depricated_path = os.path.join(cache_dir, 'deprecated_hosts.dat')
train_dataset_path = os.path.join(dataset_dir, '79_rows_text_depth_2.csv')

max_pages_per_site = 50


# depricated resources
sources = [
  'http://www.gov.ru/main/regions/regioni-44.html',
  [
    'gov.ru', 'kremlin.ru', 'itunes.apple.com', 'play.google.com', 'www.rostrud.ru', 't.me', 'vk.com', 'ok.ru', 'www.youtube.com',
    'mintrud.gov.ru',
    'government.ru',
    'www.gosuslugi.ru',
    'www.vcot.info',
    'spravochnik.rosmintrud.ru',
    'www.migrakvota.gov.ru',
    'rutube.ru',
    'www.tiktok.com',
  ],
]

# parser
docs_extensions = [
  'pdf',
  'doc',
  'docx',
  'odt',
  'xls',
  'xlsx',
  'csv',
]
