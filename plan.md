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
User(id, username)
Round(id, theme, users)
Song(id, round_id, submitter_user, title, artist, url)
#RoundVotes(id, round_id, user_id, )
Vote(round_id, song_id, user_id, rank)
