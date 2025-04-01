# backend package
Handles operations on User's OS .
For example - Executing shell commands, doing files operations on user's files.
Reading OS specific configurations like disk partitions.

Name 'backend' isn't the most appropriate, open to suggestions.

This would modularize the code, and would enable us to write unit tests 
where we would 'mock' the backend functionality with dummy interface 
implementations

# Dependencies
Should not import any "ui" package
Can import common and its subpackages

# Implementation specifications
Implement everything via interfaces, so that we can easily write unit tests
