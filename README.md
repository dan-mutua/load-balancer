 # SimpleLB

SimpleLB is a lightweight, efficient load balancer designed to distribute incoming requests across multiple backend servers using a RoundRobin algorithm. This ensures even distribution of traffic and supports retries for failed requests, enhancing overall reliability.

In addition to load balancing, SimpleLB actively monitors backend health and performs necessary cleanup and recovery operations to maintain optimal performance. It assumes a backend is available if its root URL (/) is reachable.

## Features

- RoundRobin Load Balancing: Distributes incoming requests evenly across backend servers.
- Retry Mechanism: Automatically retries requests if a backend fails.
- Health Checks: Regularly pings backends to ensure they are available and recovers them when they are back online.
- Assumption of Availability: Considers a backend available if the root path (/) is reachable.


# How to use
```bash
Usage:
  -backends string
        Load balanced backends, use commas to separate
  -port int
        Port to serve (default 3030)
        
```

Example:

To add followings as load balanced backends
- http://localhost:3031
- http://localhost:3032
- http://localhost:3033
- http://localhost:3034
```bash
simple-lb.exe --backends=http://localhost:3031,http://localhost:3032,http://localhost:3033,http://localhost:3034
```