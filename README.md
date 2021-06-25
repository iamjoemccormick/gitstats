# Overview

GitStats runs on-demand or periodically and pulls traffic and potentially other statistics from GitHub, writing them into a time series database. It will intelligently merge existing and updated data as follows: 

* Clones and page views are tracked on a daily basis, existing entries for a given day will be updated as needed. 
  * Note: While in theory associated statistics should only ever increase, the GitHub API is treated as the source of truth and will not prevent decreases.
* GitHub does not track top referral paths and sources over time, so GitStats will log (in UTC+0) when these statistics were captured.

# Reference: 

* [GitHub REST API Docs: Traffic](https://docs.github.com/en/rest/reference/repos#traffic).
* [GitHub Docs: Viewing Traffic to a Repo](https://docs.github.com/en/github/visualizing-repository-data-with-graphs/accessing-basic-repository-data/viewing-traffic-to-a-repository#accessing-the-traffic-graph).

# InfluxDB Interactions 

For each repository GitStats will query all supported endpoints writing the results into InfluxDB using [line protocol](https://docs.influxdata.com/influxdb/v2.0/reference/syntax/line-protocol/).

## Supported Endpoints

### traffic/clones

Measurement: LifetimeClones
* Tag Set:
  * Repository: String
* Field Set: 
  * Count: Integer
  * Uniques: Integer
* Timestamp: Generated

Measurement: DailyClones
* Tag Set:
  * Repository: String
* Field Set: 
  * Count: Integer
  * Uniques: Integer
* Timestamp: Inherited

### traffic/popular/paths

Measurement: VisitedPaths
* Tag Set:
  * Repository: String
  * Path: String
  * Title: String //Note: this can be useful for interpreting dynamic content, for example understanding what people are searching for.
* Field Set: 
  * Count: Integer
  * Uniques: Integer
* Timestamp: Generated

### traffic/popular/referrers

Measurement: Referrers
* Tag Set:
  * Repository: String
  * Referrer: String
* Field Set: 
  * Count: Integer
  * Uniques: Integer
* Timestamp: Generated

### traffic/views

Measurement: LifetimePageViews
* Tag Set:
  * Repository: String
* Field Set: 
  * Count: Integer
  * Uniques: Integer
* Timestamp: Generated

Measurement: DailyPageViews
* Tag Set:
  * Repository: String
* Field Set: 
  * Count: Integer
  * Uniques: Integer
* Timestamp: Inherited