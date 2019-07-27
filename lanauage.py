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
    global ReplaceCount, ReplaceTotal, NowReplaceFileSrc, WaitReplace
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
            fileNewObj = open(fileNewSrc, "w+", encoding='UTF-8')
            try:
                fileData = fileObj.read()
                # replace content
                ReplaceCount = 0
                newFileData = ReplaceContent(fileData)
                print('replace count: ' + str(ReplaceCount) + ', total: ' + str(ReplaceTotal))
                # save new file
                fileNewObj.write(newFileData)
            finally:
                fileObj.close()
                fileNewObj.close()
    return True

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