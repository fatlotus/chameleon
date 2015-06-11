## Chameleon

This project is a proof of concept, demonstrating the capability to make one
TCP stream look like another. Once finalized, Chameleon will likely live inside
one of the Tor pluggable transports, enabling protection against particular
types of HTTP timing side channels.

As a research prototype, this code is public domain.

### Evaluation

When emulating a stream of a user watching a YouTube video, Chameleon is able
to reproduce the packet sizes and times with reasonably good accuracy:

![Performance graph of Chameleon versus native YouTube](https://raw.githubusercontent.com/fatlotus/chameleon/master/Evaluation.png)

The above graph showed the performance of loading [https://www.youtube.com/watch?v=al2DFQEZl4M], without adblocking, for 24 seconds; this evaluation was done in Chrome, which negotiates the `ECDHE-RSA-RC4-SHA` cipher suite.