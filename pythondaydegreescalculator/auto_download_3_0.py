# ------------------------------------------------------------------------------------
#               Weather Data
#-------------------------------------------------------------------------------------
import urllib.request
import webbrowser
import csv
import math

# General Format of the link: http://www.bom.gov.au/climate/dwo/ 201710 /text/IDCJDW5010. 201710.csv {year 2017; month 10; station code: IDCJDW5010}

# User Inputs:
#selected_station    = input('Enter a station from the list above: ')
user_lat  = float(input('Enter the Latitude of the field: '))
user_long = float(input('Enter the Longitude of the field: '))
year      = input('Enter Year in the format YYYY: ')
month     = input('Enter Month in the format MM: ')



#Convert the user input to radians:

user_lat_rad    = math.radians(user_lat)
user_long_rad   = math.radians(user_long)


# Import list of station names, codes and co-ordinates:

loaded_stationInfo = []
with open('stations_codes.csv', newline='') as inputfile:
    for row in csv.reader(inputfile):
        loaded_stationInfo.append(row)


# Find the index of all necessary index
location_index  = loaded_stationInfo[0].index('Location')
station_index   = loaded_stationInfo[0].index('Station')
code_index      = loaded_stationInfo[0].index('Code')
site_index      = loaded_stationInfo[0].index('Site')
lat_index       = loaded_stationInfo[0].index('Lat')
long_index      = loaded_stationInfo[0].index('Lon')
#id_index        = loaded_stationInfo[0].index('ID')


# Create the list of rows of station

rowsOfStation = list(range(1,len(loaded_stationInfo)))

# Segregate the station number in one list
statNo_list = []

for row1 in rowsOfStation:
    statNo_list.append(loaded_stationInfo[row1][site_index])


# Segregate the  name of all the  stations:
statName_list = []

for row2 in rowsOfStation:
    statName_list.append(loaded_stationInfo[row2][station_index])

# Segregate the name of all locations
location_list = []

for row3 in rowsOfStation:
    location_list.append(loaded_stationInfo[row3][location_index])


# Segregate the codes for all stations:
code_list = []

for row4 in rowsOfStation:
    code_list.append(loaded_stationInfo[row4][code_index])

# Segregare the list of ID's

#id_list = []

#for row7 in rowsOfStation:
#    id_list.append(loaded_stationInfo[row7][id_index])


# Segregate the latitude of the  stations

# Find which col is lat
statLat_list = []

for row5 in rowsOfStation:
    statLat_list.append(float(loaded_stationInfo[row5][lat_index]))


# Segregate the longitude of the station:
statLong_list = []

for row6 in rowsOfStation:
    statLong_list.append(float(loaded_stationInfo[row6][long_index]))


# Convert the longitutes and lattitudes into radians:
statLong_rad    = []
statLat_rad     = []


for row6 in range(0,len(statLong_list)):
    statLong_rad.append(math.radians(statLong_list[row6]))

for row7 in range(0,len(statLat_list)):
    statLat_rad.append(math.radians(statLat_list[row7]))

# Find the distance between the user input lat and long to all the lat and long in the database
del_lat     = [(user_lat_rad - this_lat) for this_lat in statLat_rad]
del_long    = [(user_long_rad - this_long) for this_long in statLong_rad]

phi = [this_delLat/2 for this_delLat in  del_lat]
lamb = [this_delLong/2 for this_delLong in del_long]


# Calculate the distance
a1  = [math.sin(this_phi)**2 for this_phi in phi]
a2  = [(math.cos(user_lat_rad) * math.cos(this_lat2)) for this_lat2 in statLat_rad]
a3  = [math.sin(this_lamb)**2 for this_lamb in lamb]

a   = [this_a1 + (this_a2*this_a3) for this_a1,this_a2,this_a3 in zip(a1,a2,a3)]

c   = [2 * math.atan2(math.sqrt(this_a),(math.sqrt(1-this_a))) for this_a in a]

R = 6371 * 10**3

distance = [R * this_c for this_c in c]

# Find the minimum in the distance list
min_distance = min(distance)
min_distance_ind = distance.index(min_distance)

closest_station = statName_list[min_distance_ind]

# Call min_distance another name
selected_ind = min_distance_ind

print(closest_station + ' is the closest weather station and is ' + str(round(min_distance)/1000) + ' KM away from the desired farm')


# Constant part of the link string for weather data:

w_strng_1 = 'http://www.bom.gov.au/climate/dwo/'
w_strng_2 = '/text/'
w_strng_3 = '.'
w_strng_4 = '.csv'

#Variable parts of weather data


# Combine the constant and variable part of the link to form a downloadable url

download_url_weather    = w_strng_1 + year + month + w_strng_2 + code_list[selected_ind] + w_strng_3 + year + month + w_strng_4


req = urllib.request.urlopen(download_url_weather)
webbrowser.open(download_url_weather)

