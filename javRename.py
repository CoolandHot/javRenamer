import os
from bs4 import BeautifulSoup
import re
import requests
# from requests.exceptions import RequestException
import urllib3
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

reNameOrder = ['actress', 'javCode', 'title', 'publishDate']
downimg = False
useJavBus = True # if False, use avmoo


PROXY = {"http": "socks5://127.0.0.1:1099", "https": "socks5://127.0.0.1:1099"}
HEADER = {
    "user-agent":
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36",
    "accept": "application/json, text/javascript, */*; q=0.01",
}

class basic_class: # this class applies to both javBus and avmoo
    def __init__(self, javCode):
        self.javID = javCode
        self.title, self.publishDate, self.heroine = None, None, None
        self.img, self.request_url= None, None
        self.finalstr = None # store the new name
    def searchTitle(self, link):
        try:
            response = requests.get(link, headers=HEADER, allow_redirects=True, proxies=PROXY, verify=False)
            if response.status_code == 200:
                soup = BeautifulSoup(response.text, 'html.parser')            
                self.javID = soup.select_one('span.header').findNext('span').text
                self.title = soup.select_one('h3').text.split(self.javID+' ')[1]
                self.publishDate = re.search(r'\d{4}-\d{2}-\d{2}', soup.select('span.header')[1].parent.text).group(0)
                self.findActress(soup) # find heroines
                if downimg:
                    imgurl = soup.select_one('.bigImage')['href'] # find the cover
                    imgResponse = requests.get(imgurl, headers=HEADER, allow_redirects=True, proxies=PROXY) # download cover
                    self.img = imgResponse.content
                renameStr = {'actress': self.heroine,
                            'javCode': self.javID,
                            'title': self.title,
                            'publishDate': self.publishDate}
                neededOrder = [renameStr[i] for i in reNameOrder] # rename by the set order
                self.finalstr = neededOrder[0]+'-['+']-['.join(neededOrder[1:4])+']' # heroine's name is not covered by [ ]
        except Exception as e:
            print(e)
    def searchID(self):
        try:
            response = requests.get(self.request_url, headers=HEADER, allow_redirects=True, proxies=PROXY, verify=False)
            if response.status_code == 200:
                soup = BeautifulSoup(response.text, 'html.parser')
                links = soup.select('.movie-box')
                if len(links)==1: # request returns only one result
                    self.searchTitle(links[0]['href'])
                else:
                    javCode = self.javID.replace('-', '') # replace any '-' for later comparison
                    for link in links:
                        javID = link.select_one('date').text.replace('-', '') # acquire each result's javID
                        if javCode==javID: # find the one with exact match
                            self.searchTitle(link['href'])
        except Exception as e:
            print(e)
    def findActress(self, soup):
        return

class javBus(basic_class):
    def __init__(self, javCode):
        super().__init__(javCode)
        self.request_url = 'https://www.javbus.com/search/'+self.javID
    def findActress(self, soup):
        avatars = soup.select('.star-name')
        actresses = [avatar.select_one('a').text for avatar in avatars] # there could be multiple heroines
        if len(actresses):
            self.heroine = ' '.join(actresses)
        else:
            self.heroine = 'unknown'

class avmoo(basic_class):
    def __init__(self, javCode):
        super().__init__(javCode)
        self.request_url = 'https://avmoo.cyou/cn/search/'+self.javID
    def findActress(self, soup):
        avatars = soup.select('.avatar-box')
        actresses = [avatar.select_one('span').text for avatar in avatars]
        if len(actresses):
            self.heroine = ' '.join(actresses)
        else:
            self.heroine = 'unknown'



def javRe(fullpath):
    splitstr = os.path.split(fullpath.rstrip(os.path.sep))
    basepath, filename0 = splitstr[0], splitstr[1]
    matchObj = re.search(r"\.[A-Za-z0-9]{3,10}$", filename0)
    suffix = matchObj.group(0)
    filename = filename0[0:matchObj.span()[0]]
    javCode = re.search(r'(?![^A-Za-z])?[A-Za-z]{2,5}-\d{3,5}(?=\D)?', filename)
    if javCode is None:
        javCode = re.search(r'(?![^A-Za-z])?[A-Za-z]{2,5}-?\d{3,5}(?=\D)?', filename)

    if useJavBus:
        searchResult = javBus(javCode.group(0)) # use javBus website
    else:
        searchResult = avmoo(javCode.group(0))
    searchResult.searchID()
    newName, imgResp = searchResult.finalstr, searchResult.img

    if newName is not None:
        try:
            os.rename(fullpath, os.path.join(basepath, newName+suffix))
            if downimg:
                open(os.path.join(basepath, newName+'.jpg'), 'wb').write(imgResp)
        except Exception as e:
            print(e)
            print('fail on', filename)
        else:
            print('successful on', filename)
    else:
        print('fail on', filename)

if __name__ == '__main__':
    import sys
    if len(sys.argv)==1:
        print("use command: python javRename.py inputFile.mp4")
        input()
    else:
        javRe(sys.argv[1])