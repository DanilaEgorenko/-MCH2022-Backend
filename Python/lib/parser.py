import mimetypes
from pprint import pprint
import re
import aiohttp
import datefinder
from operator import itemgetter

from lib.utils import get_host
from lib.models import Page, Document


async def net_worker(session, root, level, deprecated, docs_extensions, depth, max_pages_per_site, search_docs, silent, docs, pages, visited):
  for inner_url in root.inner_urls:
    # Проверка, чтобы избежать цикла
    if inner_url in visited:
      continue

    # скипаем мусор
    if inner_url == '#' or inner_url == '/' or inner_url.startswith('mailto') or inner_url.startswith('tel'):
      continue

    if inner_url == '':
      continue

    # относительный путь
    if inner_url[0] == '/':
      inner_url = root.url + inner_url

    if get_host(inner_url) in deprecated:
      if not silent: print('ignoring', inner_url, 'as deprecated')
      continue
    
    # не делать head, если не ищем документы
    if not search_docs and level + 1 == depth:
      continue

    if not silent: print('[' + str(level) + '] visiting', inner_url)
    timeout = aiohttp.ClientTimeout(total=10)   # 10 сек на все операции
    try:
      response = await session.head(inner_url, timeout=timeout)
      visited.add(inner_url)
      # проблема с text/html; charset=UTF-8
      # https://stackoverflow.com/questions/40412228/file-extension-from-mime-type-with-charset-utf-8
      content_type = response.headers['content-type'].split(";")[0]
      extension = mimetypes.guess_extension(content_type)

      if extension == '.html':
        if len(visited) >= max_pages_per_site:
          print('pages per site limit reached (' + str(max_pages_per_site) + ')')
          break
        if level+1 < depth :
          async with session.get(inner_url, timeout=timeout) as response_get:
            page = Page(inner_url, await response_get.text(), root.url, level+1)
            pages.append(page)
            await net_worker(session, page, level+1, deprecated, docs_extensions, depth, max_pages_per_site, search_docs, silent, docs, pages, visited)
        continue

      # собираем только документы из docs_extensions ([1:] - обрезает точку)
      if extension[1:] not in docs_extensions:
        continue

      is_transfer_encoding_chunked = False
      if 'Transfer-Encoding' in response.headers and response.headers['Transfer-Encoding'] == 'chunked':
        is_transfer_encoding_chunked = True
      size = 0 if is_transfer_encoding_chunked else response.headers['content-length'] if 'content-length' in response.headers else 0

      if not silent: print(f'got resource {inner_url}\n\tcontent-type: {content_type}\n\textension: {extension}\n\tsize: {size}')
      if extension:
        docs.append(Document(inner_url, root.url, extension, size, level))
    except Exception as e:
      if not silent: print('Error in getting resource:', inner_url, '\n\tMsg: ', e)

async def parse_site(
  url,
  deprecated = [],
  docs_extensions = [],
  depth = 2,
  search_docs = True,
  silent = False,
  max_pages_per_site = 10,
):
  docs = []
  pages = []
  visited = set()
  async with aiohttp.ClientSession() as session:
    async with session.get(url) as response:
      root = Page(url, await response.text(), None, 0)
      pages.append(root)
      visited.add(root)
      await net_worker(session, root, 0, deprecated, docs_extensions, depth, max_pages_per_site, search_docs, silent, docs, pages, visited)
  
  return (pages, docs)

def regex_to_phone(reg):
  reg0, reg1 = reg
  if reg0 == '' or reg0 == '+7':
    reg0 = '8'
  return reg0 + '-' + re.sub(' ', '-', reg1)

phoneRegex = '(8|\+7)?\ ?(\d{3}\ \d{3}\ \d{2}\ ?\d{2})'

def proccess_page(page):
  body = page.soup.find('body')
  desciption = page.soup.head.find('meta', attrs={'name': 'description', 'content': True})
  body_text = ''
  desciption_text = ''
  if not body:
    print('Сайт без тела (' + page.url + ')')
  else:
    body_text = body.get_text()
  if not desciption:
    print('нет meta desciption')
  else:
    desciption_text = re.sub('[\d]', '', desciption['content']).lower()

  body_text = re.sub('[\W|\s|\\n]+', ' ', body_text)

  numbers = 'a'.join(re.findall('\d+[\s\d]+', body_text))
  dates = ' '.join(map(lambda d: f'{d.day}.{d.month}.{d.year}', filter(lambda d: d.year > 1990, list(datefinder.find_dates(numbers)))))
  phones = ' '.join(map(regex_to_phone, re.findall(phoneRegex, numbers)))

  body_text = re.sub('[\d]', '', body_text).lower()

  return {'dates': dates, 'phones': phones, 'text': body_text + ' ' + desciption_text}

def merge_page_proccess_results(results):
  if type(results) == list:
    return {
      'dates': ' '.join(map(itemgetter('dates'), results)),
      'phones': ' '.join(map(itemgetter('phones'), results)),
      'text': ' '.join(map(itemgetter('text'), results))
    }
  return results
