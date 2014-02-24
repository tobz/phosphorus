[![Coverage Status](https://coveralls.io/repos/tobz/phosphorus/badge.png?branch=master)](https://coveralls.io/r/tobz/phosphorus?branch=master) [![Build Status](https://drone.io/github.com/tobz/phosphorus/status.png)](https://drone.io/github.com/tobz/phosphorus/latest)

phosphorus
===========

A Dark Age of Camelot emulator written in Go.  


Details
===========

Phosphorus strives to stay close to idiomatic Go.  Code is split into multiple packages with loose coupling provided by interfaces.  There's extensive use of goroutines and channels, building in concurrency from the get go.

Phosphorus also strives for nearly-zero-tolerance error handling.  Client errors are bubbled up and aren't buried or thrown away.  This helps us quickly discover bugs and avoid broken behavior being left lying around, only to come up and bite us in the ass later.

The only external dependency for Phosphorus is a database.  The database layer supports SQLite 3, Postgres and MySQL, so you can get off the ground running right after compiling Phosphorus.  Phosphorus will create the database schema no matter what so you need only set up a database and the proper grants - no manual SQL hndling required.

Phosphorus also supports direct integration with InfluxDB for detailed runtime metrics.  You can see Go runtime statistics and Phosphorus statistics, such as network throughput, timer performance and more.
