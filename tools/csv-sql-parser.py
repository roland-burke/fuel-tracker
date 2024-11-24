import csv

# This script reads a csv file and converts it to a sql script

columnNames = ["users_id", "description", "date_time", "price_per_liter_euro", "total_liter", "price_per_liter", "currency", "mileage", "license_plate"]
tableName = "refuel"
insertQuery = "INSERT INTO {} ({},\"{}\",{},{},{},{},\"{}\",\"{}\",{}) VALUES({},'{}','{} 0:0:0',{},{},0,'',{},'{}');\n"

content = ""

def getDate(date: str):
    return date[6:10] + "-" + date[3:5] + "-" + date[:2]

f = open("insertData.sql", "w")

with open('Spritutf16.csv', newline='', encoding='utf16') as csvfile:
    fuelreader = csv.reader(csvfile, delimiter=';', quotechar='|')
    for row in fuelreader:
        content = content + insertQuery.format(tableName, columnNames[0], columnNames[1], columnNames[2], columnNames[3], columnNames[4], columnNames[5], columnNames[6], columnNames[7], columnNames[8], 3, row[4], getDate(row[0]), row[2].replace(",", "."), row[1].replace(",", "."), row[3], "Kennzeichen")


f.write(content)
f.close()

print("done")
