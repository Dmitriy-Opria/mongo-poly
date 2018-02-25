
#Monthly methorological csv files. Example files here are from July to October 2017
CSVs_files<-list.files("C:/Flurosat/Modelling/Weather_download/Data", "csv$",full.names = T)

#Aggregate the separate monthly files into one file

whole_file<-""#name of new big file
counter<-1
for(csv_file in CSVs_files){
 file<-read.csv(csv_file,skip=5, header=T)[,2:5]#subset useful columns only
 if(counter==1){
   whole_file<-file}
 else{
   whole_file<-rbind(whole_file,file)
 }
 
 counter<-counter+1
}


#User to supply sowing date from which to start calculating the growing degree days (GDD)
#Format is in yyyy, mm,dd
sowing_date<-readline("Enter sowing date as yyyy-mm-d:")                 #   

#User should supply a point in time in the growing season to which the GDD will be commulated to

in_season_date<-readline("Enter in season date as yyyy-mm-d:")

#Use a  regular expression to obtain the sowing and in-season dates row index 
#from the date column of the combined metrological data
sow_date_row_index<-grep(paste("^",sowing_date,"$",sep=""),whole_file$Date)
in_season_row_index<-grep(paste("^",in_season_date,"$",sep=""),whole_file$Date)

#Subset the combined files using the row index of sowing and in-season dates in the combined file

subset_meteorological_data<-whole_file[sow_date_row_index:in_season_row_index,]

#Obtain the minimum and maximum temperature values in the data subset

Min_temperature<-subset_meteorological_data$Minimum.temperature...C.

Max_temperature<-subset_meteorological_data$Maximum.temperature...C.

#Iterate through the minimum and maximum temperature values to computed cummulated growing degree days

Commulative_day_degrees_values<-c()
for(i in 1:length(Min_temperature)){
  Tmin<-Min_temperature[i]
  if(Tmin<12){
    Tmin<-12
  }
  
  Tmax<-Max_temperature[i]
  
  #Compute day degrees
  DD<-((Tmax-12)+(Tmin-12))/2
  
  #calculate the cummulative day degrees and store in a vector (python-list equivalent)
  if(i==1){
    Commulative_day_degrees_values<-c(Commulative_day_degrees_values,DD)
    
  }
  
  else{
    Commulative_day_degrees_values<-c(Commulative_day_degrees_values,DD+Commulative_day_degrees_values[i-1])
  }
  
}


#Add Cummulative day degrees to the subset of the meterological data 

subset_meteorological_data$Cummulative_DD<-Commulative_day_degrees_values


#Output file to disk

write.csv(subset_meteorological_data,"C:/Flurosat/Modelling/Weather_download/Data/SA_winter_met.csv",row.names = F)





