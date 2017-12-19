# JailBot ques
Holds global singleton que system. This is useful for handling alerts.

Whenever someone is banned, add it to the alerts que, and broadcast the message to every server out there. The element used to hold a report should be a structure that also keeps track of the number of servers notified(?) so when the total amount is met (all servers handled) the element will be removed from the stack. This should be threadsafe as services runs on another thread.
