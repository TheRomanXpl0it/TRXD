# TODO
 - change name of chall containers from chall_PORT to chall_TID_CID
 - change name from _TID_CID to _CID_TID
 - change name of challs, network and so on, with $PROJECT_NAME prefix
 - single user mode, send user id instead of team id (ex scoreboard)
 - optimize badge queries
 - tests with removal of trxd-shared (and fault tolerance)
 - https://doc.traefik.io/traefik/expose/docker/
 - password reset for players
 - change image pfp system

## Frontend
 - convert .pngs to webp
 - Better tags handling in /challenges page
 - Use the same style with tabs in the "Update profile" page
 - Country selection sucks/does not work + it's buggy
 - Improve "Team join" page on smaller screens
 - Add custom buttons for first/second/third bloods?
 - Improve challenge creation procedure
 - Write tests
 - Add frontend testing to the gh workflow

## QoS
 - flush cache flag
 - verify no comunication between instances on standard default bridge network
 - change csrf token name
 - error generator for tests
 - split main.go flag function
 - resolve TODOs left into the code

## Alpha
 - config to start and end ctf
 - integration tests (for generic behaviour)
 - tests for distributed functioning
 - role configuraion endpoint
 - dropdown menu categories
 - dropdown menu container
 - DOCUMENTATION

## Extra
 - editable homepage & theme
 - deletable instances (from admin)
 - deletable submissions (from admin)
 - link to join team
 - endpoint + cli cmd to parse a category/challenge file (maybe yaml) with all infos
 - extract data for ctftime
 - scoreboard freeze (idea: use a table or a view to take a snapshot of the scoreboard, update it every time someone solves if not in freeze time (or just prefetch and compute if not exists))
 - telegram bot for first bloods
 - when finished endpoints, review needed data and revome `RETURNING *;` and similar when not needed
 - hash verifier flag or script

## Bonus Ideas
 - maybe set default starting points for challenges as global config
 - ctf stats page
 - flag submit that checks for all challenges and mark solved the one of the correct flag (credits to: midnightsun)
 - writeups for challs into the platform (visible after ctf ends)
 - signed flags (aes ecb with team id (4 byte) and chall id (4 byte))
 - tls instances: https://github.com/inconshreveable/slt
 - likes and dislikes for challs (only for who actually solved)
 - kick user from team
 - timeout on endpoints
 - optional mail server
 - endpoint to store image files (like badges and pfp)
