#!/usr/bin/python
# -*- coding: utf-8 -*-

# This script is used to automatically replace all corresponding characters in the Excel table.
# warning! Replace all file contents in the target directory after running.If you make a mistake, you can retrieve the previous code in the backup directory.

import xlrd
import types
import glob
import time
import json
import os
import re

##########################################################
# Config
##########################################################
# Language
Language = 'Simplified Chinese'
# error list
ErrorList = []
# column list
Columns = ['A','B','C','D','E','F','G','H','I','J','K','L','M','N','O','P','Q','R','S','T','U','V','W','X','Y','Z','AA','AB','AC','AD','AE','AF','AG','AH','AI','AJ','AK','AL','AM','AN','AO','AP','AQ','AR','AS','AT','AU','AV','AW','AX','AY','AZ']
# src dir
SrcDir = './game-src'
# dest dir
DestDir = './game'
# allow english take
AllowWaitEnglish = False
# allow out location tag
AllowWaitLocation = False

##########################################################
# libs
##########################################################
# get excel
def GetExcel(fileSrc):
    if fileSrc.find('~$') > -1:
        return False
    excelDom = xlrd.open_workbook(fileSrc)
    if not excelDom:
        SendError('无法读取excel文件：' + fileSrc)
    print('读取excel文件成功：' + fileSrc)
    return excelDom

# get value by sheet key
# eg: key = 'A1'
def GetExcelSheetKey(excel,key):
    try:
        return excel.sheet_by_index(key)
    except:
        return False

# send error
def SendError(message):
    print('error: ' + message)
    ErrorList.append(message)
    time.sleep(300)
    exit()

# get data from excel sheet
def ReadSheetValues(sheet):
    if not sheet:
        return False
    newData = {}
    for h in range(sheet.nrows):
        for c in range(sheet.ncols):
            key = Columns[c] + str(h+1)
            newData[key] = sheet.cell(h,c).value
    return newData

# delete dir
def DeleteDir(path):
    if not os.path.exists(path):
        return True
    if os.path.isfile(path):
        try:
            os.remove(path)
            print('delete file: ' + path)
            return True
        except IOError:
            SendError('delete file: ' + path)
            return False
    dirs = os.listdir(path)
    for fileName in dirs:
        fileSrc = path + '/' + fileName
        if os.path.isfile(fileSrc):
            DeleteDir(fileSrc)
        if os.path.isdir(fileSrc):
            if not DeleteDir(fileSrc):
                SendError('delete dir or file: ' + fileSrc)
                return False
    os.rmdir(path)
    print('delete dir(path): ' + path)
    return True

# copy dir
def CopyDir(src,dest):
    dirs = os.listdir(src)
    for fileName in dirs:
        fileSrc = src + '/' + fileName
        fileNewSrc = dest + '/' + fileName
        # is dir
        if os.path.isdir(fileSrc):
            # create dir
            os.makedirs(fileNewSrc)
            if not SearchListFromGame(fileSrc, fileNewSrc):
                return False
        # is file, replace content
        if os.path.isfile(fileSrc):
            # open file
            fileObj = open(fileSrc, "r+", encoding='UTF-8')
            fileNewObj = open(fileNewSrc, "w+", encoding='UTF-8')
            try:
                fileData = fileObj.read()
                # save new file
                fileNewObj.write(fileData)
            finally:
                fileObj.close()
                fileNewObj.close()
    return True

##########################################################
# run
##########################################################
# now replace count
ReplaceCount = 0
# replace total
ReplaceTotal = 0
# Now file src
NowReplaceFileSrc = ''
# need replace content
# eg: [{'src': 'english', 'dest': 'chinese'}]
WaitReplace = []
# wait write english data
WaitEnglishData = []
# skip word
WaitEnglishSkip = ['<', '>', '你', '我', '它']
# wait location tag
WaitLocationTags = []

# get language file data
excel = GetExcel('./languages/' + Language + '.xlsx')
print('load file...')
if not excel:
    SendError('load excel file.')

# load add need replace content
for sheetKey in range(0, len(excel.sheets())):
    sheet = GetExcelSheetKey(excel, sheetKey)
    rowCount = sheet.nrows
    for rowKey in range(0, rowCount):
        srcLang = sheet.cell_value(rowKey, 0)
        if srcLang == '':
            continue
        newLang = sheet.cell_value(rowKey, 1)
        if newLang == '':
            continue
        WaitReplace.append({'src': srcLang, 'dest': newLang})

# find game-src file list
def SearchListFromGame(path, newPath):
    global NowReplaceFileSrc
    dirs = os.listdir(path)
    for fileName in dirs:
        fileSrc = path + '/' + fileName
        fileNewSrc = newPath + '/' + fileName
        NowReplaceFileSrc = fileNewSrc
        # is dir
        if os.path.isdir(fileSrc):
            # create dir
            os.makedirs(fileNewSrc)
            print('create dir: ' + fileNewSrc)
            if not SearchListFromGame(fileSrc, fileNewSrc):
                return False
        # is file, replace content
        if os.path.isfile(fileSrc):
            # open file
            fileObj = open(fileSrc, "r+", encoding='UTF-8')
            newFileData = ''
            try:
                fileData = fileObj.read()
                # get new data
                newFileData = GetNewContent(fileData)
            finally:
                fileObj.close()
            # save new file
            fileNewObj = open(fileNewSrc, "w+", encoding='UTF-8')
            try:
                # save new file data
                fileNewObj.write(newFileData)
            finally:
                fileNewObj.close()
    return True

# get new content
def GetNewContent(content):
    global ReplaceCount, ReplaceTotal
    # replace content
    ReplaceCount = 0
    newFileData = ReplaceContent(content)
    # strip
    newFileData2 = newFileData.lstrip()
    newFileData2 = newFileData2.rstrip()
    # get english data
    if AllowWaitEnglish:
        GetEnglishWord(newFileData2)
    # search location
    if AllowWaitLocation:
        GetLocationWord(newFileData)
    print('replace count: ' + str(ReplaceCount) + ', total: ' + str(ReplaceTotal))
    return newFileData

# replace content
# A: src language
# B: dest language
def ReplaceContent(content):
    global sheet, ReplaceCount, ReplaceTotal, NowReplaceFileSrc
    # range rows
    for key in range(0, len(WaitReplace)):
        if content.find(WaitReplace[key]['src']) == -1:
            continue
        content = content.replace(WaitReplace[key]['src'], WaitReplace[key]['dest'])
        ReplaceCount += 1
        ReplaceTotal += 1
        # print('replace file: ' + NowReplaceFileSrc + ', src: ' + WaitReplace[key]['src'] + ', new: ' + WaitReplace[key]['dest'])
    return content

# wait english pattern list
WaitEnglishWordPattern = [r"([A-Z]+)([\w \,\']+...)\."]
# check AllowWaitEnglish
def GetEnglishWord(content):
    global WaitEnglishData, WaitEnglishWordPattern, WaitEnglishSkip
    englishDataCount = 0
    for patternVal in WaitEnglishWordPattern:
        pattern = re.compile(patternVal)
        englishData = pattern.findall(content)
        if not englishData:
            continue
        englishDataCount = len(englishData)
        for mathKey in range(0, len(englishData)):
            if not englishData[mathKey]:
                continue
            isFind = False
            for mathFindKey in range(0, len(WaitEnglishSkip)):
                vFind = re.match(WaitEnglishSkip[mathFindKey], englishData[mathKey])
                if vFind:
                    isFind = True
                    break
            if isFind:
                continue
            englishData[mathKey] = englishData[mathKey].lstrip()
            englishData[mathKey] = englishData[mathKey].rstrip()
            if WaitEnglishData:
                for waitEnglishKey in range(0, len(WaitEnglishData)):
                    if WaitEnglishData[waitEnglishKey] == englishData[mathKey]:
                        isFind = True
                        break
                if isFind:
                    continue
            WaitEnglishData.append(englishData[mathKey])
    if englishDataCount > 0:
        print('find englishDataCount: ' + str(englishDataCount))
    return True

# get location word
def GetLocationWord(content):
    global WaitLocationTags
    pattern = re.compile(r'(\[{2})([A-Za-z ]+)\|')
    locationData = pattern.findall(content)
    locationDataLen = 0
    if locationData:
        locationDataLen = len(locationData)
        for locationKey in range(0, len(locationData)):
            if len(locationData[locationKey]) > 1:
                WaitLocationTags.append(locationData[locationKey][1])
        # print(locationData)
    if locationDataLen > 0:
        print('find locationDataLen: ' + str(locationDataLen))

# run
# copy src to dest dir
if not os.path.exists(SrcDir):
    if not CopyDir(DestDir, SrcDir):
        SendError('cannot copy game dir to src dir.')
    print('copy game dir to game-src.')
# delete game
if not DeleteDir(DestDir):
    SendError('cannot delete game dir.')
# search dirs
print('start replace lang data...')
SearchListFromGame(SrcDir, DestDir)
# save WaitEnglishData
fileObj = open('./languages/out.txt', "w+", encoding='UTF-8')
fileData = ''
for key in range(0, len(WaitEnglishData)):
    fileData = str(WaitEnglishData[key]) + '\n' + fileData
try:
    # save new file
    fileObj.write(fileData)
finally:
    fileObj.close()
print('find english word count: ' + str(len(WaitEnglishData)))
# save location tag
fileObj = open('./languages/location.txt', "w+", encoding='UTF-8')
fileData = ''
for key in range(0, len(WaitLocationTags)):
    fileData = str(WaitLocationTags[key]) + '\n' + fileData
try:
    # save new file
    fileObj.write(fileData)
finally:
    fileObj.close()
print('find location tag count: ' + str(len(WaitLocationTags)))