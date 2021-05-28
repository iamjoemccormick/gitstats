# Overview

GitStats runs on-demand or periodically and pulls traffic and potentially other statistics from GitHub, writing them into a time series database. It will intelligently merge existing and updated data as follows: 

* Clones and page views are tracked on a daily basis, existing entries for a given day will be updated as needed. 
  * Note: While in theory associated statistics should only ever increase, the GitHub API is treated as the source of truth and will not prevent decreases.
* GitHub does not track top referral paths and sources over time, so GitStats will log (in UTC+0) when these statistics were captured.

# Reference: 

* [GitHub REST API Docs: Traffic](https://docs.github.com/en/rest/reference/repos#traffic).
* [GitHub Docs: Viewing Traffic to a Repo](https://docs.github.com/en/github/visualizing-repository-data-with-graphs/accessing-basic-repository-data/viewing-traffic-to-a-repository#accessing-the-traffic-graph).