import sqlite3 as sq
import argparse

# Add commandline arguments of:
# specify target sqlite3 database name, default to database.db
# specify whether to enable testing mode, default to False
parser = argparse.ArgumentParser()
parser.add_argument('-d', '--database', default='database.db', help='specify target sqlite3 database name, default to database.db')
parser.add_argument('-t', '--testing', default=False, help='specify whether to enable testing mode, default to False')
args = parser.parse_args()

# Connect to the database
conn = sq.connect(args.database)

# Create voters table schema with id, unique name, grouptype, public_key
conn.execute('''CREATE TABLE IF NOT EXISTS voters (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE,
    grouptype TEXT,
    public_key TEXT
)''')

if (args.testing):
    pass

# Close db connection
conn.close()
