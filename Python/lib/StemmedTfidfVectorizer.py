import nltk
from nltk.corpus import stopwords
from nltk.stem import SnowballStemmer
nltk.download('stopwords')
from sklearn.feature_extraction.text import TfidfVectorizer

def tokenize(raw):
  return [w.lower() for w in nltk.word_tokenize(raw) if w.isalpha()]

class StemmedTfidfVectorizer(TfidfVectorizer):
  stemmer = SnowballStemmer(language="russian")
  
  def build_analyzer(self):
    analyzer = super(StemmedTfidfVectorizer, self).build_analyzer()
    return lambda doc: (StemmedTfidfVectorizer.stemmer.stem(w) for w in analyzer(doc))
