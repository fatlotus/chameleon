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

The above graph showed the performance of loading [this video][video], without
adblocking, for 24 seconds; this evaluation was done in Chrome, which
negotiates the `ECDHE-RSA-RC4-SHA` cipher suite.

[video]: https://www.youtube.com/watch?v=al2DFQEZl4M

As can be seen in the above data, there is still work to be done. In
particular, the code has a race condition when packets are sent faster than
500 Hz; as a result, packets sent faster than that are put into a queue, 
smearing the resulting packet stream somewhat. We also do not yet support 
actually sending data over the resulting stream.

The [Jupyter notebook][ipynb] for this example is also available for reference.
When capturing data, we used `tcpdump -i iface -w File.pcap`.

[ipynb]: http://nbviewer.ipython.org/github/fatlotus/chameleon/blob/master/Analysis.ipynb