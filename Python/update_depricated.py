import requests
from pprint import pprint
import joblib

from lib.utils import get_host, load_deprecated
from lib.models import Page
from lib.consts import sources, depricated_path

# deprecated = set()
deprecated = load_deprecated()

for source in sources:
  chunk = source
  if type(source) == list:
    chunk = source
  else:
    chunk = []
    if source in deprecated:
      continue
    try:
      r = requests.get(source)
      root = Page(source, r.text, None, 0)
      chunk.extend(map(get_host, root.inner_urls))
      chunk.append(source)
    except Exception as e:
      print('Error in getting resource:', source, '\n\tMsg: ', e)
  for el in chunk:
    if el == None:
      continue
    deprecated.add(el)

joblib.dump(deprecated, depricated_path)

print('updated cache:') 
pprint(deprecated)
