# sockmapd

**sockmapd** is a [**socketmap**](http://www.postfix.org/socketmap_table.5.html) server daemon for Mail Transfer Agents
(MTA) that uses configurable database repositories for [lookup tables](http://www.postfix.org/DATABASE_README.html).

**sockmapd** is written in Go Lang

## Purpose
Mail Transfer Agents like Sendmail and [Postfix](http://www.postfix.org) can use either hash files, database queries, or
[socketmaps](http://www.postfix.org/socketmap_table.5.html) to lookup for valid recipients' eMail addresses,
hosts or eMail addresses in blacklists or whitelists, or even header policies, just to mention a few examples.

Often we have this information stored in database servers that are not directly accessible from our DMZ where the MTAs
are usually deployed. Although some MTAs (like Postifx) do support direct database access and queries, we want to keep
access to our internal database well protected. 

This is the case, for example, if we use [dbmail](https://github.com/dbmail/dbmail), an IMAP, POP server. **dbmail**
maintains mailboxes and their eMail aliases on the protected database server, but the MTA in the DMZ is prevented 
direct access to them. If we want our MTA to validate a recipient and reject, or accept the eMail for further delivery
at the time it is contacted by the sender's MTA, then we need a service like sockmapd, that takes queries from the MTA
and verifies whether the recipient is valid. The same applies when we want our MTA to check if a sender's eMail address
or the sender MTA's host name is in our blacklist.


## Description 
**sockmapd** is a configurable service daemon that accepts [socketmap](http://www.postfix.org/socketmap_table.5.html)
protocol queries (TCP protocol) on a specified TCP Port number and returns a response to the calling service client.
Sockmapd queries database tables that are mapped to a service map via a configuration file that can be specified in the
command line.

## Feature: socketmap request/response
Socketmap request queries have the format of a [netstring](http://cr.yp.to/proto/netstrings.txt) as follows:

    "[len]:[query],"
    
Example: 

if my request is to lookup for **someone@somewhere.com** in the **recipient** map, the request string will be:

    "32:recipient someone@somewhere.com,"
    
where 32 is the length of "recipient somebody@somewhere.com"
in front of this request, **sockmapd** queries the database table mapped to the socketmap in the configuration file
and returns a response as follows:

    "9:NOTFOUND ,"
or

    "24:OK someone@somewhere.com,"

The response strings are documented in the [socketmap table](http://www.postfix.org/socketmap_table.5.html) man page.

## Feature: configuration and database mapping
Database mapping is achieved through a json format configuration file. The configuration file includes 3 sections:

**sysconfig** includes the tcp port number and the path for an alternative log file.

**database** includes host, port, username, and password, these are common configuration parameters for the connection
to a database server. Note that **host** is a string array that can contain multiple nodes of a database cluster 
(or replication nodes). **sockmapd** initializes a database connection for each database trying to connect to each node
until a successful connection is achieved. In case the database connection fails when a request is being processed,
**sockmapd** tries to re-initialize the database connections.

**postmaps** is an array of an object including service, database, table, key, value and reason. **sockmapd** maintains
a database connection for each database specified (as we can have different databases holding map data), and creates
the query using the **key** column and taking the result from the **value** column (if specified, otherwise the query
checks only that **key** is present). If specified, the content of **reason** is added to the response.


