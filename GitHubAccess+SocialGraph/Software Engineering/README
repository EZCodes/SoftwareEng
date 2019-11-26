Social Graph And Github Access Project

This folder contains backend and front end for this project.

A bit about implementation.

BackEnd:
BackEnd first connects to MongoDB and GitHub API using APIs.
Then all repositories for Google and Microsoft are fetched. After that program splits into 4 threads, 
while main thread will wait for completion of them. Two threads are fetching Google and Microsoft committs respectively, for every repository
that we fetched. Therefore, it will take at least a month, so my fetch is probably up on AWS cloud machine atm.
Because of that, every commit we fetch is processed, so that we get only the data that we require and then this data is passed
to uploader thread, that updates our data on MongoDB. Forth thread is caching thread. Since, in order to update our information on MongoDB,
I delete the old and insert the new entry, for stablility of FrontEnd, I created the caching thread, which works on separate collection
and updates the information only once per hour. Fetching process will not break on some of the expected errors and will wait an hour if it
hits the Rate Limit of github API.

FrontEnd:
FrontEnd is independent from BackEnd. It creates the client for MongoDB and fetches the cached data from there. After that it processes it
and creates a Struct that can be converted into JSON that is accepted by a graph. That JSON is written to a file, that is imported by chart.js
Graph is made using D3js library and is hosted on the web server, executed by frontEnd.go on localhost with socket 8080.

That's mostly it about implementation.
For a bit more explanation and demo of it working watch the demo video.

SETUP:
In case you are interested in setting it up yourself.

Prerequisites: 
1. Github Personal Token - for fetching information from github.
2. MongoDB Cluster - for stoing the data.
3. GoLang installed

After that you can download my project and execute them in go.
BackEnd and FrontEnd are executed separately. For login details provide text file/s with credentials and tokens (unless leaking it with the code is ok for you)
Once BackEnd starts, after a bit, initial information will be uploaded and can be used by FrontEnd.
As I mentioned, it's a lot of information, so it will fetch for a long time for full metrics, but FrontEnd will work with partial ones.
Once FrontEnd is executed, access localhost:8080/static to see the results.


