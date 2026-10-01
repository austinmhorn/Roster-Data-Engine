import gspread
from oauth2client.service_account import ServiceAccountCredentials
import csv
import os
import time

# Define scope and credentials
scope = ["https://spreadsheets.google.com/feeds", "https://www.googleapis.com/auth/drive"]
creds = ServiceAccountCredentials.from_json_keyfile_name('credentials.json', scope)
client = gspread.authorize(creds)

# Open the main Google Sheet
spreadsheet = client.open("Roster Data Feed")

# Get all CSV files in the current directory
csv_files = [f for f in os.listdir('.') if f.endswith('.csv')]

def retry_google_api(func, *args, max_retries=5, **kwargs):
    """Retry temporary Google API failures using exponential backoff."""

    for attempt in range(1, max_retries + 1):
        try:
            return func(*args, **kwargs)

        except gspread.exceptions.APIError as error:
            status_code = error.response.status_code

            # Only retry temporary/server-side errors
            if status_code not in (429, 500, 502, 503, 504):
                raise

            if attempt == max_retries:
                print(
                    f"❌ Google API still unavailable after "
                    f"{max_retries} attempts."
                )
                raise

            wait_time = 2 ** attempt

            print(
                f"⚠️ Google API returned {status_code}. "
                f"Retrying in {wait_time} seconds "
                f"({attempt}/{max_retries})..."
            )

            time.sleep(wait_time)

for csv_file in csv_files:
    sheet_name = os.path.splitext(csv_file)[0]  # Remove .csv extension
    print(f"Uploading {csv_file} to Google Sheet: {sheet_name}")
    
    # Try to open the specific worksheet, create if it doesn't exist
    try:
        sheet = spreadsheet.worksheet(sheet_name)
    except gspread.exceptions.WorksheetNotFound:
        print(f"Worksheet '{sheet_name}' not found. Creating a new one.")
        sheet = spreadsheet.add_worksheet(title=sheet_name, rows="1000", cols="40")
    
    # Read the CSV file
    with open(csv_file, 'r', encoding='utf-8') as csvfile:
        reader = csv.reader(csvfile)
        data = list(reader)
    
    # Clear and resize the existing sheet
    retry_google_api(sheet.clear)

    retry_google_api(
        sheet.resize,
        rows=len(data),
        cols=max(len(row) for row in data)
    )

    retry_google_api(
        sheet.update,
        "A1",
        data,
        value_input_option="USER_ENTERED"
    )

    print(f"{csv_file} synced to Google Sheet '{sheet_name}' successfully.")
