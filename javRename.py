import os
from bs4 import BeautifulSoup
import re
import requests
# from requests.exceptions import RequestException
import urllib3
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

SEARCHURL = 'https://avmoo.host/cn/search/'
reNameOrder = ['actress', 'javCode', 'title', 'publishDate']
downimg = False


PROXY = {"http": "socks5://127.0.0.1:1099", "https": "socks5://127.0.0.1:1099"}
HEADER = {
    "user-agent":
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36",
    "accept": "application/json, text/javascript, */*; q=0.01",
    "accept-language": "zh-CN,zh;q=0.9",
    "dnt": "1",
    "Connection": "keep-alive",
    # "content-type": "application/x-www-form-urlencoded; charset=UTF-8",
    # "x-requested-with": "XMLHttpRequest",
    'content-length': '0',
    'sec-fetch-dest': 'empty',
    'sec-fetch-site': 'same-origin',
    'sec-fetch-mode': 'cors'
}


def findActress(soup):
    avatars = soup.select('.avatar-box')
    actresses = [avatar.select_one('span').text for avatar in avatars]
    if len(actresses):
        actress = ' '.join(actresses)
    else:
        actress = 'unknown'
    return actress
    
def searchTitle(link):
    finalstr = None
    imgResponse = None
    try:
        PROXY
    except NameError:
        PROXY = None
    try:
        response = requests.get(link, headers=HEADER, allow_redirects=True, proxies=PROXY, verify=False)
        if response.status_code == 200:
            soup = BeautifulSoup(response.text, 'html.parser')            
            javCode = soup.select_one('span.header').findNext('span').text
            title = soup.select_one('h3').text.split(javCode+' ')[1]
            publishDate = re.search(r'\d{4}-\d{2}-\d{2}', soup.select('span.header')[1].parent.text).group(0)
            actress = findActress(soup)
            imgurl = soup.select_one('.bigImage')['href']
            imgResponse = requests.get(imgurl, headers=HEADER, allow_redirects=True, proxies=PROXY)
            renameStr = {'actress': actress, 'javCode': javCode, 'title': title, 'publishDate': publishDate}
            neededOrder = [renameStr[i] for i in reNameOrder]
            finalstr = neededOrder[0]+'-['+']-['.join(neededOrder[1:4])+']'
    except Exception as e:
        print(e)
    return finalstr, imgResponse

def searchID(javCode):
    request_url = SEARCHURL+javCode
    try:
        PROXY
    except NameError:
        PROXY = None
    try:
        response = requests.get(request_url, headers=HEADER, allow_redirects=True, proxies=PROXY, verify=False)
        if response.status_code == 200:
            soup = BeautifulSoup(response.text, 'html.parser')
            links = soup.select('.movie-box')
            if len(links)==1:
                return searchTitle(links[0]['href'])
            else:
                javCode = javCode.replace('-', '')
                for link in links:
                    javID = link.select_one('date').text.replace('-', '')
                    if javCode==javID:
                        return searchTitle(link['href'])
                return None, None
    except Exception as e:
        print(e)

def javRe(fullpath):
    splitstr = os.path.split(fullpath.rstrip(os.path.sep))
    basepath, filename0 = splitstr[0], splitstr[1]
    matchObj = re.search(r"\.[A-Za-z0-9]{3,10}$", filename0)
    suffix = matchObj.group(0)
    filename = filename0[0:matchObj.span()[0]]
    javCode = re.search(r'(?![^A-Za-z])?[A-Za-z]{3,5}-?\d{3,5}(?=\D)?', filename)
    if javCode is not None:
        newName, imgResponse = searchID(javCode.group(0))
    else:
        newName = None
    if newName is not None:
        try:
            os.rename(fullpath, os.path.join(basepath, newName+suffix))
            if downimg:                
                open(os.path.join(basepath, newName+'.jpg'), 'wb').write(imgResponse.content)
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
    else:
        javRe(sys.argv[1])