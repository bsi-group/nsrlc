nsrlc (NSRL Client)
===================

nsrlc (NSRL Client) is a command line application designed to work in conjunction
with an **nsrls** server instance. The client uses the servers HTTP API service 
to lookup large sets of hash data. 

By default the **nsrlc** application sends the hashes in batches of 1000. The 
number of hashes in a batch is configured via the **-b** command line parameter

## Output ##

The application outputs directly to a file. The output format can be defined by 
the command line parameters. The options for the output format are:
 
- i: Outputs only the identified hashes
- u: Outputs only the unidentified hashes
- a: Outputs both the identified and unidentified hashes, along with a status column