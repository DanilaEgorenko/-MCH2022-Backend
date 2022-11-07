from operator import itemgetter
from bs4 import BeautifulSoup

class Document:
  parent = None
  url = ''
  extension = ''
  level = 0
  size = 0

  def __init__(self, url, parent, extension, size, level=0):
    self.url = url
    self.extension = extension,
    self.level = level
    self.parent = parent
    self.size = size

  def __str__(self):
    return f'Document [.{self.extension}] - {self.url}'

  def __repr__(self):
    return self.__str__()

class Page:
  soup = None
  parent = None
  title = ''
  url = ''
  level = 0
  inner_urls = []
  inner_docs = []

  def __init__(self, url, content, parent, level=0):
    self.soup = BeautifulSoup(content, 'html.parser')
    self.url = url
    self.level = level
    self.parent = parent
    self.title = self.soup.title.string
    self.inner_urls = list(map(itemgetter('href'), self.soup.find_all('a', href=True)))

  def __str__(self):
    return f'page [{self.level}]; {self.url}'

  def __repr__(self):
    return self.__str__()
    