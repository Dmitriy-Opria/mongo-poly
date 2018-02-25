# -*- coding: utf-8 -*-
"""
Calculate the degree days from input csv, return as csv
"""

import pandas as pd
import csv
import re
import os

def getFileList(input_location):
    """ return a list of files in a directory with name format: yyyy_mm_Data.csv"""
    out = []
    for file in os.listdir(input_location):
        if re.fullmatch(r'\d{4}_\d{1,2}_Data\.csv', file):
            out.append(file)
    return(out)

def readCSVData(fileName):
    """ read data from a BOM input CVSfile, columns: Date, Min Temp, Max Temp"""
    with open(fileName, 'r') as f:
        reader = csv.reader(f, dialect='excel')
        outRows=[]
        for row in reader:
            try:
                if re.fullmatch(r'\d{4}-\d{1,2}-\d{1,2}', row[1]):
                    outRows.append(row[1:4])
            except IndexError:
                pass
    return outRows

def appendCSVData(outFile, data):
    """ append a list of row data to an existing CSV file"""
    with open(outFile, 'a') as f:
        writer = csv.writer(f, dialect='excel')
        writer.writerows(data)

def mergeCSVData(fileList, outFile):
    """ merge a list of CSVs into a single CSV"""

    # reset file contents
    with open(outFile, 'w') as f:
        writer = csv.writer(f, dialect='excel')
        writer.writerows([])
    
    # read each CSV then append to outFile
    for inFile in fileList:
        data = readCSVData(inFile)
        appendCSVData(outFile, data)

def getCSVData(fileName):
    """
    retrieve the relevant csv data
    """
    df = pd.read_csv(fileName,
             header=None,
             names='Date minT maxT'.split(),
             index_col=0,
             infer_datetime_format=True,
             skipinitialspace=True,
             parse_dates=True,    
             )
    df = df.sort_index(ascending=True)
    return df

def makeDegDay(df):
    """
    return a dataframe, with additional "cum. degree days" column
    degree days = (min temp - 12) + (max temp - 12) / 2
    """ 
    df['degDay'] = df.apply(degDayRow, axis =1)
    df['cumSum'] = df['degDay'].cumsum()
    
    return df


def degDayRow(row):
    """ calculate the degree day for a single DataFrame row """
    minT = row['minT']
    maxT = row['maxT']
    
    if minT < 12:
        minT = 12
    d =  (minT-12 + maxT-12)/2
    if d > 0:
        return d
    return 0



if __name__ == '__main__':
    
    start_date = input("Enter sowing date as yyyy-mm-d:")
    end_date = input("Enter in season date as yyyy-mm-d:")
    input_location = input("Enter filepath to csv files (default ./):") or './'
    output_file = input("Enter filename for output (default ./outFile.csv):") or 'outFile.csv'

    fileList = getFileList(input_location)
    
    mergeCSVData(fileList, output_file)
    df = getCSVData(output_file)
    df = makeDegDay(df)
    
    mask = (df.index >= start_date) & (df.index <= end_date)
    df = df.loc[mask]
    df.to_csv(output_file,
              header = ['Minimum temperature (C)',
                         'Maximum temperature (C)',
                         'Degree day',
                         'Cumulative degree days']
              )


    
