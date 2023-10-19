# Voting

This is the voting system for HackSoc committee elections. It uses single transferrable vote (or more specifically, instant-runoff voting, which is the single-winner equivalent of STV), and is designed so that nobody can possibly detect who cast which ballot.

It is made of three web servers. Two are user-facing and the third is an admin panel for us to check the winners.

 - The registration server asks members to enter their student ID, verifies it against the database, and then generates them an empty ballot.
 - The ballot server lets members fill in their ballot and submit it to the database.
 - The admin panel server lets us see the winners at any point during the election, and also the amount of votes and ballots given out.

## Holding an Election

A few things need to be in place to run an election.

 - Each of the three servers should be running separately, ideally on different machines. The registration, ballot, and admin servers run on ports `:10000`, `:10001`, and `:10002` respectively.
 - A reverse proxy should be set up so that the three servers can be accessed at `/register`, `/ballot`, and `/admin`.
 - The MongoDB database credentials need to be provided as environment variables `MONGO_USER` and `MONGO_PASS` and `MONGO_HOST`.
   - There should be one database, "hacksoc", with four collections, `members`, `ballots`, `candidates`, and `members_voted`.
   - `members` should be a collection of all members with at least an `ID` field for the student ID.
   - `candidates` should be a collection of all committee positions to fill, e.g.
     ```json
     {
        "role": "President",
        "candidates": ["Candidate 1", "Candidate 2", "Candidate 3", "Re-open Nominations"],
        "index": 0
     }
     ```
    - The other two collections should start off empty.

## Contributing

I've added a VS code workspace file in the root of this repository to make it easier to work with in VS code at least - the VS code Go plugin is funny about multiple Go modules existing in the same workspace folder.