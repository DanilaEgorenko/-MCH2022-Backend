import joblib
import urllib3

from lib.consts import depricated_path

def convert_bytes(num):
  """
  this function will convert bytes to MB.... GB... etc
  """
  for x in ['bytes', 'KB', 'MB', 'GB', 'TB']:
    if num < 1024.0:
        return "%3.1f %s" % (num, x)
    num /= 1024.0

def get_host(url):
  try:
   return urllib3.util.parse_url(url).host
  except:
    return None

def load_deprecated():
  return joblib.load(depricated_path)
