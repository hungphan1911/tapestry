# Kattis

A dedicated service for tracking Kattis problem-solving progress and making Kattis problem sets easier to discover.

Kattis is great for competitive programming, especially ICPC contests, but its built-in tooling for tracking long-term progress and discovering recent problem sets is fairly limited. In particular, there is no convenient way to search through available problem sets or get a consolidated view of personal statistics, which isn't quite ideal for individual practices.

This service exists to fill those gaps. There are a couple goals in mind when designing this service:
- Tracking problem-solving progress over time
- Visualize Kattis activity and performance
- Discovering and searching available problem sets

Conceptually, the architecture should look something like this:
```text
                           +----------------------+
                           |        Kattis        |
                           +----------+-----------+
                                      |
                                      | scrape / sync
                                      v
                           +----------------------+
                           |   Kattis Worker /    |
                           |       Scraper        |
                           +----------+-----------+
                                      |
                                      | write / update
                                      v
                           +----------------------+
                           |   Kattis Database    |
                           +----------+-----------+
                                      ^
                                      |
                                      | read / write
                                      |
                           +----------+-----------+
                           |      Kattis API      |
                           +----------+-----------+
                                      ^
                                      |
                                      | HTTP
                                      |
                           +----------+-----------+
                           |    Main Frontend     |
                           +----------------------+
```