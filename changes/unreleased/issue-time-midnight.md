## Fixed

- `bill`: an `issue_time` of midnight (`00:00:00`) is no longer treated as unset. Previously the calculators replaced it, along with the `issue_date`, with the current date and time. The issue time is now always kept as provided; only a missing `issue_date` is still filled in with today's date.
