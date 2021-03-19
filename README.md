# River

River is a tool for publishing data from data-base into RabbitMQ.

```
 _______ _________          _______  _______ 
(  ____ )\__   __/|\     /|(  ____ \(  ____ )
| (    )|   ) (   | )   ( || (    \/| (    )|
| (____)|   | |   | |   | || (__    | (____)|
|     __)   | |   ( (   ) )|  __)   |     __)
| (\ (      | |    \ \_/ / | (      | (\ (   
| ) \ \_____) (___  \   /  | (____/\| ) \ \__
|/   \__/\_______/   \_/   (_______/|/   \__/

```

## How to use it?

First download the appropriate binary with respect to your OS or compile it with makefile rule. Then run the program 
like this:

```shell
./river stream --db-uri <MONGODB_CONNECTION_URL> --db-name <MONGODB_DATABASE_NAME> --db-col <MONGODB_COLLECTION_NAME>
```
