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

## Supported routes

The service exposes the following endpoints:

- `/problemsets`: This exposes the problem sets available on Kattis, with all of the metadata associated with the user. By problem set, I mean the problem sets available on [this page](https://open.kattis.com/problem-sources). This does not intend to pull every single problem on Kattis down but just the problem sources. The filter features is something to be done from the frontend side.

- `/stats/solved`: This endpoint exposes the number of problems the user has solved so far.

- `/stats/summary`: Return the scalar facts of the user. In particular, the number of solved problems, the current total scure, and average difficulty. 

- `/stats/difficulty-distribution`: This endpoint exposes the difficulty distribution among the solved problems. For bar chart visualizations.

- `/stats/activity`: Number of solved/submissions over time. This is for a GitHub-style heatmap

- `/submissions/recent`: This shows the recent submissions to Kattis

- `/problemsets/completed`: The list of all completed problem sets.

- `/dashboard`: Basically returns everything above, except for `/problemsets`