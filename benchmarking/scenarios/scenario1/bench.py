import os
import logging
import MySQLdb
import time
import progressbar
import sys
import Queue
import threading

createUserSQL   = "INSERT IGNORE INTO users (name) VALUES (%s);"
getUserByUsernameSQL   = "SELECT * FROM users WHERE name=%s;"
getAuthorByNameSQL   = "SELECT * FROM authors WHERE name=%s;"
createAuthorSQL = "INSERT IGNORE INTO authors (userid, name) VALUES (%s, %s);"
createBookSQL = "INSERT IGNORE INTO books (name, author, price) VALUES (%s, %s, %s);"

firstNames = sorted(["Kenia ", "Randal", "Shawnna ", "Rey ", "Cordia", "Kendal",
    "Alina", "Dianna", "Misti", "Chelsie", "Gracia", "Teena", "Ronny", "Willy",
    "Betsy", "Kenisha", "Elsy", "Cheryle", "Lurline ", "Karina", "Luba", "Vita",
    "Lu", "Frances", "Lavenia", "Nereida", "Zetta", "Melony", "Eloise",
    "Nickolas", "Ericka", "Cecilia", "Jenni", "Sofia", "Nobuko", "Trudy",
    "Petronila", "Donnette", "Santos", "Viola", "Jessika", "Chere", "Azalee",
    "Meggan", "Floyd", "Liberty", "Tabitha", "Juliana", "Pamila", "Blondell"])

lastNames = sorted(["Watterson", "Lawler", "Walt", "Birch", "Bryd", "Speight",
    "Monroy", "Milledge", "Davilla", "Behrendt", "Mustain", "Blythe", "Gandhi",
    "Brady", "Gooden", "Jellison", "Hager", "Selders", "Seaton", "Wind",
    "Jelinek", "Reiser", "Lacour", "Maginnis", "Baggs", "Crossno", "Shadley",
    "Bramer", "Mento", "Manigault", "Jacobi", "Deckman", "Spikes", "Duncan",
    "Ackman", "Hornick", "Bourbeau", "Riehl", "Sena", "Rolon", "Pereira",
    "Mikula", "Luk", "Albaugh", "Akin", "Bradburn", "Houlihan", "Frisina",
    "Funnell", "Keister"])

def connect():
    return MySQLdb.connect(host="mysql",    # your host, usually localhost
                         user="root",         # your username
                         passwd="password",  # your password
                         db="bench")        # name of the data base

createUserThreads = []
def createUsers(name):
    logging.debug("Creating... "+name)
    sys.stdout.flush()
    db = connect();
    cur = db.cursor()
    for j in lastNames:
        for k in range(0, 10):
            myname = name + " " + j + "(" + str(k) + ")"
            sys.stdout.flush()
            cur.execute(createUserSQL, (myname,))
            cur.execute(getUserByUsernameSQL, (myname, ))
            row = cur.fetchone()
            if not row == None:
                cur.execute(createAuthorSQL, [str(row[0]), ("Author "+myname)])
            else:
                print "Could not create ", myname

    db.commit()
    db.close()
    logging.debug("Created! "+name)
    sys.stdout.flush()

createBookThreads = []
def createBook(username):
    logging.debug("Creating books... "+username)
    sys.stdout.flush()
    db = connect()
    cur = db.cursor()
    for j in lastNames:
        for k in range(0, 3):
            myname = "Author " + username + " " + j + "(" + str(k) + ")"

            cur.execute(getAuthorByNameSQL, (myname, ))
            row = cur.fetchone()
            if not row == None:
                for i in range(0,2):
                    bookname = myname+"'s book "+str(i)
                    cur.execute(createBookSQL, [bookname, str(row[0]), i * 5])

            else:
                print "Could not find ", myname

    db.commit()
    db.close()
    logging.debug("Created books! "+username)
    sys.stdout.flush()

def initilizeUsers():
    logging.debug("Initilizing users...")
    start = time.time();
    for i in firstNames:
        name = i + " " + hostname
        t = threading.Thread(target=createUsers, args = (name, ))
        t.daemon = True
        createUserThreads.append(t)

    # Start all the threads
    for x in createUserThreads:
        x.start()

    # Wait for them to complete
    for x in createUserThreads:
        x.join()

    # Return the time it took to run
    logging.debug("Creating users took: "+str(time.time() - start))
    return time.time() - start;

def initilizeBooks():
    logging.debug("Initilizing books...")
    start = time.time();
    for i in firstNames:
        name = i + " " + hostname
        t = threading.Thread(target=createBook, args = (name, ))
        t.daemon = True
        createBookThreads.append(t)

    # Start all the threads
    for x in createBookThreads:
        x.start()

    # Wait for them to complete
    for x in createBookThreads:
        x.join()

    # Return the time it took to run
    logging.debug("Creating books took: "+str(time.time() - start))
    return time.time() - start;

def main():
    logging.debug("Starting...")
    db = connect();
    intUserTime = initilizeUsers();
    intBookTime = initilizeBooks();

    # cur.execute("SELECT * FROM users")
    # # print all the first cell of all the rows
    # for row in cur.fetchall():
    #     logging.debug(row[1])
    #
    # cur.execute("SELECT * FROM authors")
    # # print all the first cell of all the rows
    # for row in cur.fetchall():
    #     logging.debug(row[2])
    # db.close()

    logging.info("[Create Users] " + str(intUserTime))
    logging.info("[Create Books] " + str(intBookTime))


if __name__ == '__main__':
    hostname = os.uname()[1]
    logging.basicConfig(format=hostname + ' %(asctime)s %(levelname)s: %(message)s', level=logging.DEBUG)
    main()
