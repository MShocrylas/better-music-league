# Prototype Plan
## Ideal user flow
1. User "logs in" to their account (Just by a username right now, no auth/password)
2. Submit a song in plain text: Title, Artist. Maybe an optional link to the song (youtube or spotify or bandcamp or whatever)
3. View the "league" page with every user's submission
4. Rank the submissions from best to worst, and submit your vote
5. Once everyone has voted, you can view a results page

## For the prototype, one page
The voting page is probably the most important one, start with that
- A list of songs queried from the database for that round
- Draggable (or just setting numbers beside each element)
- Submit vote for a user

### Schema
Users(id, username)
Leagues(id, name)
Rounds(id, theme)
Songs(id, title, artist, url)
Submissions(id, user, round, song)
Votes(round_id, song_id, user_id, song_rank)

# TODO:
- Add `goose` for sql migrations management (only really important once you have a persistent db)
    - Goose is pretty simple to add in, just switching away from a schema.sql file to a schema directory with numbered files. And adding Up and Down comments to the sql files.
- Maybe look into using [Bun](https://bun.uptrace.dev/guide/#quick-start) as a lightweight ORM if getting sick of writing sql 
- Add constraint in Go code that only members of a league can submit songs to, or vote in, that league’s rounds (This is not constrained by sql schema)
