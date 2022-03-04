### Thoughts while porting the flag_tool
- Now that we have seperated flag_tool from the database. 
  - Is it possible and does it makes sense to send the tweets as a stream, fetching and sending all the tweets in one batch for the -i flag will take a long time
  - Right now anyone with this tool can flag messages, is this an issue? we discussed adding the end points to an authorized group. 
  - copy pasting the message and user struct from minitwit. Doesn't seem like the best solution.
 
