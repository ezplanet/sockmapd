# sockmapd

**sockmapd** is a [**socketmap**](http://www.postfix.org/socketmap_table.5.html) server daemon for Mail Transfer Agents
(MTA) that uses configurable database repositories for [lookup tables](http://www.postfix.org/DATABASE_README.html).

**sockmapd** is written in Go Lang

## Purpose
Mail Transfer Agents like Sendmail and [Postfix](http://www.postfix.org) can use hash files or sockemaps to lookup
for valid eMail addresses recipients or eMail address and hosts in blacklists and whitelists or header policies just to 
mention a few examples.

Often we have this information stored in database servers that are not directly accessible from our DMZ where the MTAs
are usually deployed, neither, though possible, we want to give our MTA access to query directly our Databases as it
would defeat the purpose of protecting them behind a firewall.

This is the case, for example, if we use [dbmail](https://github.com/dbmail/dbmail), an IMAP, POP server. In this case
mailboxes and their eMail aliases are maintained on the backend database server, but the MTA in the DMZ does not know
about them. If we want our MTA to validate a recipient and reject/accept the eMail for further delivery when contacted
by the foreign MTA, then we need an agent like sockmapd, that takes queries from the MTA and verifies whether the
recipient is valid.

## Description 
**sockmapd** is a configurable service daemon that accepts socketmap protocol queries (TCP protocol) on a specified
TCP Port number and returns a response to the calling service client. Sockmapd queries database tables that are mapped
to a service map via a configuration file that can be specified in the command line.

## Feature: socketmap request/response
Socketmap request queries have the format of a [netstring](http://cr.yp.to/proto/netstrings.txt) as follows:

    [len]":"[query]","
    
Example: 

if my request is for **someone@somewhere.com** in the **recipient** map, the request string will be:

    "32:recipient someone@somewhere.com,"
    
where 32 is the length of "recipient somebody@somewhere.com"
in front of this request, **sockmapd** queries the database table mapped to the socket map in the configuration file
and returns a response as follows:

    "9:NOTFOUND ,"

The response strings are documented in the [socketmap table](http://www.postfix.org/socketmap_table.5.html) man page.

## Feature: configuration and database mapping
Database mapping is achieved through a json format configuration file. The configuration file includes 3 sections:

**sysconfig** includes the tcp port number and the path for an alternative log file.

**database** includes host, port, username, and passord, these are common configuration parameter for the connection
to a database server. Note that **host** is a string array that can contain multiple nodes of a database cluster 
(or replication nodes). **sockmapd** will retry the connection to all the specified database nodes in case of a
connection failure.

**postmaps** is an array of an object including service, database, table, key, value and reason. **sockmapd** will
maintain a database connection for each database specified (as we can have different databases holding map data), and
will create the query using the **key** column and taking the result from the **value** column (if specified, otherwise
the query checks only that **key** is present). If specified, the **reason** is added to the response.


